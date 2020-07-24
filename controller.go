package main

import (
	samplecrdv1 "github.com/Mathew857/democrd/pkg/apis/democrd/v1"
	mydemoscheme "github.com/Mathew857/democrd/pkg/client/clientset/versioned/scheme"
	informers "github.com/Mathew857/democrd/pkg/client/informers/externalversions/democrd/v1"
	listers "github.com/Mathew857/democrd/pkg/client/listers/democrd/v1"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"k8s.io/kubernetes/staging/src/k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/kubernetes/staging/src/k8s.io/apimachinery/pkg/util/runtime"
)

const controllerAgentName = "mydemo-controller"

const (
	SuccessSynced         = "Synced"
	MessageResourceSynced = "Mydemo synced successfully"
)

type Controller struct {
	//kubeclientset is a standard kubernetes clientset
	kubeclientset   kubernetes.Interface
	mydemoslientset clientset.Interface

	demoInformer  listers.MydemoLister
	mydemosSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as sonn as a change happens. This  means
	// we can ensure we only process a fixed amount of resources at a time, and makes
	// it eay to ensure we are never processing the same item simultaneously in two
	// different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder  is an event recorder for recording Event resources to the
	// Kubenretes API
	recorder record.EventRecorder
}

// NewController returns a new mydemo controller
func NewController(
	kubeclientset kubernetes.Interface,
	mydemoslientset clientset.Interface,
	mydemoInformer informers.MydemoInformer) *Controller {

	// Create event broadcaster
	// Add sample-controller types to the default kubernetes scheme so events can be
	// loggd for sample-controller types.
	utilruntime.Must(mydemoscheme.AddToScheme(scheme.Scheme))
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:   kubeclientset,
		mydemoslientset: mydemoslientset,
		demoInformer:    mydemoInformer.Lister(),
		mydemosSynced:   mydemoInformer.Informer().HasSynced,
		workqueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Mydemos"),
		recorder:        recorder,
	}

	glog.Info("Setting up mydemo event handlers")
	mydemoInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueMydemo,
		UpdateFunc: func(old, new interface{}) {
			oldMydemo := old.(*samplecrdv1.Mydemo)
			newMydemo := new.(*samplecrdv1.Nydemo)
			if oldMydemo.ResourceVersion == newMydemo.ResourceVersion {
				return
			}
			controller.enqueueMydemo(new)
		},
		DeleteFunc: controller.enqueueMydemoForDelete,
	})

	return controller
}

// enqueueMydemo takes a Mydemo resource and converts it into a namespace/name
// string wrhich is then put onto the queue. This method should *not* be
// passed resources of any type other than Mydemo.

func (c *Controller) enqueueMydemo(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

func (c *Controller) enqueueMydemoForDelete(obj interface{}) {
	var key string
	var err error
	if key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}
