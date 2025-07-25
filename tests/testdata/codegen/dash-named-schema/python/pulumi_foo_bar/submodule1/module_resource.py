# coding=utf-8
# *** WARNING: this file was generated by test. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import builtins as _builtins
import warnings
import sys
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
if sys.version_info >= (3, 11):
    from typing import NotRequired, TypedDict, TypeAlias
else:
    from typing_extensions import NotRequired, TypedDict, TypeAlias
from .. import _utilities
from .. import _inputs as _root_inputs
from .. import outputs as _root_outputs

__all__ = ['ModuleResourceArgs', 'ModuleResource']

@pulumi.input_type
class ModuleResourceArgs:
    def __init__(__self__, *,
                 thing: Optional[pulumi.Input['_root_inputs.TopLevelArgs']] = None):
        """
        The set of arguments for constructing a ModuleResource resource.
        """
        if thing is not None:
            pulumi.set(__self__, "thing", thing)

    @_builtins.property
    @pulumi.getter
    def thing(self) -> Optional[pulumi.Input['_root_inputs.TopLevelArgs']]:
        return pulumi.get(self, "thing")

    @thing.setter
    def thing(self, value: Optional[pulumi.Input['_root_inputs.TopLevelArgs']]):
        pulumi.set(self, "thing", value)


@pulumi.type_token("foo-bar:submodule1:ModuleResource")
class ModuleResource(pulumi.CustomResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 thing: Optional[pulumi.Input[Union['_root_inputs.TopLevelArgs', '_root_inputs.TopLevelArgsDict']]] = None,
                 __props__=None):
        """
        Create a ModuleResource resource with the given unique name, props, and options.
        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: Optional[ModuleResourceArgs] = None,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        Create a ModuleResource resource with the given unique name, props, and options.
        :param str resource_name: The name of the resource.
        :param ModuleResourceArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(ModuleResourceArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 thing: Optional[pulumi.Input[Union['_root_inputs.TopLevelArgs', '_root_inputs.TopLevelArgsDict']]] = None,
                 __props__=None):
        opts = pulumi.ResourceOptions.merge(_utilities.get_resource_opts_defaults(), opts)
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.id is None:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = ModuleResourceArgs.__new__(ModuleResourceArgs)

            __props__.__dict__["thing"] = thing
        super(ModuleResource, __self__).__init__(
            'foo-bar:submodule1:ModuleResource',
            resource_name,
            __props__,
            opts)

    @staticmethod
    def get(resource_name: str,
            id: pulumi.Input[str],
            opts: Optional[pulumi.ResourceOptions] = None) -> 'ModuleResource':
        """
        Get an existing ModuleResource resource's state with the given name, id, and optional extra
        properties used to qualify the lookup.

        :param str resource_name: The unique name of the resulting resource.
        :param pulumi.Input[str] id: The unique provider ID of the resource to lookup.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        opts = pulumi.ResourceOptions.merge(opts, pulumi.ResourceOptions(id=id))

        __props__ = ModuleResourceArgs.__new__(ModuleResourceArgs)

        __props__.__dict__["thing"] = None
        return ModuleResource(resource_name, opts=opts, __props__=__props__)

    @_builtins.property
    @pulumi.getter
    def thing(self) -> pulumi.Output[Optional['_root_outputs.TopLevel']]:
        return pulumi.get(self, "thing")

