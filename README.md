# k8s-observability-demo

Learning cluster observability tools

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

## Test Services via Docker

```bash
docker compose up -d
curl -X POST -d '{"name":"Finish observability demo"}' -H "Content-Type: application/json" http://localhost:8080/tasks
# {"id":1,"name":"Finish observability demo"}
curl http://localhost:8080/tasks
# [{"id":1,"name":"Finish observability demo"}]
curl http://localhost:8081/summary
# {"task_count":1}
docker compose down
```

## Helm

- (Optional) Set up your chart scaffolding (skip if you already have the files):
```bash
helm create charts/service1
helm create charts/service2
```
- Run the helm install script:
```bash
chmod +x helm-install.sh
./helm-install ../charts/service1 ../charts/service2
```
- Check pod status (can take a little bit):
```bash
kubectl get pods
NAME                        READY   STATUS    RESTARTS   AGE
service1-6c6bcf64dd-gzmv7   1/1     Running   0          30s
service2-bc87db7b8-npjqr    1/1     Running   0          30s
```
> **_TIP:_** If there are issues, use `kubectl describe pod <pod-name>` to get detailed pod info and troubleshoot.