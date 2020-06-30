See the original document at https://github.com/nuweba/faasbenchmark.

# Setup

Follow these steps before building `faasbenchmark` Docker image.

1. Run `terraform` to create EKS clusters and deploy required FaaS platforms.

    ```sh
    # Start from the project root directory
    cd clusters/eks
    terraform apply

    cp kubeconfig_* ../../benchmarks/faastest/kubeconfigs

    cd ../..
    ./setup.sh knative -r # For example, this deploys Knative on its cluster
    ```

2. Copy an ingress URL of each platform to `provider/providers.json`. Repeat the following commands to see the URL for each platform. There will be only one value (hostname) on column `EXTERNAL-IP`.

    ```sh
    # Stay at the project root directory
    export KUBECONFIG="$(pwd)/clusters/eks/kubeconfig_knative"
    kubectl get svc -A
    ```

3. Build and push Docker images (if you want to use your own repository).

    ```sh
    # Go to the directory of the function you want to deploy inside `arsenal`

    # Knative
    # Run these and update `service.yml` with the correct image URL
    docker build -t <repo>/<tag> .
    docker push <repo>/<tag>

    # OpenFaaS
    # Run this and update `openfaas.yml` with the correct image tag
    faas-cli build -f openfaas.yml
    faas-cli push -f openfaas.yml
    ```

Build `faasbenchmark` by running:

```sh
docker build -t faasbenchmark .
```

Run a container by running:

```sh
docker run -it faasbenchmark /bin/bash
```

Inside the container, follow the required steps described below and run tests. The optional flag `-c` or `--config` is for a filepath to a config file containing values to override the `httpConfig` used by the test. Default values are used for missing fields. Keep the test's default behavior by not specifying the flag.

```sh
./faasbenchmark run knative BurstLvl1 -c config.json
```

Example of the config file:

TBA

## Platform Credentials

These need to be run before running tests on the corresponding platforms.

### Knative

```sh
export AWS_ACCESS_KEY_ID=<your key id> \
  AWS_SECRET_ACCESS_KEY=<your secret key>
```

### OpenFaaS

```sh
export AWS_ACCESS_KEY_ID=<your key id> \
  AWS_SECRET_ACCESS_KEY=<your secret key>

OPENFAAS_URL="http://$(kubectl get svc -n openfaas gateway-external -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' --kubeconfig /app/kubeconfigs/kubeconfig_openfaas):8080"
PASSWORD=$(kubectl -n openfaas get secret basic-auth -o jsonpath="{.data.basic-auth-password}" --kubeconfig /app/kubeconfigs/kubeconfig_openfaas | base64 --decode)
echo -n $PASSWORD | faas-cli login -g $OPENFAAS_URL -u admin --password-stdin
```

### OpenWhisk

```sh
export AWS_ACCESS_KEY_ID=<your key id> \
  AWS_SECRET_ACCESS_KEY=<your secret key>

wsk property set --apihost "$(kubectl get svc -n openwhisk owdev-nginx -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' --kubeconfig /app/kubeconfigs/kubeconfig_openwhisk):443"
wsk property set --auth 23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP
```

### Kubeless

```sh
export AWS_ACCESS_KEY_ID=<your key id> \
  AWS_SECRET_ACCESS_KEY=<your secret key> \
  KUBECONFIG=/app/kubeconfigs/kubeconfig_kubeless
```

### Fission

```sh
export AWS_ACCESS_KEY_ID=<your key id> \
  AWS_SECRET_ACCESS_KEY=<your secret key> \
  KUBECONFIG=/app/kubeconfigs/kubeconfig_fission
```
