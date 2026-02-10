package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/owulveryck/onnx-go"
	"github.com/owulveryck/onnx-go/backend/x/gorgonnx"
)

//go:embed models/**/*.onnx
var modelsFS embed.FS

func main() {
	models := []string{
		"models/realcugan/realcugan-pro.onnx",
		"models/lsdir/4xLSDIR.onnx",
	}

	for _, modelPath := range models {
		fmt.Printf("\n🔍 Testando: %s\n", modelPath)

		modelData, err := modelsFS.ReadFile(modelPath)
		if err != nil {
			log.Printf("  ❌ Não encontrado: %v", err)
			continue
		}

		backend := gorgonnx.NewGraph()
		model := onnx.NewModel(backend)

		if err := model.UnmarshalBinary(modelData); err != nil {
			log.Printf("  ❌ INCOMPATÍVEL: %v", err)
			log.Printf("     Provável causa: operadores não suportados")
			continue
		}

		fmt.Printf("  ✅ COMPATÍVEL (modelo carregado com sucesso)\n")
		fmt.Printf("     Inputs: %d\n", len(model.GetInputTensors()))
		outputs, _ := model.GetOutputTensors()
		fmt.Printf("     Outputs: %d\n", len(outputs))
	}

	fmt.Println("\n═══════════════════════════════════════════")
	fmt.Println("Se todos compatíveis: prossiga com migração")
	fmt.Println("Se incompatível: use solução com DLL ou GoMLX")
	fmt.Println("═══════════════════════════════════════════")
}
