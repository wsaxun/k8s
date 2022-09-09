# kubernetes集群部署工具
- - -  

## 简介  
- - -   

本项目为kubernetes集群***一键部署工具***，通过配置yaml文件定义集群组件及分布，底层通过*Ansible*执行  
实际安装动作，目前只支持centos，redhat系列系统  
  
支持的kubernetes版本: *v1.20* / *v1.21* / *v.1.22* / *v.1.23* / *v.1.24*  

## 功能  
- - -   
* *kubernetes*集群安装   
* *node*节点安装

## 用法   
- - -   
### 前提条件  
1. 部署的节点能通外部网络  
2. 执行本工具的主机能免密登录到各个节点 
3. 提前安装好*ansible*工具 
4. 所有节点配置好*hostname*  
5. 所有节点配置好*/etc/hosts*文件  

### 命令  
k8s-install -h  
打印帮助   
``` 
[root@master1 tmp]# ./k8s-install -h
Usage:
  k8s-install [OPTIONS]

Application Options:
  -p, --PrintDefault  print install default config
  -i, --install=      k8s or node
  -f, --file=         install config file

Help Options:
  -h, --help          Show this help message

```
k8s-install -p >  /path/to/install.yml  
打印配置样式清单  
```yaml
k8s:
  master:
    components:
      - name: etcd
        hosts:
          - 192.168.58.129
          - 192.168.58.130
          - 192.168.58.131
        dir: /opt/etcd
        dataDir: /data/etcd
      - name: api-server
        hosts:
          - 192.168.58.129
          - 192.168.58.130
          - 192.168.58.131
        dir: /opt/api-server
      - name: scheduler
        hosts:
          - 192.168.58.129
          - 192.168.58.130
          - 192.168.58.131
        dir: /opt/scheduler
      - name: controller-manager
        hosts:
          - 192.168.58.129
          - 192.168.58.130
          - 192.168.58.131
        dir: /opt/controller-manager
  node:
    hosts:
      - 192.168.58.129
      - 192.168.58.130
      - 192.168.58.131
    components:
      - name: kubelet
        dir: /opt/kubelet
      - name: kubproxy
        dir: /opt/kube-proxy
  plugin:
    - name: coreDns
      dns: 172.16.0.10
#    如果使用calico插件，注释掉flannel部分，反之亦然
    - name: calico
      podCIDR: 10.100.0.0/16
      calicoUrl: https://raw.githubusercontent.com/projectcalico/calico/v3.24.1/manifests/calico-typha.yaml
#    - name: flannel
#      podCIDR: 10.100.0.0/16
#      flannelUrl: https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
  CIDR:
    serviceCIDR: 172.16.0.0/16
  certificate: 876000h

# k8s v1.24.0 版本后不支持dockershim
# docker 和 containerd 保留一个
cri:
#  docker:
#    yumRepo: http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
#    dataRoot: /var/lib/docker
#    registryMirrors: https://mvaav0ar.mirror.aliyuncs.com
#    pod-infra-container-image: registry.aliyuncs.com/google_containers/pause:3.5
  containerd:
    yumRepo: http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
    dataRoot: /var/lib/containerd
    registryMirrors:  https://mvaav0ar.mirror.aliyuncs.com
    sandboxImage: registry.aliyuncs.com/google_containers/pause:3.5

haproxy:
  hosts:
    - 192.168.58.129
    - 192.168.58.130
    - 192.168.58.131
  frontendPort: 16443

keepalived:
  hosts:
    - 192.168.58.129
    - 192.168.58.130
    - 192.168.58.131
  vip: 192.168.58.110
  interface: ens33

# 如果手工下载的时，请确保下文件和此yml定义url中名称一致，如
# https://pkg.cfssl.org/R1.2/cfssl_linux-amd64 的文件名为 cfssl_linux-amd64
packages:
  downloadDir: /tmp/packages
  url:
    - https://dl.k8s.io/v1.23.9/bin/linux/amd64/kube-controller-manager
    - https://dl.k8s.io/v1.23.9/bin/linux/amd64/kube-apiserver
    - https://dl.k8s.io/v1.23.9/bin/linux/amd64/kube-scheduler
    - https://github.com/etcd-io/etcd/releases/download/v3.4.20/etcd-v3.4.20-linux-amd64.tar.gz
    - https://dl.k8s.io/v1.23.9/bin/linux/amd64/kube-proxy
    - https://dl.k8s.io/v1.23.9/bin/linux/amd64/kubelet
    - https://dl.k8s.io/v1.23.9/bin/linux/amd64/kubectl
    - https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
    - https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64

```
*建议* 在*跑此工具的主机*里提前下载好相应的包放入 *packages.downloadDir* 定义的目录（calico or flannel清单也可以提前下载好放入此目录）  

k8s-install -i k8s -f /path/to/install.yml  
部署kubernetes集群，此命令也会部署node节点  
   
    
k8s-install -i node -f  /path/to/install.yml  
部署node节点，通常用于已有集群后，新增node节点的部署  
  
  
 ## 安装文件说明 
 ```
k8s.master.components   master节点组件
    components.name     组件
    components.hosts    部署的主机
    components.dir      组件安装目录
    components.dataDir 数据目录

k8s.node.components   node节点组件
    hosts             部署的主机
    components.name   组件
    components.dir    组件安装目录

k8s.plugin            插件
    name              插件名
    dns               DNS
    podCIDR           pod网段
    calicoUrl         calico的清单URL
    flannelUrl        flannel的清单URL

k8s.CIDR              service的网段
    serviceCIDR       service的网段

k8s.certificate       证书有效期

cir.docker            docker 配置
    yumRepo           yum repo地址
    dataRoot          数据目录
    registryMirrors   国内镜像源
    pod-infra-container-image  基础容器镜像地址

cri.containerd       containerd 配置
    yumRepo           yum repo地址
    dataRoot          数据目录
    registryMirrors   国内镜像源
    sandboxImage      基础容器镜像地址

haproxy               haproxy配置
    hosts             部署的节点
    frontendPort      监听的端口

keepalived            keepalived配置
    hosts             部署的节点
    vip               虚拟vip
    interface         网口

packages              k8s 组件存放信息，建议安装后保留此目录方便后续新增node节点时使用相关文件
    downloadDir       存放路径
    url               下载url
```