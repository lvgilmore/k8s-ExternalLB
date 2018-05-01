package haproxyCluster

import (
	"net/http"
	"io/ioutil"
	"log"
	"github.com/SchSeba/k8s-ExternalLB/pkg/loadbalancer"
	"io"
	"encoding/json"
	"strconv"
	"os"
)

type Agent struct {
	KeepalivedConfig KeepalivedConfig
	HaproxyConfig HaproxyConfig
	SyncTime int64
}

func (a *Agent) Create(w http.ResponseWriter, r *http.Request) {
	log.Println(" Get create post start working on it")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	var serviceInstance = loadbalancer.ServiceForAgentStruct{}
	json.Unmarshal(body,&serviceInstance)
	log.Println("Unmarshal data for the creation")
	log.Println(serviceInstance)
	a.SyncTime = serviceInstance.SyncTime

	a.KeepalivedConfig.AddNewVirtualInterface(serviceInstance)
	a.KeepalivedConfig.CreateConfigFile()
	a.KeepalivedConfig.ReloadKeepAliveDConfig()

	a.HaproxyConfig.AddNewFarms(serviceInstance)

	a.HaproxyConfig.CreateConfigFile()
	a.HaproxyConfig.ReloadHaproxyConfig()

	io.WriteString(w, "OK")
}

func (a *Agent) Update(w http.ResponseWriter, r *http.Request) {
	log.Println(" Get update post start working on it")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	var serviceInstance = loadbalancer.ServiceForAgentStruct{}
	json.Unmarshal(body,&serviceInstance)
	log.Println("Unmarshal data for the creation")
	log.Println(serviceInstance)

	if a.SyncTime != serviceInstance.SyncTime {
		a.SyncTime = serviceInstance.SyncTime

		a.KeepalivedConfig.AddNewVirtualInterface(serviceInstance)
		a.KeepalivedConfig.CreateConfigFile()
		a.KeepalivedConfig.ReloadKeepAliveDConfig()

		a.HaproxyConfig.UpdateFarms(serviceInstance)

		a.HaproxyConfig.CreateConfigFile()
		a.HaproxyConfig.ReloadHaproxyConfig()
	}
	io.WriteString(w, "OK")
}

func (a *Agent) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println(" Get delete post start working on it")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	var serviceInstance = loadbalancer.ServiceForAgentStruct{}
	json.Unmarshal(body,&serviceInstance)
	log.Println("Unmarshal data for the creation")
	log.Println(serviceInstance)
	a.SyncTime = serviceInstance.SyncTime

	a.HaproxyConfig.DeleteFarms(serviceInstance)

	a.HaproxyConfig.CreateConfigFile()
	a.HaproxyConfig.ReloadHaproxyConfig()

	a.KeepalivedConfig.DeleteVirtualInterface(serviceInstance)

	a.KeepalivedConfig.CreateConfigFile()
	a.KeepalivedConfig.ReloadKeepAliveDConfig()

	io.WriteString(w, "OK")
}

func (a *Agent)Nodes(w http.ResponseWriter, r *http.Request) {
	log.Println(" Get nodes post start working on it")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	var nodes = []string{}
	json.Unmarshal(body,&nodes)

	a.HaproxyConfig.UpdateNodes(nodes)
	a.HaproxyConfig.CreateConfigFile()
	a.HaproxyConfig.ReloadHaproxyConfig()
}

func (a *Agent)Sync(w http.ResponseWriter, r *http.Request) {
	log.Println(" Get sync post start working on it")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}
	var serviceInstance = []loadbalancer.ServiceForAgentStruct{}
	json.Unmarshal(body,&serviceInstance)
	log.Println("Unmarshal data for the creation")
	log.Println(serviceInstance)

	a.HaproxyConfig.Services = make(map[string]Services)
	a.KeepalivedConfig.VirtualInterface = make(map[string]VirtualInterface)

	for _, value := range serviceInstance {
		a.HaproxyConfig.AddNewFarms(value)

		a.KeepalivedConfig.AddNewVirtualInterface(value)

	}

	a.HaproxyConfig.CreateConfigFile()
	a.HaproxyConfig.ReloadHaproxyConfig()

	a.KeepalivedConfig.CreateConfigFile()
	a.KeepalivedConfig.ReloadKeepAliveDConfig()

}

func (a *Agent)SyncCheck(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	syncTime, _ := strconv.ParseInt(string(body),10,64)
	if syncTime != a.SyncTime {
		io.WriteString(w,"true")
	} else {
		io.WriteString(w,"false")
	}

	a.SyncTime = syncTime
}

func (a *Agent)StartProcess() {
	// Run for the stats page
	a.HaproxyConfig.CreateConfigFile()
	a.HaproxyConfig.RunHaproxy()

	a.KeepalivedConfig.CreateConfigFile()
	a.KeepalivedConfig.RunKeepAliveD()
}

func (a *Agent)TermSignal(c chan os.Signal) {
	log.Println("Wating for signal")
	<-c
	log.Println("Get TERMSIGNAL clear virtual interfaces")
	a.KeepalivedConfig.Stop()
	a.KeepalivedConfig.ClearInterfaces()
	a.HaproxyConfig.Stop()
	log.Println("Finnish clearing virtual interfaces")

	os.Exit(0)
}

func CreateAgentInstance(Interface,State string) (Agent) {
	return Agent{KeepalivedConfig:KeepalivedConfig{Interface:Interface,
												   State:State,
												   VirtualInterface:make(map[string]VirtualInterface)},
												   HaproxyConfig:HaproxyConfig{Services:make(map[string]Services)},
												   SyncTime:0}
}

