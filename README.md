# k8s-observability-demo
to learn cluster observability tools

Current stack idea:
- Docker Desktop - local K8s cluster
- Istio - service mesh for routing
- Prometheus - metric collection
- Grafana - dashboard and visualization
- Kiali (maybe) - service mesh topology and heatlh viewer
- Go Microservices - to generate traffic and metrics

## Local Setup

- Install Docker Desktop and enable Kubernetes in settings
- Ensure `kubectl` is installed and available in your terminal (`kubectl version` to check)
- Switch your kubectl context to Docker Desktop:
```bash
kubectl config use-context docker-desktop
``` 
- Verify Kubernetes is up and running:
```bash
kubectl get nodes
```
> **_TIP:_** If this errors, make sure Docker Desktop is open and Kubernetes is enabled
- For a clean workspace, create a dedicated namespace for this project:
```bash
kubectl create namespace observability-demo
kubectl get namespaces
```
- Set this namespace as the default for all kubectl commands in this context:
```bash
kubectl config set-context --current --namespace=observability-demo
```
- Double-check that the namespace is set:
```bash
kubectl config view --minify | grep namespace:
```