# coding=utf-8
# *** WARNING: this file was generated by pulumi-language-python. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import builtins as _builtins
from . import _utilities
import typing
# Export this package's modules as members:
from .component import *
from .custom import *
from .provider import *
_utilities.register(
    resource_modules="""
[
 {
  "pkg": "component-property-deps",
  "mod": "index",
  "fqn": "pulumi_component_property_deps",
  "classes": {
   "component-property-deps:index:Component": "Component",
   "component-property-deps:index:Custom": "Custom"
  }
 }
]
""",
    resource_packages="""
[
 {
  "pkg": "component-property-deps",
  "token": "pulumi:providers:component-property-deps",
  "fqn": "pulumi_component_property_deps",
  "class": "Provider"
 }
]
"""
)
