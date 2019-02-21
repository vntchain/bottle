package main

type ExportType int

const (
	ExportTypeMutable ExportType = iota
	ExportTypeUnmutable
	ExportTypeNone
)

type FunctionInfo struct {
	Name      string
	Signature string
	Export    ExportType
	Payable   bool
	Location  Location
}

type Function struct {
	Hash uint32
	Name string
	Call []CallFunction
	Info FunctionInfo
}

type FunctionTree struct {
	Root map[uint32]*Function
}

type CallFunction struct {
	Name     string
	RootHash uint32
	Location Location
}

var functionTree *FunctionTree

func NewFunction(hash uint32, name string, info FunctionInfo) *Function {
	return &Function{
		Hash: hash,
		Name: name,
		Call: []CallFunction{},
		Info: info,
	}
}

func NewFunctionTree() *FunctionTree {
	return &FunctionTree{
		Root: map[uint32]*Function{},
	}
}

func (t *FunctionTree) AddFunction(f *Function) {
	t.Root[f.Hash] = f
}

func (t *FunctionTree) AddCall(hash uint32, name string, path string, offset int, size int) {
	t.Root[hash].AddCall(name, hash, path, offset, size)
}

func (f *Function) AddCall(name string, root uint32, path string, offset int, size int) {
	f.Call = append(f.Call, CallFunction{
		Name:     name,
		RootHash: root,
		Location: NewLocation(path, offset, size),
	})
}
