package haproxyCluster

import (
	"os"
	"text/template"
	"github.com/SchSeba/k8s-ExternalLB/pkg/loadbalancer"
	"bytes"
	"fmt"
	"log"
	"strings"
	"os/exec"
	"strconv"
)

type HaproxyConfig struct {
	Services map[string]Services
	Config string

}

type Services struct {
	Farms map[int32]Farm
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
	farmName := serviceInstance.Name
	services, ok := h.Services[farmName]
	if !ok {
		h.Services[farmName] = Services{make(map[int32]Farm)}
		services = h.Services[farmName]
	}

	for _,value := range serviceInstance.Ports {
		realServers := make([]RealServer, len(serviceInstance.Nodes))
		for index,serverValue := range serviceInstance.Nodes {
			realServers[index] = RealServer{serverValue,value.NodePort}
		}

		services.Farms[value.Port] = Farm{Name:farmName + "-" + strconv.Itoa(int(value.Port)),
										  Protocol:strings.ToLower(serviceInstance.Protocol),
										  BindAddr:serviceInstance.VirtualIp,
										  BindPort:value.Port,
										  Servers:realServers}
	}
}

func (h *HaproxyConfig)UpdateFarms(serviceInstance loadbalancer.ServiceForAgentStruct) {
	farmName := serviceInstance.Name
	delete(h.Services, farmName)
	h.AddNewFarms(serviceInstance)

}

func (h *HaproxyConfig)CreateConfigFile() {
	h.Config = haproxyGlobalConfig
	tmpl, _ := template.New("HaproxyFarmConfig").Parse(haproxyTemplate)
	for _,service := range h.Services {
		for _, value := range service.Farms {
			farmConfig := new(bytes.Buffer)
			tmpl.Execute(farmConfig, value)
			h.Config += farmConfig.String()
		}
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