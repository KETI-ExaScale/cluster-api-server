package rest

import (
	"cluster-api-server/pkg/client"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

var kubeClient kubernetes.Interface

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ClusterResponse struct {
	ClusterName string `json:"clusterName"`
	MasterNode  string `json:"masterNode"`
	Nodes       string `json:"nodes"`
	TotalGPU    string `json:"totalGPU"`
}

type NodeResponse struct {
	ClusterName string         `json:"clusterName"`
	VirtualGPU  string         `json:"virtualGPU"`
	Age         string         `json:"age"`
	GpuPods     map[string]int `json:"gpuPods"`
}
type Resource interface {
	Uri() string
	Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response
	Post(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response
	Put(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response
	Delete(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response
}

type (
	GetNotSupported    struct{}
	PostNotSupported   struct{}
	PutNotSupported    struct{}
	DeleteNotSupported struct{}
)

func (GetNotSupported) Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{405, "", nil}
}

func (PostNotSupported) Post(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{405, "", nil}
}

func (PutNotSupported) Put(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{405, "", nil}
}

func (DeleteNotSupported) Delete(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{405, "", nil}
}

func abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
}

func HttpResponse(rw http.ResponseWriter, req *http.Request, res Response) {
	content, err := json.Marshal(res)

	if err != nil {
		abort(rw, 500)
	}

	rw.WriteHeader(res.Code)
	rw.Write(content)
}

func AddResource(router *httprouter.Router, resource Resource) {
	fmt.Println("\"" + resource.Uri() + "\" api is registerd")
	kubeClient = client.NewClient()
	router.GET(resource.Uri(), func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := resource.Get(rw, r, ps)
		HttpResponse(rw, r, res)
	})
	router.POST(resource.Uri(), func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := resource.Post(rw, r, ps)
		HttpResponse(rw, r, res)
	})
	router.PUT(resource.Uri(), func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := resource.Put(rw, r, ps)
		HttpResponse(rw, r, res)
	})
	router.DELETE(resource.Uri(), func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := resource.Delete(rw, r, ps)
		HttpResponse(rw, r, res)
	})
}

// /cluster/:clusterName
type ClusterResource struct {
	PostNotSupported
	PutNotSupported
	DeleteNotSupported
}

func (ClusterResource) Uri() string {
	return "/cluster/:clusterName"
}

func (ClusterResource) Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	clusterName := ps.ByName("clusterName")

	nodeList, err := kubeClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Errorln(err)
	}

	for _, node := range nodeList.Items {
		node.Labels["clusterName"] = clusterName
		_, updateErr := kubeClient.CoreV1().Nodes().Update(context.TODO(), &node, metav1.UpdateOptions{})
	}
	return Response{200, "", mres}
}

// /node/:nodeName
type NodeResource struct {
	PostNotSupported
	PutNotSupported
	DeleteNotSupported
}

func (NodeResource) Uri() string {
	return "/node/:nodeName"
}

func (NodeResource) Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	st := time.Now().Format(time.ANSIC)
	mres := &MigrationResponse{}
	mres.MigrationStartTime = st
	mres.MigrationSource = make([]string, 0)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		klog.Errorln(err)
	}
	req := SetRequest(body)
	scn := req.Source.ClusterName + "/" + req.Source.NodeName
	fcn := req.Target.ClusterName + "/" + req.Target.NodeName
	// 디플로이먼트 탐색
	mres.MigrationSource = setSource(req.Source.ClusterName, req.Source.DepName)
	// request Migration 사용해야함
	migrationRequest(*req, mres.MigrationSource, req.Source.DepName)
	// time.Sleep(time.Second * 1)
	time.Sleep(time.Second * 5)

	ft := time.Now().Format(time.ANSIC)
	mres.MigrationStartCN = scn
	mres.MigrationFinishCN = fcn
	mres.MigrationFinishTime = ft
	// resb, err := json.Marshal(mres)
	// if err != nil {
	// 	klog.Errorln(err)
	// }

	return Response{200, "", mres}
}
