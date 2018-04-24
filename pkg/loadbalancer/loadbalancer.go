package loadbalancer

import (
	"github.com/coreos/etcd/client"
	"time"
	"log"
	"net"
	"net/http"
	"context"
	"errors"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"github.com/SchSeba/k8s-ExternalLB/pkg/k8s/loadbalance"
	"io/ioutil"
	"sync"
	"k8s.io/api/core/v1"
	"bytes"
)

type Node struct {
	IpAddr string
	IsOnline bool
}


type LBController struct {
	agents []Node
	DbConnection client.Client
	kapi client.KeysAPI
	HostsList []string
	mutex *sync.Mutex
}

type ServiceDataStruct struct {
	ServiceData *v1.Service `json:"service_data"`
	Nodes []string `json:"nodes"`
	IsCreated bool `json:"is_created"`
	RouterID int `json:"router_id"`
}

type ServiceForAgentStruct struct {
	Name string      `json:"name"`
	NameSpace string `json:"name_space"`
	VirtualIp string `json:"virtual_ip"`
	Nodes []string   `json:"nodes"`
	Protocol string  `json:"protocol"`
	Ports []Port     `json:"ports"`
	RouterID int `json:"router_id"`
}

type Port struct {
	Name string `json:"name"`
	Port int32    `json:"port"`
	NodePort int32 `json:"node_port"`
}

func convertToAgentStruct(serviceDataStruct ServiceDataStruct) ([]byte,error) {
	ports := make([]Port, len(serviceDataStruct.ServiceData.Spec.Ports))

	for index, value := range serviceDataStruct.ServiceData.Spec.Ports {
		ports[index] = Port{Name:value.Name,Port:value.Port,NodePort:value.NodePort}
	}

	serviceForAgentInstance := ServiceForAgentStruct{VirtualIp:serviceDataStruct.ServiceData.Spec.ExternalIPs[0],
	                                                 Nodes:serviceDataStruct.Nodes,
	                                                 Protocol:string(serviceDataStruct.ServiceData.Spec.Ports[0].Protocol),
													 Ports:ports,
													 RouterID:serviceDataStruct.RouterID,
													 Name:getServiceName(serviceDataStruct),
													 NameSpace:serviceDataStruct.ServiceData.Namespace}

	body,err := json.Marshal(&serviceForAgentInstance)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (l *LBController)unMarshalDataStruct(body []byte) ( loadbalance.LbDataStruct) {
	var lbDataStruct = loadbalance.LbDataStruct{}
	json.Unmarshal(body,&lbDataStruct)

	return  lbDataStruct
}

func getServiceName(serviceData ServiceDataStruct) string {
	return serviceData.ServiceData.Namespace+"-"+serviceData.ServiceData.Name
}

func (l *LBController)GetEmptyIp() (string, error) {
	for _,ip := range l.HostsList {
		_, err := l.kapi.Get(context.Background(), "/mylb/allocatedIps/"+ ip, nil)
		if err != nil {
			err := err.(client.Error)
			if err.Code == 100 {
				return ip, nil
			} else {
				log.Fatal(err)
			}
		}
	}
	return "", errors.New("Fail to find any Key")
}

func (l *LBController)getVirtualRouterID() int {
	for i :=1; i< 255; i++  {
		_, err := l.kapi.Get(context.Background(), "/mylb/VirtualRouterID/"+ string(i), nil)
		if err != nil {
			err := err.(client.Error)
			if err.Code == 100 {
				return i
			}
		}
	}

	return -1
}

func (l *LBController)ClearDB() {
	for _,ip := range l.HostsList {
		_, err := l.kapi.Get(context.Background(), "/mylb/allocatedIps/"+ip, nil)
		if err == nil {
			l.kapi.Delete(context.Background(), "/mylb/allocatedIps/"+ip, nil)
		}
	}
}

func (l *LBController)sendDataToAgents(serviceDataStruct ServiceDataStruct, commandType string) (error) {
	IsAnyAgentAlive := false
	jsonToAgents, err := convertToAgentStruct(serviceDataStruct)

	if err != nil {
		return err
	}

	for _,agent := range l.agents {
		_, err = http.Post("http://" + agent.IpAddr+"/"+commandType, "application/json",bytes.NewBuffer(jsonToAgents))

		if err != nil {
			log.Print(err)
		} else {
			IsAnyAgentAlive = true
		}

	}

	if !IsAnyAgentAlive {
		return errors.New("0 Agents Alive")
	}

	return nil
}

func (l *LBController)Create(w http.ResponseWriter, r *http.Request) {
	body, errIO := ioutil.ReadAll(r.Body)
	if errIO != nil {
		log.Panic(errIO)
	}
	lbDataStruct := l.unMarshalDataStruct(body)

	var ip string
	var err error

	l.mutex.Lock()

	if len(lbDataStruct.ServiceData.Spec.ExternalIPs) == 0 {
		ip , err = l.GetEmptyIp()
	} else {
		ip = lbDataStruct.ServiceData.Spec.ExternalIPs[0]
		err = nil
	}

	// TODO: REMOVE FOR DEBUG ONLY
	ip , err = l.GetEmptyIp()

	if err != nil {
		log.Print(err)
		l.mutex.Unlock()
		io.WriteString(w, "")
	} else {
		l.kapi.Set(context.Background(), "/mylb/allocatedIps/"+ip, lbDataStruct.ServiceData.Namespace+"-"+lbDataStruct.ServiceData.Name, nil)
		routerID := l.getVirtualRouterID()
		l.mutex.Unlock()

		lbDataStruct.ServiceData.Spec.ExternalIPs = []string{ip}
		serviceDataStruct := ServiceDataStruct{ServiceData:lbDataStruct.ServiceData,Nodes:lbDataStruct.NodeList,IsCreated:false,RouterID:routerID}
		err = l.sendDataToAgents(serviceDataStruct, "Create")
		if err == nil {
			serviceDataStruct.IsCreated = true
		}

		body,err = json.Marshal(&serviceDataStruct)
		l.kapi.Set(context.Background(), "/mylb/services/"+getServiceName(serviceDataStruct), string(body), nil)
		io.WriteString(w, ip)
	}
}

func (l *LBController)Update(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	lbDataStruct := l.unMarshalDataStruct(body)
	log.Print(lbDataStruct)

}

func (l *LBController)Delete(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	lbDataStruct:= l.unMarshalDataStruct(body)
	log.Print(lbDataStruct)

}

func (l *LBController)Nodes(w http.ResponseWriter, r *http.Request) {

}



func LBControllerInitializer(Endpoints []string,Nodes []string,cidr string) *LBController {

	cfg := client.Config{
		Endpoints:               Endpoints,
		Transport:               client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	etdConnection, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi := client.NewKeysAPI(etdConnection)

	agents := make([]Node,len(Nodes))

	for index,node := range Nodes {
		agents[index] =  Node{node,true}
	}

	ips, err := GetsHosts(cidr)
	if err != nil {
		log.Fatal(err)
	}

	return &LBController{agents,etdConnection,kapi,ips,&sync.Mutex{}}

}

func GetsHosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}