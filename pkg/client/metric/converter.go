package metric

import "fmt"

type MetricResponse struct {
	GPUCore   string `json:"gpuCore,omitempty"`
	GPUMemory string `json:"gpuMemory,omitempty"`
	GPUPower  string `json:"gpuPower,omitempty"`
	CPUCore   string `json:"cpuCore,omitempty"`
	Memory    string `json:"memory,omitempty"`
	Storage   string `json:"storage,omitempty"`
	NetworkRX string `json:"networkRX,omitempty"`
	NetworkTX string `json:"networkTX,omitempty"`
}

const (
	GPUCORETAG   = "Host_GPU_Core_Gauge"
	GPUMEMORYTAG = "Host_GPU_Memory_Gauge"
	GPUPOWERTAG  = "Host_GPU_Power_Gauge"
	CPUCORETAG   = "CPU_Core_Gauge"
	MEMORYTAG    = "Memory_Gauge"
	STORAGETAG   = "Storage_Gauge"
	NETWORKRXTAG = "Network_Gauge"
	NETWORKTXTAG = "Network_Counter"
)

func Convert(grpcres *Response) *MetricResponse {
	gpuCoreVal := grpcres.Message[GPUCORETAG].Metric[0]
	gpuMemoryVal := grpcres.Message[GPUMEMORYTAG].Metric[0]
	gpuPowerVal := grpcres.Message[GPUPOWERTAG].Metric[0]
	cpuCoreVal := grpcres.Message[CPUCORETAG].Metric[0]
	memoryVal := grpcres.Message[MEMORYTAG].Metric[0]
	storageVal := grpcres.Message[STORAGETAG].Metric[0]
	networkRXVal := grpcres.Message[NETWORKRXTAG].Metric[0]
	networkTXVal := grpcres.Message[NETWORKTXTAG].Metric[0]

	return &MetricResponse{
		GPUCore:   fmt.Sprint(gpuCoreVal.Gauge.GetValue()),
		GPUMemory: fmt.Sprint(gpuMemoryVal.Gauge.GetValue()),
		GPUPower:  fmt.Sprint(gpuPowerVal.Gauge.GetValue()),
		CPUCore:   fmt.Sprint(cpuCoreVal.Gauge.GetValue()),
		Memory:    fmt.Sprint(memoryVal.Gauge.GetValue()),
		Storage:   fmt.Sprint(storageVal.Gauge.GetValue()),
		NetworkRX: fmt.Sprint(networkRXVal.Gauge.GetValue()),
		NetworkTX: fmt.Sprint(networkTXVal.Gauge.GetValue()),
	}
}
