package loadbalance

import (
	"fmt"
	"time"

	"github.com/golang/glog"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"github.com/SchSeba/k8s-ExternalLB/pkg/k8s/node"
	"k8s.io/client-go/util/retry"
	"k8s.io/client-go/kubernetes"
	"log"

)

type Controller struct {
	clientset *kubernetes.Clientset
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
	lbControllerInstance *lbController
	nodeController *node.NodeController
}


type EventData struct {
	EventType string
	Name string
	Data *v1.Service
}

func NewController(clientset *kubernetes.Clientset,
				   queue workqueue.RateLimitingInterface,
	               indexer cache.Indexer,
	               informer cache.Controller,
	               ipAddr string,
	               port int,
	               nodeController *node.NodeController) *Controller {

	lbControllerInstance := NewLBController(ipAddr,port)

	return &Controller{
		clientset: clientset,
		informer: informer,
		indexer:  indexer,
		queue:    queue,
		lbControllerInstance: lbControllerInstance,
		nodeController: nodeController,
	}
}

func (c *Controller) NodeChannelWatch() {
	for true{
		<- c.nodeController.LBChannel
		c.lbControllerInstance.NodesChange(c.nodeController.GetNodes())
	}
}

func (c *Controller) processNextItem() bool {
	// Wait until there is a new item in the working queue
	eventData, quit := c.queue.Get()
	event := eventData.(EventData)
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two pods with the same key are never processed in
	// parallel.
	defer c.queue.Done(event)

	// Invoke the method containing the business logic
	err := c.syncLoadBalancer(event)
	// Handle the error if something went wrong during the execution of the business logic
	c.handleErr(err, event)
	return true
}

// syncLoadBalancer is the business logic of the controller.
func (c *Controller) syncLoadBalancer(event EventData) error {
	nodes := c.nodeController.GetNodes()
	var ipAddr string
	var err error
	fmt.Printf("%s for Service %s\n",event.EventType, event.Data.GetName())

	if event.EventType == "Create" {
		ipAddr, err =c.lbControllerInstance.Create(event.Data,nodes)
	} else if event.EventType == "Update" {
		ipAddr, err = c.lbControllerInstance.Update(event.Data, nodes)
	} else if event.EventType == "Delete" {
		ipAddr, err = c.lbControllerInstance.Delete(event.Data, nodes)
	}


	if err != nil {
		fmt.Printf("Fail sending data to LB Controller error: %v\n",err)
		return err
	} else {
		return c.updateService(event.Data,ipAddr)
	}

}


// handleErr checks if an error happened and makes sure we will retry later.
func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < 5 {
		glog.Infof("Error syncing pod %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	glog.Infof("Dropping pod %q out of the queue: %v", key, err)
}

func (c *Controller) Run(threadiness int, stopCh chan struct{}) {

	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer c.queue.ShutDown()
	glog.Info("Starting Pod controller")

	go c.informer.Run(stopCh)

	go c.NodeChannelWatch()

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	// TODO: Need to fix this
	//go c.Sync()
	<-stopCh
	glog.Info("Stopping Pod controller")
}

func (c *Controller)updateService(service *v1.Service,ipAddr string) error {
	// Update Service ip address
	lbclient := c.clientset.CoreV1().Services(service.ObjectMeta.Namespace)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		// TODO: Remove this its For Debug Only
		//ipAddr = "10.0.0.0"

		service.Spec.ExternalIPs = []string{ipAddr}
		_, updateErr := lbclient.Update(service)
		if updateErr != nil {
			return updateErr
		}
		return nil
	})

	if retryErr != nil {
		fmt.Printf("Update failed: %v\n", retryErr)
		return retryErr
	}

	return nil
}

func (c *Controller)Sync() {
	tick := time.Tick(30 * time.Second)
	for {
		<-tick
		log.Println("Syncing all services with the controller")
		nodes := c.nodeController.GetNodes()
		for _,value := range c.indexer.List() {
			if value.(*v1.Service).Spec.Type == "LoadBalancer" {
				ipAddr, err := c.lbControllerInstance.Update(value.(*v1.Service), nodes)
				if err != nil {
					fmt.Printf("Fail sending data to LB Controller error: %v\n",err)
				} else if ipAddr != ""{
					c.updateService(value.(*v1.Service),ipAddr)
				}
			}
		}
	}
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}