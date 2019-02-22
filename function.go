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
