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
	"fmt"
	"strings"

	"github.com/vntchain/go-vnt/accounts/abi"
)

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

var structLists abi.Root
var varLists abi.Root

type value struct {
	Path     string
	TypeInfo string
}

type symbol struct {
	Key          string
	KeyType      string
	IsArrayIndex bool
}

type ValueSymbol struct {
	ValueType   string
	ValueSymbol []symbol
}

var valueLists map[string]value
var lengthLists map[string]value
var keyLists map[string]value
var root *KVTree



func initList(node map[string]*abi.Node) {
	valueLists = make(map[string]value)
	lengthLists = make(map[string]value)
	keyLists = make(map[string]value)
	root = NewKVTree()
	for k, _ := range node {
		s := strings.Split(k, ".")
		keyLists[s[0]] = value{Path: s[0], TypeInfo: "pointer"}
	}
}

func RecursiveVarLists(node map[string]*abi.Node, path string, ty string) {
	for _, v := range node {
		var p string
		var t string
		if path == "" {
			p = v.FieldName
		} else {
			p = fmt.Sprintf("%s.%s", path, v.FieldName)
		}

		if ty == "" {
			t = v.FieldType
		} else {
			t = fmt.Sprintf("%s.%s", ty, v.FieldType)
		}

		if len(v.Children) != 0 {
			// fmt.Printf("RecursiveVarLists node 【%v】\n", v)
			root.AddNode(v.FieldName, v.StorageType, v.FieldType, p)
			RecursiveVarLists(v.Children, p, t)
		} else {

			root.AddNode(v.FieldName, v.StorageType, v.FieldType, p)
			if v.StorageType != abi.MappingKeyTy && v.StorageType != abi.ArrayIndexTy {
				valueLists[p] = value{Path: p, TypeInfo: t}
				if v.StorageType == abi.LengthTy {
					lengthLists[p] = value{Path: p, TypeInfo: t}
				}
			} else {
				keyLists[p] = value{Path: p, TypeInfo: t}
			}
		}
	}

}

func parseKey() map[string]ValueSymbol {
	valueSymbol := make(map[string]ValueSymbol)
	for _, v := range valueLists {
		var ts []string
		var ps []string
		if _, ok := lengthLists[v.Path]; ok {
			ts = strings.Split(v.TypeInfo, ".")
			ts = append(ts[0:len(ts)-2], ts[len(ts)-1])
			ps = strings.Split(v.Path, ".")
		} else {
			ts = strings.Split(v.TypeInfo, ".")
			ps = strings.Split(v.Path, ".")
		}

		//【a.value.value】
		//   【mapping.mapping.int32】
		sym := []symbol{}
		sym = append(sym, symbol{Key: ps[0], KeyType: "pointer"})

		for i, v := range ts {
			// fmt.Printf("i %d ts %s\n", i, v)
			res := ""
			switch v {
			case "mapping":
				tmp := make([]string, len(ps))
				copy(tmp, ps)
				tmp[i+1] = "key"
				res = strings.Join(tmp[0:i+2], ".")
			case "array":
				tmp := make([]string, len(ps))
				copy(tmp, ps)
				tmp[i+1] = "index"
				res = strings.Join(tmp[0:i+2], ".")
			case "struct":
				tmp := make([]string, len(ps))
				copy(tmp, ps)
				// fmt.Printf("i %d struct %v\n", i, tmp)
				res = strings.Join(tmp[0:i+2], ".")
				sym = append(sym, symbol{Key: res, KeyType: "pointer"})
			}
			if v, ok := keyLists[res]; ok {
				split := strings.Split(v.TypeInfo, ".")
				if split[len(split)-2] == "array" {
					sym = append(sym, symbol{Key: v.Path, KeyType: split[len(split)-1], IsArrayIndex: true})
				} else {
					sym = append(sym, symbol{Key: v.Path, KeyType: split[len(split)-1], IsArrayIndex: false})
				}

			}
		}

		valueSymbol[v.Path] = ValueSymbol{
			ValueType:   ts[len(ts)-1],
			ValueSymbol: sym,
		}
	}
	// for _, v := range lengthLists {
	// 	ts := strings.Split(v.TypeInfo, ".")
	// 	ps := strings.Split(v.Path, ".")
	// 	//[d length]
	// 	//[array uint64]
	// 	fmt.Printf("***lengthLists ts %v***\n", ts)
	// 	fmt.Printf("***lengthLists ps %v***\n", ps)
	// }
	return valueSymbol
}
