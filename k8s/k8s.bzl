load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:utils.bzl", "maybe")
load("//toolchains/kubectl:kubectl_configure.bzl", "kubectl_configure")
load(":with-defaults.bzl", _k8s_defaults = "k8s_defaults")

k8s_defaults = _k8s_defaults

def k8s_repositories():
    """Download dependencies of k8s rules."""

    # Register the default kubectl toolchain targets for supported platforms
    # note these work with the autoconfigured toolchain
    native.register_toolchains(
        "@io_bazel_rules_k8s//toolchains/kubectl:kubectl_linux_amd64_toolchain",
        "@io_bazel_rules_k8s//toolchains/kubectl:kubectl_linux_arm64_toolchain",
        "@io_bazel_rules_k8s//toolchains/kubectl:kubectl_linux_s390x_toolchain",
        "@io_bazel_rules_k8s//toolchains/kubectl:kubectl_macos_x86_64_toolchain",
        "@io_bazel_rules_k8s//toolchains/kubectl:kubectl_macos_arm64_toolchain",
        "@io_bazel_rules_k8s//toolchains/kubectl:kubectl_windows_toolchain",
    )

    maybe(
        http_archive,
        name = "io_bazel_rules_go",
        integrity = "sha256-kP6PtALe6VejdfPrhRFFW9c4x+1WJpX03RF6x9LYM7E=",
        urls = [
            "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.52.0/rules_go-v0.52.0.zip",
            "https://github.com/bazelbuild/rules_go/releases/download/v0.52.0/rules_go-v0.52.0.zip",
        ],
    )

    maybe(
        http_archive,
        name = "bazel_gazelle",
        integrity = "sha256-XYDmKnAxTznMdkwcPqqADFk2yfHqkWJQBiJ85NIM0IY=",
        urls = [
            "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.42.0/bazel-gazelle-v0.42.0.tar.gz",
            "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.42.0/bazel-gazelle-v0.42.0.tar.gz",
        ],
    )

    maybe(
        http_archive,
        name = "io_bazel_rules_docker",
        integrity = "sha256-A5jYBCeY4s0eJNYKrFYV6pnRSFekMRSilmr+xVF/Q+s=",
        strip_prefix = "rules_docker-d517338f5a4e29a11b6077ec39e533a518424b53",
        # This is our own fork of rules_docker with bzlmod support.
        # Diff: https://github.com/bazelbuild/rules_docker/compare/master...buildbuddy-io:rules_docker:sluongng/bzlmod-enable
        urls = ["https://github.com/buildbuddy-io/rules_docker/archive/d517338f5a4e29a11b6077ec39e533a518424b53.tar.gz"],
    )

    maybe(
        http_archive,
        name = "bazel_skylib",
        sha256 = "bc283cdfcd526a52c3201279cda4bc298652efa898b10b4db0837dc51652756f",
        urls = [
            "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.7.1/bazel-skylib-1.7.1.tar.gz",
            "https://github.com/bazelbuild/bazel-skylib/releases/download/1.7.1/bazel-skylib-1.7.1.tar.gz",
        ],
    )

    # WORKSPACE target to configure the kubectl tool
    maybe(
        kubectl_configure,
        name = "k8s_config",
        build_srcs = False,
    )
