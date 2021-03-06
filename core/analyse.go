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
	"regexp"
	"strconv"
	"strings"

	"github.com/vntchain/go-vnt/accounts/abi"
)

//key 类型检查<complete>，如果key写在其他方法之间，报错并提示  未完成，当前只判断key写在constructor之前，
//没有constructor <complete>
//多个constructor <complete>
//call,第一个参数检查,参数类型检查及返回值检查,event,参数类型检查 <complete>
//mutable
//unmutable，sendfromcontract,transferfromcontract检查,mutable方法调用检查 <complete>
//unmutable，key的写入检查
//unmutable调用payable或mutable <complete>
//payable,input的检查
//unpayable 调用getvalue的提示
//导出方法的类型和参数检查 <complete>
//uint256及address类型

//方法中KEY写操作的判断：
//1)直接通过key变量进行写，
//2)临时变量通过指针指向key变量，并进行了写操作
//3)方法中调用了其他方法，其他方法中包含上面两种情况，所以必须对每个方法添加一个是否包含写操作的字段
//4)方法中调用了其他方法，其他方法返回了KEY的变量
//4.1)返回KEY变量的成员，若成员为array类型的index和mapping类型的key，则该方法为非写操作，
//4.2)返回KEY变量的成员，若成员为非array类型的index和mapping类型的key，需要从上层方法判断是否进行了写操作
//4.3)返回KEY变量，非成员，需要从上层方法判断是否对成员进行写操作

//1)临时变量指向一个方法的返回值
//2)临时变量指向另一个变量
//2.1)临时变量指向另一个临时变量，递归获取直到最终变量是否是KEY变量，
//2.1.1)若最终变量为KEY变量，染色为KEY变量，否则，为临时变量
//2.2)临时变量指向KEY变量，则把该临时变量染色为KEY变量
//2.3）临时变量指向一个方法的返回值

type Hint struct {
	Path           string
	Code           []byte
	ConstructorPos int
}

type HintType int

const (
	HintTypeWarning HintType = iota
	HintTypeError
)

type HintMessage struct {
	Name     string
	Message  string
	Type     HintType
	Location Location
}
type HintMessages []HintMessage

func (msgs HintMessages) ToString() string {
	m := ""
	for i, v := range msgs {
		m = m + fmt.Sprintf("%s:%d:%d: error: %s", v.Location.Path, v.Location.Line, v.Location.Offset, v.Message)
		if i < len(msgs)-1 {
			m = m + "\n"
		}
	}
	return m
}

func newHint(path string, code []byte) *Hint {

	return &Hint{
		Path: path,
		Code: code,
	}
}

func (h *Hint) contructorCheck() (HintMessages, error) {
	var msgs HintMessages
	reg := regexp.MustCompile(constructorReg)
	idx := reg.FindAllStringIndex(string(h.Code), -1)
	if len(idx) == 0 {
		offset := 0
		size := 0
		line, lineOffset := GetLineNumber(offset, fileContent[h.Path])
		msg := HintMessage{
			Message:  "必须定义一个Contructor方法",
			Type:     HintTypeError,
			Location: NewLocation(h.Path, line, lineOffset, size),
		}
		msgs = append(msgs, msg)
		return msgs, nil
	}

	h.ConstructorPos = idx[0][0]
	if len(idx) > 1 {
		for i := 0; i < len(idx); i++ {
			offset := idx[i][0]
			size := idx[i][1] - idx[i][0]
			line, lineOffset := GetLineNumber(offset, fileContent[h.Path])
			msg := HintMessage{
				Message:  "重复定义Contructor方法",
				Type:     HintTypeError,
				Location: NewLocation(h.Path, line, lineOffset, size),
			}
			msgs = append(msgs, msg)
		}
		return msgs, nil
	}

	constructorReg := `(constructor)[^(;|\r\n|\r|\n|\{|\})]*(\s+)(%s)(\s*)(\({1})([a-zA-Z0-9_\*\s,]*)(\){1})`
	for _, v := range functionTree.Root {
		call := fmt.Sprintf(constructorReg, escape(v.Name))
		reg := regexp.MustCompile(call)
		stridx := reg.FindAllStringIndex(string(h.Code), -1)
		if len(stridx) != 0 {
			//find
			left, right := removeSpaceAndParen(v.Info.Signature)
			offset := stridx[0][0]
			size := stridx[0][1] - stridx[0][0]
			line, lineOffset := GetLineNumber(offset, fileContent[h.Path])
			if len(left) != 1 {
				msg := HintMessage{
					Message:  "Constructor方法的返回值为不支持的类型：" + strings.Join(left, " "),
					Type:     HintTypeError,
					Location: NewLocation(h.Path, line, lineOffset, size),
				}
				msgs = append(msgs, msg)
			} else {
				if !isSupportedType(left[0]) {
					msg := HintMessage{
						Message:  "Constructor方法的返回值为不支持的类型：" + left[0],
						Type:     HintTypeError,
						Location: NewLocation(h.Path, line, lineOffset, size),
					}
					msgs = append(msgs, msg)
				}
			}
			var msgStr = "Constructor方法的参数为不支持的类型："
			types := []string{}
			for i := 0; i < len(right); i++ {
				if !isSupportedType(right[i]) {
					types = append(types, right[i])
				}
			}
			if len(types) != 0 {
				msg := HintMessage{
					Message:  msgStr + strings.Join(types, " ,"),
					Type:     HintTypeError,
					Location: NewLocation(h.Path, line, lineOffset, size),
				}
				msgs = append(msgs, msg)
			}
		}
	}

	return msgs, nil
}

func (h *Hint) keyCheck() (HintMessages, error) {
	var msgs HintMessages
	for _, v := range varLists.Root {
		loc := v.FieldLocation
		offset, err := strconv.Atoi(loc)
		if err != nil {
			return nil, err
		}
		size := len(v.FieldName)
		line, lineOffset := GetLineNumber(offset, fileContent[h.Path])
		if offset >= h.ConstructorPos {
			msg := HintMessage{
				Message:  "KEY必须定义在construct之前",
				Type:     HintTypeWarning,
				Location: NewLocation(h.Path, line, lineOffset, size),
			}
			msgs = append(msgs, msg)
		}
		// 类型判断
		// key类型分基本类型和复杂类型，
		// 基本类型包括int32/64,uint32/64,uint256,bool,string,address
		// 复杂类型包括mapping,array,由基本类型构成的struct
		types := Traversal(v)
		unsupported := []string{}
		for _, t := range types {
			if !isSupportedKeyType(t) {
				unsupported = append(unsupported, t)
			}
		}
		if len(unsupported) != 0 {
			msg := HintMessage{
				Message:  "KEY为不支持的类型：" + strings.Join(unsupported, " ,"),
				Type:     HintTypeError,
				Location: NewLocation(h.Path, line, lineOffset, size),
			}
			msgs = append(msgs, msg)
		}
	}

	return msgs, nil
}

func Traversal(node *abi.Node) []string {
	types := []string{node.FieldType}
	if len(node.Tables) != 0 {
		for _, v := range node.Tables {
			types = append(types, Traversal(v)...)
		}
	}
	return types
}

func (h *Hint) callCheck() (HintMessages, error) {
	var msgs HintMessages
	callReg := `(CALL)[^(;|\r\n|\r|\n|\{|\})]*(\s+)(%s)(\s*)(\({1})([a-zA-Z0-9_\*\s,]*)(\){1})`
	for _, v := range functionTree.Root {
		call := fmt.Sprintf(callReg, escape(v.Name))
		reg := regexp.MustCompile(call)
		stridx := reg.FindAllStringIndex(string(h.Code), -1)
		if len(stridx) != 0 {
			offset := stridx[0][0]
			size := stridx[0][1] - stridx[0][0]
			line, lineOffset := GetLineNumber(offset, fileContent[h.Path])
			//find call
			left, right := removeSpaceAndParen(v.Info.Signature)
			if len(left) != 1 {
				msg := HintMessage{
					Message:  "CALL方法的返回值为不支持的类型：" + strings.Join(left, " "),
					Type:     HintTypeError,
					Location: NewLocation(h.Path, line, lineOffset, size),
				}
				msgs = append(msgs, msg)
			} else {
				if !isSupportedType(left[0]) {
					msg := HintMessage{
						Message:  "CALL方法的返回值为不支持的类型：" + left[0],
						Type:     HintTypeError,
						Location: NewLocation(h.Path, line, lineOffset, size),
					}
					msgs = append(msgs, msg)
				}
			}
			if len(right) == 0 {
				msg := HintMessage{
					Message:  "CALL方法至少需要一个类型为CallParams的参数",
					Type:     HintTypeError,
					Location: NewLocation(h.Path, line, lineOffset, size),
				}
				msgs = append(msgs, msg)
			} else {
				types := []string{}
				for i := 0; i < len(right); i++ {
					if i == 0 { //callprams
						if right[i] != "CallParams" {
							msg := HintMessage{
								Message:  "CALL方法的第一个参数为不支持的类型：" + right[i] + ", 支持的类型为CallParams",
								Type:     HintTypeError,
								Location: NewLocation(h.Path, line, lineOffset, size),
							}
							msgs = append(msgs, msg)
						}
					} else { //input type
						if !isSupportedType(right[i]) {
							types = append(types, right[i])
						}
					}
				}
				if len(types) != 0 {
					msg := HintMessage{
						Message:  "CALL方法的参数为不支持的类型：" + strings.Join(types, " ,"),
						Type:     HintTypeError,
						Location: NewLocation(h.Path, line, lineOffset, size),
					}
					msgs = append(msgs, msg)
				}
			}
		}
	}

	return msgs, nil
}

//EVENT event_name(indexed[option] param_type param_name)
func (h *Hint) eventCheck() (HintMessages, error) {
	var msgs HintMessages
	eventReg := `(EVENT)[^(;|\r\n|\r|\n|\{|\})]*(\s+)(%s)(\s*)(\({1})([a-zA-Z0-9_\*\s,]*)(\){1})`
	for _, v := range functionTree.Root {
		event := fmt.Sprintf(eventReg, escape(v.Name))
		reg := regexp.MustCompile(event)
		stridx := reg.FindAllStringIndex(string(h.Code), -1)
		if len(stridx) != 0 {
			offset := stridx[0][0]
			size := stridx[0][1] - stridx[0][0]
			line, lineOffset := GetLineNumber(offset, fileContent[h.Path])
			//find event
			left, right := removeSpaceAndParen(v.Info.Signature)

			if len(left) != 1 {
				msg := HintMessage{
					Message:  "EVENT方法的返回值为不支持的类型：" + strings.Join(left, " "),
					Type:     HintTypeError,
					Location: NewLocation(h.Path, line, lineOffset, size),
				}
				msgs = append(msgs, msg)
			} else {
				if !isSupportedType(left[0]) {
					msg := HintMessage{
						Message:  "EVENT方法的返回值为不支持的类型：" + left[0],
						Type:     HintTypeError,
						Location: NewLocation(h.Path, line, lineOffset, size),
					}
					msgs = append(msgs, msg)
				}
			}
			if len(right) == 0 {
				msg := HintMessage{
					Message:  "EVENT方法至少需要一个参数",
					Type:     HintTypeError,
					Location: NewLocation(h.Path, line, lineOffset, size),
				}
				msgs = append(msgs, msg)
			} else {
				//类型判断
				types := []string{}
				for i := 0; i < len(right); i++ {
					if !isSupportedType(right[i]) {
						types = append(types, right[i])

					}
				}
				if len(types) != 0 {
					msg := HintMessage{
						Message:  "EVENT方法的参数为不支持的类型：" + strings.Join(types, " ,"),
						Type:     HintTypeError,
						Location: NewLocation(h.Path, line, lineOffset, size),
					}
					msgs = append(msgs, msg)
				}
				//indexed位置
				sym := splitArgs(string(h.Code[stridx[0][0]:stridx[0][1]]))
				sym = sym[1:]
				for i := 0; i < len(sym); i++ {
					irregular := true
					idx := strings.Index(sym[i], KWIndexed)
					if idx == 0 || idx == -1 {
						irregular = false
					}
					if irregular {
						msg := HintMessage{
							Message:  "EVENT方法中的indexed写法不规范",
							Type:     HintTypeError,
							Location: NewLocation(h.Path, line, lineOffset, size),
						}
						msgs = append(msgs, msg)
					}
				}
			}
		}
	}
	return msgs, nil
}

func (h *Hint) payableCheck() (HintMessages, error) {
	var msgs HintMessages
	//payable和unmutable不能共存
	for _, v := range functionTree.Root {
		if v.Info.Export == ExportTypeNone {
			//ignore
			continue
		}
		if v.Info.Payable {
			if v.Info.Export != ExportTypeMutable {
				msg := HintMessage{
					Message:  "Payable方法必须使用关键字MUTABLE进行导出",
					Type:     HintTypeError,
					Location: NewLocation(v.Info.Location.Path, v.Info.Location.Line+1, 1, v.Info.Location.Size),
				}
				msgs = append(msgs, msg)
			}
		}
	}
	return msgs, nil
}

func (h *Hint) exportCheck() (HintMessages, error) {
	var msgs HintMessages
	exportReg := `[^(;|\r\n|\r|\n|\{|\})]*(\s+)(%s)(\s*)(\({1})([a-zA-Z0-9_\*\s,]*)(\){1})`
	for _, v := range functionTree.Root {
		if v.Info.Export == ExportTypeNone {
			//ignore
			continue
		}
		export := fmt.Sprintf(exportReg, escape(v.Name))
		reg := regexp.MustCompile(export)
		stridx := reg.FindAllStringIndex(string(h.Code), -1)
		if len(stridx) != 0 {
			//find method
			offset := stridx[0][0]
			size := stridx[0][1] - stridx[0][0]
			line, lineOffset := GetLineNumber(offset, fileContent[h.Path])
			left, right := removeSpaceAndParen(v.Info.Signature)
			if len(left) != 1 {
				msg := HintMessage{
					Message:  "方法的返回值为不支持的类型：" + strings.Join(left, " "),
					Type:     HintTypeError,
					Location: NewLocation(h.Path, line, lineOffset, size),
				}
				msgs = append(msgs, msg)
			} else {
				if !isSupportedType(left[0]) {
					msg := HintMessage{
						Message:  "方法的返回值为不支持的类型：" + left[0],
						Type:     HintTypeError,
						Location: NewLocation(h.Path, line, lineOffset, size),
					}
					msgs = append(msgs, msg)
				}
			}
			types := []string{}
			for i := 0; i < len(right); i++ {
				if !isSupportedType(right[i]) {
					types = append(types, right[i])
				}
			}
			if len(types) != 0 {
				msg := HintMessage{
					Message:  "方法的参数为不支持的类型：" + strings.Join(types, " ,"),
					Type:     HintTypeError,
					Location: NewLocation(h.Path, line, lineOffset, size),
				}
				msgs = append(msgs, msg)
			}
		}
	}
	return msgs, nil
}

func (h *Hint) checkUnmutableFunction() (HintMessages, error) {
	var msgs HintMessages
	for _, v := range functionTree.Root {
		if v.Info.Export == ExportTypeUnmutable {
			for _, call := range v.Call {
				for _, f := range functionTree.Root {
					if f.Name == call.Name && (f.Info.Export == ExportTypeMutable || f.Info.Payable) {
						msg := HintMessage{
							Message:  "UNMUTABLE的方法中调用了MUTABLE或PAYABLE方法：" + call.Name,
							Type:     HintTypeError,
							Location: call.Location,
						}
						msgs = append(msgs, msg)
					}
					if f.Name == call.Name && (f.Name == "SendFromContract" || f.Name == "TransferFromContract") {
						msg := HintMessage{
							Message:  "UNMUTABLE的方法中调用了SendFromContract或TransferFromContract方法",
							Type:     HintTypeError,
							Location: call.Location,
						}
						msgs = append(msgs, msg)
					}
				}
			}
		}
	}
	return msgs, nil
}

func (h *Hint) typeCheck() {

}
