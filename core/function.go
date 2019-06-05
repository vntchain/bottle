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

import "sync"

type ExportType int

const (
	ExportTypeMutable ExportType = iota
	ExportTypeUnmutable
	ExportTypeNone
)

type DeclRefType int

const (
	DeclRefTypeFunction ExportType = iota
	DeclRefTypeVar
)

type FunctionInfo struct {
	Name      string
	Signature string
	Export    ExportType
	Payable   bool
	Location  Location
	WriteKey  bool
	ReturnKey bool
}

type FunctionVar struct {
	Name    string
	Type    DeclRefType
	Returns []string
	Chilren map[string]*FunctionVar
}

type Function struct {
	Hash uint32
	Name string
	Call []CallFunction
	Vars map[string]*FunctionVar
	Info FunctionInfo
}

type FunctionTree struct {
	Mutex *sync.RWMutex
	Root  map[uint32]*Function
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
		Vars: map[string]*FunctionVar{},
		Info: info,
	}
}

func NewFunctionTree() *FunctionTree {
	return &FunctionTree{
		Mutex: new(sync.RWMutex),
		Root:  map[uint32]*Function{},
	}
}

func (t *FunctionTree) AddFunction(f *Function) {
	t.Root[f.Hash] = f
}

func (t *FunctionTree) AddCall(hash uint32, name string, path string, line int, offset int, size int) {
	t.Root[hash].AddCall(name, hash, path, line, offset, size)
}

func (f *Function) AddCall(name string, root uint32, path string, line int, offset int, size int) {
	f.Call = append(f.Call, CallFunction{
		Name:     name,
		RootHash: root,
		Location: NewLocation(path, line, offset, size),
	})
}

// func (f *Function) AddVar(name string, ref string, writed bool) {
// 	isKey := false
// 	for _, v := range f.Vars {
// 		if v.Name == ref {
// 			isKey = true
// 		}
// 	}
// 	f.Vars = append(f.Vars, FunctionVar{
// 		Name:  name,
// 		IsKey: isKey,
// 	})
// }

// func (t *FunctionTree) ChangeWriteKeyState(hash uint32, state bool) {
// 	t.Root[hash].Info.WriteKey = state
// }
