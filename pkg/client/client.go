package client

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

func NewClient() kubernetes.Interface {
	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Errorln(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Errorln(err.Error())
	}
	return clientset
}
