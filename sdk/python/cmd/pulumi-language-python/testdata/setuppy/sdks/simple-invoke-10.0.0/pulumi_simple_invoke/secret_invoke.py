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

__all__ = [
    'SecretInvokeResult',
    'AwaitableSecretInvokeResult',
    'secret_invoke',
    'secret_invoke_output',
]

@pulumi.output_type
class SecretInvokeResult:
    def __init__(__self__, response=None, secret=None):
        if response and not isinstance(response, str):
            raise TypeError("Expected argument 'response' to be a str")
        pulumi.set(__self__, "response", response)
        if secret and not isinstance(secret, bool):
            raise TypeError("Expected argument 'secret' to be a bool")
        pulumi.set(__self__, "secret", secret)

    @_builtins.property
    @pulumi.getter
    def response(self) -> _builtins.str:
        return pulumi.get(self, "response")

    @_builtins.property
    @pulumi.getter
    def secret(self) -> _builtins.bool:
        return pulumi.get(self, "secret")


class AwaitableSecretInvokeResult(SecretInvokeResult):
    # pylint: disable=using-constant-test
    def __await__(self):
        if False:
            yield self
        return SecretInvokeResult(
            response=self.response,
            secret=self.secret)


def secret_invoke(secret_response: Optional[_builtins.bool] = None,
                  value: Optional[_builtins.str] = None,
                  opts: Optional[pulumi.InvokeOptions] = None) -> AwaitableSecretInvokeResult:
    """
    Use this data source to access information about an existing resource.
    """
    __args__ = dict()
    __args__['secretResponse'] = secret_response
    __args__['value'] = value
    opts = pulumi.InvokeOptions.merge(_utilities.get_invoke_opts_defaults(), opts)
    __ret__ = pulumi.runtime.invoke('simple-invoke:index:secretInvoke', __args__, opts=opts, typ=SecretInvokeResult).value

    return AwaitableSecretInvokeResult(
        response=pulumi.get(__ret__, 'response'),
        secret=pulumi.get(__ret__, 'secret'))
def secret_invoke_output(secret_response: Optional[pulumi.Input[_builtins.bool]] = None,
                         value: Optional[pulumi.Input[_builtins.str]] = None,
                         opts: Optional[Union[pulumi.InvokeOptions, pulumi.InvokeOutputOptions]] = None) -> pulumi.Output[SecretInvokeResult]:
    """
    Use this data source to access information about an existing resource.
    """
    __args__ = dict()
    __args__['secretResponse'] = secret_response
    __args__['value'] = value
    opts = pulumi.InvokeOutputOptions.merge(_utilities.get_invoke_opts_defaults(), opts)
    __ret__ = pulumi.runtime.invoke_output('simple-invoke:index:secretInvoke', __args__, opts=opts, typ=SecretInvokeResult)
    return __ret__.apply(lambda __response__: SecretInvokeResult(
        response=pulumi.get(__response__, 'response'),
        secret=pulumi.get(__response__, 'secret')))
