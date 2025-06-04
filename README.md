# k8s-observability-demo

Learning cluster observability tools

Current stack:
- Docker Desktop - local K8s cluster
- Istio - service mesh for routing
- Prometheus - metric collection
- Grafana - dashboard and visualization
- Kiali - service mesh topology and heatlh viewer
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
./helm-install-services ../charts/app/service1 ../charts/app/service2
```
- Check pod status (can take a little bit):
```bash
kubectl get pods
NAME                        READY   STATUS    RESTARTS   AGE
service1-6c6bcf64dd-gzmv7   1/1     Running   0          30s
service2-bc87db7b8-npjqr    1/1     Running   0          30s
```
> **_TIP:_** If there are issues, use `kubectl describe pod <pod-name>` to get detailed pod info and troubleshoot

## Istio CLI (nice to have, not required)

- Install `istioctl`:
```bash
# via Homebrew
brew install istioctl
# via istio.io
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.26.1
export PATH=$PWD/bin:$PATH
```
- Confirm the installation:
```bash
istioctl version
Istio is not present in the cluster: no running Istio pods in namespace "istio-system"
client version: 1.26.1
```

## Istio

- I'd recommend following the docs since they are very thorough: https://istio.io/latest/docs/setup/install/helm/
- At the end make sure to label the namespace for sidecar injection:
```bash
kubectl label namespace observability-demo istio-injection=enabled --overwrite
```
- Then restart your deployments to trigger injection (or run helm uninstall and install scripts):
```bash
kubectl rollout restart deployment service1
kubectl rollout restart deployment service2
```
- And check your pods:
```bash
kubectl get pods
NAME                        READY   STATUS    RESTARTS   AGE
service1-6c6bcf64dd-d5vn4   2/2     Running   0          11s
service2-bc87db7b8-6xctq    2/2     Running   0          11s
```

## Kiali

- Once again I'd recommend following the docs: https://kiali.io/docs/installation/installation-guide/install-with-helm/ (install both the operator and server)
- Check Kiali pod:
```bash
kubectl get pods -n istio-system | grep kiali
kiali-549fddf87c-6l57x    1/1     Running   0          66s
```
- Check Kiali service:
```bash
kubectl get svc -n istio-system | grep kiali
kiali    ClusterIP   10.104.22.45    <none>        20001/TCP,9090/TCP                      4m7s
```
- View the Kiali dashboard via port-forwarding:
```bash
kubectl port-forward svc/kiali -n istio-system 20001:20001
```

## Prometheus

- Get Prometheus (or run the setup script):
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/kube-prometheus-stack -n istio-system --create-namespace
```
- Check that it worked:
```bash
kubectl get pods -n istio-system
NAME                                                     READY   STATUS             RESTARTS         AGE
alertmanager-prometheus-kube-prometheus-alertmanager-0   2/2     Running            0                166m
istiod-5d4b7d89bb-xjpj9                                  1/1     Running            1 (3h46m ago)    26h
kiali-58fb7f6674-m4g7v                                   1/1     Running            0                3h46m
prometheus-grafana-76cd8bb66b-hpkqp                      3/3     Running            0                166m
prometheus-kube-prometheus-operator-5ccb5b5fb9-jb5jg     1/1     Running            0                166m
prometheus-kube-state-metrics-74b7dc4795-zt5vr           1/1     Running            0                166m
prometheus-prometheus-kube-prometheus-prometheus-0       2/2     Running            0                166m
# Node exporter has issues running within Docker Desktop Kubernetes, but don't worry about it
prometheus-prometheus-node-exporter-679dw                0/1     CrashLoopBackOff   25 (2m36s ago)   166m
```
- By default, Prometheus will hit `/metrics` so make sure it's an existing endpoint (or you'll probably get a 404)
- You will also need to make sure your services have a `ServiceMonitor` available and deployed, check the service templates for an example
```bash
kubectl get servicemonitor
NAME                      AGE
service1-servicemonitor   13m
```
- View the Prometheus dashboard via port-forwarding:
```bash
kubectl port-forward svc/prometheus-kube-prometheus-prometheus -n istio-system 9090:9090
```
> **_TIP:_** Go has a nice prometheus library for for helping set up metrics, check the service code for an example

## Grafana 

- Grafana should be installed alongside prometheus stack!
- Check Grafana pod:
```bash
kubectl get pods -n istio-system | grep grafana
prometheus-grafana-76cd8bb66b-hpkqp                      3/3     Running
```
- View the Grafana dashboard via port-forwarding:
```bash
kubectl port-forward -n istio-system svc/prometheus-grafana 3000:80
```
- Create your own dashboard:
    1. From `Home` switch to the `Dashboard` tab
    2. Create your query, for our microservices we can do something like: `http_requests_total{job="service1"}` since that's what we defined in the service
    3. Or if you choose to use the builder you can do: `https_requests_total` for Metric and `job = service1` for Label filters
    4. Then you can be creative!

## TODO

- Get Istio Dashboards into Grafana so metrics can be properly viewed in Kiali
- Add cleanup scripts
- Add service2 metric monitoring, but that could also be a reader challenge so not sure how to go about that yet
- Add pictures!