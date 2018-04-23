package haproxyCluster

import (
	"os"
	"text/template"
	"github.com/SchSeba/k8s-ExternalLB/pkg/loadbalancer"
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"os/exec"
)

type HaproxyConfig struct {
	Farms map[string]Farm
	Config string

}

type Farm struct {
	Name string
	BindAddr string
	BindPort int32
	Protocol string
	Servers []RealServer
}

type RealServer struct {
	IpAddr string
	Port int32
}

const (
	haproxyGlobalConfig = `
global
        log 127.0.0.1   local0
        log 127.0.0.1   local1 notice
        #log loghost    local0 info
        maxconn 4096
        #chroot /usr/share/haproxy
        user haproxy
        group haproxy
        daemon
        #debug
        #quiet

defaults
        log     global
        option  dontlognull
        retries 3
        option redispatch
        maxconn 2000
        timeout connect      5000
        timeout client      50000
        timeout server      50000

`

	haproxyTemplate  = `listen {{.Name}}
	bind {{.BindAddr}}:{{.BindPort}}
	mode {{.Protocol}}
	balance roundrobin
	{{range .Servers}}server {{.IpAddr}}:{{.Port}} {{.IpAddr}}:{{.Port}} 
{{end}}
`
)

func (h *HaproxyConfig)AddNewFarms(serviceInstance loadbalancer.ServiceForAgentStruct) {
	for _,value := range serviceInstance.Ports {
		realServers := make([]RealServer, len(serviceInstance.Nodes))
		for index,serverValue := range serviceInstance.Nodes {
			realServers[index] = RealServer{serverValue,value.NodePort}
		}
		farmName := serviceInstance.Name+"-"+strconv.Itoa(int(value.Port))

		h.Farms[farmName] = Farm{Name:farmName,
								 Protocol:strings.ToLower(serviceInstance.Protocol),
								 BindAddr:serviceInstance.VirtualIp,
								 BindPort:value.Port,
								 Servers:realServers}
	}
}

func (h *HaproxyConfig)CreateConfigFile() {
	h.Config = haproxyGlobalConfig
	tmpl, _ := template.New("HaproxyFarmConfig").Parse(haproxyTemplate)
	for _,value := range h.Farms {
		farmConfig := new(bytes.Buffer)
		tmpl.Execute(farmConfig, value)
		h.Config += farmConfig.String()
	}

	f, err := os.Create("/etc/haproxy/haproxy.cfg")
	if err != nil {
		log.Print(err)
	}

	defer f.Close()

	f.WriteString(h.Config)

}

func (h *HaproxyConfig)ReloadHaproxyConfig() {
	if _, err := os.Stat("/run/haproxy.pid"); os.IsNotExist(err) {
		cmd := exec.Command("haproxy", "-f", "/etc/haproxy/haproxy.cfg", "-p", "/run/haproxy.pid")
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			log.Print(err)
		}

		fmt.Printf("%s\n", stdoutStderr)

	} else {
		cmd := exec.Command("cat", "/run/haproxy.pid")
		pid, _ := cmd.Output()
		cmd = exec.Command("haproxy", "-f", "/etc/haproxy/haproxy.cfg", "-p", "/run/haproxy.pid", "-sf", string(pid))
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			log.Print(err)
		}

		fmt.Printf("%s\n", stdoutStderr)
	}
}

func (h *HaproxyConfig)RunHaproxy(){
	if _, err := os.Stat("/run/haproxy.pid"); os.IsNotExist(err) {
		cmd := exec.Command("haproxy", "-f", "/etc/haproxy/haproxy.cfg", "-p", "/run/haproxy.pid")
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			log.Print(err)
		}

		fmt.Printf("%s\n", stdoutStderr)

	} else {
		h.ReloadHaproxyConfig()
	}
}

func Temp() {
	farm := Farm{Name:"test-Farm",Protocol:"TCP",BindAddr:"10.0.0.0",BindPort:1000,Servers:[]RealServer{{IpAddr:"10.0.0.1",Port:12345},
		{IpAddr:"10.0.0.2",Port:12345}}}
	tmpl, err := template.New("test").Parse(haproxyTemplate)
	if err != nil { panic(err) }
	err = tmpl.Execute(os.Stdout, farm)
	if err != nil { panic(err) }
}