package providers

import (
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/shirou/gopsutil/v3/mem"
)

type ProviderType string

const (
	ProviderCUDA   ProviderType = "CUDAExecutionProvider"
	ProviderCoreML ProviderType = "CoreMLExecutionProvider"
	ProviderCPU    ProviderType = "CPUExecutionProvider"
)

type HardwareConfig struct {
	Provider    ProviderType
	MemoryLimit int64
	DeviceID    int
	Options     map[string]string
}

func DetectHardware() (*HardwareConfig, error) {
	log.Printf("🔍 Detectando hardware... (GOOS=%s, GOARCH=%s)", runtime.GOOS, runtime.GOARCH)

	config := &HardwareConfig{
		Provider: ProviderCPU,
		DeviceID: 0,
		Options:  make(map[string]string),
	}

	// Debug: verificar sistema operacional
	switch runtime.GOOS {
	case "darwin": // macOS
		log.Println("🍎 Sistema macOS detectado")
		return detectMacHardware(config)

	case "linux", "windows":
		log.Printf("🐧/🪟 Sistema %s detectado", runtime.GOOS)
		return detectCUDAHardware(config)

	default:
		log.Printf("⚠️  Sistema operacional %s não reconhecido, usando CPU", runtime.GOOS)
		config.MemoryLimit = getHalfSystemRAM()
		return config, nil
	}
}

// detectMacHardware configura CoreML para macOS
func detectMacHardware(config *HardwareConfig) (*HardwareConfig, error) {
	// Verificar se CoreML está explicitamente desabilitado
	if os.Getenv("DISABLE_COREML") == "1" {
		log.Println("⚙️  CoreML desabilitado via DISABLE_COREML=1")
		config.Provider = ProviderCPU
		config.MemoryLimit = getHalfSystemRAM()
		return config, nil
	}

	log.Println("✅ Habilitando CoreML Execution Provider")
	config.Provider = ProviderCoreML
	config.MemoryLimit = getHalfSystemRAM()

	// Configurações do CoreML para GPU + Neural Engine
	config.Options = map[string]string{
		"MLComputeUnits":           getMLComputeUnits(),
		"ModelFormat":              "MLProgram", // Melhor performance no macOS 12+
		"RequireStaticInputShapes": "0",
		"EnableOnSubgraphs":        "0",
	}

	log.Printf("🍎 CoreML configurado:")
	log.Printf("   - MLComputeUnits: %s", config.Options["MLComputeUnits"])
	log.Printf("   - ModelFormat: %s", config.Options["ModelFormat"])
	log.Printf("   - Memória disponível: %.2f GB",
		float64(config.MemoryLimit)/(1024*1024*1024))

	return config, nil
}

// getMLComputeUnits retorna o modo de computação
func getMLComputeUnits() string {
	if units := os.Getenv("COREML_COMPUTE_UNITS"); units != "" {
		// Valores: ALL, CPUOnly, CPUAndNeuralEngine, CPUAndGPU
		return units
	}
	// Padrão: usa CPU, GPU e Neural Engine
	return "ALL"
}

// detectCUDAHardware configura NVIDIA CUDA
func detectCUDAHardware(config *HardwareConfig) (*HardwareConfig, error) {
	cudaVisible := os.Getenv("CUDA_VISIBLE_DEVICES")
	if cudaVisible == "-1" {
		log.Println("⚙️  CUDA desabilitado via CUDA_VISIBLE_DEVICES=-1")
		config.Provider = ProviderCPU
		config.MemoryLimit = getHalfSystemRAM()
		return config, nil
	}

	log.Println("🚀 Habilitando CUDA Execution Provider")

	if cudaVisible != "" {
		if deviceID, err := strconv.Atoi(cudaVisible); err == nil {
			config.DeviceID = deviceID
		}
	}

	config.Provider = ProviderCUDA
	config.MemoryLimit = getHalfFreeGPUMemory(config.DeviceID)

	config.Options = map[string]string{
		"device_id":              strconv.Itoa(config.DeviceID),
		"gpu_mem_limit":          strconv.FormatInt(config.MemoryLimit, 10),
		"arena_extend_strategy":  "kSameAsRequested",
		"cudnn_conv_algo_search": "DEFAULT",
	}

	log.Printf("🚀 CUDA configurado:")
	log.Printf("   - Device ID: %d", config.DeviceID)
	log.Printf("   - Limite de memória: %.2f GB",
		float64(config.MemoryLimit)/(1024*1024*1024))

	return config, nil
}

func getHalfSystemRAM() int64 {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Println("⚠️  Erro ao detectar RAM, usando 4GB como padrão")
		return 4 * 1024 * 1024 * 1024
	}

	halfRAM := int64(v.Available / 2)

	// Limite máximo de segurança: 16GB
	maxRAM := int64(16 * 1024 * 1024 * 1024)
	if halfRAM > maxRAM {
		halfRAM = maxRAM
	}

	return halfRAM
}

func getHalfFreeGPUMemory(deviceID int) int64 {
	if limitGB := os.Getenv("GPU_MEMORY_LIMIT_GB"); limitGB != "" {
		if gb, err := strconv.ParseFloat(limitGB, 64); err == nil {
			return int64(gb * 1024 * 1024 * 1024)
		}
	}
	return 4 * 1024 * 1024 * 1024 // 4GB padrão
}

// GetProviderPriority retorna ordem de providers para fallback
func GetProviderPriority(config *HardwareConfig) []ProviderType {
	providers := []ProviderType{}

	if config.Provider != ProviderCPU {
		providers = append(providers, config.Provider)
	}

	// CPU sempre como fallback
	providers = append(providers, ProviderCPU)

	return providers
}
