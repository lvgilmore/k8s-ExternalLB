# K8S-ExternalLB
external loadbalancer for kubernetes without any cloud provider

## Deployment
### For the K8S Controller and the Haproxy Controller
* You need to have a kubernetes cluster running
* change the configuration in the docker folder for the haproxy controller and the k8s controller
* Deploy the kubernetes configurations
```
kubectl apply -f docker/k8s-controller/k8s-deployment.yaml
kubectl apply -f docker/haproxycluster/k8s-deployment.yaml
```
### For the Agent
* Need to have a docker host installed
* Run the docker run command with the variables

```$xslt
docker run -d --name lb-agent --privileged --cap-add NET_ADMIN --env Prod=TRUE --env interfaceName=<Interface> --env state=<MASTER-SLAVE> --env adminIpInterface=<HostIp> --env adminPort=<port> --network host sebassch/lb-agent
```

* Dont Run the agent on a kubernetes cluster (kubeproxy will conflict with the haproxy listener)

## Work in progress

### HaproxyCluster Controller Todo:
* Add more logging messages
* add tests with ginkgo
* create config map
* create deployment


### HaproxyCluster Agent Todo:
* create a ansible deployment
* check the cluster
* create a production ready haproxy configuration
* create a production ready keepalived configuration
* test member down

### k8s controller
* add more logging messages
* add tests with ginkgo


### Run the Grpc
```protoc externallb/externalLB.proto --go_out=plugins=grpc:```