# Set the default Kubernetes context to kind
k8s_context("kind-kind")

# Path to cloud-provider-kind binary or image
CLOUD_PROVIDER_KIND_PATH = "cloud-provider-kind"  # or docker image name

# Run it in the background
local_resource(
    "run-cloud-provider-kind",
    "cloud-provider-kind",
    allow_parallel=True,
)

include('./route/Tiltfile')
include('./maintenance/Tiltfile')
