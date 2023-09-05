package rest

import (
	"cluster-api-server/pkg/client"
	"cluster-api-server/pkg/client/metric"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	corev1 "k8s.io/api/core/v1"
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
	ClusterName string   `json:"clusterName"`
	MasterNode  string   `json:"masterNode"`
	Nodes       []string `json:"nodes"`
	TotalGPU    string   `json:"totalGPU"`
}

type NodeResponse struct {
	ClusterName    string         `json:"clusterName"`
	VirtualGPU     int32          `json:"virtualGPU"`
	Age            string         `json:"age"`
	GpuPods        map[string]int `json:"gpuPods"`
	GpuPodForPrint map[int]string `json:"gpuPodForPrint"`
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

	res := &ClusterResponse{
		ClusterName: clusterName,
		Nodes:       make([]string, 0),
	}
	totalGPU := int32(0)
	nodeList, err := kubeClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Errorln(err)
	}

	for _, node := range nodeList.Items {
		res.Nodes = append(res.Nodes, node.Name)
		if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
			res.MasterNode = node.Name
		}
		if _, ok := node.Labels["clusterName"]; !ok {
			node.Labels["clusterName"] = clusterName
			_, updateErr := kubeClient.CoreV1().Nodes().Update(context.TODO(), &node, metav1.UpdateOptions{})
			if updateErr != nil {
				klog.Errorln(updateErr)
			}
		}
		nodeIP := ""
		for _, address := range node.Status.Addresses {
			if address.Type == corev1.NodeInternalIP {
				nodeIP = address.Address
			}
		}
		conn, err := grpc.Dial(nodeIP+client.HOSTSERVERPORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			klog.Errorln(err)
		}
		defer conn.Close()
		req := &client.NodeRequest{ClusterName: clusterName}
		travelClient := client.NewTravelerClient(conn)
		res, err := travelClient.Node(context.Background(), req)
		if err != nil {
			klog.Errorln(err)
		}
		if res.GPU > 0 {
			node.Labels["gpu"] = "on"
		}
		totalGPU += res.GPU
	}
	return Response{200, "", res}
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
	nodeName := ps.ByName("nodeName")
	clusterName := ""

	res := &NodeResponse{
		GpuPods: make(map[string]int),
	}
	node, err := kubeClient.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		klog.Errorln(err)
	}
	res.ClusterName = node.Labels["clusterName"]

	nodeIP := ""
	for _, address := range node.Status.Addresses {
		if address.Type == corev1.NodeInternalIP {
			nodeIP = address.Address
		}
	}
	conn, err := grpc.Dial(nodeIP + client.HOSTSERVERPORT)
	if err != nil {
		klog.Errorln(err)
	}
	defer conn.Close()
	req := &client.NodeRequest{ClusterName: clusterName}
	travelClient := client.NewTravelerClient(conn)
	nodegpures, err := travelClient.Node(context.Background(), req)
	if err != nil {
		klog.Errorln(err)
	}
	res.Age = time.Since(node.ObjectMeta.CreationTimestamp.Time).String()
	res.VirtualGPU = nodegpures.GPU * 20

	podList, err := kubeClient.CoreV1().Pods(corev1.NamespaceAll).List(context.Background(), metav1.ListOptions{FieldSelector: "spec.nodeName=" + nodeName})
	if err != nil {
		klog.Errorln(err)
	}
	currentGPU := 0
	for _, pod := range podList.Items {
		for _, container := range pod.Spec.Containers {
			quantity := container.Resources.Limits["nvidia.com/gpu"]
			gpuCount, _ := quantity.AsInt64()
			res.GpuPods[pod.Name] += int(gpuCount)
			for i := 0; i < int(gpuCount); i++ {
				res.GpuPodForPrint[currentGPU] = pod.Name
				currentGPU++
			}
		}
	}
	return Response{200, "", res}
}

// /node/:nodeName/metrics
type NodeMetricResource struct {
	PostNotSupported
	PutNotSupported
	DeleteNotSupported
}

func (NodeMetricResource) Uri() string {
	return "/node/:nodeName/metrics"
}

func (NodeMetricResource) Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	nodeName := ps.ByName("nodeName")
	podIP := ""

	podList, err := kubeClient.CoreV1().Pods("keti-system").List(context.Background(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		klog.Errorln(err)
	}

	for _, pod := range podList.Items {
		if strings.Contains(pod.Name, "metric-collector") {
			podIP = strings.ReplaceAll(pod.Status.PodIP, ".", "-")
		}
	}

	conn, err := grpc.Dial(podIP + ".keti-system.pod.cluster.local:50051")
	if err != nil {
		klog.Errorln(err)
	}
	defer conn.Close()
	req := &metric.Request{}
	metricClient := metric.NewMetricGathererClient(conn)
	node, err := kubeClient.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		klog.Errorln(err)
	}

	var grpcRes *metric.Response

	if strings.Compare(node.Labels["gpu"], "on") == 0 {
		grpcRes, err = metricClient.GPU(context.Background(), req)
		if err != nil {
			klog.Errorln(err)
		}
	} else {
		grpcRes, err = metricClient.Node(context.Background(), req)
		if err != nil {
			klog.Errorln(err)
		}
	}

	return Response{200, "", metric.Convert(grpcRes)}
}
