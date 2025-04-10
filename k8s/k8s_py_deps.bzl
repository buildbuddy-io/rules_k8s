load("@rules_python//python:pip.bzl", "pip_parse")

def k8s_py_deps():
    pip_parse(
        name = "pip",
        requirements_lock = Label("//:requirements_lock.txt"),
    )
