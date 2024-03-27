servername=user
serverport=9801

# 申请配置
kubectl apply -f ${servername}.config.yaml
kubectl apply -f ${servername}.deployment.yaml
kubectl apply -f ${servername}.server.yaml
# 开放端口
kubectl port-forward service/api-server-user ${serverport}:${serverport} -n default