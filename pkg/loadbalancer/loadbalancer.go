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
	"strconv"
	"fmt"
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
	SyncTime int64 `json:"sync_time"`
}

type Port struct {
	Name string `json:"name"`
	Port int32    `json:"port"`
	NodePort int32 `json:"node_port"`
}

func getAgentStruct(serviceDataStruct ServiceDataStruct,syncTime int64) ServiceForAgentStruct {
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
													 NameSpace:serviceDataStruct.ServiceData.Namespace,SyncTime:syncTime}

	return serviceForAgentInstance
}

func convertToAgentStruct(serviceDataStruct ServiceDataStruct,syncTime int64) ([]byte,error) {
	serviceForAgentInstance := getAgentStruct(serviceDataStruct,syncTime)

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

func (l *LBController)getVirtualRouterID(nameSpace string) int {
	value, err := l.kapi.Get(context.Background(), "/mylb/VirtualRouterID/Namespace/"+ nameSpace, nil)
	if err == nil {
		i , _ := strconv.Atoi(value.Node.Value)
		return i
	} else {
		for i := 1; i < 255; i++ {
			_, err := l.kapi.Get(context.Background(), "/mylb/VirtualRouterID/ID/"+strconv.Itoa(i), nil)
			if err != nil {
				log.Println(err)
				errType := fmt.Sprintf("%T", err)
				if errType == "client.Error" && err.(client.Error).Code == 100 {
					l.kapi.Set(context.Background(),"/mylb/VirtualRouterID/Namespace/"+ nameSpace,strconv.Itoa(i),nil)
					l.kapi.Set(context.Background(),"/mylb/VirtualRouterID/ID/"+strconv.Itoa(i),strconv.Itoa(i),nil)
					return i
				}
			}
		}

		return -1
	}
}

func (l *LBController)ClearDB() {
	//for _,ip := range l.HostsList {
	//	_, err := l.kapi.Get(context.Background(), "/mylb/allocatedIps/"+ip, nil)
	//	if err == nil {
	//		l.kapi.Delete(context.Background(), "/mylb/allocatedIps/"+ip, nil)
	//	}
	//}

	_ ,err := l.kapi.Delete(context.Background(), "/mylb/allocatedIps/",&client.DeleteOptions{Recursive:true,Dir:true})
	if err != nil {
		log.Println(err)
	}

	_,err = l.kapi.Delete(context.Background(), "/mylb/VirtualRouterID/",&client.DeleteOptions{Recursive:true,Dir:true})
	if err != nil {
		log.Println(err)
	}
	//for i := 1; i < 255; i++ {
	//	_, err := l.kapi.Get(context.Background(), "/mylb/VirtualRouterID/ID/"+string(i), nil)
	//	if err == nil {
	//		l.kapi.Delete(context.Background(), "/mylb/VirtualRouterID/ID/"+string(i),nil)
	//	}
	//}
}

func (l *LBController)sendDataToAgents(serviceDataStruct ServiceDataStruct, commandType string) (error) {
	syncTime := time.Now().Unix()
	l.kapi.Set(context.Background(), "/mylb/SyncTime",strconv.FormatInt(syncTime,10),nil)

	IsAnyAgentAlive := false
	jsonToAgents, err := convertToAgentStruct(serviceDataStruct,syncTime)

	if err != nil {
		return err
	}

	for _,agent := range l.agents {
		_, err = http.Post("http://" + agent.IpAddr+"/"+commandType, "application/json",bytes.NewBuffer(jsonToAgents))

		if err != nil {
			log.Println("Agent " + agent.IpAddr + " is down")
			agent.IsOnline = false
			log.Print(err)
		} else {
			log.Println("Agent " + agent.IpAddr + " is alive")
			agent.IsOnline = true
			IsAnyAgentAlive = true
		}

	}

	if !IsAnyAgentAlive {
		return errors.New("0 Agents Alive")
	}

	return nil
}

func (l *LBController)DeleteDataFromDB(serviceDataStruct ServiceDataStruct) {
	// Remove Ip
	_, err := l.kapi.Delete(context.Background(), "/mylb/allocatedIps/"+serviceDataStruct.ServiceData.Spec.ExternalIPs[0], nil)
	if err != nil {
		log.Println(err)
	}

	// remove Virtual Router
	_, err = l.kapi.Delete(context.Background(), "/mylb/VirtualRouterID/Namespace/"+ serviceDataStruct.ServiceData.Namespace, nil)
	if err != nil {
		log.Println(err)
	}

	_, err = l.kapi.Delete(context.Background(), "/mylb/VirtualRouterID/ID/"+strconv.Itoa(serviceDataStruct.RouterID), nil)
	if err != nil {
		log.Println(err)
	}

	// Remove Service Data
	_, err = l.kapi.Delete(context.Background(), "/mylb/services/"+getServiceName(serviceDataStruct), nil)
	if err != nil {
		log.Println(err)
	}

}

func (l *LBController)Create(w http.ResponseWriter, r *http.Request) {
	body, errIO := ioutil.ReadAll(r.Body)
	if errIO != nil {
		log.Panic(errIO)
	}
	lbDataStruct := l.unMarshalDataStruct(body)
	resourceVersion := lbDataStruct.ServiceData.ObjectMeta.ResourceVersion

	resp, err := l.kapi.Get(context.Background(),"/mylb/services/"+lbDataStruct.ServiceData.Namespace +"-"+ lbDataStruct.ServiceData.Name,nil)
	if err != nil {
		log.Println(err)
	} else {
		var serviceDataStruct ServiceDataStruct
		err := json.Unmarshal([]byte(resp.Node.Value), &serviceDataStruct)

		if err != nil {
			log.Println(err)
		} else if serviceDataStruct.ServiceData.ObjectMeta.ResourceVersion != resourceVersion {
			log.Println("Need to update service " + getServiceName(serviceDataStruct))

			var ip string
			var err error

			l.mutex.Lock()

			if len(lbDataStruct.ServiceData.Spec.ExternalIPs) == 0 {
				ip, err = l.GetEmptyIp()
			} else {
				ip = lbDataStruct.ServiceData.Spec.ExternalIPs[0]
				err = nil
			}

			// TODO: REMOVE FOR DEBUG ONLY
			//ip , err = l.GetEmptyIp()

			if err != nil {
				log.Print(err)
				l.mutex.Unlock()
				io.WriteString(w, "")
			} else {
				l.kapi.Set(context.Background(), "/mylb/allocatedIps/"+ip, lbDataStruct.ServiceData.Namespace+"-"+lbDataStruct.ServiceData.Name, nil)
				routerID := l.getVirtualRouterID(lbDataStruct.ServiceData.Namespace)
				l.mutex.Unlock()

				lbDataStruct.ServiceData.Spec.ExternalIPs = []string{ip}
				serviceDataStruct := ServiceDataStruct{ServiceData: lbDataStruct.ServiceData, Nodes: lbDataStruct.NodeList, IsCreated: false, RouterID: routerID}
				err = l.sendDataToAgents(serviceDataStruct, "Create")
				if err == nil {
					serviceDataStruct.IsCreated = true
				}

				body, err = json.Marshal(&serviceDataStruct)
				l.kapi.Set(context.Background(), "/mylb/services/"+getServiceName(serviceDataStruct), string(body), nil)
				io.WriteString(w, ip)
			}
		} else {
			log.Println("No update needed for " + getServiceName(serviceDataStruct))
		}
	}
}

func (l *LBController)Update(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	lbDataStruct := l.unMarshalDataStruct(body)
	resourceVersion := lbDataStruct.ServiceData.ObjectMeta.ResourceVersion

	resp, err := l.kapi.Get(context.Background(),"/mylb/services/"+lbDataStruct.ServiceData.Namespace +"-"+ lbDataStruct.ServiceData.Name,nil)
	if err != nil {
		log.Println(err)
	} else {
		var serviceDataStruct ServiceDataStruct
		err := json.Unmarshal([]byte(resp.Node.Value),&serviceDataStruct)

		if err != nil {
			log.Println(err)
		} else if serviceDataStruct.ServiceData.ObjectMeta.ResourceVersion != resourceVersion {
			log.Println("Need to update service " + getServiceName(serviceDataStruct))
			routerID := l.getVirtualRouterID(lbDataStruct.ServiceData.Namespace)
			serviceDataStruct = ServiceDataStruct{ServiceData: lbDataStruct.ServiceData, Nodes: lbDataStruct.NodeList, IsCreated: false, RouterID: routerID}
			err = l.sendDataToAgents(serviceDataStruct, "Update")
			if err == nil {
				serviceDataStruct.IsCreated = true
			}

			body, err = json.Marshal(&serviceDataStruct)
			l.kapi.Set(context.Background(), "/mylb/services/"+getServiceName(serviceDataStruct), string(body), nil)

		} else {
			log.Println("No update need for service " + getServiceName(serviceDataStruct))
		}
	}

	io.WriteString(w, lbDataStruct.ServiceData.Spec.ExternalIPs[0])
}

func (l *LBController)Delete(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	lbDataStruct:= l.unMarshalDataStruct(body)
	value, err := l.kapi.Get(context.Background(), "/mylb/services/"+ lbDataStruct.ServiceData.Namespace+ "-" + lbDataStruct.ServiceData.Name, nil)
	if err == nil {
		var serviceDataStruct = ServiceDataStruct{}
		json.Unmarshal([]byte(value.Node.Value),&serviceDataStruct)
		err = l.sendDataToAgents(serviceDataStruct, "Delete")
		l.DeleteDataFromDB(serviceDataStruct)
		io.WriteString(w, lbDataStruct.ServiceData.Spec.ExternalIPs[0])
	} else {
		log.Println("Service" + lbDataStruct.ServiceData.Namespace+ "-" + lbDataStruct.ServiceData.Name + " is not in the database")
	}

	io.WriteString(w, "")
}

func (l *LBController)Nodes(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	lbDataStruct:= l.unMarshalDataStruct(body)
	jsonToAgents, err := json.Marshal(lbDataStruct.NodeList)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Update Nodes")
		for _, agent := range l.agents {
			_, err = http.Post("http://"+agent.IpAddr+"/Nodes", "application/json", bytes.NewBuffer(jsonToAgents))

			if err != nil {
				log.Println("Agent " + agent.IpAddr + " is down")
				agent.IsOnline = false
				log.Print(err)
			} else {
				log.Println("Agent " + agent.IpAddr + " is alive")
				agent.IsOnline = true
			}

		}

		directory, err := l.kapi.Get(context.Background(), "/mylb/services/", nil)
		if err != nil {
			log.Println(err)
			io.WriteString(w, "")
		}

		for _, node := range directory.Node.Nodes {
			var serviceDataInstance = ServiceDataStruct{}
			json.Unmarshal([]byte(node.Value), &serviceDataInstance)
			serviceDataInstance.Nodes = lbDataStruct.NodeList
			body, err = json.Marshal(&serviceDataInstance)
			l.kapi.Set(context.Background(), "/mylb/services/"+getServiceName(serviceDataInstance), string(body), nil)
		}

		io.WriteString(w, "OK")
	}
}

func (l *LBController)SyncAgents() {
	c := time.Tick(10 * time.Second)

	for {
		<-c
		log.Println("Syncing with agents")
		NeedToSyncAgents := []Node{}

		node, err := l.kapi.Get(context.Background(),"/mylb/SyncTime",nil)
		if err == nil {

			syncTime := []byte(node.Node.Value)
			syncTimeInt, _ := strconv.ParseInt(node.Node.Value,10,64)
			for _, value := range l.agents {
				r, err := http.Post("http://"+value.IpAddr+"/SyncCheck", "application/json", bytes.NewBuffer(syncTime))

				if err != nil {
					log.Println("Agent " + value.IpAddr + " is down")
					value.IsOnline = false
					log.Print(err)
				} else {
					log.Println("Agent " + value.IpAddr + " is alive")
					value.IsOnline = true
					resp, _ := ioutil.ReadAll(r.Body)
					if b, _ := strconv.ParseBool(string(resp)); b == true {
						NeedToSyncAgents = append(NeedToSyncAgents, value)
					}
				}
			}

			if len(NeedToSyncAgents) > 0 {
				syncData := []ServiceForAgentStruct{}
				services, err := l.kapi.Get(context.Background(), "/mylb/services/", nil)
				if err != nil {
					log.Println(err)
				} else {
					for _, value := range services.Node.Nodes {
						var serviceData ServiceDataStruct
						json.Unmarshal([]byte(value.Value), &serviceData)
						syncData = append(syncData, getAgentStruct(serviceData, syncTimeInt))
					}
				}

				syncDataByte, err := json.Marshal(&syncData)
				if err != nil {
					log.Println(err)
				} else {
					for _, value := range NeedToSyncAgents {
						_, err = http.Post("http://"+value.IpAddr+"/Sync", "application/json", bytes.NewBuffer(syncDataByte))
					}
				}
			}
		}
	}

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