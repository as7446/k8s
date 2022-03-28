1 安装jaeger
kubectl apply -f jaeger.yaml
kubectl edit configmap istio -n istio-system
set tracing.sampling=100

kubectl create ns tracing
kubectl label ns tracing istio-injection=enabled
kubectl -n tracing apply -f service0-deploy.yaml
kubectl -n tracing apply -f service1-deploy.yaml
kubectl -n tracing apply -f service2-deploy.yaml
kubectl apply -f istio-specs.yaml -n tracing

