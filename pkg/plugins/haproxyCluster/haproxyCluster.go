package haproxyCluster

import (
	"net/http"
	"io/ioutil"
	"log"
	"github.com/SchSeba/k8s-ExternalLB/pkg/loadbalancer"
	"io"
	"encoding/json"
)

type Agent struct {
	KeepalivedConfig KeepalivedConfig
	HaproxyConfig HaproxyConfig
}

func (a *Agent) Create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	var serviceInstance = loadbalancer.ServiceForAgentStruct{}
	json.Unmarshal(body,&serviceInstance)
	a.HaproxyConfig.AddNewFarms(serviceInstance)

	a.HaproxyConfig.CreateConfigFile()
	a.HaproxyConfig.ReloadHaproxyConfig()

	a.KeepalivedConfig.AddNewVirtualInterface(serviceInstance)
	a.KeepalivedConfig.CreateConfigFile()
	a.KeepalivedConfig.ReloadKeepAliveDConfig()



	io.WriteString(w, "OK")
}

func (a *Agent) Update(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	var serviceInstance = loadbalancer.ServiceForAgentStruct{}
	json.Unmarshal(body,&serviceInstance)
	a.HaproxyConfig.UpdateFarms(serviceInstance)

	a.HaproxyConfig.CreateConfigFile()
	a.HaproxyConfig.ReloadHaproxyConfig()

	io.WriteString(w, "OK")
}

func (a *Agent) Delete(w http.ResponseWriter, r *http.Request) {
	
}

func (a *Agent) Nodes(w http.ResponseWriter, r *http.Request) {
}

func (a *Agent)StartProcess() {
	a.HaproxyConfig.CreateConfigFile()
	a.HaproxyConfig.RunHaproxy()

	a.KeepalivedConfig.CreateConfigFile()
	a.KeepalivedConfig.RunKeepAliveD()
}

func CreateAgentInstance(Interface,State string) (Agent) {
	return Agent{KeepalivedConfig:KeepalivedConfig{Interface:Interface,
												   State:State,
												   VirtualInterface:make(map[string]VirtualInterface)},
												   HaproxyConfig:HaproxyConfig{Services:make(map[string]Services)}}
}

