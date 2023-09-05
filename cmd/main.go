package main

import (
	"cluster-api-server/pkg/rest"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"k8s.io/klog"
)

func main() {
	router := httprouter.New()
	rest.AddResource(router, new(rest.NodeMetricResource))
	rest.AddResource(router, new(rest.NodeResource))
	rest.AddResource(router, new(rest.ClusterResource))

	klog.Fatal(http.ListenAndServe(":30850", router))
}
