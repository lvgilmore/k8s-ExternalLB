package main

import (
	"net/http"
	"github.com/SchSeba/k8s-ExternalLB/pkg/loadbalancer"
	"log"
	"os"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
)

var (
	etcdEndPointsConst =  []string{"http://127.0.0.1:6666"}
	cidrConst = "192.168.1.32/27"
	agentsConst = []string{"192.168.1.124:9090"}
)

type Variabels struct {
	EtcdEndPoints []string `json:"etcd_end_points"`
	Cidr string `json:"cidr"`
	Agents []string `json:"agents"`
} 

func loadVariables() Variabels{
	var variables Variabels

	raw, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Panic(err)
	}

	json.Unmarshal(raw,&variables)

	return variables
}

func main() {
	log.Println("Starting HaproxyCluster controller")
	log.Println("Check for enviroment variables")
	var variables Variabels

	if h := os.Getenv("Prod"); h == "TRUE" {
		log.Println("Load enviroment variables")
		variables = loadVariables()
	} else {
		log.Println("Enviroment variables not found use const data (for development only!)")
		variables = Variabels{EtcdEndPoints:etcdEndPointsConst,Cidr:cidrConst,Agents:agentsConst}
	}

	lbController := loadbalancer.LBControllerInitializer(variables.EtcdEndPoints,variables.Agents,variables.Cidr)

	// TODO: For Debug
	//lbController.ClearDB()

	go lbController.SyncAgents()


	http.HandleFunc("/Create", lbController.Create)
	http.HandleFunc("/Update", lbController.Update)
	http.HandleFunc("/Delete", lbController.Delete)
	http.HandleFunc("/Nodes", lbController.Nodes)


	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
