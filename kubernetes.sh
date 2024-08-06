#! /bin/bash
# default
servername=user
kubernetesServername="$servername"-server # *.server.yaml.metadata.name
podName="$servername"-server # *.deployment.yaml.metadata.name
serverPort=9801
targetPort=9801
defaultNamespace="default" # 默认命名空间
imageSecretName="ali-image-key" # 阿里云镜像密钥, 修改请注意:所有deployment.imagePullSecrets.name 引用的名字
attempt=1
maxAttempts=5

function createDockerImageSecrete() {
    imageSecretName=$1
    namespace=$2
    defaultDockerImageServer="registry.cn-hangzhou.aliyuncs.com"
    echo -n -e "\033[1;32mPlease enter your docker server (default $defaultDockerImageServer): \033[0m"
    read -r dockerImageServer
    if [ "$dockerImageServer" = "" ]; then
      dockerImageServer=${defaultDockerImageServer}
    fi
    echo -n -e "\033[1;32mPlease enter your docker username: \033[0m"
    read -r dockerUsername
    echo -n -e "\033[1;32mPlease enter your docker password: \033[0m"
    read -r dockerPssword

    kubectl create secret docker-registry "$imageSecretName" \
    --docker-server=$dockerImageServer \
    --docker-username="$dockerUsername" \
    --docker-password="$dockerPssword" \
    -n "$namespace"

}

# istio
echo -n -e "\033[1;32mDo you want to inject istio-injection[y/n] (default y): \033[0m"
read -r selectInjectIstio
case $selectInjectIstio in
  N | n)
    selectInjectIstio=false;;
	*)
	  selectInjectIstio=true;;
esac
# port-forward
echo -n -e "\033[1;32mDo you want to use port-forward[y/n] (default y): \033[0m"
read -r selectPortForward
case $selectPortForward in
  N | n)
    selectPortForward=false;;
	*)
	  selectPortForward=true;;
esac
# namespace
echo -n -e "\033[1;32mPlease enter your namespace (default $defaultNamespace): \033[0m"
read -r namespace
if [ "$namespace" = "" ]; then
  namespace=${defaultNamespace}
fi

allNamespaces=$(kubectl get namespaces -o jsonpath="{.items[*].metadata.name}")
if echo "$allNamespaces" | grep -q "$namespace"; then
    echo "Namespace $namespace exists"
else
    echo "Create Namespace $namespace"
    kubectl create namespace $namespace
fi

# 申请密钥
secretExists=$(kubectl get secrets -n "$namespace" "$imageSecretName" --ignore-not-found=true --no-headers -o custom-columns=:.metadata.name 2>/dev/null)
# 检查Secret的名称是否与期望的匹配
if [[ "$secretExists" == "$imageSecretName" ]]; then
    echo -n -e "\033[1;32mDo you want reset $imageSecretName secret[y/n] (default n): \033[0m"
    read -r needRestSecret
    case $needRestSecret in
      Y | y) # 重制
        kubectl delete secret "$imageSecretName" -n $namespace
        createDockerImageSecrete "$imageSecretName" "$namespace"
    esac
else # 创建
    createDockerImageSecrete "$imageSecretName" "$namespace"
fi

# 申请配置与资源
kubectl apply -f ${servername}.secret.yaml -n $namespace
kubectl apply -f ${servername}.config.yaml -n $namespace


# 检查Deployment的状态，确保至少有一个可用的副本
deploymentStatus=$(kubectl get deployments "$podName" --namespace="$namespace" -o jsonpath='{.status.availableReplicas}')
if [ "$deploymentStatus" -eq 0 ]; then # 无成功运行货不存在申请pods
    kubectl apply -f ${servername}.deployment.yaml -n $namespace
    echo "Pods status:"
    kubectl get pods --selector=app=$podName --namespace="$namespace"
    exit 1
else # 有正在运行的pods重启
    kubectl rollout restart deployment $podName
    # 可选：列出Pods作为确认
    echo "Pods status:"
    kubectl get pods --selector=app=$podName --namespace="$namespace"
fi

kubectl apply -f ${servername}.server.yaml -n $namespace


if [ "$selectPortForward" = true ];then
  while true; do
    # 尝试执行的命令 开放端口
    kubectl port-forward service/${kubernetesServername} ${serverPort}:${targetPort} -n $namespace && exitCode=$? || exitCode=$?

    # 如果命令成功执行，退出循环
    if [ $exitCode -eq 0 ]; then
      break
    fi

    # 如果已达到最大尝试次数，退出循环
    if [ $attempt -eq $maxAttempts ]; then
      echo "Max attempts reached, giving up."
      exit 1
    fi

    # 等待一段时间后再次尝试
    sleepTime=5  # 等待5秒
    echo "Attempt failed. Retrying in $sleepTime seconds (Attempt $attempt of $maxAttempts)."
    sleep $sleepTime
    attempt=$((attempt + 1))
  done
fi