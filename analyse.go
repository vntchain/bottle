package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//key 类型检查，如果key写在其他方法之间，报错并提示  未完成，当前只判断key写在constructor之前，
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

func newHint(path string, code []byte) *Hint {

	return &Hint{
		Path: path,
		Code: code,
	}
}

func (h *Hint) contructorCheck() ([]HintMessage, error) {
	var msgs []HintMessage
	reg := regexp.MustCompile(constructorReg)
	idx := reg.FindAllStringIndex(string(h.Code), -1)

	if len(idx) == 0 {
		msg := HintMessage{
			Message:  "必须定义一个Contructor方法",
			Type:     HintTypeError,
			Location: NewLocation(h.Path, 0, 0),
		}
		msgs = append(msgs, msg)
		return msgs, nil
	}
	if len(idx) > 1 {
		for i := 0; i < len(idx); i++ {
			msg := HintMessage{
				Message:  "重复定义Contructor方法",
				Type:     HintTypeError,
				Location: NewLocation(h.Path, idx[i][0], idx[i][1]-idx[i][0]),
			}
			msgs = append(msgs, msg)
		}

		return msgs, nil
	}
	return msgs, nil
}

func (h *Hint) keyCheck() ([]HintMessage, error) {
	var msgs []HintMessage
	for _, v := range varLists.Root {
		loc := v.FieldLocation
		offset, err := strconv.Atoi(loc)
		if err != nil {
			return nil, err
		}
		if offset >= h.ConstructorPos {
			msg := HintMessage{
				Message:  "KEY必须定义在construct之前",
				Type:     HintTypeWarning,
				Location: NewLocation(h.Path, offset, len(v.FieldName)),
			}
			msgs = append(msgs, msg)
		}
	}
	if len(msgs) != 0 {
		return msgs, nil
	}
	return msgs, nil
}

func (h *Hint) callCheck() ([]HintMessage, error) {
	var msgs []HintMessage
	callReg := `(CALL)[^(;|\r|\n|\{|\})]*(%s)(\s*)(\({1})([a-zA-Z0-9_\$\s,]*)(\){1})`
	for _, v := range functionTree.Root {
		call := fmt.Sprintf(callReg, escape(v.Name))
		reg := regexp.MustCompile(call)
		stridx := reg.FindAllStringIndex(string(h.Code), -1)
		if len(stridx) != 0 {
			//find call
			left, right := removeSpaceAndParen(v.Info.Signature)
			if len(left) != 1 {
				msg := HintMessage{
					Message:  "CALL方法的返回值为不支持的类型：" + strings.Join(left, " "),
					Type:     HintTypeError,
					Location: NewLocation(h.Path, stridx[0][0], stridx[0][1]-stridx[0][0]),
				}
				msgs = append(msgs, msg)
			} else {
				if !isSupportedType(left[0]) {
					msg := HintMessage{
						Message:  "CALL方法的返回值为不支持的类型：" + left[0],
						Type:     HintTypeError,
						Location: NewLocation(h.Path, stridx[0][0], stridx[0][1]-stridx[0][0]),
					}
					msgs = append(msgs, msg)
				}
			}
			if len(right) == 0 {
				msg := HintMessage{
					Message:  "CALL方法至少需要一个类型为CallParams的参数",
					Type:     HintTypeError,
					Location: NewLocation(h.Path, stridx[0][0], stridx[0][1]-stridx[0][0]),
				}
				msgs = append(msgs, msg)
			} else {
				for i := 0; i < len(right); i++ {
					if i == 0 { //callprams
						if right[i] != "CallParams" {
							msg := HintMessage{
								Message:  "CALL方法的第一个参数为不支持的类型：" + right[i] + ", 支持的类型为CallParams",
								Type:     HintTypeError,
								Location: NewLocation(h.Path, stridx[0][0], stridx[0][1]-stridx[0][0]),
							}
							msgs = append(msgs, msg)
						}
					} else { //input type
						if !isSupportedType(right[i]) {
							msg := HintMessage{
								Message:  "CALL方法的参数为不支持的类型：" + right[i],
								Type:     HintTypeError,
								Location: NewLocation(h.Path, stridx[0][0], stridx[0][1]-stridx[0][0]),
							}
							msgs = append(msgs, msg)
						}
					}
				}
			}
		}
	}

	return msgs, nil
}

func (h *Hint) eventCheck() ([]HintMessage, error) {
	var msgs []HintMessage
	eventReg := `(EVENT)[^(;|\r|\n|\{|\})]*(%s)(\s*)(\({1})([a-zA-Z0-9_\$\s,]*)(\){1})`
	for _, v := range functionTree.Root {
		event := fmt.Sprintf(eventReg, escape(v.Name))
		reg := regexp.MustCompile(event)
		stridx := reg.FindAllStringIndex(string(h.Code), -1)
		if len(stridx) != 0 {
			//find event
			_, right := removeSpaceAndParen(v.Info.Signature)
			if len(right) == 0 {
				msg := HintMessage{
					Message:  "EVENT方法至少需要一个参数",
					Type:     HintTypeError,
					Location: NewLocation(h.Path, stridx[0][0], stridx[0][1]-stridx[0][0]),
				}
				msgs = append(msgs, msg)
			} else {
				//类型判断
				//indexed位置
				//string index, address,addres indexed
				for i := 0; i < len(right); i++ {
					if !isSupportedType(right[i]) {
						msg := HintMessage{
							Message:  "EVENT方法的参数为不支持的类型：" + right[i],
							Type:     HintTypeError,
							Location: NewLocation(h.Path, stridx[0][0], stridx[0][1]-stridx[0][0]),
						}
						msgs = append(msgs, msg)
						continue
					}
				}
				sym := removeSymbol(string(h.Code[stridx[0][0]:stridx[0][1]]))
				sym = sym[2:]
				for i := 0; i < len(sym); i++ {
					irregular := false
					if sym[i] == KWIndexed { //match indexed
						if i-1 < 0 {
							irregular = true
						} else {
							if !isSupportedType(sym[i-1]) {
								irregular = true
							}
						}
						if irregular {
							msg := HintMessage{
								Message:  "EVENT方法中的indexed写法不规范",
								Type:     HintTypeError,
								Location: NewLocation(h.Path, stridx[0][0], stridx[0][1]-stridx[0][0]),
							}
							msgs = append(msgs, msg)
						}
					}
				}
			}

		}

	}
	return msgs, nil
}

func (h *Hint) payableCheck() ([]HintMessage, error) {
	var msgs []HintMessage
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
					Location: v.Info.Location,
				}
				msgs = append(msgs, msg)
			}
		}
	}
	return msgs, nil
}

func (h *Hint) exportCheck() ([]HintMessage, error) {
	var msgs []HintMessage
	exportReg := `[^(;|\r|\n|\{|\})]*(%s)(\s*)(\({1})([a-zA-Z0-9_\$\s,]*)(\){1})`
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
			left, right := removeSpaceAndParen(v.Info.Signature)
			if len(left) != 1 {
				msg := HintMessage{
					Message:  "方法的返回值为不支持的类型：" + strings.Join(left, " "),
					Type:     HintTypeError,
					Location: NewLocation(h.Path, stridx[0][0], stridx[0][1]-stridx[0][0]),
				}
				msgs = append(msgs, msg)
			} else {
				if !isSupportedType(left[0]) {
					msg := HintMessage{
						Message:  "方法的返回值为不支持的类型：" + left[0],
						Type:     HintTypeError,
						Location: NewLocation(h.Path, stridx[0][0], stridx[0][1]-stridx[0][0]),
					}
					msgs = append(msgs, msg)
				}
			}
			for i := 0; i < len(right); i++ {
				if !isSupportedType(right[i]) {
					msg := HintMessage{
						Message:  "方法的参数为不支持的类型：" + right[i],
						Type:     HintTypeError,
						Location: NewLocation(h.Path, stridx[0][0], stridx[0][1]-stridx[0][0]),
					}
					msgs = append(msgs, msg)
				}
			}
		}
	}
	return msgs, nil
}

func (h *Hint) checkUnmutableFunction() ([]HintMessage, error) {
	var msgs []HintMessage
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
