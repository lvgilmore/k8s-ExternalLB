package loadbalance

import (
	"k8s.io/api/core/v1"
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"fmt"
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


func (l *lbController) Create(serviceObject *v1.Service, nodeList []string) (string,error) {
	url := fmt.Sprintf("http://%s:%d/Create",l.ipAddr,l.port)

	return l.sendData(url,l.marshalData(serviceObject,nodeList))
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