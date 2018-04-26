package haproxyCluster

import (
	"os"
	"github.com/SchSeba/k8s-ExternalLB/pkg/loadbalancer"
	"text/template"
	"bytes"
	"log"
	"os/exec"
	"fmt"

)

type KeepalivedConfig struct {
	Interface string
	State string
	VirtualInterface map[string]VirtualInterface
	Config string
}

type VirtualInterface struct {
	Name string
	State string
	Interface string
	RouterID int
	Servers map[string]struct{}

}
const (
	keepalivedGlobalConfig = `global_defs {
    tlived process identifier
lvs_id haproxy_DH
}
# Script used to check if HAProxy is running
vrrp_script check_haproxy {
script "killall -0 haproxy"
interval 2
weight 2
}
`
	keepalivedTemplate = `vrrp_instance {{.Name}} {
state {{.State}}
interface {{.Interface}}
virtual_router_id {{.RouterID}}
priority 101
# The virtual ip address shared between the two loadbalancers
virtual_ipaddress_excluded {
{{range $key, $value := .Servers}}{{$key}}
{{end}}
}
track_script {
check_haproxy
}
}`
)

func (k *KeepalivedConfig)AddNewVirtualInterface(serviceInstance loadbalancer.ServiceForAgentStruct) {
	virtualInterface, ok := k.VirtualInterface[serviceInstance.NameSpace]

	if !ok {
		k.VirtualInterface[serviceInstance.NameSpace] = VirtualInterface{Name: serviceInstance.NameSpace,
			Interface: k.Interface,
			State: k.State,
			RouterID: serviceInstance.RouterID,Servers:make(map[string]struct{})}
		virtualInterface = 	k.VirtualInterface[serviceInstance.NameSpace]
	}

	virtualInterface.Servers[serviceInstance.VirtualIp] = struct{}{}

}

func (k *KeepalivedConfig)CreateConfigFile() {
	k.Config = keepalivedGlobalConfig
	tmpl, _ := template.New("KeepalivedConfigTemaplte").Parse(keepalivedTemplate)
	for _,value := range k.VirtualInterface {
		vip := new(bytes.Buffer)
		tmpl.Execute(vip, value)
		k.Config += vip.String()
	}

	f, err := os.Create("/etc/keepalived/keepalived.conf")
	if err != nil {
		log.Print(err)
	}

	defer f.Close()

	f.WriteString(k.Config)
}

func (k *KeepalivedConfig)DeleteVirtualInterface(serviceInstance loadbalancer.ServiceForAgentStruct) {
	virtualInterface, ok := k.VirtualInterface[serviceInstance.NameSpace]

	if ok {
		if len(virtualInterface.Servers) == 1 {
			delete(k.VirtualInterface, serviceInstance.NameSpace)
		} else {
			delete(virtualInterface.Servers, serviceInstance.VirtualIp)
		}
	}
}

func (k *KeepalivedConfig)ReloadKeepAliveDConfig() {
	if _, err := os.Stat("/var/run/keepalived.pid"); os.IsNotExist(err) {
		cmd := exec.Command("/usr/sbin/keepalived")
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			log.Print(err)
		}

		fmt.Printf("%s\n", stdoutStderr)
	} else {
		cmd := exec.Command("cat", "/var/run/keepalived.pid")
		pid, _ := cmd.Output()

		cmd = exec.Command("kill", "-HUP", string(pid[:len(pid)-1]))
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			log.Print(err)
		}

		fmt.Printf("%s\n", stdoutStderr)
	}
}

func (k *KeepalivedConfig)RunKeepAliveD(){
	if _, err := os.Stat("/var/run/keepalived.pid"); os.IsNotExist(err) {
		cmd := exec.Command("/usr/sbin/keepalived")
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			log.Print(err)
		}

		fmt.Printf("%s\n", stdoutStderr)
	} else {
		k.ReloadKeepAliveDConfig()
	}
}