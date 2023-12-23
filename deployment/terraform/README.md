
```
$env:GOOGLE_APPLICATION_CREDENTIALS="C:\Users\MAV\IdeaProjects\game\deployment\terraform\credentials.json" 
terraform init
```


```
helm upgrade --install gitlab gitlab/gitlab --set global.hosts.domain=gl.kube.atpstealer.com --set global.hosts.externalIP=34.38.168.94 --set certmanager-issuer.email=forpubmail@gmail.com -n gitlab --create-namespace
```