Ingress 

```agsl
helm install ingress ingress-nginx/ingress-nginx --namespace ingress-nginx  --create-namespace  --set controller.service.loadBalancerIP=34.79.170.192
```
34.79.170.192 is Region IP

Cluster issuer
```agsl
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.8.2/cert-manager.yaml
kubectl apply -f issuer-lets-encrypt.yaml
```

GitLab runner access
```
kubectl apply -f gitlab-rolebinding.yaml
```