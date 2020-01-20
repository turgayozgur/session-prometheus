# Session Prometheus

```bash
# minikube
minikube start
minikube addons enable ingress

# Build app.
eval $(minikube docker-env) # to build container for minikube.
docker build --rm=false -t turgayozgur/session-prometheus:0.1 .

# Create namespaces
kubectl create namespace dev
kubectl create namespace monitoring

# Install or upgrade the prometheus operator.
helm upgrade -i prom -n monitoring --version 8.5.11 -f deployment/prometheus-values.yml stable/prometheus-operator

# Push the container to the k8s cluster.
kubectl apply -f deployment/app.yml
kubectl port-forward -n dev svc/session-prometheus 8080:80
# Change image version if required.
kubectl set image deployment/session-prometheus -n dev session-prometheus=turgayozgur/session-prometheus:0.2

# Enable service monitor.
kubectl apply -f deployment/servicemonitor.yml

# Make some load.
kubectl run -i --restart=Never --rm fortio --image=fortio/fortio -n dev -- load -qps 4 -c 2 -t 0 -H "Content-Type: application/json" --payload '{"total": 200, "bankType": "A", "recordedcardkey": "testcard1"}' http://session-prometheus.dev/pay

# Enable prometheus rule.
kubectl apply -f deployment/prometheusrule.yml
```
