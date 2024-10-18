# Tiltfile

# Connect to the 'codema-hub-dev' cluster
k8s_context('codema-hub-dev')

# Set up a local registry
local_registry_name = 'codema-registry'
local_registry_port = 30051

# Deploy the local registry
k8s_yaml("k8s/registry.yaml")
k8s_resource(                 
  'registry',                 
  port_forwards='30051:5000', 
  labels=['external-service'],
)

def get_registry():
    return 'localhost:' + str(local_registry_port)

default_registry(get_registry())

# Build the CSS
local_resource(
    'build-css',
    cmd='npm run build-css',
    deps=['src/input.css', 'tailwind.config.js'],
    labels=['css']
)

# Build the Go server
local_resource(
    'build-go',
    cmd='go build -o codema-server .',
    deps=['*.go'],
    labels=['go']
)

# Deploy MongoDB
k8s_yaml('k8s/mongo.yaml')

# Deploy S3-compatible storage (e.g., MinIO)
k8s_yaml('k8s/minio.yaml')

k8s_yaml('k8s/mongo-express.yaml')

# Deploy the server
docker_build(
    get_registry() + '/codema-server',
    '.',
    dockerfile='Dockerfile',
)
k8s_yaml('k8s/codema-server.yaml')

# Port forward for local development
k8s_resource('codema-server', port_forwards='8090:8090')
k8s_resource('mongodb', port_forwards='27017:27017')
k8s_resource('minio', port_forwards='9000:9000')
k8s_resource('mongo-express', port_forwards='8081:8081')

