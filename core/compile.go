// Copyright 2019 The bottle Authors
// This file is part of the bottle library.
//
// The bottle library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The bottle library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the bottle library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"bytes"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/wavm"
	"github.com/vntchain/vnt-wasm/validate"
	"github.com/vntchain/vnt-wasm/wasm"
)

//clang -Xclang -ast-dump -fsyntax-only main3.cpp

//VNT_WASM_EXPORT
//uint64 init   (   uint64 totalsupply   )
//todo \s不仅匹配了空格符，还匹配了其他符号
//	匹配任何空白字符，包括空格、制表符、换页符等等。等价于 [ \f\n\r\t\v]。注意 Unicode 正则表达式会匹配全角空格符。
const (
	//methodReg     = `(VNT_WASM_EXPORT\n)(\s*)(int(|32|64)|uint(|32|64|256)|address|string|bool|void)(\s+)([a-zA-Z0-9_]+)(\s*)(\({1})([a-zA-Z0-9_\s,]*)(\){1})`
	methodReg     = `(MUTABLE|UNMUTABLE\n)(\s*)(int(|32|64)|uint(|32|64|256)|address|string|bool|void)(\s+)([a-zA-Z0-9_\$\*]+)(\s*)(\({1})([a-zA-Z0-9_\*\s,]*)(\){1})([^{]*)({){1}`
	openParenReg  = `(\s*)(\()(\s*)`
	closeParenReg = `(\s*)(\))(\s*)`
	commaReg      = `(\s*),(\s*)`
	spaceReg      = `(\s+)`
	letterReg     = `[a-zA-Z0-9_\$\*]{1,}`

	functionReg = `(mutable|unmutable|)(\s*)(int(|32|64)|uint(|32|64|256)|address|string|bool|void)(\s+)([a-zA-Z0-9_\$\*]+)(\s*)(\({1})([a-zA-Z0-9_\*\s,]*)(\){1})(\s*)({){1}`
	// constructorAndMethodReg = `(VNT_WASM_EXPORT\n)(\s*)(int(|32|64)|uint(|32|64|256)|address|string|bool|void|constructor)(\s+)([a-zA-Z0-9_]+)(\s*)(\({1})([a-zA-Z0-9_\s,]*)(\){1})([^{]*)({){1}`
)

//KEY uint64 aaa_a;
//KEY mapping(address, address) mapping_g;
//KEY array(int32) array_a;
const (
// ([\s\S]*)(\()[^}]*(\){1,})
// keyReg = `(KEY)([ ]+)(int(|32|64)|uint(|32|64|256)|address|string|bool|mapping|array|struct)[^;]*(;{1})`
// keyReg = `(KEY)[^(;|\r|\n)]*(;{1})`
// keyReg = `(KEY)([ ]+)( int(|32|64)|uint(|32|64|256)|address|string|bool|mapping|array)([ ]+)([a-zA-Z0-9_]+)`
)

//event transfer_event(address _from,address _to,uint64 _amount);
const (
	eventReg = `(EVENT)(\s+)([a-zA-Z0-9_\$]+)(\s*)(\({1})([a-zA-Z0-9_\*\s,]*)(\){1})`
)

//call uint64 GetAmount(address _contractaddress,uint64 _amount, address _addr);
const (
	callReg = `(CALL)(\s+)(int(|32|64)|uint(|32|64|256)|address|string|bool|void)(\s+)([a-zA-Z0-9_\$\*]+)(\s*)(\({1})(\s*)(CallParams)(\s+)([a-zA-Z0-9_\*\s,]*)(\){1})`
)

//construct token   (   uint64 totalsupply   ) {
const (
	constructorReg = `(constructor)(\s+)([a-zA-Z0-9_\$\*]+)(\s*)(\({1})([a-zA-Z0-9_\*\s,]*)(\){1})([^{]*)({){1}`
)

const (
	fallBackReg = `(\$*)(_\(\))(\s*)({){1}`
)

//处理代码中的注释
//event transfer_event(address _from,/*address _to,*/uint64 _amount);
//todo 处理 //
const (
	commandReg = `/\*([^*]|[\r\n]|(\*+([^*/]|[\r\n])))*\*+/|//(.*)`
)

const (
	mutable   = "MUTABLE"
	unmutable = "UNMUTABLE"
)

// {"name":"init","constant":false,"inputs":[{"name":"totalsupply","type":"uint64"}],"outputs":[{"name":"bbb","type":"uint64"}],"payable":false,"stateMutability":"nonpayable","type":"function"},

type ABI struct {
	Constructor Method            `json:"constructor"`
	Methods     map[string]Method `json:"methods"`
	Events      map[string]Event  `json:"events"`
	Calls       map[string]Method `json:"calls"`
	// Keys        map[string]Key    `json:"keys"`
}

type Method struct {
	Name    string    `json:"name"`
	Const   bool      `json:"constant"`
	Inputs  Arguments `json:"inputs"`
	Outputs Arguments `json:"outputs"`
	Type    string    `json:"type"`
}

type Event struct {
	Name      string    `json:"name"`
	Anonymous bool      `json:"anonymous"`
	Inputs    Arguments `json:"inputs"`
	Type      string    `json:"type"`
}

type Key struct {
	Name   string      `json:"name"`
	Tables []*abi.Node `json:"tables"`
	Type   string      `json:"type"`
}

type Argument struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Indexed bool   `json:"indexed"`
}

type Arguments []Argument

type Contract struct {
	ContractName string              `json:"contractName"`
	Abi          interface{}         `json:"abi"`
	Bytecode     string              `json:"bytecode"`
	SourcePath   string              `json:"sourcePath"`
	UpdatedAt    string              `json:"updatedAt"`
	Networks     map[string]Networks `json:"networks,omitempty"`
}
type Networks struct {
	Events          map[string]Events `json:"events"`
	Address         string            `json:"address"`
	TransactionHash string            `json:"transactionHash"`
}

type Events struct {
}

type abiGen struct {
	Code []byte
	abi  ABI
}

var fileContent = map[string][]ContentPerLine{}

var (
	codePath          string
	outputDir         string
	includeDir        string
	abiPath           string
	wasmPath          string
	compressPath      string
	DecompressDirPath string
)

func newAbiGen(code []byte) *abiGen {

	return &abiGen{
		Code: code,
		abi: ABI{
			Methods: make(map[string]Method),
			Events:  make(map[string]Event),
			Calls:   make(map[string]Method),
			// Keys:    make(map[string]Key),
		},
	}
}

func (gen *abiGen) removeComment() {
	gen.Code = removeComment(string(gen.Code))
}

//将显示声明的struct和typedef声明的struct替换成隐示声明
// struct node{
// int32 a,b;
// int32 c,d;
// };
// KEY struct node s4;
// typedef struct node{
// int32 a,b;
// int32 c,d;
// } _node;
// 复杂情况
// KEY _node s4;
// struct node1{
// 	int32 a,b;
// 	struct {
// int c,d
// } e;
// } ;

// KEY struct node1 s6;
//=================>>
// KEY struct{int32 a,b;int32 c,d;} s4;
//

const (
	anonymousStruct = `(struct)([ ]+)([a-zA-Z0-9_\$]+)`
)

func removeSymbol(input string) []string {
	s := strings.TrimSpace(input)
	re := regexp.MustCompile(openParenReg)
	str := re.ReplaceAllString(s, "(")

	re = regexp.MustCompile(closeParenReg)
	str = re.ReplaceAllString(str, ")")

	re = regexp.MustCompile(spaceReg)
	str = re.ReplaceAllString(str, " ")
	re = regexp.MustCompile(letterReg)
	final := re.FindAllString(str, -1)
	return final
}

func splitArgs(input string) []string {
	s := strings.TrimSpace(input)
	re := regexp.MustCompile(openParenReg)
	str := re.ReplaceAllString(s, "(")

	re = regexp.MustCompile(closeParenReg)
	str = re.ReplaceAllString(str, ")")

	re = regexp.MustCompile(spaceReg)
	str = re.ReplaceAllString(str, " ")

	re = regexp.MustCompile(commaReg)
	str = re.ReplaceAllString(str, ",")

	re = regexp.MustCompile(`(\(|\)|,)`)
	final := re.Split(str, -1)
	final = final[0 : len(final)-1]
	return final
}

func removeSpaceAndParen(input string) ([]string, []string) {
	s := strings.TrimSpace(input)
	re := regexp.MustCompile(openParenReg)
	str := re.ReplaceAllString(s, "(")

	re = regexp.MustCompile(closeParenReg)
	str = re.ReplaceAllString(str, "")

	re = regexp.MustCompile(spaceReg)
	str = re.ReplaceAllString(str, "")
	sp := strings.Split(str, "(")
	leftFinal := []string{sp[0]}
	right := sp[1]
	var rightFinal []string
	if right == "" {
		rightFinal = []string{}
	} else {
		rightFinal = strings.Split(right, ",")
	}
	return leftFinal, rightFinal
}

func removeComment(code string) []byte {
	codeBytes := []byte(code)
	reg := regexp.MustCompile(commandReg)
	idx := reg.FindAllStringIndex(code, -1)
	for _, v := range idx {
		for i := v[0]; i < v[1]; i++ {
			if codeBytes[i] != byte('\n') {
				codeBytes[i] = byte(' ')
			}
		}
	}
	return codeBytes
}

func (gen *abiGen) parseMethod() {
	reg := regexp.MustCompile(methodReg)
	res := reg.FindAllString(string(gen.Code), -1)
	for _, v := range res {
		s1 := strings.Split(strings.Replace(v, "\r\n", "\n", -1), "\n")
		if len(s1) < 2 {
			panic("Irregular method structure")
		}
		s2 := ""
		for i := 1; i < len(s1); i++ {
			s2 = s2 + s1[i]
		}

		final := removeSymbol(s2)

		var method Method
		mutable := s1[0]
		method.Const = isConstant(mutable)

		name := final[1]
		funcType := final[0]
		method.Name = name
		if funcType == "void" {
			method.Outputs = Arguments{}

		} else {
			output := Argument{
				Name:    "output",
				Type:    funcType,
				Indexed: false,
			}
			method.Outputs = Arguments{output}
		}
		inputs := Arguments{}
		for i := 2; i < len(final); i = i + 2 {
			input := Argument{
				Name:    final[i+1],
				Type:    final[i],
				Indexed: false,
			}
			inputs = append(inputs, input)
		}
		method.Inputs = inputs
		method.Type = "function"
		gen.abi.Methods[name] = method
	}
}

// func (gen *abiGen) parseKey() {

// 	for _, v := range varLists.Root {
// 		key := Key{
// 			Name: v.FieldName,
// 			Type: "key",
// 			Tables: []*abi.Node{
// 				&abi.Node{
// 					FieldName: v.FieldName,
// 					FieldType: v.FieldType,
// 					Tables:    v.Tables,
// 				},
// 			},
// 		}
// 		gen.abi.Keys[v.FieldName] = key
// 	}

// 	fmt.Printf("abi key 111%+v", gen.abi.Keys)
// }

func (gen *abiGen) normalType(inputs []string) {
	if len(inputs) < 3 {
		panic("Illegal input key")
	}
	switch inputs[1] {
	case "int32", "int64", "uint32", "uint64", "uint256", "string", "bool", "address":
	case "mapping":
	case "array":
	case "struct":
	}
}

func (gen *abiGen) parseEvent() {
	reg := regexp.MustCompile(eventReg)
	res := reg.FindAllString(string(gen.Code), -1)
	for _, v := range res {
		final := removeSymbol(v)
		//makeAbi.Method[final[1]]=
		var event Event
		name := final[1]
		event.Name = name
		inputs := Arguments{}
		for i := 2; i < len(final); i = i + 2 {
			if final[i] == "indexed" {
				i = i + 1
				input := Argument{
					Name:    final[i+1],
					Type:    final[i],
					Indexed: true,
				}
				inputs = append(inputs, input)
			} else {
				input := Argument{
					Name:    final[i+1],
					Type:    final[i],
					Indexed: false,
				}
				inputs = append(inputs, input)
			}

		}
		event.Inputs = inputs
		event.Type = "event"
		gen.abi.Events[name] = event
	}
}

//todo 参数1和参数2类型判断
func (gen *abiGen) parseCall() {
	reg := regexp.MustCompile(callReg)
	res := reg.FindAllString(string(gen.Code), -1)
	for _, v := range res {
		final := removeSymbol(v)
		var call Method
		name := final[2]
		call.Name = name

		if final[1] == "void" {
			call.Outputs = Arguments{}

		} else {
			output := Argument{
				Name:    "output",
				Type:    final[1],
				Indexed: false,
			}
			call.Outputs = Arguments{output}
		}

		inputs := Arguments{}

		for i := 5; i < len(final); i = i + 2 { //忽略第一个参数
			input := Argument{
				Name:    final[i+1],
				Type:    final[i],
				Indexed: false,
			}
			inputs = append(inputs, input)

		}
		call.Inputs = inputs
		call.Type = "call"
		gen.abi.Calls[name] = call
	}
}

func (gen *abiGen) parseConstructor() {
	reg := regexp.MustCompile(constructorReg)
	res := reg.FindAllString(string(gen.Code), -1)
	if len(res) == 0 {
		panic("Can't find Contructor function")
	}
	for _, v := range res {
		final := removeSymbol(v)
		var method Method
		method.Const = false
		name := final[1]
		method.Name = name
		method.Outputs = Arguments{}

		inputs := Arguments{}
		for i := 2; i < len(final); i = i + 2 {
			input := Argument{
				Name:    final[i+1],
				Type:    final[i],
				Indexed: false,
			}
			inputs = append(inputs, input)
		}
		method.Inputs = inputs
		method.Type = "constructor"
		gen.abi.Constructor = method
	}
}

//Registery(value_address,value_type,key_address,key_type,is_array_index)
//Registery(&a.key,"mapping","int")
const regFmt = "AddKeyInfo( &%s, %d, &%s, %d, %t);\n"
const funcFmt = "\nvoid %s(){%s}\n"
const initializeVariables = "\nInitializeVariables();"

// 插入AddKeyInfo和init代码
// AddKeyInfo代码用于构建key信息
// InitializeVariables用于在constructor方法里存储key的初始化值
func (gen *abiGen) insertRegistryCode() []byte {
	initList(varLists.Root)
	RecursiveVarLists(varLists.Root, "", "")
	sym := parseKey()
	insert := "\n"
	for k, v1 := range sym {
		for _, v2 := range v1.ValueSymbol {
			// fmt.Printf("key2222 %s val2 %s StorageType %s \n", k, v1.ValueType, v2.Key, v2.KeyType)
			insert = insert + fmt.Sprintf(regFmt, k, abi.KeyType(v1.ValueType), v2.Key, abi.KeyType(v2.KeyType), v2.IsArrayIndex)
		}
	}
	funcName := fmt.Sprintf("%s%s", "key", GetRandomString(8))
	insertFuncBody := fmt.Sprintf(funcFmt, funcName, insert)

	reg := regexp.MustCompile(methodReg)
	res := reg.FindAllStringIndex(string(gen.Code), -1)
	reg = regexp.MustCompile(constructorReg)
	consres := reg.FindAllStringIndex(string(gen.Code), -1)
	reg = regexp.MustCompile(fallBackReg)
	fbres := reg.FindAllStringIndex(string(gen.Code), -1)
	res = append(append(res, consres...), fbres...)
	sort.Sort(Index(res))
	var code = gen.Code
	originCode := make([]byte, len(code))
	copy(originCode, code)
	insertFunc := fmt.Sprintf("\n%s();", funcName)
	insertFuncBodyBytes := []byte(insertFuncBody)
	initLen := 0
	for i, v := range res {
		f := originCode[v[0]:v[1]]
		flag := isConstructor(string(f))
		if i == 0 {
			codeInter := common.Insert(code, v[0], insertFuncBodyBytes)
			code = codeInter.([]byte)
		}

		insertBytes := []byte(insertFunc)
		//code = append(code[0:v[1]+i*len(insertBytes)+len(insertFuncBodyBytes)], append(insertBytes, code[v[1]+i*len(insertBytes)+len(insertFuncBodyBytes):]...)...)
		insertRes := common.Insert(code, v[1]+i*len(insertBytes)+len(insertFuncBodyBytes)+initLen, insertBytes)
		code = insertRes.([]byte)
		if flag {
			initBytes := []byte(initializeVariables)
			initLen = len(initBytes)
			insertRes := common.Insert(code, v[1]+(i+1)*len(insertBytes)+len(insertFuncBodyBytes), initBytes)
			code = insertRes.([]byte)
		}
	}
	return code
}

func BuildWasm(input string, output string) {
	buildCFile("-g -O3", input, output)
}

func ValidWasm(wasmcode []byte, abi abi.ABI) error {
	wasm.SetDebugMode(false)
	buf := bytes.NewReader(wasmcode)
	wavm := &wavm.Wavm{
		ChainContext: wavm.ChainContext{
			Code: wasmcode,
			Abi:  abi,
		},
	}
	m, err := wasm.ReadModule(buf, wavm.ResolveImports)
	if err != nil {
		newErr := fmt.Errorf("%s\n%s", err.Error(), "Please check that the wasm trigger an invalid ENV function (ex: call c standard library function).")
		return newErr
	}
	if err := validate.VerifyModule(m); err != nil {
		newErr := fmt.Errorf("%s\n%s", err.Error(), "Please check that the wasm trigger an invalid ENV function (ex: call c standard library function).")
		return newErr
	}
	return nil
}

type Index [][]int

func (idx Index) Len() int           { return len(idx) }
func (idx Index) Less(i, j int) bool { return idx[i][0]-idx[j][0] < 0 }
func (idx Index) Swap(i, j int)      { idx[i], idx[j] = idx[j], idx[i] }

func isConstructor(input string) bool {
	regStr := `constructor`
	reg := regexp.MustCompile(regStr)
	res := reg.FindAllString(input, -1)
	if len(res) == 0 {
		return false
	} else {
		return true
	}
}

func isConstant(input string) bool {
	inputs := strings.Split(input, " ")
	if inputs[0] == mutable {
		return false
	} else if inputs[0] == unmutable {
		return true
	} else {
		panic("Unsupport keyword " + inputs[0])
	}
}

func mustCFile(codefilepath string) {
	res := path.Ext(codefilepath)
	if strings.Compare(res, ".c") != 0 {
		panic("合约文件的后缀名必须为.c")
	}
}
