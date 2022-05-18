### 1 kubectl 设置Pod CPU、内存资源
```
  resources:
      limits:
        cpu: 200m
        memory: 512Mi
      requests:
        cpu: 100m
        memory: 256Mi
```
对已有deployment设置
```
kubectl set resources deployment nginx-deployment -c=nginx --limits=cpu=200m,memory=512Mi --requests=cpu=100m,memory=256Mi
```
* -c 容器名
