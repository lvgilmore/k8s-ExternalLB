package main

import (
	"net/http"
	"github.com/SchSeba/k8s-ExternalLB/pkg/plugins/haproxyCluster"
	"os"
	"log"
	"os/signal"
	"syscall"
)

const (
	interfaceNameConst = "ens33"
	stateConst = "MASTER"
	adminIpInterfaceConst = "192.168.1.124"
	adminPortConst = "9090"
)

var (
	interfaceName string
	state string
	adminIpInterface string
	adminPort string
)

func loadVariables() {
	if h := os.Getenv("interfaceName"); h != "" {
		interfaceName = h
	} else {
		log.Panic("no interface name found")
	}

	if h := os.Getenv("state"); h == "MASTER" || h == "SLAVE" {
		state = h
	} else {
		log.Panic("no state (need to be MASTER or SLAVE)")
	}

	if h := os.Getenv("adminIpInterface"); h != "" {
		adminIpInterface = h
	} else {
		log.Panic("no adminIpInterface")
	}

	if h := os.Getenv("adminPort"); h != "" {
		adminPort = h
	} else {
		log.Panic("no adminPort")
	}
}

func main() {
	log.Println("Starting agent")
	log.Println("Check for enviroment variables")

	if h := os.Getenv("Prod"); h == "TRUE" {
		log.Println("Load enviroment variables")
		loadVariables()
	} else {
		log.Println("Enviroment variables not found use const data (for development only!)")
		interfaceName = interfaceNameConst
		state = stateConst
		adminIpInterface = adminIpInterfaceConst
		adminPort = adminPortConst
	}


	agent := haproxyCluster.CreateAgentInstance(interfaceName,state)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go agent.TermSignal(c)


	agent.StartProcess()


	http.HandleFunc("/Create", agent.Create)
	http.HandleFunc("/Update", agent.Update)
	http.HandleFunc("/Delete", agent.Delete)
	http.HandleFunc("/Nodes", agent.Nodes)
	http.HandleFunc("/SyncCheck", agent.SyncCheck)
	http.HandleFunc("/Sync", agent.Sync)

	if err := http.ListenAndServe(adminIpInterface+ ":" + adminPort, nil); err != nil {
		panic(err)
	}
}
