# coding=utf-8
# *** WARNING: this file was generated by test. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import builtins as _builtins
import pulumi
from enum import Enum

__all__ = [
    'EnumA',
]


@pulumi.type_token("myPkg:myMod/childA:EnumA")
class EnumA(_builtins.str, Enum):
    A1 = "a1"
    A2 = "a2"
