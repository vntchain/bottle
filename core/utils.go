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
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/dollarshaveclub/line"
)

type ContentPerLine struct {
	Content []byte
	Offset  int
}

func cutUSR(t string) string {
	pt := t
	idx := strings.LastIndex(t, "@FI@")
	if idx != -1 {
		pt = t[0:idx]
	}
	return pt
}

func readfile(fi *os.File) []ContentPerLine {
	// contents := []string{}
	contensPerLine := []ContentPerLine{}
	br := bufio.NewReader(fi)
	offset := 0
	for {
		a, c := br.ReadBytes('\n')
		// contents = append(contents, string(a))
		contensPerLine = append(contensPerLine, ContentPerLine{
			Content: a,
			Offset:  offset,
		})
		offset += len(a)
		if c == io.EOF {
			break
		}
	}
	// for i, v := range contensPerLine {
	// 	fmt.Printf("line %d content %s\n", i, v)
	// }
	return contensPerLine
}

func GetLineNumber(offset int, filecontent []ContentPerLine) (int, int) {
	for i := 1; i < len(filecontent); i++ {
		if offset >= filecontent[i-1].Offset && offset < filecontent[i].Offset {
			return i, offset - filecontent[i-1].Offset + 1
		} else if i == len(filecontent)-1 && offset >= filecontent[i].Offset {
			return i + 1, offset - filecontent[i].Offset + 1
		}
	}
	return 1, 1
}

//KEY _complex s3;
var astKeyReg = `([ ]*)(KEY)([ ]+)(int(|32|64)|uint(|32|64|256)|address|string|bool|mapping|array|struct)([\s\S]*)`
var astKeyRegFmt = `([ ]*)(KEY)([ ]+)(int(|32|64)|uint(|32|64|256)|address|string|bool|mapping|array|struct|%s)([\s\S]*)`

const keyTypeReg = `(int(|32|64)|uint(|32|64|256)|address|string|bool|mapping|array|struct)`

func isKey(input string, structnames string) bool {
	keyReg := ""
	if structnames == "" {
		keyReg = astKeyReg
	} else {
		keyReg = fmt.Sprintf(astKeyRegFmt, structnames)
	}

	reg, err := regexp.Compile(keyReg)
	flag := false
	if err != nil {
		return flag
	}
	idx := reg.FindAllStringIndex(input, -1)
	if len(idx) == 0 {
		return flag
	}
	flag = true
	return flag
}

func isSupportKeyType() bool {
	return false
}

func DeepCopy(dst, src interface{}) error {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(src); err != nil {
		return err
	}

	return gob.NewDecoder(&buffer).Decode(dst)
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func writeFile(file string, content []byte) error {
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return err
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}
	f.Close()
	return os.Rename(f.Name(), file)
}

func unpackZip(dst string, src string, depth int) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		split := strings.Split(file.Name, "/")
		if len(split) <= depth {
			continue
		}
		filename := path.Join(dst, strings.Join(split[depth:], "/"))
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		var content []byte
		_, err = io.ReadFull(rc, content)
		if err != nil {
			return err
		}
		err = writeFile(filename, content)
		if err != nil {
			return err
		}
	}
	return nil
}

func isEmpty(dir string) bool {
	files, _ := ioutil.ReadDir(dir)
	if len(files) == 0 {
		return true
	} else {
		return false
	}
}

func isSupportedType(tp string) bool {
	types := []string{"int32", "int64", "uint32", "uint64", "uint256", "string", "address", "bool", "_Bool", "void"}
	typesmap := map[string]bool{}
	for _, v := range types {
		typesmap[v] = true
	}
	return typesmap[tp]
}

func isSupportedKeyType(tp string) bool {
	types := []string{"int32", "int64", "uint32", "uint64", "uint256", "string", "address", "bool", "_Bool", "struct", "mapping", "array"}
	typesmap := map[string]bool{}
	for _, v := range types {
		typesmap[v] = true
	}
	return typesmap[tp]
}

func escape(input string) string {
	if len(input) == 0 {
		return input
	}
	if input[0:1] == "$" {
		input = `\` + input
	}
	return `[\*]?` + input
}

func deployText(abi, code string) string {
	return fmt.Sprintf(`
	var projectContract = vnt.core.contract(%s);
	var project = projectContract.new(
    {
     	from: vnt.core.accounts[0], 
     	data: '%s', 
     	gas: '4000000'
    }, function (e, contract){
    	console.log(e, contract);
    	if (typeof contract.address !== 'undefined') {
        	console.log('Contract address: ' + contract.address + ' transactionHash: ' + contract.transactionHash);
   	 	}
 	})
	`, abi, code)
}

func PrintfHeader(output *line.Line, format string, a ...interface{}) {
	li := output.Prefix("").White()
	li.Printf(format, a...)
	li.Printf("%s\n", strings.Repeat("=", len(fmt.Sprintf(format, a...))))
}

func PrintfBody(output *line.Line, before string, after string) {
	li := output.Prefix("   > ").White()
	pad := rightPadBytes([]byte(before), 25)
	text := fmt.Sprintf("%s%s\n", string(pad), after)
	li.Printf("%s", text)
}

func rightPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}
	padded := []byte(strings.Repeat(" ", l))
	copy(padded, slice)
	return padded
}

// merges extraHeaders into headers and returns headers
func merge(contract1, contrac2 Contract) Contract {
	contract1.ContractName = contrac2.ContractName
	contract1.Abi = contrac2.Abi
	contract1.Bytecode = contrac2.Bytecode
	contract1.SourcePath = contrac2.SourcePath
	contract1.UpdatedAt = contrac2.UpdatedAt
	return contract1
}
