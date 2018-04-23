package main

import (
	"flag"

	"github.com/SchSeba/k8s-ExternalLB/pkg/k8s/loadbalance"

	"k8s.io/client-go/tools/clientcmd"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path/filepath"
	"os"
	//"fmt"
	"github.com/SchSeba/k8s-ExternalLB/pkg/k8s/node"
)

const  (
	ipAddr = "192.168.1.124"
	port = 8080
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func NodeControllerStart(clientset *kubernetes.Clientset) *node.NodeController {
	podListWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(),"Nodes","", fields.Everything())

	// create the workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the pod key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the Pod than the version which was responsible for triggering the update.
	indexer, informer := cache.NewIndexerInformer(podListWatcher, &v1.Node{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				event := node.EventData{EventType:"Create",Name:key,Data:obj.(*v1.Node)}
				queue.Add(event)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				event := node.EventData{EventType:"Update",Name:key,Data:new.(*v1.Node)}
				queue.Add(event)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				event := node.EventData{EventType:"Delete",Name:key,Data:obj.(*v1.Node)}
				queue.Add(event)
			}
		},
	}, cache.Indexers{})

	nodeController := node.NewNodeController(queue, indexer, informer,clientset)

	return nodeController
}

func ServiceControllerStart(clientset *kubernetes.Clientset, nodeController *node.NodeController) *loadbalance.Controller {
	// create the pod watcher
	podListWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(),"services","", fields.Everything())

	// create the workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the pod key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the Pod than the version which was responsible for triggering the update.
	indexer, informer := cache.NewIndexerInformer(podListWatcher, &v1.Service{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil &&  obj.(*v1.Service).Spec.Type == "LoadBalancer" {
				event := loadbalance.EventData{EventType:"Create",Name:key,Data:obj.(*v1.Service)}
				queue.Add(event)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil &&  new.(*v1.Service).Spec.Type == "LoadBalancer" {
				event := loadbalance.EventData{EventType:"Update",Name:key,Data:new.(*v1.Service)}
				queue.Add(event)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil &&  obj.(*v1.Service).Spec.Type == "LoadBalancer" {
				event := loadbalance.EventData{EventType:"Delete",Name:key,Data:obj.(*v1.Service)}
				queue.Add(event)
			}
		},
	}, cache.Indexers{})

	lbController := loadbalance.NewController(clientset,queue, indexer, informer, ipAddr, port,nodeController)

	return lbController
}



func main() {
	//var kubeconfig string
	var master string

	//flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	//flag.StringVar(&master, "master", "", "master url")
	//flag.Parse()

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// creates the connection
	config, err := clientcmd.BuildConfigFromFlags(master, *kubeconfig)
	if err != nil {
		glog.Fatal(err)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatal(err)
	}


	nodeController := NodeControllerStart(clientset)
	lbController := ServiceControllerStart(clientset,nodeController)


	// Now let's start the controller
	stop := make(chan struct{})
	defer close(stop)
	go nodeController.Run(1, stop)
	go lbController.Run(1, stop)
	// Wait forever
	select {}
}