package runtime

import (
	"fmt"
	"log"
	"sync"

	ort "github.com/yalue/onnxruntime_go"
)

// ONNXRuntime gerencia o ciclo de vida do ONNX Runtime como singleton
type ONNXRuntime struct {
	initialized bool
	libPath     string
	mu          sync.RWMutex
}

var (
	instance *ONNXRuntime
	once     sync.Once
)

// GetInstance retorna a instância singleton do ONNX Runtime
func GetInstance() *ONNXRuntime {
	once.Do(func() {
		instance = &ONNXRuntime{
			initialized: false,
		}
	})
	return instance
}

// Initialize inicializa o ONNX Runtime com o caminho da biblioteca
// Pode ser chamado múltiplas vezes, mas só inicializa uma vez
func (r *ONNXRuntime) Initialize(libPath string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.initialized {
		log.Printf("ℹ️  ONNX Runtime já inicializado")
		return nil
	}

	if libPath == "" {
		return fmt.Errorf("caminho da biblioteca ONNX não pode ser vazio")
	}

	// Define o caminho da biblioteca
	ort.SetSharedLibraryPath(libPath)
	r.libPath = libPath

	// Inicializa o ambiente
	if err := ort.InitializeEnvironment(); err != nil {
		return fmt.Errorf("falha ao inicializar ONNX Runtime: %w", err)
	}

	r.initialized = true
	log.Printf("✅ ONNX Runtime inicializado")
	log.Printf("   Biblioteca: %s", libPath)

	return nil
}

// IsInitialized verifica se o runtime foi inicializado
func (r *ONNXRuntime) IsInitialized() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.initialized
}

// GetLibPath retorna o caminho da biblioteca configurado
func (r *ONNXRuntime) GetLibPath() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.libPath
}

// Destroy limpa os recursos do ONNX Runtime
// Deve ser chamado ao finalizar a aplicação
func (r *ONNXRuntime) Destroy() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.initialized {
		return
	}

	ort.DestroyEnvironment()
	r.initialized = false
	log.Println("✅ ONNX Runtime destruído")
}

// MustInitialize é como Initialize mas entra em pânico se falhar
// Útil para inicialização na main()
func (r *ONNXRuntime) MustInitialize(libPath string) {
	if err := r.Initialize(libPath); err != nil {
		log.Fatalf("❌ Falha crítica ao inicializar ONNX Runtime: %v", err)
	}
}

// WithRuntime é um helper para garantir que o runtime está inicializado
// antes de executar uma função
func WithRuntime(libPath string, fn func() error) error {
	runtime := GetInstance()
	if err := runtime.Initialize(libPath); err != nil {
		return err
	}
	return fn()
}
