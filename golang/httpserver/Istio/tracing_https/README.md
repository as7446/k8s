1 安装jaeger
```
kubectl apply -f jaeger.yaml
kubectl edit configmap istio -n istio-system
set tracing.sampling=100

```
2 生成ssl证书
```
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=httserver Inc./CN=*.httserver.io' -keyout httserver.io.key -out httserver.io.crt
kubectl create -n istio-system secret tls httserver-credential --key=httserver.io.key --cert=httserver.io.crt
```
3 部署httpserver、istio-ingress
```
kubectl create ns tracing
kubectl label ns tracing istio-injection=enabled
kubectl -n tracing apply -f service0-deploy.yaml
kubectl -n tracing apply -f service1-deploy.yaml
kubectl -n tracing apply -f service2-deploy.yaml
kubectl apply -f istio-specs.yaml -n tracing
```
