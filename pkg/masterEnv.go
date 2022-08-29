package pkg

import "k8s/pkg/utils"

func InitEnv(host string) {
	utils.Exec(host, "shell", "yum install -y yum-utils   device-mapper-persistent-data  sed lvm2  jq ipvsadm ipset  conntrack libseccomp")
	utils.Exec(host, "service", "name=firewalld enable=no state=stopped")
	utils.Exec(host, "service", "name=NetworkManager enable=no state=stopped")
	utils.Exec(host, "shell", "setenforce 0 && sed -i 's@SELINUX=enforcing@SELINUX=disabled@g'  /etc/selinux/config")
	utils.Exec(host, "shell", "swapoff -a && sysctl -w vm.swappiness=0 && sed -ri '/^[^#]*swap/s@^@#@' /etc/fstab")



}
