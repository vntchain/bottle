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
	"strings"

	"github.com/vntchain/go-vnt/accounts/abi"
)

type KVTree struct {
	Root map[string]*KVNode
}

type KVNode struct {
	Name        string
	StorageType abi.StorageType
	Type        string
	Children    map[string]*KVNode
}

func NewKVTree() *KVTree {
	return &KVTree{
		Root: make(map[string]*KVNode),
	}
}

func NewKVNode(name string, styp abi.StorageType, typ string) *KVNode {
	return &KVNode{
		Name:        name,
		StorageType: styp,
		Type:        typ,
		Children:    make(map[string]*KVNode),
	}
}

func (node *KVNode) AddNode(name string, styp abi.StorageType, typ string) {
	if _, ok := node.Children[name]; !ok {
		n := NewKVNode(name, styp, typ)
		node.Children[name] = n
	}
}

func (tree *KVTree) AddNode(name string, styp abi.StorageType, typ string, path string) {
	keys := strings.Split(path, ".")
	if len(keys) <= 1 {
		root := NewKVNode(name, styp, typ)
		tree.Root[name] = root
	} else {
		node := tree.Root[keys[0]]
		for i := 1; i < len(keys)-1; i++ {
			node = node.Children[keys[i]]
		}

		node.AddNode(name, styp, typ)
	}
}
