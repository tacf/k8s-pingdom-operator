/*
Copyright 2019 github.com/tacf.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"time"

	//corev1 "k8s.io/api/core/v1"

	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

	//"k8s.io/client-go/kubernetes/scheme"
	//typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	//"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	pingdomv1alpha1 "github.com/tacf/k8s-pingdom-operator/pkg/apis/pingdomcontroller/v1alpha1"
	clientset "github.com/tacf/k8s-pingdom-operator/pkg/generated/clientset/versioned"

	//pingdomoperatorscheme "github.com/tacf/k8s-pingdom-operator/pkg/generated/clientset/versioned/scheme"
	informers "github.com/tacf/k8s-pingdom-operator/pkg/generated/informers/externalversions/pingdomcontroller/v1alpha1"
	listers "github.com/tacf/k8s-pingdom-operator/pkg/generated/listers/pingdomcontroller/v1alpha1"

	pingdom "github.com/tacf/k8s-pingdom-operator/pkg/pingdom"
)

const controllerAgentName = "pingdom-controller"

// Controller is the controller implementation for PingdomOperator resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// pingdomclientset is a clientset for our own API group
	pingdomclientset clientset.Interface

	pingdomOperatorsLister listers.PingdomOperatorLister
	pingdomOperatorsSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface

	pingdomCredentials pingdom.Credentials
}

const (
	add    = iota
	update = iota
	delete = iota
)

// QueueObject allows for aditional information to be added
// to the queue for processing
type QueueObject struct {
	opType uint8
	value  *pingdomv1alpha1.PingdomOperator
}

// NewController returns a new sample controller
func NewController(
	kubeclientset kubernetes.Interface,
	pingdomclientset clientset.Interface,
	pingdomCredentials pingdom.Credentials,
	pingdomOperatorInformer informers.PingdomOperatorInformer) *Controller {

	controller := &Controller{
		kubeclientset:          kubeclientset,
		pingdomclientset:       pingdomclientset,
		pingdomOperatorsLister: pingdomOperatorInformer.Lister(),
		pingdomOperatorsSynced: pingdomOperatorInformer.Informer().HasSynced,
		workqueue:              workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "PingdomOperators"),
		pingdomCredentials:     pingdomCredentials,
	}

	enqueueObject := func(keyFunc func(interface{}) (string, error), obj interface{}, op uint8) {
		if _, err := keyFunc(obj); err == nil {
			controller.workqueue.Add(QueueObject{opType: op, value: obj.(*pingdomv1alpha1.PingdomOperator)})
		}
	}

	klog.Info("Setting up event handlers")
	pingdomOperatorInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			enqueueObject(cache.MetaNamespaceKeyFunc, obj, add)
		},
		UpdateFunc: func(old, new interface{}) {
			enqueueObject(cache.MetaNamespaceKeyFunc, new, update)
		},
		DeleteFunc: func(obj interface{}) {
			enqueueObject(cache.DeletionHandlingMetaNamespaceKeyFunc, obj, delete)
		},
	})
	return controller
}

// Run executes the controller
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	klog.Info("Starting PingdomOperator controller")

	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.pingdomOperatorsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var qObj QueueObject
		var ok bool

		if qObj, ok = obj.(QueueObject); !ok {
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		if err := c.syncHandler(obj.(QueueObject)); err != nil {
			c.workqueue.AddRateLimited(qObj)
			return fmt.Errorf("error syncing '%s': %s, requeuing", getResourceNamespaceKey(qObj), err.Error())
		}

		c.workqueue.Forget(obj)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func getResourceNamespaceKey(obj QueueObject) string {
	var key string
	//TODO FIX ERROR HANDLING
	if (obj.opType == add) && (obj.opType == update) {
		key, _ = cache.MetaNamespaceKeyFunc(obj.value)
	} else {
		key, _ = cache.DeletionHandlingMetaNamespaceKeyFunc(obj.value)
	}

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return ""
	}

	return fmt.Sprintf("%s/%s", namespace, name)
}

func (c *Controller) syncHandler(obj QueueObject) error {
	// Create Pingdom client
	client, err := pingdom.NewClient(c.pingdomCredentials)
	if err != nil {
		klog.Error("Could not create Pingdom client")
		return nil
	}

	// Convert Custom resource spec int a list of Pingdom Checks
	var checksList []pingdom.BasicHTTPCheck
	for i, check := range obj.value.Spec.Checks {
		checksList = append(checksList, pingdom.BasicHTTPCheck{Name: fmt.Sprintf("%s_%d", obj.value.GetName(), i), URL: check.URL, Interval: check.Interval})
	}

	switch obj.opType {
	case add:
		client.AddCheck(checksList)
	case update:
		client.UpdateCheck(checksList)
	case delete:
		client.DeleteCheck(checksList)
	}

	return nil
}
