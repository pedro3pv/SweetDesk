//go:build windows

package main

import _ "embed"

//go:embed lib/windows/onnxruntime.dll
var onnxLib []byte

func getONNXLibrary() ([]byte, string) {
	return onnxLib, "onnxruntime.dll"
}
