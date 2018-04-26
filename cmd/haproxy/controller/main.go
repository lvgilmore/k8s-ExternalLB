package main

import (
	"net/http"
	"github.com/SchSeba/k8s-ExternalLB/pkg/loadbalancer"
)

var (
	etcdEndPoints =  []string{"http://127.0.0.1:6666"}
	cidr = "192.168.1.32/27"
	agents = []string{"192.168.1.124:9090"}
)


func main() {

	//conn, _ :=net.Dial("tcp", "127.0.0.1:8081")
	//conn.Write([]byte("asdasdasd"))

	lbController := loadbalancer.LBControllerInitializer(etcdEndPoints,agents,cidr)

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
