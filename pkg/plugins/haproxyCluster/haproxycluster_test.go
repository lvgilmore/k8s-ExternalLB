package haproxyCluster_test

import (
	. "github.com/onsi/ginkgo"

	. "github.com/SchSeba/k8s-ExternalLB/pkg/plugins/haproxyCluster"
	. "github.com/onsi/gomega"

	"github.com/SchSeba/k8s-ExternalLB/pkg/loadbalancer"
	"time"

)

var _ = Describe("Haproxycluster", func() {
	var (
		interfaceName = "ens33"
		state = "MASTER"

		haproxyCluster Agent
		ports = []loadbalancer.Port{loadbalancer.Port{Name:"Test-Port",Port:80,NodePort:32000},loadbalancer.Port{Name:"Test-Port-1",Port:81,NodePort:32001}}
		agentDataStruct = loadbalancer.ServiceForAgentStruct{SyncTime:time.Now().Unix(),
															 Nodes:[]string{"10.0.0.1"},
															 VirtualIp:"10.0.0.10",
															 NameSpace:"default",
															 Name:"default-Test-Service",
															 RouterID:1,
															 Protocol:"TCP",Ports:ports}

	)


	Describe("loading from JSON", func() {
		Context("Check haproxy objects", func() {
			haproxyCluster = CreateAgentInstance(interfaceName, state)
			haproxyCluster.HaproxyConfig.AddNewFarms(agentDataStruct)
			haconfig := haproxyCluster.HaproxyConfig

			It("Should create one service", func() {
				Expect(len(haconfig.Services)).To(Equal(1))
			})

			service := haconfig.Services["default-Test-Service"]
			It("Should create one farm", func() {
				Expect(len(service.Farms)).To(Equal(len(ports)))
			})

			newNodes := []string{"10.0.0.2"}
			haconfig.UpdateNodes(newNodes)
			It("Should change the nodes on the farms", func() {
				for _,value := range service.Farms[80].Servers{
					Expect(value.IpAddr).To(Equal("10.0.0.2"))
				}
			})
		})
	})
})
