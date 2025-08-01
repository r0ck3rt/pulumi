// Copyright 2016-2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploy

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/pulumi/pulumi/pkg/v3/util/gsync"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/promise"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
)

const (
	// Dummy workerID for synchronous operations.
	synchronousWorkerID = -1
	infiniteWorkerID    = -2

	// Utility constant for easy debugging.
	stepExecutorLogLevel = 4
)

// StepApplyFailed is a sentinel error for errors that arise when step application fails.
// We (the step executor) are not responsible for reporting those errors so this sentinel ensures
// that we don't do so.
type StepApplyFailed struct {
	Err error
}

func (saf StepApplyFailed) Error() string {
	return fmt.Sprintf("step application failed: %s", saf.Err)
}

func (saf StepApplyFailed) Unwrap() error {
	return saf.Err
}

// The step executor operates in terms of "chains" and "antichains". A chain is set of steps that are totally ordered
// when ordered by dependency; each step in a chain depends directly on the step that comes before it. An antichain
// is a set of steps that is completely incomparable when ordered by dependency. The step executor is aware that chains
// must be executed serially and antichains can be executed concurrently.
//
// See https://en.wikipedia.org/wiki/Antichain for more complete definitions. The below type aliases are useful for
// documentation purposes.

// A Chain is a sequence of Steps that must be executed in the given order.
type chain = []Step

// An Antichain is a set of Steps that can be executed in parallel.
type antichain = []Step

// A CompletionToken is a token returned by the step executor that is completed when the chain has completed execution.
// Callers can use it to optionally wait synchronously on the completion of a chain.
type completionToken struct {
	channel chan bool
}

// Wait blocks until the completion token is signalled or until the given context completes, whatever occurs first.
func (c completionToken) Wait(ctx context.Context) {
	select {
	case <-c.channel:
	case <-ctx.Done():
	}
}

// incomingChain represents a request to the step executor to execute a chain.
type incomingChain struct {
	Chain          chain     // The chain we intend to execute
	CompletionChan chan bool // A completion channel to be closed when the chain has completed execution
}

// stepExecutor is the component of the engine responsible for taking steps and executing
// them, possibly in parallel if requested. The step generator operates on the granularity
// of "chains", which are sequences of steps that must be executed exactly in the given order.
// Chains are a simplification of the full dependency graph DAG within Pulumi programs. Since
// Pulumi language hosts can only invoke the resource monitor once all of their dependencies have
// resolved, we (the engine) can assume that any chain given to us by the step generator is already
// ready to execute.
type stepExecutor struct {
	// The deployment currently being executed.
	deployment *Deployment
	// Resources that have been created but are pending a RegisterResourceOutputs.
	pendingNews gsync.Map[resource.URN, Step]

	// True if errors should be ignored completely, without any handling or
	// reporting. This is used in the case of imports and refreshes. It is _not_
	// the same as ContinueOnError, which allows execution to continue in the face
	// of errors that may occur during updates. If both ignoreErrors and
	// ContinueOnError are set, ignoreErrors takes precedence.
	ignoreErrors bool

	// Lock protecting the running of workers. This can be used to synchronize with step executor.
	workerLock sync.RWMutex

	workers        sync.WaitGroup     // WaitGroup tracking the worker goroutines that are owned by this step executor.
	incomingChains chan incomingChain // Incoming chains that we are to execute

	ctx    context.Context    // cancellation context for the current deployment.
	cancel context.CancelFunc // CancelFunc that cancels the above context.

	// async promise indicating an error seen by the step executor, if multiple errors are seen this will only
	// record the first.
	sawError promise.CompletionSource[struct{}]

	erroredStepLock sync.RWMutex
	erroredSteps    []Step

	// ExecuteRegisterResourceOutputs will save the event for the stack resource so that the stack outputs
	// can be finalized at the end of the deployment. We do this so we can determine whether or not the
	// deployment succeeded. If there were errors, we update any stack outputs that were updated, but don't delete
	// any old outputs.
	stackOutputsEvent RegisterResourceOutputsEvent
}

//
// The stepExecutor communicates with a stepGenerator by listening to a channel. As the step generator
// generates new chains that need to be executed, the step executor will listen to this channel to execute
// those steps.
//

// Execute submits a Chain for asynchronous execution. The execution of the chain will begin as soon as there
// is a worker available to execute it.
func (se *stepExecutor) ExecuteSerial(chain chain) completionToken {
	// The select here is to avoid blocking on a send to se.incomingChains if a cancellation is pending.
	// If one is pending, we should exit early - we will shortly be tearing down the engine and exiting.

	completion := make(chan bool)
	select {
	case se.incomingChains <- incomingChain{Chain: chain, CompletionChan: completion}:
	case <-se.ctx.Done():
		close(completion)
	}

	return completionToken{channel: completion}
}

// Locks the step executor from executing any more steps. This is used to synchronize with the step executor.
func (se *stepExecutor) Lock() {
	se.workerLock.Lock()
}

// Unlocks the step executor to allow it to execute more steps. This is used to synchronize with the step executor.
func (se *stepExecutor) Unlock() {
	se.workerLock.Unlock()
}

func (se *stepExecutor) GetErroredSteps() []Step {
	se.erroredStepLock.RLock()
	defer se.erroredStepLock.RUnlock()
	return se.erroredSteps
}

// ExecuteParallel submits an antichain for parallel execution. All of the steps within the antichain are submitted for
// concurrent execution.
func (se *stepExecutor) ExecuteParallel(antichain antichain) completionToken {
	var wg sync.WaitGroup

	// ExecuteParallel is implemented in terms of ExecuteSerial - it executes each step individually and waits for all
	// of the steps to complete.
	wg.Add(len(antichain))
	for _, step := range antichain {
		tok := se.ExecuteSerial(chain{step})
		go func() {
			defer wg.Done()
			tok.Wait(se.ctx)
		}()
	}

	done := make(chan bool)
	go func() {
		wg.Wait()
		close(done)
	}()

	return completionToken{channel: done}
}

// ExecuteRegisterResourceOutputs services a RegisterResourceOutputsEvent synchronously on the calling goroutine.
func (se *stepExecutor) ExecuteRegisterResourceOutputs(e RegisterResourceOutputsEvent) error {
	return se.executeRegisterResourceOutputs(e, false /* errored */, false /* finalizingStackOutputs */)
}

func (se *stepExecutor) executeRegisterResourceOutputs(
	e RegisterResourceOutputsEvent,
	errored,
	finalizingStackOutputs bool,
) error {
	urn := e.URN()

	if finalizingStackOutputs {
		contract.Assertf(urn.QualifiedType() == resource.RootStackType, "expected a stack resource urn, got %v", urn)
	}

	// If we're not finalizing and we've received an event for the stack's outputs, save the event for finalization
	// later. We finalize stack outputs at the end of the deployment, so we can determine whether or not the
	// deployment succeeded. If the deployment was successful, we use the new stack outputs. If there was an error,
	// we only replace outputs that have new outputs, but to otherwise leave old outputs untouched.
	if !finalizingStackOutputs && urn.QualifiedType() == resource.RootStackType {
		se.stackOutputsEvent = e
		e.Done()
		return nil
	}

	// Look up the final state in the pending registration list.
	reg, has := se.pendingNews.Load(urn)
	if !has {
		return fmt.Errorf("cannot complete a resource '%v' whose registration isn't pending", urn)
	}
	contract.Assertf(reg != nil, "expected a non-nil resource step ('%v')", urn)
	se.pendingNews.Delete(urn)
	// Unconditionally set the resource's outputs to what was provided.  This intentionally overwrites whatever
	// might already be there, since otherwise "deleting" outputs would have no affect.
	outs := e.Outputs()
	se.log(synchronousWorkerID,
		"registered resource outputs %s: old=#%d, new=#%d", urn, len(reg.New().Outputs), len(outs))

	old := se.deployment.Olds()[urn]
	var oldOuts resource.PropertyMap
	if old != nil {
		oldOuts = old.Outputs
	}

	// If we're finalizing stack outputs and there was an error, the absence of an output can't safely be assumed to
	// mean it was deleted, so we keep old outputs, overwriting new ones.
	if finalizingStackOutputs && errored {
		outs = oldOuts.Copy()
		for k, v := range e.Outputs() {
			outs[k] = v
		}
	}

	// If a plan is present check that these outputs match what we recorded before
	if se.deployment.plan != nil {
		resourcePlan, ok := se.deployment.plan.ResourcePlans[urn]
		if !ok {
			return fmt.Errorf("no plan for resource %v", urn)
		}

		if err := resourcePlan.checkOutputs(oldOuts, outs); err != nil {
			return fmt.Errorf("resource violates plan: %w", err)
		}
	}

	reg.New().Lock.Lock()
	reg.New().Outputs = outs
	reg.New().Lock.Unlock()

	// If we're generating plans save these new outputs to the plan
	if se.deployment.opts.GeneratePlan {
		if resourcePlan, ok := se.deployment.newPlans.get(urn); ok {
			resourcePlan.Goal.OutputDiff = NewPlanDiff(oldOuts.Diff(outs))
			resourcePlan.Outputs = outs
		} else {
			return fmt.Errorf(
				"resource should already have a plan from when we called register resources [urn=%v]", urn)
		}
	}

	if !reg.Res().Custom {
		// Run the after hooks for Create/Update/Delete steps for
		// `ComponentResources `now that `RegisterResourceOutputs `is about to
		// complete.
		if s, ok := reg.(*CreateStep); ok {
			if err := s.Deployment().RunHooks(
				s.new.ResourceHooks[resource.AfterCreate],
				false, /* isBeforeHook */
				s.new.ID,
				s.new.URN,
				s.new.URN.Name(),
				s.Type(),
				s.new.Inputs,
				nil, /* oldInputs */
				s.new.Outputs,
				nil, /* oldOutputs */
			); err != nil {
				return err
			}
		} else if s, ok := reg.(*UpdateStep); ok {
			if err := s.Deployment().RunHooks(
				s.new.ResourceHooks[resource.AfterUpdate],
				false, /* isBeforeHook */
				s.new.ID,
				s.new.URN,
				s.new.URN.Name(),
				s.Type(),
				s.new.Inputs,
				s.old.Inputs,
				s.new.Outputs,
				s.old.Outputs,
			); err != nil {
				return err
			}
		}
	}

	// If there is an event subscription for finishing the resource, execute them.
	if e := se.deployment.events; e != nil {
		if eventerr := e.OnResourceOutputs(reg); eventerr != nil {
			se.log(synchronousWorkerID, "register resource outputs failed: %s", eventerr)

			// This is a bit of a kludge, but ExecuteRegisterResourceOutputs is an odd duck
			// in that it doesn't execute on worker goroutines. Arguably, it should, but today it's
			// not possible to express RegisterResourceOutputs as a step. We could 1) more generally allow
			// clients of stepExecutor to do work on worker threads by e.g. scheduling arbitrary callbacks
			// or 2) promote RRE to be step-like so that it can be scheduled as if it were a step. Neither
			// of these are particularly appealing right now.
			outErr := fmt.Errorf("resource complete event returned an error: %w", eventerr)
			diagMsg := diag.RawMessage(reg.URN(), outErr.Error())
			se.deployment.Diag().Errorf(diagMsg)
			se.cancelDueToError(eventerr, nil)
			return nil
		}
	}
	if !finalizingStackOutputs {
		e.Done()
	}
	return nil
}

// Errored returns whether or not this step executor saw a step whose execution ended in failure.
func (se *stepExecutor) Errored() error {
	// See if the sawError promise has been rejected yet
	_, err, _ := se.sawError.Promise().TryResult()
	// err will be nil if the promise has not been rejected yet
	return err
}

// SignalCompletion signals to the stepExecutor that there are no more chains left to execute. All worker
// threads will terminate as soon as they retire all of the work they are currently executing.
func (se *stepExecutor) SignalCompletion() {
	close(se.incomingChains)
}

// WaitForCompletion blocks the calling goroutine until the step executor completes execution of all in-flight
// chains.
func (se *stepExecutor) WaitForCompletion() {
	se.log(synchronousWorkerID, "StepExecutor.waitForCompletion(): waiting for worker threads to exit")
	se.workers.Wait()
	se.log(synchronousWorkerID, "StepExecutor.waitForCompletion(): worker threads all exited")
}

//
// As calls to `Execute` submit chains for execution, some number of worker goroutines will continuously
// read from `incomingChains` and execute any chains that are received. The core execution logic is in
// the next few functions.
//

// executeChain executes a chain, one step at a time. If any step in the chain fails to execute, or if the
// context is canceled, the chain stops execution.
func (se *stepExecutor) executeChain(workerID int, chain chain) {
	for _, step := range chain {
		select {
		case <-se.ctx.Done():
			se.log(workerID, "step %v on %v canceled", step.Op(), step.URN())
			return
		default:
		}

		// Take the work lock before executing the step, this uses the "read" side of the lock because we're ok with as
		// many workers as possible executing steps in parallel.
		se.workerLock.RLock()
		err := se.executeStep(workerID, step)
		// Regardless of error we need to release the lock here.
		se.workerLock.RUnlock()

		if err != nil {
			se.log(workerID, "step %v on %v failed, signalling cancellation", step.Op(), step.URN())
			se.cancelDueToError(err, step)

			var saf StepApplyFailed
			if !errors.As(err, &saf) {
				// Step application errors are recorded by the OnResourceStepPost callback. This is confusing,
				// but it means that at this level we shouldn't be logging any errors that came from there.
				//
				// The StepApplyFailed sentinel signals that the error that failed this chain was a step apply
				// error and that we shouldn't log it. Everything else should be logged to the diag system as usual.
				diagMsg := diag.RawMessage(step.URN(), err.Error())
				se.deployment.Diag().Errorf(diagMsg)
			}
			return
		}
	}
}

func (se *stepExecutor) cancelDueToError(err error, step Step) {
	set := se.sawError.Reject(err)
	if !set {
		logging.V(10).Infof("StepExecutor already recorded an error then saw: %v", err)
	}

	continueOnError := se.deployment.opts.ContinueOnError
	// We special case deleteProtectedError to always continue past them, this allows us to report back all
	// protect delete errors in a single run. We know we didn't attempt to do any actual delete, _and_ as
	// these are deletes there's nothing replied to the user program, it's safe to continue past them. There
	// are no observable effects in the user program, its just the CLI will be able to continue past and
	// report any other issues.
	if errors.As(err, &deleteProtectedError{}) {
		continueOnError = true
	}

	if se.ignoreErrors {
		// Do nothing.
	} else if continueOnError {
		step.Fail()
		// Record the failure, but allow the deployment to continue.
		se.erroredStepLock.Lock()
		defer se.erroredStepLock.Unlock()
		se.erroredSteps = append(se.erroredSteps, step)
	} else {
		se.cancel()
	}
}

//
// The next few functions are responsible for executing individual steps. The basic flow of step
// execution is
//   1. The pre-step event is raised, if there are any attached callbacks to the engine
//   2. If successful, the step is executed (if not a preview)
//   3. The post-step event is raised, if there are any attached callbacks to the engine
//
// The pre-step event returns an interface{}, which is some arbitrary context that must be passed
// verbatim to the post-step event.
//

// executeStep executes a single step, returning true if the step execution was successful and
// false if it was not.
func (se *stepExecutor) executeStep(workerID int, step Step) error {
	var payload interface{}
	events := se.deployment.events

	// DiffSteps are special, we just use them for step worker parallelism but they shouldn't be passed to the rest of
	// the system.
	_, isDiff := step.(*DiffStep)

	if events != nil && !isDiff {
		var err error
		payload, err = events.OnResourceStepPre(step)
		if err != nil {
			se.log(workerID, "step %v on %v failed pre-resource step: %v", step.Op(), step.URN(), err)
			return fmt.Errorf("pre-step event returned an error: %w", err)
		}
	}

	return se.continueExecuteStep(payload, workerID, step)
}

func (se *stepExecutor) continueExecuteStep(payload interface{}, workerID int, step Step) error {
	events := se.deployment.events

	se.log(workerID, "applying step %v on %v (preview %v)", step.Op(), step.URN(), se.deployment.opts.DryRun)
	status, stepComplete, err := step.Apply()

	// DiffSteps are special, we just use them for step worker parallelism but they shouldn't be passed to the rest of
	// the system.
	if _, isDiff := step.(*DiffStep); isDiff {
		return nil
	}

	if err == nil {
		// If we have a state object, and this is a create or update, remember it, as we may need to update it later.
		if step.Logical() && step.New() != nil {
			se.pendingNews.Store(step.URN(), step)
		}
	}

	// Ensure that any secrets properties in the output are marked as such and that the resource is tracked in the set
	// of registered resources. We skip this for replace steps because while they _do_ have a "new" side to them that
	// state may have already been added to the snapshot manager (in the case of create before delete replacements
	// because the Create step is run before the Replace step) and mutating the state again causes dataraces (see
	// https://github.com/pulumi/pulumi/issues/14994).
	if step.New() != nil && step.Op() != OpReplace {
		newState := step.New()
		newState.Lock.Lock()

		for _, k := range newState.AdditionalSecretOutputs {
			if k == "id" {
				se.deployment.Diag().Warningf(&diag.Diag{
					URN:     step.URN(),
					Message: "The 'id' property cannot be made secret. See pulumi/pulumi#2717 for more details.",
				})
			} else {
				if v, has := newState.Outputs[k]; has && !v.IsSecret() {
					newState.Outputs[k] = resource.MakeSecret(v)
				} else if !has { //nolint:staticcheck // https://github.com/pulumi/pulumi/issues/9926
					// TODO (https://github.com/pulumi/pulumi/issues/9926): We want to re-enable this warning
					// but it requires that providers always return back _every_ output even in preview. We
					// might need to add a new "unset" PropertyValue to do this as there might be optional
					// secret outputs and the engine needs to be able to tell the difference between "this
					// isn't a valid output of the resource" and "this value just hasn't been set in this
					// instance". Arguably for user side additionalSecretOutputs that distinction probably
					// doesn't matter (if you ask for an optional output to be made secret but then the
					// provider doesn't return it maybe you want the warning that nothing is actually being
					// affected?). But for SDK generated we always send the same list and the user doesn't
					// control it so we need to make sure that if there is an optional output that this
					// warning doesn't get triggered.

					// User asked us to make k a secret, but we don't have a property k. This is probably a
					// mistake (mostly likely due to casing, eg my_prop vs myProp) but warn the user so they know
					// the key didn't do anything.
					// msg := fmt.Sprintf("Could not find property '%s' listed in additional secret outputs.", k)
					// se.deployment.Diag().Warningf(diag.RawMessage(step.URN(), msg))
				}
			}
		}

		// If an input secret is potentially leaked as an output, preemptively mark it as secret.
		for k, out := range newState.Outputs {
			if !out.IsSecret() {
				in, has := newState.Inputs[k]
				if !has {
					continue
				}
				if in.IsSecret() {
					newState.Outputs[k] = resource.MakeSecret(out)
				}
			}
		}

		newState.Lock.Unlock()

		// At this point we'll track the resource so that it can be looked up later by e.g. consumers resolving resource
		// references using pulumi:pulumi:getResource. If it's a resource managed by Pulumi (i.e. it's the result of a
		// resource registration with a goal state), we'll record it in news. If it's not managed by Pulumi (i.e. it's the
		// result of a Read, perhaps caused by a .get in an SDK, for instance), we'll record it in reads.
		if step.Op() == OpRead || step.Op() == OpReadReplacement {
			se.deployment.reads.Store(newState.URN, newState)
		} else {
			// Read tells us a state was for an external, but we can't check deployment.goals to see if a
			// resource is managed because untargetted but removed from program resources are still considered
			// managed but don't have a goal state.
			se.deployment.news.Store(newState.URN, newState)
		}

		// If we're generating plans update the resource's outputs in the generated plan.
		if se.deployment.opts.GeneratePlan {
			if resourcePlan, ok := se.deployment.newPlans.get(newState.URN); ok {
				resourcePlan.Outputs = newState.Outputs
			}
		}
	}

	// Executes view steps serially, in the order they were published to the resource status server.
	executeViewSteps := func() error {
		steps := se.deployment.resourceStatus.ReleaseToken(step.URN())
		// Execute the steps in the order they were published.
		var errs []error
		for _, step := range steps {
			if err := se.continueExecuteStep(step.payload, workerID, step.step); err != nil {
				errs = append(errs, err)
			}
		}
		return errors.Join(errs...)
	}

	_, isDelete := step.(*DeleteStep)

	// For delete operations, execute view steps before the post-step event. This way, any delete steps for child views
	// are executed and saved to state before OnResourceStepPost for this resource saves the results of the step in the
	// snapshot.
	if isDelete {
		if err := executeViewSteps(); err != nil {
			return err
		}
	}

	if events != nil {
		if postErr := events.OnResourceStepPost(payload, step, status, err); postErr != nil {
			se.log(workerID, "step %v on %v failed post-resource step: %v", step.Op(), step.URN(), postErr)
			return fmt.Errorf("post-step event returned an error: %w", postErr)
		}
	}

	// For non-delete operations, execute view steps after the post-step event.
	if !isDelete {
		if err := executeViewSteps(); err != nil {
			return err
		}
	}

	// Calling stepComplete allows steps that depend on this step to continue. OnResourceStepPost saved the results
	// of the step in the snapshot, so we are ready to go.
	if stepComplete != nil {
		se.log(workerID, "step %v on %v retired", step.Op(), step.URN())
		stepComplete()
	}

	if err != nil {
		se.log(workerID, "step %v on %v failed with an error: %v", step.Op(), step.URN(), err)
		return StepApplyFailed{err}
	}

	return nil
}

// log is a simple logging helper for the step executor.
func (se *stepExecutor) log(workerID int, msg string, args ...interface{}) {
	if logging.V(stepExecutorLogLevel) {
		message := fmt.Sprintf(msg, args...)
		logging.V(stepExecutorLogLevel).Infof("StepExecutor worker(%d): %s", workerID, message)
	}
}

//
// The step executor owns a number of goroutines that it considers to be "workers", responsible for
// executing steps. By default, as we ease into the waters of parallelism, there is at most one worker
// active.
//
// Workers continuously pull from se.incomingChains, executing chains as they are provided to the executor.
// There are two reasons why a worker would exit:
//
//  1. A worker exits if se.ctx is canceled. There are two ways that se.ctx gets canceled: first, if there is
//     a step error in another worker, it will cancel the context. Second, if the deployment executor experiences an
//     error when generating steps or doing pre or post-step events, it will cancel the context.
//  2. A worker exits if it experiences an error when running a step.
//

// worker is the base function for all step executor worker goroutines. It continuously polls for new chains
// and executes any that it gets from the channel. If `launchAsync` is true, worker launches a new goroutine
// that will execute the chain so that the execution continues asynchronously and this worker can proceed to
// the next chain.
func (se *stepExecutor) worker(workerID int, launchAsync bool) {
	se.log(workerID, "worker coming online")
	defer se.workers.Done()

	oneshotWorkerID := 0
	for {
		se.log(workerID, "worker waiting for incoming chains")
		select {
		case request := <-se.incomingChains:
			contract.Assertf(request.CompletionChan != nil || request.Chain == nil,
				"worker received a non-nil chain with a nil completion channel: %v", request.Chain,
			)

			if request.CompletionChan == nil {
				se.log(workerID, "worker received nil completion channel, exiting")
				return
			}

			se.log(workerID, "worker received chain for execution")
			if !launchAsync {
				se.executeChain(workerID, request.Chain)
				close(request.CompletionChan)
				continue
			}

			// If we're launching asynchronously, make up a new worker ID for this new oneshot worker and record its
			// launch with our worker wait group.
			se.workers.Add(1)
			newWorkerID := oneshotWorkerID
			go func() {
				defer se.workers.Done()
				se.log(newWorkerID, "launching oneshot worker")
				se.executeChain(newWorkerID, request.Chain)
				close(request.CompletionChan)
			}()

			oneshotWorkerID++
		case <-se.ctx.Done():
			se.log(workerID, "worker exiting due to cancellation")
			return
		}
	}
}

func newStepExecutor(
	ctx context.Context,
	cancel context.CancelFunc,
	deployment *Deployment,
	ignoreErrors bool,
) *stepExecutor {
	exec := &stepExecutor{
		deployment:     deployment,
		ignoreErrors:   ignoreErrors,
		incomingChains: make(chan incomingChain),
		ctx:            ctx,
		cancel:         cancel,
	}

	// If we're being asked to run as parallel as possible, spawn a single worker that launches chain executions
	// asynchronously.
	if deployment.opts.InfiniteParallelism() {
		exec.workers.Add(1)
		go exec.worker(infiniteWorkerID, true /*launchAsync*/)
		return exec
	}

	// Otherwise, launch a worker goroutine for each degree of parallelism.
	fanout := deployment.opts.DegreeOfParallelism()
	for i := 0; i < int(fanout); i++ {
		exec.workers.Add(1)
		go exec.worker(i, false /*launchAsync*/)
	}

	return exec
}
