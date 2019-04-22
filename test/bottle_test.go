package tests

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

var (
	gopath        = os.Getenv("GOPATH")
	bottlePath    = path.Join(gopath, "src/github.com/vntchain/bottle/build/bin/bottle")
	contractsPath = path.Join(gopath, "src/github.com/vntchain/bottle/test/contracts")
)

func TestBottle(t *testing.T) {
	tests := []struct {
		contract string
		wanted   string
	}{
		{
			contract: "test_constructor_1.c",
			wanted:   "必须定义一个Contructor方法",
		},
		{
			contract: "test_constructor_2.c",
			wanted:   "重复定义Contructor方法,重复定义Contructor方法",
		},
		{
			contract: "test_constructor_3.c",
			wanted:   "Constructor方法的参数为不支持的类型：string* ,address* ,int32* ,int64* ,uint32* ,uint64* ,_Bool* ,uint256*",
		},
		{
			contract: "test_constructor_4.c",
			wanted:   "Constructor方法的返回值为不支持的类型：void*",
		},
		{
			contract: "test_constructor_5.c",
			wanted:   "Constructor方法的返回值为不支持的类型：void*,Constructor方法的参数为不支持的类型：int*",
		},
		{
			contract: "test_call_1.c",
			wanted:   "CALL方法至少需要一个类型为CallParams的参数",
		},
		{
			contract: "test_call_2.c",
			wanted:   "CALL方法的返回值为不支持的类型：int32*",
		},
		{
			contract: "test_call_3.c",
			wanted:   "CALL方法的第一个参数为不支持的类型：CallParams*, 支持的类型为CallParams",
		},
		{
			contract: "test_call_4.c",
			wanted:   "CALL方法的参数为不支持的类型：int32* ,int64* ,int64*",
		},
		{
			contract: "test_call_5.c",
			wanted:   "CALL方法的返回值为不支持的类型：int*,CALL方法的参数为不支持的类型：int64*",
		},
		{
			contract: "test_event_1.c",
			wanted:   "EVENT方法至少需要一个参数",
		},
		{
			contract: "test_event_2.c",
			wanted:   "EVENT方法的返回值为不支持的类型：void*,EVENT方法至少需要一个参数",
		},
		{
			contract: "test_event_3.c",
			wanted:   "EVENT方法的返回值为不支持的类型：void*",
		},
		{
			contract: "test_event_4.c",
			wanted:   "EVENT方法的参数为不支持的类型：int*",
		},
		{
			contract: "test_event_5.c",
			wanted:   "EVENT方法的参数为不支持的类型：string*",
		},
		{
			contract: "test_event_6.c",
			wanted:   "EVENT方法的参数为不支持的类型：string*,EVENT方法中的indexed写法不规范,EVENT方法中的indexed写法不规范",
		},
		{
			contract: "test_key_1.c",
			wanted:   "KEY为不支持的类型：int32 *",
		},
		{
			contract: "test_key_2.c",
			wanted:   "KEY为不支持的类型：int32 * ,int64 * ,int",
		},
		{
			contract: "test_key_3.c",
			wanted:   "KEY为不支持的类型：int ,int *",
		},
		{
			contract: "test_key_4.c",
			wanted:   "KEY为不支持的类型：int *",
		},
		{
			contract: "test_key_5.c",
			wanted:   "KEY为不支持的类型：int32 * ,int64 * ,string * ,int",
		},
		{
			contract: "test_key_6.c",
			wanted:   "KEY必须定义在construct之前,KEY为不支持的类型：int32 *",
		},
		{
			contract: "test_key_7.c",
			wanted:   "KEY为不支持的类型：int32 *",
		},
		{
			contract: "test_key_8.c",
			wanted:   "KEY为不支持的类型：int32 * ,int64 * ,int",
		},
		{
			contract: "test_export_1.c",
			wanted:   "方法的返回值为不支持的类型：int",
		},
		{
			contract: "test_export_2.c",
			wanted:   "方法的返回值为不支持的类型：int*",
		},
		{
			contract: "test_export_3.c",
			wanted:   "方法的参数为不支持的类型：string* ,int",
		},
		{
			contract: "test_export_4.c",
			wanted:   "方法的参数为不支持的类型：string* ,int",
		},
		{
			contract: "test_payable_1.c",
			wanted:   "Payable方法必须使用关键字MUTABLE进行导出",
		},
		{
			contract: "test_payable_2.c",
			wanted:   "方法的返回值为不支持的类型：int32*",
		},
		{
			contract: "test_unmutable_1.c",
			wanted:   "UNMUTABLE的方法中调用了MUTABLE或PAYABLE方法：test_export_1",
		},
		{
			contract: "test_unmutable_2.c",
			wanted:   "UNMUTABLE的方法中调用了SendFromContract或TransferFromContract方法",
		},
		{
			contract: "test_unmutable_3.c",
			wanted:   "UNMUTABLE的方法中调用了SendFromContract或TransferFromContract方法",
		},
	}
	for _, v := range tests {
		res := execBottle(path.Join(contractsPath, v.contract))
		errStrs := strings.Split(res, "\n")
		errs := []string{}
		for _, err := range errStrs {
			spl := strings.Split(err, ": ")
			if spl[len(spl)-1] != "" {
				errs = append(errs, spl[len(spl)-1])
			}
		}
		errStr := strings.Join(errs, ",")
		if errStr != v.wanted {
			t.Fatalf("test [%s] result mismatch, got %s, want %s", v.contract, errStr, v.wanted)
		}
	}

}

func execBottle(contract string) string {
	cmdPath := bottlePath
	cmdArgs := []string{"hint", "-code", contract}
	cmd := exec.Command(cmdPath, cmdArgs...)
	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return stderr.String()
	}
	return ""
}
