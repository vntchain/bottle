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
	"bytes"
	"os/exec"
	"path"
	"strings"
)

var llvmDir = ""
var sysrootDir = ""

func SetEnvPath() {
	wasmceptionDir := wasmCeptionFlag

	llvmDir = path.Join(wasmceptionDir, "dist")
	sysrootDir = path.Join(wasmceptionDir, "sysroot")
}

func getClangOptions(options string) []string {
	clangFlags := `--target=wasm32-unknown-unknown-wasm --sysroot=` + sysrootDir + ` -fdiagnostics-print-source-range-info -fno-exceptions`
	if options == "" {
		return strings.Split(clangFlags, " ")
	}
	availableOptions := []string{
		"-O0", "-O1", "-O2", "-O3", "-O4", "-Os", "-fno-exceptions", "-fno-rtti",
		"-ffast-math", "-fno-inline", "-std=c99", "-std=c89", "-std=c++14",
		"-std=c++1z", "-std=c++11", "-std=c++98", "-g"}
	safeOptions := "-c"
	for _, o := range availableOptions {
		if strings.Contains(options, o) {
			safeOptions += " " + o
		} else if strings.Contains(o, "-std=") && strings.Contains(strings.ToLower(options), o) {
			safeOptions += " " + o
		}
	}
	return strings.Split(clangFlags+" "+safeOptions, " ")
}

func getLldOptions(options string) []string {
	clangFlags := `--target=wasm32-unknown-unknown-wasm --sysroot=` + sysrootDir + ` -nostartfiles -Wl,--allow-undefined,--no-entry,--no-threads`
	if options == "" {
		return strings.Split(clangFlags, " ")
	}
	availableOptions := []string{"--import-memory", "-g", "-O0", "-O1", "-O2", "-O3", "-O4", "-Os"}
	safeOptions := ""
	for _, o := range availableOptions {
		if strings.Contains(options, o) {
			safeOptions += " -Wl, " + o
		}
	}
	return strings.Split(clangFlags+safeOptions, " ")
}

func buildCFile(options string, input string, output string) {
	cmdPath := llvmDir + "/bin/clang"
	cmdArgs := append(getLldOptions(options), []string{input, "-o", output, "-I", includeDir}...)
	cmd := exec.Command(cmdPath, cmdArgs...)
	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		panic(cmd.Stderr)
	}
}
