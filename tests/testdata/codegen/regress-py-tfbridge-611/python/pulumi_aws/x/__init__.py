# coding=utf-8
# *** WARNING: this file was generated by test. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import builtins as _builtins
from .. import _utilities
import typing
from ._inputs import *
from . import outputs

# Make subpackages available:
if typing.TYPE_CHECKING:
    import pulumi_aws.x.iam as __iam
    iam = __iam
else:
    iam = _utilities.lazy_import('pulumi_aws.x.iam')

