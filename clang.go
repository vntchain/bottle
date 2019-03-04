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

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/go-clang/bootstrap/clang"
	"github.com/vntchain/go-vnt/accounts/abi"
)

var KeyPos [][]int

var index = 0

func cmd(args []string) int {
	// idx := clang.NewIndex(0, 0)
	idx := clang.NewIndex(0, 1) //显示diagnostics
	defer idx.Dispose()
	var tu clang.TranslationUnit
	// tu = idx.ParseTranslationUnit("", []string{"-I", includeDir, "-x", "c", "-"}, nil, 0)
	tu = idx.ParseTranslationUnit(args[0], []string{"-I", includeDir}, nil, 0)
	if args[0] == "<stdin>" { //stdin
		fmt.Printf("<stdin> \n")
		tu = idx.ParseTranslationUnit("", []string{"-I", includeDir, "-x", "c", "-"}, nil, 0)
	} else {
		tu = idx.ParseTranslationUnit(args[0], []string{"-I", includeDir}, nil, 0)
	}

	defer tu.Dispose()

	diagnostics := tu.Diagnostics()
	for _, d := range diagnostics {
		// fmt.Printf("d %+v\n", d)
		// fmt.Println(d.Spelling(), d.CategoryText())
		// fmt.Println("PROBLEM:", d.Spelling(), " LEVEL:", d.Severity())
		if d.Severity() == clang.Diagnostic_Error || d.Severity() == clang.Diagnostic_Fatal {
			// return err
		}
	}

	cursor := tu.TranslationUnitCursor()

	// cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
	// 	fmt.Printf("cursor type %s\n", cursor.Kind().Spelling())
	// 	fmt.Printf("cursor %p parent %p\n", cursor, parent)
	// 	return clang.ChildVisit_Recurse
	// })

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.IsNull() {
			return clang.ChildVisit_Continue
		}
		// fmt.Printf("\n******          %s: %s (%s) (%s)\n", cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR(), cursor.Type().Spelling())
		// fmt.Printf("******parent    %s: %s (%s) (%s)\n", parent.Kind().Spelling(), parent.Spelling(), parent.USR(), parent.Type().Spelling())
		createFileContent(cursor, parent)
		createStructList(cursor, parent)
		switch cursor.Kind() {
		case clang.Cursor_ClassDecl, clang.Cursor_EnumDecl, clang.Cursor_StructDecl, clang.Cursor_Namespace:
			return clang.ChildVisit_Recurse
		}
		return clang.ChildVisit_Continue
	})
	structLists.Fulling()

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.IsNull() {
			return clang.ChildVisit_Continue
		}
		getGlobalVarDecl(cursor, parent)
		getFunc(cursor, parent)
		// getVarInFunction(cursor, parent)
		// switch cursor.Kind() {
		// case clang.Cursor_ClassDecl, clang.Cursor_EnumDecl, clang.Cursor_StructDecl, clang.Cursor_Namespace:
		// 	return clang.ChildVisit_Recurse
		// }
		// return clang.ChildVisit_Continue
		return clang.ChildVisit_Recurse
	})
	structLists.Fulling()
	// jsonres, _ := json.Marshal(structLists)
	// fmt.Printf("structLists %s\n", jsonres)

	if len(diagnostics) > 0 {
		// fmt.Println("NOTE: There were problems while analyzing the given file")
	}

	return 0
}

func createStructList(cursor, parent clang.Cursor) {

	decl := cursor.Kind()
	cursorname := cursor.Spelling()
	cursortype := cursor.Type().Spelling()
	usr := cursor.USR()

	pdecl := parent.Kind()
	pcursorname := parent.Spelling()
	// //pusr := parent.USR()
	pcursortype := parent.Type().Spelling()
	// fmt.Printf("\n******          %s: %s (%s) (%s)\n", cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR(), cursor.Type().Spelling())
	// fmt.Printf("******parent    %s: %s (%s) (%s)\n", parent.Kind().Spelling(), parent.Spelling(), parent.USR(), parent.Type().Spelling())
	if decl == clang.Cursor_StructDecl && pdecl == clang.Cursor_TranslationUnit { //声明结构体
		if cursorname == "" {
			if strings.Contains(cursortype, "struct (anonymous") { //匿名结构
				// fmt.Printf("匿名结构体\n")
				index = index + 1
				fieldtype := fmt.Sprintf("%s@@@%d", usr, index)
				node := abi.NewNode(cursorname, fieldtype, "")
				structStack = append(structStack, node)
				structLists.Root[fieldtype] = node
				// fmt.Printf("cursorname %s fieldtype %s cursor", cursorname, fieldtype)
			} else {
				// fmt.Printf("typedef 匿名结构体\n")
				// fmt.Printf("cursortype %s\n", cursortype)
				node := abi.NewNode(cursorname, cursortype, "")
				structLists.Root[cursortype] = node
				// fmt.Printf("typedef 匿名结构体 end\n")
			}
		} else {
			node := abi.NewNode(cursorname, cursortype, "")
			structLists.Root[cursortype] = node
			if strings.Contains(cursortype, "struct") {
				structLists.Root[cursortype[7:]] = node
			} else {
				structLists.Root[fmt.Sprintf("struct %s", cursortype)] = node
			}

			//structStack = append(structStack, node)
		}
	} else if decl == clang.Cursor_TypedefDecl && pdecl == clang.Cursor_TranslationUnit { //使用typedef定义的结构体解析
		// fmt.Printf("cursor.TypedefDeclUnderlyingType().Spelling() %s\n", cursor.TypedefDeclUnderlyingType().Spelling())
		// fmt.Printf("cursor.Type() %s\n", cursor.Type().Spelling())
		if strings.Contains(cursor.TypedefDeclUnderlyingType().Spelling(), "struct") {
			fieldname := cursor.TypedefDeclUnderlyingType().Spelling()[7:]
			// fmt.Printf("cursorname %s fieldname %s\n", cursorname, fieldname)
			if fieldname == cursorname { //匿名结构体
			} else {
				node := abi.NewNode(cursorname, fieldname, "")
				structLists.Root[cursorname] = node
			}
		}
	} else if decl == clang.Cursor_StructDecl && pdecl == clang.Cursor_StructDecl { //结构体内部定义的结构体
		if cursorname == "" {
			index = index + 1
			fieldtype := fmt.Sprintf("%s@@@%d", usr, index)
			node := abi.NewNode(cursorname, fieldtype, "")
			structStack = append(structStack, node)
			structLists.Root[fieldtype] = node
		} else {
			node := abi.NewNode(cursorname, cursortype, "")
			structLists.Root[cursorname] = node
			structLists.Root[fmt.Sprintf("struct %s", cursorname)] = node
		}

	} else if decl == clang.Cursor_FieldDecl && pdecl == clang.Cursor_StructDecl {
		// fmt.Printf(" decl == clang.Cursor_FieldDecl && pdecl == clang.Cursor_StructDecl \n")
		if len(structStack) != 0 {
			node := structStack[len(structStack)-1]
			if cutUSR(usr) == strings.Split(node.FieldType, "@@@")[0] {
				// fmt.Printf("cutUSR(usr) == strings.Split(node.FieldType, `@@@`)[0]\n")
				//fmt.Printf("struct element %s\n", strings.Split(node.FieldType, "@@@")[1])
				node.Add(cursorname, cursortype, "", node.FieldType)
			} else {
				// fmt.Printf("node.FieldType %s\n", node.FieldType)
				// fmt.Printf("struct element %v\n", strings.Split(node.FieldType, "@@@"))
				if pcursorname == "" {
					childnode := structStack[len(structStack)-1]
					structStack = structStack[0 : len(structStack)-1]
					node := structStack[len(structStack)-1]
					node.Add(cursorname, childnode.FieldType, "", node.FieldType)
				} else if strings.Contains(cursortype, "anonymous struct") {
					childnode := structStack[len(structStack)-1]
					node = structLists.Root[pcursorname]
					node.Add(cursorname, childnode.FieldType, "", pcursorname)
				} else {
					node = structLists.Root[pcursorname]
					node.Add(cursorname, cursortype, "", pcursorname)
				}
			}
		} else {
			if pcursorname != "" {
				node := structLists.Root[pcursorname]
				node.Add(cursorname, cursortype, "", pcursorname)
			} else {
				node := structLists.Root[pcursortype]
				// fmt.Printf("node %+v\n", node)
				node.Add(cursorname, cursortype, "", pcursortype)
				// fmt.Printf("cursorname %s cursortype %s  pcursortype %s\n", cursorname, cursortype, pcursortype)
			}

		}

	} else if decl == clang.Cursor_VarDecl && strings.Contains(cursortype, "anonymous struct") {
		// fmt.Printf(` decl == clang.Cursor_VarDecl && strings.Contains(cursortype, "anonymous struct")\n`)
		node := structStack[len(structStack)-1]
		structStack = structStack[0 : len(structStack)-1]
		structLists.Root[cursortype] = node
	}

}

//c:main6.cpp@S@main6.cpp@8255
func getGlobalVarDecl(cursor, parent clang.Cursor) {
	decl := cursor.Kind()
	cursortype := cursor.Type().Spelling()
	cursorname := cursor.Spelling()
	allstruct := []string{}
	for k, _ := range structLists.Root {
		allstruct = append(allstruct, k)
	}
	structnames := strings.Join(allstruct, "|")
	if decl == clang.Cursor_VarDecl {
		file, _, _, offset := cursor.Location().FileLocation()
		// fmt.Printf("\n******          %s: %s (%s) (%s)\n", cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR(), cursor.Type().Spelling())
		// fmt.Printf("******parent    %s: %s (%s) (%s)\n", parent.Kind().Spelling(), parent.Spelling(), parent.USR(), parent.Type().Spelling())
		// if strings.Contains(cursortype, "struct") || strings.Contains(cursortype, "volatile _S") {
		// fmt.Printf("!!!cursortype %s\n", cursortype)
		if !strings.Contains(cursortype, "volatile") {
			return
		}
		if strings.Contains(strings.Join(allstruct, ""), cursortype[9:]) {
			// fmt.Printf("======cursortype=======%s\n", cursortype)
			// fmt.Printf("\n******          %s: %s (%s) (%s)\n", cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR(), cursor.Type().Spelling())
			// fmt.Printf("******parent    %s: %s (%s) (%s)\n", parent.Kind().Spelling(), parent.Spelling(), parent.USR(), parent.Type().Spelling())
			if strings.Contains(cursortype, "struct (anonymous") {
				contents := strings.Split(cursortype, ":")
				num, err := strconv.Atoi(contents[len(contents)-2])
				if err != nil {
					panic(err)
				}
				if !isKey(string(fileContent[file.Name()][num-1].Content), structnames) {
					return
				}
			} else {
				_, x1, _, _ := cursor.Location().FileLocation()
				if x1-1 < 0 {
					return
				}
				if !isKey(string(fileContent[file.Name()][x1-1].Content), structnames) {
					return
				}
			}
			for k, v := range structLists.Root {
				//volatile struct (anonymous xxxx
				if k == cursortype || cursortype == "volatile "+k {
					varLists.Root[cursorname] = v
					v.FieldName = cursorname
					v.FieldType = ""
					v.FieldLocation = fmt.Sprintf("%d", offset)

					//todo 优化 key.go:84
					var allkey = ""
					for k, _ := range v.Children {
						allkey = allkey + k
					}
					if strings.Contains(allkey, "mapping1537182776") {
						v.FieldType = "mapping"
					} else if strings.Contains(allkey, "array1537182776") {
						v.FieldType = "array"
					} else {
						v.FieldType = "struct"
					}
					for k, c := range v.Children {
						if v.FieldType == "mapping" {
							if k == "key" {
								c.StorageType = abi.MappingKeyTy
							} else if k == "value" {
								c.StorageType = abi.MappingValueTy
							} else {
								delete(v.Children, k)
							}
						} else if v.FieldType == "array" {
							if k == "index" {
								c.StorageType = abi.ArrayIndexTy
							} else if k == "value" {
								c.StorageType = abi.ArrayValueTy
							} else if k == "length" {
								c.StorageType = abi.LengthTy
							} else {
								delete(v.Children, k)
							}
						} else if v.FieldType == "struct" {
							c.StorageType = abi.StructValueTy
						}
					}
				}
			}
		} else {
			sourceFile, x1, _, offset := cursor.Location().FileLocation()
			ext := path.Ext(sourceFile.Name())
			if strings.Compare(ext, ".h") == 0 {
				return
			}
			if x1-1 < 0 {
				return
			}
			if !isKey(string(fileContent[sourceFile.Name()][x1-1].Content), structnames) {
				return
			}
			offsetStr := fmt.Sprintf("%d", offset)
			if strings.Contains(cursortype, "volatile") {
				//volatile int64
				node := abi.NewNode(cursorname, cursortype[9:], offsetStr)
				node.StorageType = abi.NormalTy
				varLists.Root[cursorname] = node
			} else {
				node := abi.NewNode(cursorname, cursortype, offsetStr)
				node.StorageType = abi.NormalTy
				varLists.Root[cursorname] = node
			}

		}
	}
}

var currentFunctionHash uint32

func getFunc(cursor, parent clang.Cursor) {
	// fmt.Printf("func          %s: %s (%s) (%s)\n", cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR(), cursor.Type().Spelling())
	// fmt.Printf("func parent    %s: %s (%s) (%s)\n", parent.Kind().Spelling(), parent.Spelling(), parent.USR(), parent.Type().Spelling())
	if cursor.Kind() == clang.Cursor_FunctionDecl && parent.Kind() == clang.Cursor_TranslationUnit {
		// fmt.Printf("=============Cursor_FunctionDecl==============\n")
		currentFunctionHash = cursor.HashCursor()
		// fmt.Printf("function\n")
		// fmt.Printf("func          %s: %s (%s) (%s)\n", cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR(), cursor.Type().Spelling())
		// fmt.Printf("func parent    %s: %s (%s) (%s)\n", parent.Kind().Spelling(), parent.Spelling(), parent.USR(), parent.Type().Spelling())
		// fmt.Printf("hash %d\n ", currentFunctionHash)
		info := getFunctionInfo(cursor, parent)
		function := NewFunction(cursor.HashCursor(), cursor.Spelling(), info)
		if functionTree == nil {
			functionTree = NewFunctionTree()
		}
		functionTree.AddFunction(function)
		// fmt.Printf("function %+v\n parent %d\n", function, parent.HashCursor())
	}
	if cursor.Kind() == clang.Cursor_CallExpr {
		file, x1, _, _ := cursor.Location().FileLocation()
		// fmt.Printf("call\n")
		// fmt.Printf("hash %d function %s\n ", currentFunctionHash, cursor.Spelling())
		offset := fileContent[file.Name()][x1-1].Offset
		size := len(fileContent[file.Name()][x1-1].Content)
		functionTree.AddCall(currentFunctionHash, cursor.Spelling(), file.Name(), int(x1-1), offset, size)
		// fmt.Printf("func           %s: %s (%s) (%s)\n", cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR(), cursor.Type().Spelling())
		// fmt.Printf("func parent    %s: %s (%s) (%s)\n", parent.Kind().Spelling(), parent.Spelling(), parent.USR(), parent.Type().Spelling())
	}

}

func createFileContent(cursor, parent clang.Cursor) {
	file, _, _, _ := cursor.Location().FileLocation()
	if _, ok := fileContent[file.Name()]; !ok {
		if file.Name() == "<stdin>" {
			fileContent[file.Name()] = readfile(os.Stdin)
		} else {
			fi, err := os.Open(file.Name())
			if err != nil {
				panic(err.Error())
			}
			fileContent[file.Name()] = readfile(fi)
			fi.Close()
		}
	}
}

func getFunctionInfo(cursor, parent clang.Cursor) FunctionInfo {
	file, x1, _, _ := cursor.Location().FileLocation()
	// fmt.Println("func =======================")
	// fmt.Printf("cursor %s %d %d %d \n", file.Name(), x1, x2, x3)
	// fmt.Printf("func          %s: %s (%s) (%s)\n", cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR(), cursor.Type().Spelling())
	// fmt.Printf("func parent    %s: %s (%s) (%s)\n", parent.Kind().Spelling(), parent.Spelling(), parent.USR(), parent.Type().Spelling())
	var export ExportType
	var payable bool
	if x1-2 >= 0 {
		cont := removeSymbol(string(fileContent[file.Name()][x1-2].Content))
		if len(cont) == 0 {
			export = ExportTypeNone
		} else {
			cont = strings.Split(cont[0], "\n")
			if cont[0] == KWMutable {
				export = ExportTypeMutable
			} else if cont[0] == KWUnmutable {
				export = ExportTypeUnmutable
			} else {
				export = ExportTypeNone
			}
		}
	}
	if cursor.Spelling()[0:1] == "$" {
		payable = true
	} else {
		payable = false
	}
	// fmt.Printf("content %s x1 %d\n", fileContent[file.Name()][x1-1].Content, x1)
	return FunctionInfo{
		Name:      cursor.Spelling(),
		Signature: cursor.Type().Spelling(),
		Export:    export,
		Payable:   payable,
		Location:  NewLocation(file.Name(), int(x1-1), fileContent[file.Name()][x1-1].Offset, len(fileContent[file.Name()][x1-1].Content)),
	}
}

var CompoundStmtHash uint32
var varDeclSpell string
var memberExpr string

func getVarInFunction(cursor, parent clang.Cursor) {
	fmt.Printf("func          %s: %s (%s) (%s)\n", cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR(), cursor.Type().Spelling())
	fmt.Printf("func parent    %s: %s (%s) (%s)\n", parent.Kind().Spelling(), parent.Spelling(), parent.USR(), parent.Type().Spelling())

	if parent.HashCursor() == CompoundStmtHash { //开始一条新的语句
		if varDeclSpell != "" {
			//

			//初始化
			varDeclSpell = ""
			memberExpr = ""

		}
	}

	switch cursor.Kind() {
	case clang.Cursor_CompoundStmt:
		CompoundStmtHash = cursor.HashCursor()
	case clang.Cursor_VarDecl:
		varDeclSpell = cursor.Spelling()
		fmt.Printf("VarDecl Spell %s\n", varDeclSpell)
	case clang.Cursor_ReturnStmt: //返回值是否是key
	}
	// if cursor.Kind() == clang.Cursor_CompoundStmt {
	// 	CompoundStmtHash = cursor.HashCursor()
	// }
	// if cursor.Kind() == clang.Cursor_VarDecl {
	// 	varDeclSpell = cursor.Spelling()
	// 	fmt.Printf("VarDecl Spell %s\n", varDeclSpell)
	// }

}
