
This repo demonstrates how to create a minimal control loop in Go that monitors
for successful Deployment rollouts and triggers a post-deploy test suite that
executes as a Job.

## Prerequisites

- Docker
- Kind
- Helm

## Setup

First build the container images for the three components:

```
docker build -t server --build-arg="CMD=server" .
docker build -t controller --build-arg="CMD=controller" .
docker build -t tests --build-arg="CMD=tests" --build-arg="MODE=test" .
```

Then create the cluster:

```
kind create cluster
```

Then load the container images into your kind cluster:

```
kind load docker-image server controller tests
```

Now deploy the helm chart:

```
helm upgrade test-after-deploy ./chart --install --create-namespace --namespace test-after-deploy
```

You can now edit the deployment and every time it finishes deploying it will
create a Job to run the tests.

Finally, cleanup:

```
kind delete cluster
```
