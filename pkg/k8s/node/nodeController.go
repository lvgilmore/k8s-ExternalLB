package node

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"github.com/golang/glog"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

type NodeController struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
	clientset *kubernetes.Clientset
	NodeList map[string]struct{}
	LBChannel chan int
}

type EventData struct {
	EventType string
	Name string
	Data *v1.Node
}

func NewNodeController(queue workqueue.RateLimitingInterface,
					   indexer cache.Indexer,
	                   informer cache.Controller,
	                   clientset *kubernetes.Clientset) *NodeController {

	nodeController := NodeController{queue:queue,
									indexer:indexer,
									informer:informer,
									clientset:clientset,
									NodeList:make(map[string]struct{}),
									LBChannel:make(chan int, 5)}

	nodeController.SyncNodeData()
	return &nodeController
}

func (n *NodeController)GetNodes() []string {
	keys := make([]string, len(n.NodeList))
	i := 0
	for k := range n.NodeList {
		keys[i] = k
		i++
	}

	return keys
}

func (n *NodeController)SyncNodeData() {
	nodes, err := n.clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, node := range nodes.Items{
		for _, IpAddr := range node.Status.Addresses {
			if IpAddr.Type == "InternalIP" {
				n.NodeList[IpAddr.Address] = struct{}{}
			}
		}
	}
}

func (c *NodeController) processNextItem() bool {
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
	err := c.syncNode(event)
	// Handle the error if something went wrong during the execution of the business logic
	c.handleErr(err, event)
	return true
}

// syncLoadBalancer is the business logic of the controller.
func (c *NodeController) syncNode(event EventData) error {
	if event.EventType == "Create" {
		for _, value := range event.Data.Status.Addresses {
			if value.Type == "InternalIP" {
				if _,ok := c.NodeList[value.Address]; !ok {
					c.SyncNodeData()
					c.LBChannel <- 1

					fmt.Printf("%s for Node %s\n", event.EventType, event.Data.GetName())

					return nil
				}
			}
		}
	} else {
		c.SyncNodeData()
		c.LBChannel <- 1
	}

	fmt.Printf("%s for Node %s\n", event.EventType, event.Data.GetName())

	return nil
}

// handleErr checks if an error happened and makes sure we will retry later.
func (c *NodeController) handleErr(err error, key interface{}) {
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

func (c *NodeController) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer c.queue.ShutDown()
	glog.Info("Starting Node controller")

	go c.informer.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	glog.Info("Stopping Node controller")
}

func (c *NodeController) runWorker() {
	for c.processNextItem() {
	}
}