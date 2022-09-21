version_settings(constraint='>=0.22.2')

# https://docs.tilt.dev/api.html#api.docker_build
# https://docs.tilt.dev/live_update_reference.html
custom_build(
    'jkremser/cert-manager-webhook-externaldns',
    'docker build -f ./Dockerfile.multistage -t $EXPECTED_REF . && \
     k3d image import $EXPECTED_REF',
    ['Dockerfile.multistage', 'main.go'],
    # live_update=[
    #     sync('./Dockerfile.multistage', '/main.go'),
    # ],
    skips_local_docker=True,
    disable_push=True,
)
# docker_build(
#     'absaoss/cert-manager-webhook-externaldns',
#     context='.',
#     dockerfile='./Dockerfile' 
# )

# k8s_yaml automatically creates resources in Tilt for the entities
# and will inject any images referenced in the Tiltfile when deploying
# https://docs.tilt.dev/api.html#api.k8s_yaml
k8s_yaml([
    './manifest/0-namespace.yaml',
    './manifest/rbac.yaml',
    './manifest/svc.yaml',
    './manifest/deployment.yaml',
    './manifest/issuers.yaml',
    './manifest/ca-cert.yaml',
    './manifest/serving-cert.yaml',
    './manifest/api-svc.yaml',
])

# # k8s_resource allows customization where necessary such as adding port forwards and labels
# # https://docs.tilt.dev/api.html#api.k8s_resource
k8s_resource(
    workload='cert-manager-externaldns-webhook',
    port_forwards=[
        port_forward(1443, 443, "app"),
  ]
)

# config.main_path is the absolute path to the Tiltfile being run
# there are many Tilt-specific built-ins for manipulating paths, environment variables, parsing JSON/YAML, and more!
# https://docs.tilt.dev/api.html#api.config.main_path
tiltfile_path = config.main_path

# print writes messages to the (Tiltfile) log in the Tilt UI
# the Tiltfile language is Starlark, a simplified Python dialect, which includes many useful built-ins
# config.tilt_subcommand makes it possible to only run logic during `tilt up` or `tilt down`
# https://github.com/bazelbuild/starlark/blob/master/spec.md#print
# https://docs.tilt.dev/api.html#api.config.tilt_subcommand
if config.tilt_subcommand == 'up':
    print("""
    \033[32m\033[32mHello World from cert-manager-webhook for externaldns!\033[0m

    Check the API on https://localhost:1443/

    """.format(tiltfile=tiltfile_path))