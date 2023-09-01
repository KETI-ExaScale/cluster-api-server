package deviceplugin

import (
	"context"
	"flag"
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func main() {
	var kubeconfig *string
	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()

	config, err := getRestConfig(kubeconfig)
	if err != nil {
		klog.Fatalf("Error creating kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error creating clientset: %v", err)
	}

	nodeName := os.Getenv("NODE_NAME") // Assuming you're running the code inside a pod

	deviceInfo, err := getDevicePluginInfo(clientset, nodeName)
	if err != nil {
		klog.Fatalf("Error getting device plugin info: %v", err)
	}

	printDevicePluginInfo(deviceInfo)
}

func getRestConfig(kubeconfigPath *string) (*rest.Config, error) {
	if kubeconfigPath == nil || *kubeconfigPath == "" {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	config, err := rest.BuildConfigFromFlags("", *kubeconfigPath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func getDevicePluginInfo(clientset *kubernetes.Clientset, nodeName string) (*v1beta1.DevicePluginList, error) {
	devicePluginClient := clientset.DevicepluginV1beta1().DevicePlugins(nodeName)

	devicePluginList, err := devicePluginClient.List(context.Background(), v1beta1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return devicePluginList, nil
}

func printDevicePluginInfo(devicePluginList *v1beta1.DevicePluginList) {
	for _, plugin := range devicePluginList.Items {
		fmt.Printf("Device Plugin Name: %s\n", plugin.ObjectMeta.Name)
		fmt.Printf("Device Plugin Endpoint: %s\n", plugin.Status.Endpoint)
		fmt.Printf("Device Plugin ResourceName: %s\n", plugin.Status.ResourceName)
		fmt.Println("Devices:")
		for _, dev := range plugin.Status.Devices {
			fmt.Printf("  %s: %s\n", dev.ID, dev.Health)
		}
		fmt.Println("======================")
	}
}
