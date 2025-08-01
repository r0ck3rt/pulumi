"""
@generated by mypy-protobuf.  Do not edit manually!
isort:skip_file
Copyright 2016-2018, Pulumi Corporation.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import builtins
import collections.abc
import google.protobuf.descriptor
import google.protobuf.internal.containers
import google.protobuf.internal.enum_type_wrapper
import google.protobuf.message
import google.protobuf.struct_pb2
import sys
import typing

if sys.version_info >= (3, 10):
    import typing as typing_extensions
else:
    import typing_extensions

DESCRIPTOR: google.protobuf.descriptor.FileDescriptor

class _EnforcementLevel:
    ValueType = typing.NewType("ValueType", builtins.int)
    V: typing_extensions.TypeAlias = ValueType

class _EnforcementLevelEnumTypeWrapper(google.protobuf.internal.enum_type_wrapper._EnumTypeWrapper[_EnforcementLevel.ValueType], builtins.type):
    DESCRIPTOR: google.protobuf.descriptor.EnumDescriptor
    ADVISORY: _EnforcementLevel.ValueType  # 0
    """Displayed to users, but does not block deployment."""
    MANDATORY: _EnforcementLevel.ValueType  # 1
    """Stops deployment, cannot be overridden."""
    DISABLED: _EnforcementLevel.ValueType  # 2
    """Disabled policies do not run during a deployment."""
    REMEDIATE: _EnforcementLevel.ValueType  # 3
    """Remediated policies actually fixes problems instead of issuing diagnostics."""

class EnforcementLevel(_EnforcementLevel, metaclass=_EnforcementLevelEnumTypeWrapper):
    """EnforcementLevel indicates the severity of a policy violation."""

ADVISORY: EnforcementLevel.ValueType  # 0
"""Displayed to users, but does not block deployment."""
MANDATORY: EnforcementLevel.ValueType  # 1
"""Stops deployment, cannot be overridden."""
DISABLED: EnforcementLevel.ValueType  # 2
"""Disabled policies do not run during a deployment."""
REMEDIATE: EnforcementLevel.ValueType  # 3
"""Remediated policies actually fixes problems instead of issuing diagnostics."""
global___EnforcementLevel = EnforcementLevel

@typing_extensions.final
class AnalyzerStackConfigureRequest(google.protobuf.message.Message):
    """`AnalyzerStackConfigureRequest` is the message for the stack configuration of the stack being analyzed."""

    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    @typing_extensions.final
    class ConfigEntry(google.protobuf.message.Message):
        DESCRIPTOR: google.protobuf.descriptor.Descriptor

        KEY_FIELD_NUMBER: builtins.int
        VALUE_FIELD_NUMBER: builtins.int
        key: builtins.str
        value: builtins.str
        def __init__(
            self,
            *,
            key: builtins.str = ...,
            value: builtins.str = ...,
        ) -> None: ...
        def ClearField(self, field_name: typing_extensions.Literal["key", b"key", "value", b"value"]) -> None: ...

    @typing_extensions.final
    class TagsEntry(google.protobuf.message.Message):
        DESCRIPTOR: google.protobuf.descriptor.Descriptor

        KEY_FIELD_NUMBER: builtins.int
        VALUE_FIELD_NUMBER: builtins.int
        key: builtins.str
        value: builtins.str
        def __init__(
            self,
            *,
            key: builtins.str = ...,
            value: builtins.str = ...,
        ) -> None: ...
        def ClearField(self, field_name: typing_extensions.Literal["key", b"key", "value", b"value"]) -> None: ...

    STACK_FIELD_NUMBER: builtins.int
    PROJECT_FIELD_NUMBER: builtins.int
    ORGANIZATION_FIELD_NUMBER: builtins.int
    DRY_RUN_FIELD_NUMBER: builtins.int
    CONFIG_SECRET_KEYS_FIELD_NUMBER: builtins.int
    CONFIG_FIELD_NUMBER: builtins.int
    TAGS_FIELD_NUMBER: builtins.int
    stack: builtins.str
    """The stack name being analyzed."""
    project: builtins.str
    """The project name of the stack being analyzed."""
    organization: builtins.str
    """The organization name of the stack being analyzed."""
    dry_run: builtins.bool
    """True if this is a preview/dry run."""
    @property
    def config_secret_keys(self) -> google.protobuf.internal.containers.RepeatedScalarFieldContainer[builtins.str]:
        """A list of configuration keys whose values should be treated as secrets."""
    @property
    def config(self) -> google.protobuf.internal.containers.ScalarMap[builtins.str, builtins.str]:
        """The configuration of the stack being analyzed."""
    @property
    def tags(self) -> google.protobuf.internal.containers.ScalarMap[builtins.str, builtins.str]:
        """Tags for the current stack."""
    def __init__(
        self,
        *,
        stack: builtins.str = ...,
        project: builtins.str = ...,
        organization: builtins.str = ...,
        dry_run: builtins.bool = ...,
        config_secret_keys: collections.abc.Iterable[builtins.str] | None = ...,
        config: collections.abc.Mapping[builtins.str, builtins.str] | None = ...,
        tags: collections.abc.Mapping[builtins.str, builtins.str] | None = ...,
    ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal["config", b"config", "config_secret_keys", b"config_secret_keys", "dry_run", b"dry_run", "organization", b"organization", "project", b"project", "stack", b"stack", "tags", b"tags"]) -> None: ...

global___AnalyzerStackConfigureRequest = AnalyzerStackConfigureRequest

@typing_extensions.final
class AnalyzerStackConfigureResponse(google.protobuf.message.Message):
    """`AnalyzerStackConfigureResponse` is the type of responses sent by a [](pulumirpc.Analyzer.ConfigureStack) call."""

    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    def __init__(
        self,
    ) -> None: ...

global___AnalyzerStackConfigureResponse = AnalyzerStackConfigureResponse

@typing_extensions.final
class AnalyzerHandshakeRequest(google.protobuf.message.Message):
    """`AnalyzerHandshakeRequest` is the type of requests sent as part of a [](pulumirpc.Analyzer.Handshake) call."""

    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    ENGINE_ADDRESS_FIELD_NUMBER: builtins.int
    ROOT_DIRECTORY_FIELD_NUMBER: builtins.int
    PROGRAM_DIRECTORY_FIELD_NUMBER: builtins.int
    engine_address: builtins.str
    """The gRPC address of the engine handshaking with the analyzer. At a minimum, this address will expose an instance
    of the [](pulumirpc.Engine) service.
    """
    root_directory: builtins.str
    """A *root directory* where the analyzer's binary, `PulumiPolicy.yaml`, or other identifying source code is located.
    In the event that the analyzer is *not* being booted by the engine (e.g. in the case that the engine has been
    asked to attach to an existing running analyzer instance via a host/port number), this field will be empty.
    """
    program_directory: builtins.str
    """A *program directory* in which the analyzer should execute. This is generally a subdirectory of the root
    directory, though this is not required. In the event that the analyzer is *not* being booted by the engine (e.g.
    in the case that the engine has been asked to attach to an existing running analyzer instance via a host/port
    number), this field will be empty.
    """
    def __init__(
        self,
        *,
        engine_address: builtins.str = ...,
        root_directory: builtins.str | None = ...,
        program_directory: builtins.str | None = ...,
    ) -> None: ...
    def HasField(self, field_name: typing_extensions.Literal["_program_directory", b"_program_directory", "_root_directory", b"_root_directory", "program_directory", b"program_directory", "root_directory", b"root_directory"]) -> builtins.bool: ...
    def ClearField(self, field_name: typing_extensions.Literal["_program_directory", b"_program_directory", "_root_directory", b"_root_directory", "engine_address", b"engine_address", "program_directory", b"program_directory", "root_directory", b"root_directory"]) -> None: ...
    @typing.overload
    def WhichOneof(self, oneof_group: typing_extensions.Literal["_program_directory", b"_program_directory"]) -> typing_extensions.Literal["program_directory"] | None: ...
    @typing.overload
    def WhichOneof(self, oneof_group: typing_extensions.Literal["_root_directory", b"_root_directory"]) -> typing_extensions.Literal["root_directory"] | None: ...

global___AnalyzerHandshakeRequest = AnalyzerHandshakeRequest

@typing_extensions.final
class AnalyzerHandshakeResponse(google.protobuf.message.Message):
    """`AnalyzerHandshakeResponse` is the type of responses sent by a [](pulumirpc.Analyzer.Handshake) call."""

    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    def __init__(
        self,
    ) -> None: ...

global___AnalyzerHandshakeResponse = AnalyzerHandshakeResponse

@typing_extensions.final
class AnalyzeRequest(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    TYPE_FIELD_NUMBER: builtins.int
    PROPERTIES_FIELD_NUMBER: builtins.int
    URN_FIELD_NUMBER: builtins.int
    NAME_FIELD_NUMBER: builtins.int
    OPTIONS_FIELD_NUMBER: builtins.int
    PROVIDER_FIELD_NUMBER: builtins.int
    type: builtins.str
    """the type token of the resource."""
    @property
    def properties(self) -> google.protobuf.struct_pb2.Struct:
        """the full properties to use for validation."""
    urn: builtins.str
    """the URN of the resource."""
    name: builtins.str
    """the name for the resource's URN."""
    @property
    def options(self) -> global___AnalyzerResourceOptions:
        """the resource options."""
    @property
    def provider(self) -> global___AnalyzerProviderResource:
        """the resource's provider."""
    def __init__(
        self,
        *,
        type: builtins.str = ...,
        properties: google.protobuf.struct_pb2.Struct | None = ...,
        urn: builtins.str = ...,
        name: builtins.str = ...,
        options: global___AnalyzerResourceOptions | None = ...,
        provider: global___AnalyzerProviderResource | None = ...,
    ) -> None: ...
    def HasField(self, field_name: typing_extensions.Literal["options", b"options", "properties", b"properties", "provider", b"provider"]) -> builtins.bool: ...
    def ClearField(self, field_name: typing_extensions.Literal["name", b"name", "options", b"options", "properties", b"properties", "provider", b"provider", "type", b"type", "urn", b"urn"]) -> None: ...

global___AnalyzeRequest = AnalyzeRequest

@typing_extensions.final
class AnalyzerResource(google.protobuf.message.Message):
    """AnalyzerResource defines the view of a Pulumi-managed resource as sent to Analyzers. The properties
    of the resource are specific to the type of analysis being performed. See the Analyzer
    service definition for more information.
    """

    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    @typing_extensions.final
    class PropertyDependenciesEntry(google.protobuf.message.Message):
        DESCRIPTOR: google.protobuf.descriptor.Descriptor

        KEY_FIELD_NUMBER: builtins.int
        VALUE_FIELD_NUMBER: builtins.int
        key: builtins.str
        @property
        def value(self) -> global___AnalyzerPropertyDependencies: ...
        def __init__(
            self,
            *,
            key: builtins.str = ...,
            value: global___AnalyzerPropertyDependencies | None = ...,
        ) -> None: ...
        def HasField(self, field_name: typing_extensions.Literal["value", b"value"]) -> builtins.bool: ...
        def ClearField(self, field_name: typing_extensions.Literal["key", b"key", "value", b"value"]) -> None: ...

    TYPE_FIELD_NUMBER: builtins.int
    PROPERTIES_FIELD_NUMBER: builtins.int
    URN_FIELD_NUMBER: builtins.int
    NAME_FIELD_NUMBER: builtins.int
    OPTIONS_FIELD_NUMBER: builtins.int
    PROVIDER_FIELD_NUMBER: builtins.int
    PARENT_FIELD_NUMBER: builtins.int
    DEPENDENCIES_FIELD_NUMBER: builtins.int
    PROPERTYDEPENDENCIES_FIELD_NUMBER: builtins.int
    type: builtins.str
    """the type token of the resource."""
    @property
    def properties(self) -> google.protobuf.struct_pb2.Struct:
        """the full properties to use for validation."""
    urn: builtins.str
    """the URN of the resource."""
    name: builtins.str
    """the name for the resource's URN."""
    @property
    def options(self) -> global___AnalyzerResourceOptions:
        """the resource options."""
    @property
    def provider(self) -> global___AnalyzerProviderResource:
        """the resource's provider."""
    parent: builtins.str
    """an optional parent URN that this child resource belongs to."""
    @property
    def dependencies(self) -> google.protobuf.internal.containers.RepeatedScalarFieldContainer[builtins.str]:
        """a list of URNs that this resource depends on."""
    @property
    def propertyDependencies(self) -> google.protobuf.internal.containers.MessageMap[builtins.str, global___AnalyzerPropertyDependencies]:
        """a map from property keys to the dependencies of the property."""
    def __init__(
        self,
        *,
        type: builtins.str = ...,
        properties: google.protobuf.struct_pb2.Struct | None = ...,
        urn: builtins.str = ...,
        name: builtins.str = ...,
        options: global___AnalyzerResourceOptions | None = ...,
        provider: global___AnalyzerProviderResource | None = ...,
        parent: builtins.str = ...,
        dependencies: collections.abc.Iterable[builtins.str] | None = ...,
        propertyDependencies: collections.abc.Mapping[builtins.str, global___AnalyzerPropertyDependencies] | None = ...,
    ) -> None: ...
    def HasField(self, field_name: typing_extensions.Literal["options", b"options", "properties", b"properties", "provider", b"provider"]) -> builtins.bool: ...
    def ClearField(self, field_name: typing_extensions.Literal["dependencies", b"dependencies", "name", b"name", "options", b"options", "parent", b"parent", "properties", b"properties", "propertyDependencies", b"propertyDependencies", "provider", b"provider", "type", b"type", "urn", b"urn"]) -> None: ...

global___AnalyzerResource = AnalyzerResource

@typing_extensions.final
class AnalyzerResourceOptions(google.protobuf.message.Message):
    """AnalyzerResourceOptions defines the options associated with a resource."""

    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    @typing_extensions.final
    class CustomTimeouts(google.protobuf.message.Message):
        """CustomTimeouts allows a user to be able to create a set of custom timeout parameters."""

        DESCRIPTOR: google.protobuf.descriptor.Descriptor

        CREATE_FIELD_NUMBER: builtins.int
        UPDATE_FIELD_NUMBER: builtins.int
        DELETE_FIELD_NUMBER: builtins.int
        create: builtins.float
        """The create resource timeout in seconds."""
        update: builtins.float
        """The update resource timeout in seconds."""
        delete: builtins.float
        """The delete resource timeout in seconds."""
        def __init__(
            self,
            *,
            create: builtins.float = ...,
            update: builtins.float = ...,
            delete: builtins.float = ...,
        ) -> None: ...
        def ClearField(self, field_name: typing_extensions.Literal["create", b"create", "delete", b"delete", "update", b"update"]) -> None: ...

    PROTECT_FIELD_NUMBER: builtins.int
    IGNORECHANGES_FIELD_NUMBER: builtins.int
    DELETEBEFOREREPLACE_FIELD_NUMBER: builtins.int
    DELETEBEFOREREPLACEDEFINED_FIELD_NUMBER: builtins.int
    ADDITIONALSECRETOUTPUTS_FIELD_NUMBER: builtins.int
    ALIASES_FIELD_NUMBER: builtins.int
    CUSTOMTIMEOUTS_FIELD_NUMBER: builtins.int
    PARENT_FIELD_NUMBER: builtins.int
    protect: builtins.bool
    """true if the resource should be marked protected."""
    @property
    def ignoreChanges(self) -> google.protobuf.internal.containers.RepeatedScalarFieldContainer[builtins.str]:
        """a list of property names to ignore during changes."""
    deleteBeforeReplace: builtins.bool
    """true if this resource should be deleted before replacement."""
    deleteBeforeReplaceDefined: builtins.bool
    """true if the deleteBeforeReplace property should be treated as defined even if it is false."""
    @property
    def additionalSecretOutputs(self) -> google.protobuf.internal.containers.RepeatedScalarFieldContainer[builtins.str]:
        """a list of output properties that should also be treated as secret, in addition to ones we detect."""
    @property
    def aliases(self) -> google.protobuf.internal.containers.RepeatedScalarFieldContainer[builtins.str]:
        """a list of additional URNs that shoud be considered the same."""
    @property
    def customTimeouts(self) -> global___AnalyzerResourceOptions.CustomTimeouts:
        """a config block that will be used to configure timeouts for CRUD operations."""
    parent: builtins.str
    """an optional parent URN that this child resource belongs to."""
    def __init__(
        self,
        *,
        protect: builtins.bool = ...,
        ignoreChanges: collections.abc.Iterable[builtins.str] | None = ...,
        deleteBeforeReplace: builtins.bool = ...,
        deleteBeforeReplaceDefined: builtins.bool = ...,
        additionalSecretOutputs: collections.abc.Iterable[builtins.str] | None = ...,
        aliases: collections.abc.Iterable[builtins.str] | None = ...,
        customTimeouts: global___AnalyzerResourceOptions.CustomTimeouts | None = ...,
        parent: builtins.str = ...,
    ) -> None: ...
    def HasField(self, field_name: typing_extensions.Literal["customTimeouts", b"customTimeouts"]) -> builtins.bool: ...
    def ClearField(self, field_name: typing_extensions.Literal["additionalSecretOutputs", b"additionalSecretOutputs", "aliases", b"aliases", "customTimeouts", b"customTimeouts", "deleteBeforeReplace", b"deleteBeforeReplace", "deleteBeforeReplaceDefined", b"deleteBeforeReplaceDefined", "ignoreChanges", b"ignoreChanges", "parent", b"parent", "protect", b"protect"]) -> None: ...

global___AnalyzerResourceOptions = AnalyzerResourceOptions

@typing_extensions.final
class AnalyzerProviderResource(google.protobuf.message.Message):
    """AnalyzerProviderResource provides information about a resource's provider."""

    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    TYPE_FIELD_NUMBER: builtins.int
    PROPERTIES_FIELD_NUMBER: builtins.int
    URN_FIELD_NUMBER: builtins.int
    NAME_FIELD_NUMBER: builtins.int
    type: builtins.str
    """the type token of the resource."""
    @property
    def properties(self) -> google.protobuf.struct_pb2.Struct:
        """the full properties to use for validation."""
    urn: builtins.str
    """the URN of the resource."""
    name: builtins.str
    """the name for the resource's URN."""
    def __init__(
        self,
        *,
        type: builtins.str = ...,
        properties: google.protobuf.struct_pb2.Struct | None = ...,
        urn: builtins.str = ...,
        name: builtins.str = ...,
    ) -> None: ...
    def HasField(self, field_name: typing_extensions.Literal["properties", b"properties"]) -> builtins.bool: ...
    def ClearField(self, field_name: typing_extensions.Literal["name", b"name", "properties", b"properties", "type", b"type", "urn", b"urn"]) -> None: ...

global___AnalyzerProviderResource = AnalyzerProviderResource

@typing_extensions.final
class AnalyzerPropertyDependencies(google.protobuf.message.Message):
    """AnalyzerPropertyDependencies describes the resources that a particular property depends on."""

    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    URNS_FIELD_NUMBER: builtins.int
    @property
    def urns(self) -> google.protobuf.internal.containers.RepeatedScalarFieldContainer[builtins.str]:
        """A list of URNs this property depends on."""
    def __init__(
        self,
        *,
        urns: collections.abc.Iterable[builtins.str] | None = ...,
    ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal["urns", b"urns"]) -> None: ...

global___AnalyzerPropertyDependencies = AnalyzerPropertyDependencies

@typing_extensions.final
class AnalyzeStackRequest(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    RESOURCES_FIELD_NUMBER: builtins.int
    @property
    def resources(self) -> google.protobuf.internal.containers.RepeatedCompositeFieldContainer[global___AnalyzerResource]: ...
    def __init__(
        self,
        *,
        resources: collections.abc.Iterable[global___AnalyzerResource] | None = ...,
    ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal["resources", b"resources"]) -> None: ...

global___AnalyzeStackRequest = AnalyzeStackRequest

@typing_extensions.final
class AnalyzeResponse(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    DIAGNOSTICS_FIELD_NUMBER: builtins.int
    @property
    def diagnostics(self) -> google.protobuf.internal.containers.RepeatedCompositeFieldContainer[global___AnalyzeDiagnostic]:
        """information about policy violations."""
    def __init__(
        self,
        *,
        diagnostics: collections.abc.Iterable[global___AnalyzeDiagnostic] | None = ...,
    ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal["diagnostics", b"diagnostics"]) -> None: ...

global___AnalyzeResponse = AnalyzeResponse

@typing_extensions.final
class AnalyzeDiagnostic(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    POLICYNAME_FIELD_NUMBER: builtins.int
    POLICYPACKNAME_FIELD_NUMBER: builtins.int
    POLICYPACKVERSION_FIELD_NUMBER: builtins.int
    DESCRIPTION_FIELD_NUMBER: builtins.int
    MESSAGE_FIELD_NUMBER: builtins.int
    ENFORCEMENTLEVEL_FIELD_NUMBER: builtins.int
    URN_FIELD_NUMBER: builtins.int
    policyName: builtins.str
    """Name of the violated policy."""
    policyPackName: builtins.str
    """Name of the policy pack the policy is in."""
    policyPackVersion: builtins.str
    """Version of the policy pack."""
    description: builtins.str
    """Description of policy rule. e.g., "encryption enabled." """
    message: builtins.str
    """Message to display on policy violation, e.g., remediation steps."""
    enforcementLevel: global___EnforcementLevel.ValueType
    """Severity of the policy violation."""
    urn: builtins.str
    """URN of the resource that violates the policy."""
    def __init__(
        self,
        *,
        policyName: builtins.str = ...,
        policyPackName: builtins.str = ...,
        policyPackVersion: builtins.str = ...,
        description: builtins.str = ...,
        message: builtins.str = ...,
        enforcementLevel: global___EnforcementLevel.ValueType = ...,
        urn: builtins.str = ...,
    ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal["description", b"description", "enforcementLevel", b"enforcementLevel", "message", b"message", "policyName", b"policyName", "policyPackName", b"policyPackName", "policyPackVersion", b"policyPackVersion", "urn", b"urn"]) -> None: ...

global___AnalyzeDiagnostic = AnalyzeDiagnostic

@typing_extensions.final
class Remediation(google.protobuf.message.Message):
    """Remediation is a single resource remediation result."""

    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    POLICYNAME_FIELD_NUMBER: builtins.int
    POLICYPACKNAME_FIELD_NUMBER: builtins.int
    POLICYPACKVERSION_FIELD_NUMBER: builtins.int
    DESCRIPTION_FIELD_NUMBER: builtins.int
    PROPERTIES_FIELD_NUMBER: builtins.int
    DIAGNOSTIC_FIELD_NUMBER: builtins.int
    policyName: builtins.str
    """Name of the policy that performed the remediation."""
    policyPackName: builtins.str
    """Name of the policy pack the transform is in."""
    policyPackVersion: builtins.str
    """Version of the policy pack."""
    description: builtins.str
    """Description of transform rule. e.g., "auto-tag resources." """
    @property
    def properties(self) -> google.protobuf.struct_pb2.Struct:
        """the transformed properties to use."""
    diagnostic: builtins.str
    """an optional warning diagnostic to emit, if a transform failed."""
    def __init__(
        self,
        *,
        policyName: builtins.str = ...,
        policyPackName: builtins.str = ...,
        policyPackVersion: builtins.str = ...,
        description: builtins.str = ...,
        properties: google.protobuf.struct_pb2.Struct | None = ...,
        diagnostic: builtins.str = ...,
    ) -> None: ...
    def HasField(self, field_name: typing_extensions.Literal["properties", b"properties"]) -> builtins.bool: ...
    def ClearField(self, field_name: typing_extensions.Literal["description", b"description", "diagnostic", b"diagnostic", "policyName", b"policyName", "policyPackName", b"policyPackName", "policyPackVersion", b"policyPackVersion", "properties", b"properties"]) -> None: ...

global___Remediation = Remediation

@typing_extensions.final
class RemediateResponse(google.protobuf.message.Message):
    """RemediateResponse contains a sequence of remediations applied, in order."""

    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    REMEDIATIONS_FIELD_NUMBER: builtins.int
    @property
    def remediations(self) -> google.protobuf.internal.containers.RepeatedCompositeFieldContainer[global___Remediation]:
        """the list of remediations that were applied."""
    def __init__(
        self,
        *,
        remediations: collections.abc.Iterable[global___Remediation] | None = ...,
    ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal["remediations", b"remediations"]) -> None: ...

global___RemediateResponse = RemediateResponse

@typing_extensions.final
class AnalyzerInfo(google.protobuf.message.Message):
    """AnalyzerInfo provides metadata about a PolicyPack inside an analyzer."""

    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    @typing_extensions.final
    class InitialConfigEntry(google.protobuf.message.Message):
        DESCRIPTOR: google.protobuf.descriptor.Descriptor

        KEY_FIELD_NUMBER: builtins.int
        VALUE_FIELD_NUMBER: builtins.int
        key: builtins.str
        @property
        def value(self) -> global___PolicyConfig: ...
        def __init__(
            self,
            *,
            key: builtins.str = ...,
            value: global___PolicyConfig | None = ...,
        ) -> None: ...
        def HasField(self, field_name: typing_extensions.Literal["value", b"value"]) -> builtins.bool: ...
        def ClearField(self, field_name: typing_extensions.Literal["key", b"key", "value", b"value"]) -> None: ...

    NAME_FIELD_NUMBER: builtins.int
    DISPLAYNAME_FIELD_NUMBER: builtins.int
    POLICIES_FIELD_NUMBER: builtins.int
    VERSION_FIELD_NUMBER: builtins.int
    SUPPORTSCONFIG_FIELD_NUMBER: builtins.int
    INITIALCONFIG_FIELD_NUMBER: builtins.int
    name: builtins.str
    """Name of the PolicyPack."""
    displayName: builtins.str
    """Pretty name for the PolicyPack."""
    @property
    def policies(self) -> google.protobuf.internal.containers.RepeatedCompositeFieldContainer[global___PolicyInfo]:
        """Metadata about policies contained in PolicyPack."""
    version: builtins.str
    """Version of the Policy Pack."""
    supportsConfig: builtins.bool
    """Whether the Policy Pack supports config."""
    @property
    def initialConfig(self) -> google.protobuf.internal.containers.MessageMap[builtins.str, global___PolicyConfig]:
        """Map of policy name to config."""
    def __init__(
        self,
        *,
        name: builtins.str = ...,
        displayName: builtins.str = ...,
        policies: collections.abc.Iterable[global___PolicyInfo] | None = ...,
        version: builtins.str = ...,
        supportsConfig: builtins.bool = ...,
        initialConfig: collections.abc.Mapping[builtins.str, global___PolicyConfig] | None = ...,
    ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal["displayName", b"displayName", "initialConfig", b"initialConfig", "name", b"name", "policies", b"policies", "supportsConfig", b"supportsConfig", "version", b"version"]) -> None: ...

global___AnalyzerInfo = AnalyzerInfo

@typing_extensions.final
class PolicyInfo(google.protobuf.message.Message):
    """PolicyInfo provides metadata about a policy within a Policy Pack."""

    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    NAME_FIELD_NUMBER: builtins.int
    DISPLAYNAME_FIELD_NUMBER: builtins.int
    DESCRIPTION_FIELD_NUMBER: builtins.int
    MESSAGE_FIELD_NUMBER: builtins.int
    ENFORCEMENTLEVEL_FIELD_NUMBER: builtins.int
    CONFIGSCHEMA_FIELD_NUMBER: builtins.int
    name: builtins.str
    """Name of the policy."""
    displayName: builtins.str
    """Pretty name for the policy."""
    description: builtins.str
    """Description of policy rule. e.g., "encryption enabled." """
    message: builtins.str
    """Message to display on policy violation, e.g., remediation steps."""
    enforcementLevel: global___EnforcementLevel.ValueType
    """Severity of the policy violation."""
    @property
    def configSchema(self) -> global___PolicyConfigSchema:
        """Config schema for the policy."""
    def __init__(
        self,
        *,
        name: builtins.str = ...,
        displayName: builtins.str = ...,
        description: builtins.str = ...,
        message: builtins.str = ...,
        enforcementLevel: global___EnforcementLevel.ValueType = ...,
        configSchema: global___PolicyConfigSchema | None = ...,
    ) -> None: ...
    def HasField(self, field_name: typing_extensions.Literal["configSchema", b"configSchema"]) -> builtins.bool: ...
    def ClearField(self, field_name: typing_extensions.Literal["configSchema", b"configSchema", "description", b"description", "displayName", b"displayName", "enforcementLevel", b"enforcementLevel", "message", b"message", "name", b"name"]) -> None: ...

global___PolicyInfo = PolicyInfo

@typing_extensions.final
class PolicyConfigSchema(google.protobuf.message.Message):
    """PolicyConfigSchema provides the schema for a policy's configuration."""

    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    PROPERTIES_FIELD_NUMBER: builtins.int
    REQUIRED_FIELD_NUMBER: builtins.int
    @property
    def properties(self) -> google.protobuf.struct_pb2.Struct:
        """JSON schema for each property."""
    @property
    def required(self) -> google.protobuf.internal.containers.RepeatedScalarFieldContainer[builtins.str]:
        """Required properties."""
    def __init__(
        self,
        *,
        properties: google.protobuf.struct_pb2.Struct | None = ...,
        required: collections.abc.Iterable[builtins.str] | None = ...,
    ) -> None: ...
    def HasField(self, field_name: typing_extensions.Literal["properties", b"properties"]) -> builtins.bool: ...
    def ClearField(self, field_name: typing_extensions.Literal["properties", b"properties", "required", b"required"]) -> None: ...

global___PolicyConfigSchema = PolicyConfigSchema

@typing_extensions.final
class PolicyConfig(google.protobuf.message.Message):
    """PolicyConfig provides configuration for a policy."""

    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    ENFORCEMENTLEVEL_FIELD_NUMBER: builtins.int
    PROPERTIES_FIELD_NUMBER: builtins.int
    enforcementLevel: global___EnforcementLevel.ValueType
    """Enforcement level of the policy."""
    @property
    def properties(self) -> google.protobuf.struct_pb2.Struct:
        """Configuration properties of the policy."""
    def __init__(
        self,
        *,
        enforcementLevel: global___EnforcementLevel.ValueType = ...,
        properties: google.protobuf.struct_pb2.Struct | None = ...,
    ) -> None: ...
    def HasField(self, field_name: typing_extensions.Literal["properties", b"properties"]) -> builtins.bool: ...
    def ClearField(self, field_name: typing_extensions.Literal["enforcementLevel", b"enforcementLevel", "properties", b"properties"]) -> None: ...

global___PolicyConfig = PolicyConfig

@typing_extensions.final
class ConfigureAnalyzerRequest(google.protobuf.message.Message):
    """ConfigureAnalyzerRequest provides configuration information to the analyzer."""

    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    @typing_extensions.final
    class PolicyConfigEntry(google.protobuf.message.Message):
        DESCRIPTOR: google.protobuf.descriptor.Descriptor

        KEY_FIELD_NUMBER: builtins.int
        VALUE_FIELD_NUMBER: builtins.int
        key: builtins.str
        @property
        def value(self) -> global___PolicyConfig: ...
        def __init__(
            self,
            *,
            key: builtins.str = ...,
            value: global___PolicyConfig | None = ...,
        ) -> None: ...
        def HasField(self, field_name: typing_extensions.Literal["value", b"value"]) -> builtins.bool: ...
        def ClearField(self, field_name: typing_extensions.Literal["key", b"key", "value", b"value"]) -> None: ...

    POLICYCONFIG_FIELD_NUMBER: builtins.int
    @property
    def policyConfig(self) -> google.protobuf.internal.containers.MessageMap[builtins.str, global___PolicyConfig]:
        """Map of policy name to config."""
    def __init__(
        self,
        *,
        policyConfig: collections.abc.Mapping[builtins.str, global___PolicyConfig] | None = ...,
    ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal["policyConfig", b"policyConfig"]) -> None: ...

global___ConfigureAnalyzerRequest = ConfigureAnalyzerRequest
