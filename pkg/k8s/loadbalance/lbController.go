package loadbalance

import (
	"k8s.io/api/core/v1"
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"fmt"
	"google.golang.org/grpc"
	"log"
	pb "github.com/SchSeba/k8s-ExternalLB/pkg/externallb"
	"time"
	"context"
)

type lbController struct {
	ipAddr string
	port int
}

type LbDataStruct struct {
	ServiceData *v1.Service `json:"service_data"`
	NodeList []string `json:"node_list"`
}

func (l *lbController) marshalData(serviceObject *v1.Service, nodeList []string)  []byte {
	b, err := json.Marshal(&LbDataStruct{ServiceData:serviceObject,NodeList:nodeList})

	if err != nil {
		fmt.Errorf("fail to Marshal Error: %v",err)
	}

	return b
}

func (l *lbController) sendData(url string, data []byte) (string,error) {
	resp, err := http.Post(url, "application/json",bytes.NewBuffer(data))

	if err != nil {
		return "",err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), nil
}


func NewLBController(ipAddr string, port int) *lbController {
	return &lbController{ipAddr:ipAddr,port:port}
}

func createDataObject(serviceObject *v1.Service, nodeList []string) *pb.Data {
	ports := make([]*pb.Port,len(serviceObject.Spec.Ports))

	for key,value := range serviceObject.Spec.Ports {
		ports[key] = &pb.Port{Name:value.Name,Port:value.Port,NodePort:value.NodePort}
	}
	return &pb.Data{ServiceName:serviceObject.Name,
					Namespace:serviceObject.Namespace,
					Nodes:nodeList,
					ExternalIPs:serviceObject.Spec.ExternalIPs,
					Protocol:string(serviceObject.Spec.Ports[0].Protocol),Ports:ports}
}

func (l *lbController) Create(serviceObject *v1.Service, nodeList []string) (string,error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d",l.ipAddr,l.port),grpc.WithInsecure())
	if err != nil {
		log.Println(err)

		return "", err
	}

	log.Println("Connection Successfully")
	client := pb.NewExternalLBClient(conn)

	// TODO: Change this. for Debug only
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()
	dataObject := createDataObject(serviceObject,nodeList)
	log.Println("Send Object:")
	log.Println(dataObject)

	reslt, err :=client.Create(ctx,dataObject)
	if err != nil {
		log.Println(err)
	}
	return reslt.Addr, err

	//url := fmt.Sprintf("http://%s:%d/Create",l.ipAddr,l.port)
	//
	//return l.sendData(url,l.marshalData(serviceObject,nodeList))
}

func (l *lbController) Update(serviceObject *v1.Service, nodeList []string) (string,error) {
	url := fmt.Sprintf("http://%s:%d/Update",l.ipAddr,l.port)

	return l.sendData(url,l.marshalData(serviceObject,nodeList))
}

func (l *lbController) Delete(serviceObject *v1.Service, nodeList []string) (string,error) {
	url := fmt.Sprintf("http://%s:%d/Delete",l.ipAddr,l.port)

	return l.sendData(url,l.marshalData(serviceObject,nodeList))
}

func (l *lbController) NodesChange(nodeList []string) (string,error) {
	url := fmt.Sprintf("http://%s:%d/Nodes",l.ipAddr,l.port)

	emptyService := v1.Service{}

	return l.sendData(url,l.marshalData( &emptyService,nodeList))
}