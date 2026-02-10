//go:build linux

package main

import _ "embed"

//go:embed lib/linux/libonnxruntime.so
var onnxLib []byte

func getONNXLibrary() ([]byte, string) {
	return onnxLib, "libonnxruntime.so"
}
