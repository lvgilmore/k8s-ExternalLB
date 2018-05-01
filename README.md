# K8S-ExternalLB
external loadbalancer for kubernetes without any cloud provider

## Work in progress

### HaproxyCluster Controller Todo:
* check if need to send update
* Add more logging messages
* add tests with ginkgo
* create config map
* create deployment


### HaproxyCluster Agent Todo:
* check the cluster
* create a production ready haproxy configuration
* create a production ready keepalived configuration
* test member down

### k8s controller
* add more logging messages
* add tests with ginkgo

## Agent Command
```$xslt
docker run -d --name lb-agent --privileged --cap-add ALL --env Prod=TRUE --env interfaceName=ens33 --env state=MASTER --env adminIpInterface=192.168.1.123 --env adminPort=9090 --network host lb-agent
```