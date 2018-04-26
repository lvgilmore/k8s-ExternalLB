package main

import (
	"net/http"
	"github.com/SchSeba/k8s-ExternalLB/pkg/plugins/haproxyCluster"
)

const (
	interfaceName = "ens33"
	state = "MASTER"
	statsAddress = "192.168.1.124"
	statsPort = 9000
)

func main() {
	agent := haproxyCluster.CreateAgentInstance(interfaceName,state,statsAddress,statsPort)
	agent.StartProcess()

	http.HandleFunc("/Create", agent.Create)
	http.HandleFunc("/Update", agent.Update)
	http.HandleFunc("/Delete", agent.Delete)
	http.HandleFunc("/Nodes", agent.Nodes)
	http.HandleFunc("/SyncCheck", agent.SyncCheck)
	http.HandleFunc("/Sync", agent.Sync)

	if err := http.ListenAndServe(":9090", nil); err != nil {
		panic(err)
	}
}
