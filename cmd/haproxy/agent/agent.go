package main

import (
	"net/http"
	"github.com/SchSeba/k8s-ExternalLB/pkg/plugins/haproxyCluster"
)

const (
	interfaceName = "ens33"
	state = "MASTER"
)

func main() {
	agent := haproxyCluster.CreateAgentInstance(interfaceName,state)
	//agent.StartProcess()

	http.HandleFunc("/Create", agent.Create)
	http.HandleFunc("/Update", agent.Update)
	http.HandleFunc("/Delete", agent.Delete)
	http.HandleFunc("/Nodes", agent.Nodes)

	if err := http.ListenAndServe(":9090", nil); err != nil {
		panic(err)
	}
}
