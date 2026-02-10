//go:build darwin

package main

import _ "embed"

//go:embed "/opt/homebrew/lib/libonnxruntime.dylib"
var onnxLib []byte

func getONNXLibrary() ([]byte, string) {
	return onnxLib, "libonnxruntime.dylib"
}
