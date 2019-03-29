package main

import (
	"fmt"
	"os"
)

func getWasmceiptionEnv() error {
	if wasmCeptionFlag = os.Getenv("VNT_WASMCEPTION"); wasmCeptionFlag == "" {
		return fmt.Errorf("未找到VNT_WASMCEPTION的环境变量，请确认已设置VNT_WASMCEPTION")
	}
	return nil
}

func getIncludeEnv() error {
	if vntIncludeFlag = os.Getenv("VNT_INCLUDE"); vntIncludeFlag == "" {
		return fmt.Errorf("未找到VNT_INCLUDE的环境变量，请确认已设置VNT_INCLUDE")
	}
	return nil
}
