# 通过k8s启动

## 1. 创建私钥
为k8s创建镜像云密钥,免密拉私有镜像, 命令示例`[替换的变量](示例值)`
```shell
$ kubectl create secret docker-registry \
 [私钥名称](ali-image-key) \
 --docker-server=[镜像云地址](registry.cn-hangzhou.aliyuncs.com) 
 --docker-username=[用户名称] 
 --docker-password=[用户密码]
```

## 2. 上传镜像
示例中镜像上传到阿里云私有镜像，需要将`push.sh`脚本中的`docker login `命令修改为自己的仓库地址和账号，执行后输入版本后自动上传,版本自定义
```shell
$ /bin/bash ./push.sh
Please input server version
1.0.0
```

## 3. 修改配置
1. 修改`deployment.yaml`中的`image镜像地址`和`imagePullSecrets.name`第一步设置私有仓库的密钥名称
2. 修改隐私配置`secret.yaml`
```shell
$ cp [app].secret.yaml.tmp [app].secret.yaml`
```

## 4. 申请服务
`kubernetes.sh`脚本中的`kubectl port-forward`为本机调试开放端口访问。  
`serverPort`为对外访问接口,如设置9000,本机可以通过`127.0.0.1:9000`访问`http`服务。  
`targetPort`为服务定义端口,修改端口需要更新`server.yaml`文件中的`spec.ports.port`
```shell
$ /bin/bash ./kubernetes.sh
```