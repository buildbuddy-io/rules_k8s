"""Kubectl toolchain bzlmod extensions.
"""

load("@bazel_tools//tools/build_defs/repo:utils.bzl", "maybe")
load("//toolchains/kubectl:kubectl_configure.bzl", _kubectl_configure = "kubectl_configure")

#####################################################################
# toolchain
#####################################################################

toolchain_attrs = {
    "name": attr.string(
        mandatory = True,
        doc = "The name of the repo.",
    ),
    "build_srcs": attr.bool(
        doc = "Optional. Set to true to build kubectl from sources.",
        default = False,
        mandatory = False,
    ),
    "kubectl_path": attr.label(
        allow_single_file = True,
        mandatory = False,
        doc = "Optional. Path to a prebuilt custom kubectl binary file or" +
              " label. Can't be used together with attribute 'build_srcs'.",
    ),
}

_toolchain_configure_tag = tag_class(
    attrs = toolchain_attrs,
    doc = "configure kubectl toolchain",
)

def _impl(ctx):
    for mod in ctx.modules:
        for tag in mod.tags.configure:
            _kubectl_configure(
                name = tag.name,
                build_srcs = tag.build_srcs,
                kubectl_path = tag.kubectl_path,
            )

kubectl_toolchain = module_extension(
    implementation = _impl,
    tag_classes = {
        "configure": _toolchain_configure_tag,
    },
)
