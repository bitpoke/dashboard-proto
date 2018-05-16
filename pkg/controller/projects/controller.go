/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package projects

import (
	apiextenstions_util "github.com/appscode/kutil/apiextensions/v1beta1"
	"github.com/golang/glog"
	apiextenstions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	projectsApi "github.com/presslabs/dashboard/pkg/apis/projects/v1alpha1"
	"github.com/presslabs/dashboard/pkg/controller"
)

const (
	controllerName = "projects-controller"
	maxRetries     = 5
	threadiness    = 4
)

type Controller struct {
	*controller.Context
	*ProjectsContext
}

func NewController(ctx *controller.Context) (c *Controller, err error) {
	c = &Controller{}
	c.Context = ctx

	c.initProjectsWorker()

	return
}

// Run starts the control loop for the Projects Controller
func (c *Controller) Run(stopCh <-chan struct{}) {
	crds := []*apiextenstions.CustomResourceDefinition{
		projectsApi.ResourceProjectCRD,
	}

	if c.InstallCRDs {
		if err := c.installCRDs(crds); err != nil {
			glog.Fatalf(err.Error())
			return
		}
	}
	if err := c.waitForCRDs(crds); err != nil {
		glog.Fatalf(err.Error())
		return
	}

	glog.V(2).Infof("Starting shared informer factories")
	c.KubeSharedInformerFactory.Start(stopCh)
	c.DashboardSharedInformerFactory.Start(stopCh)
	// Wait for all involved caches to be synced, before processing items from the queue is started
	for t, v := range c.KubeSharedInformerFactory.WaitForCacheSync(stopCh) {
		if !v {
			glog.Fatalf("%v timed out waiting for caches to sync", t)
			return
		}
	}
	for t, v := range c.DashboardSharedInformerFactory.WaitForCacheSync(stopCh) {
		if !v {
			glog.Fatalf("%v timed out waiting for caches to sync", t)
			return
		}
	}
	glog.V(2).Infof("Informer cache synced")

	glog.Infof("Starting %s control loops", controllerName)

	c.projectsQueue.Run(stopCh)

	<-stopCh
	glog.Infof("Stopping %s control loops", controllerName)
}

func (c *Controller) installCRDs(crds []*apiextenstions.CustomResourceDefinition) error {
	glog.Info("Registering Custom Resource Definitions")

	if err := apiextenstions_util.RegisterCRDs(c.CRDClient, crds); err != nil {
		return err
	}

	return nil
}

func (c *Controller) waitForCRDs(crds []*apiextenstions.CustomResourceDefinition) error {
	glog.Info("Waiting for Custom Resource Definitions to become available")
	return apiextenstions_util.WaitForCRDReady(c.CRDClient.RESTClient(), crds)
}
