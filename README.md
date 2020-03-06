# Session Prometheus

```bash
# minikube
minikube start --kubernetes-version=v1.16.0 --memory=3g --bootstrapper=kubeadm --extra-config=kubelet.authentication-token-webhook=true --extra-config=kubelet.authorization-mode=Webhook --extra-config=scheduler.address=0.0.0.0 --extra-config=controller-manager.address=0.0.0.0
minikube addons disable metrics-server
minikube addons enable ingress

# get minikube ip
minikube ip

# add these lines to /etc/hosts file.
192.168.64.2 prometheus.dev.com # change the ip 192.168.64.2 with your own minikube ip.
192.168.64.2 alertmanager.dev.com
192.168.64.2 grafana.dev.com
192.168.64.2 session-prometheus.dev.com

# Build app.
eval $(minikube docker-env) # to build container for minikube.
docker build --rm=false -t turgayozgur/session-prometheus:0.1 .

# Create namespaces
kubectl create namespace dev && kubectl create namespace monitoring

# Install or upgrade the prometheus operator.
helm upgrade -i prom -n monitoring --version 8.5.11 -f deployment/prometheus-values.yml stable/prometheus-operator

# Push the container to the k8s cluster.
kubectl apply -f deployment/app.yml
# Change image version if required.
kubectl set image deployment/session-prometheus -n dev session-prometheus=turgayozgur/session-prometheus:0.2

# Enable service monitor.
kubectl apply -f deployment/servicemonitor.yml

# Make some load.
kubectl run -i --restart=Never --rm fortio --image=fortio/fortio -n dev -- load -qps 2 -c 2 -t 0 -H "Content-Type: application/json" --payload '{"total": 200, "bankType": "A", "recordedcardkey": "testcard1"}' http://session-prometheus.dev.com/pay

# Enable prometheus rule.
kubectl apply -f deployment/prometheusrule.yml
```
