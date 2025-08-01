# coding=utf-8
# *** WARNING: this file was generated by pulumi-language-python. ***
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
from . import _utilities

__all__ = ['ExampleArgs', 'Example']

@pulumi.input_type
class ExampleArgs:
    def __init__(__self__, *,
                 map_map_union_property: Optional[pulumi.Input[Mapping[str, pulumi.Input[Mapping[str, pulumi.Input[Union[_builtins.str, Sequence[pulumi.Input[_builtins.str]]]]]]]]] = None,
                 string_or_integer_property: Optional[pulumi.Input[Union[_builtins.str, _builtins.int]]] = None):
        """
        The set of arguments for constructing a Example resource.
        """
        if map_map_union_property is not None:
            pulumi.set(__self__, "map_map_union_property", map_map_union_property)
        if string_or_integer_property is not None:
            pulumi.set(__self__, "string_or_integer_property", string_or_integer_property)

    @_builtins.property
    @pulumi.getter(name="mapMapUnionProperty")
    def map_map_union_property(self) -> Optional[pulumi.Input[Mapping[str, pulumi.Input[Mapping[str, pulumi.Input[Union[_builtins.str, Sequence[pulumi.Input[_builtins.str]]]]]]]]]:
        return pulumi.get(self, "map_map_union_property")

    @map_map_union_property.setter
    def map_map_union_property(self, value: Optional[pulumi.Input[Mapping[str, pulumi.Input[Mapping[str, pulumi.Input[Union[_builtins.str, Sequence[pulumi.Input[_builtins.str]]]]]]]]]):
        pulumi.set(self, "map_map_union_property", value)

    @_builtins.property
    @pulumi.getter(name="stringOrIntegerProperty")
    def string_or_integer_property(self) -> Optional[pulumi.Input[Union[_builtins.str, _builtins.int]]]:
        return pulumi.get(self, "string_or_integer_property")

    @string_or_integer_property.setter
    def string_or_integer_property(self, value: Optional[pulumi.Input[Union[_builtins.str, _builtins.int]]]):
        pulumi.set(self, "string_or_integer_property", value)


@pulumi.type_token("union:index:Example")
class Example(pulumi.CustomResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 map_map_union_property: Optional[pulumi.Input[Mapping[str, pulumi.Input[Mapping[str, pulumi.Input[Union[_builtins.str, Sequence[pulumi.Input[_builtins.str]]]]]]]]] = None,
                 string_or_integer_property: Optional[pulumi.Input[Union[_builtins.str, _builtins.int]]] = None,
                 __props__=None):
        """
        Create a Example resource with the given unique name, props, and options.
        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: Optional[ExampleArgs] = None,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        Create a Example resource with the given unique name, props, and options.
        :param str resource_name: The name of the resource.
        :param ExampleArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(ExampleArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 map_map_union_property: Optional[pulumi.Input[Mapping[str, pulumi.Input[Mapping[str, pulumi.Input[Union[_builtins.str, Sequence[pulumi.Input[_builtins.str]]]]]]]]] = None,
                 string_or_integer_property: Optional[pulumi.Input[Union[_builtins.str, _builtins.int]]] = None,
                 __props__=None):
        opts = pulumi.ResourceOptions.merge(_utilities.get_resource_opts_defaults(), opts)
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.id is None:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = ExampleArgs.__new__(ExampleArgs)

            __props__.__dict__["map_map_union_property"] = map_map_union_property
            __props__.__dict__["string_or_integer_property"] = string_or_integer_property
        super(Example, __self__).__init__(
            'union:index:Example',
            resource_name,
            __props__,
            opts)

    @staticmethod
    def get(resource_name: str,
            id: pulumi.Input[str],
            opts: Optional[pulumi.ResourceOptions] = None) -> 'Example':
        """
        Get an existing Example resource's state with the given name, id, and optional extra
        properties used to qualify the lookup.

        :param str resource_name: The unique name of the resulting resource.
        :param pulumi.Input[str] id: The unique provider ID of the resource to lookup.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        opts = pulumi.ResourceOptions.merge(opts, pulumi.ResourceOptions(id=id))

        __props__ = ExampleArgs.__new__(ExampleArgs)

        __props__.__dict__["map_map_union_property"] = None
        __props__.__dict__["string_or_integer_property"] = None
        return Example(resource_name, opts=opts, __props__=__props__)

    @_builtins.property
    @pulumi.getter(name="mapMapUnionProperty")
    def map_map_union_property(self) -> pulumi.Output[Optional[Mapping[str, Mapping[str, Any]]]]:
        return pulumi.get(self, "map_map_union_property")

    @_builtins.property
    @pulumi.getter(name="stringOrIntegerProperty")
    def string_or_integer_property(self) -> pulumi.Output[Optional[Any]]:
        return pulumi.get(self, "string_or_integer_property")

