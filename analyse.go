package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//key 类型检查，如果key写在其他方法之间，报错并提示 <complete>
//没有constructor <complete>
//多个constructor
//call,第一个参数检查,参数类型检查及返回值检查,event,参数类型检查
//mutable
//unmutable，key的写入检查，sendfromcontract,transferfromcontract检查
//payable,input的检查，sendfromcontract,transferfromcontract
//导出方法的类型和参数检查
//uint256及address类型

type Hint struct {
	Code           []byte
	ConstructorPos int
}

type HintType int

const (
	HintTypeWarning HintType = iota
	HintTypeError
)

type HintMessage struct {
	Message string
	Type    HintType
	Offset  int
	Size    int
}

func newHint(code []byte) *Hint {

	return &Hint{
		Code: code,
	}
}

func (h *Hint) contructorCheck() ([]HintMessage, error) {
	var msgs []HintMessage
	reg := regexp.MustCompile(constructorReg)
	idx := reg.FindAllStringIndex(string(h.Code), -1)

	if len(idx) == 0 {
		msg := HintMessage{
			Message: "必须定义一个Contructor方法",
			Type:    HintTypeError,
			Offset:  0,
			Size:    0,
		}
		msgs = append(msgs, msg)
		return msgs, nil
	}
	if len(idx) > 1 {
		for i := 0; i < len(idx); i++ {
			msg := HintMessage{
				Message: "重复定义Contructor方法",
				Type:    HintTypeError,
				Offset:  idx[i][0],
				Size:    idx[i][1] - idx[i][0],
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
				Message: "KEY必须定义在construct之前",
				Type:    HintTypeWarning,
				Offset:  offset,
				Size:    len(v.FieldName),
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
	for _, v := range funcLists {
		call := fmt.Sprintf(callReg, v.Name)
		reg := regexp.MustCompile(call)
		stridx := reg.FindAllStringIndex(string(h.Code), -1)
		if len(stridx) != 0 {
			//find call
			left, right := removeSpaceAndParen(v.Signature)
			if len(left) != 1 {
				msg := HintMessage{
					Message: "CALL方法的返回值为不支持的类型：" + strings.Join(left, " "),
					Type:    HintTypeError,
					Offset:  stridx[0][0],
					Size:    stridx[0][1] - stridx[0][0],
				}
				msgs = append(msgs, msg)
			} else {
				if !isSupportedType(left[0]) {
					msg := HintMessage{
						Message: "CALL方法的返回值为不支持的类型：" + left[0],
						Type:    HintTypeError,
						Offset:  stridx[0][0],
						Size:    stridx[0][1] - stridx[0][0],
					}
					msgs = append(msgs, msg)
				}
			}
			if len(right) == 0 {
				msg := HintMessage{
					Message: "CALL方法至少需要一个类型为CallParams的参数",
					Type:    HintTypeError,
					Offset:  stridx[0][0],
					Size:    stridx[0][1] - stridx[0][0],
				}
				msgs = append(msgs, msg)
			}

			for i := 0; i < len(right); i++ {
				if i == 0 { //callprams
					if right[i] != "CallParams" {
						msg := HintMessage{
							Message: "CALL方法的第一个参数为不支持的类型：" + right[i] + ", 支持的类型为CallParams",
							Type:    HintTypeError,
							Offset:  stridx[0][0],
							Size:    stridx[0][1] - stridx[0][0],
						}
						msgs = append(msgs, msg)
					}
				} else { //input type
					if !isSupportedType(right[i]) {
						msg := HintMessage{
							Message: "CALL方法的参数为不支持的类型：" + right[i],
							Type:    HintTypeError,
							Offset:  stridx[0][0],
							Size:    stridx[0][1] - stridx[0][0],
						}
						msgs = append(msgs, msg)
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
	for _, v := range funcLists {
		event := fmt.Sprintf(eventReg, v.Name)
		reg := regexp.MustCompile(event)
		stridx := reg.FindAllStringIndex(string(h.Code), -1)
		if len(stridx) != 0 {
			//find event
			left, right := removeSpaceAndParen(v.Signature)
			fmt.Printf("left %v right %v\n", left, right)
			if len(right) == 0 {
				msg := HintMessage{
					Message: "EVENT方法至少需要一个参数",
					Type:    HintTypeError,
					Offset:  stridx[0][0],
					Size:    stridx[0][1] - stridx[0][0],
				}
				msgs = append(msgs, msg)
			}
		}
		//类型判断
		//indexed位置
		
		//string index, address,addres indexed 
	}
	return msgs, nil
}
