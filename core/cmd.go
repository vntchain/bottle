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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/dollarshaveclub/line"
	"github.com/vntchain/go-vnt/accounts/abi"
	cmdutils "github.com/vntchain/go-vnt/cmd/utils"
	"github.com/vntchain/go-vnt/core/wavm/utils"
	cli "gopkg.in/urfave/cli.v1"
)

var wasmCeptionFlag string
var vntIncludeFlag string

const (
	VersionMajor = 0      // Major version component of the current release
	VersionMinor = 6      // Minor version component of the current release
	VersionPatch = 0      // Patch version component of the current release
	VersionMeta  = "beta" // Version metadata to append to the version string
)

// Version holds the textual version string.
var Version = func() string {
	v := fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	if VersionMeta != "" {
		v += "-" + VersionMeta
	}
	return v
}()

// NewApp creates an app with sane defaults.
func NewApp(gitCommit, usage, includeFlag, wasmFlag string) *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Author = ""
	//app.Authors = nil
	app.Email = ""
	app.Version = Version
	if len(gitCommit) >= 8 {
		app.Version += "-" + gitCommit[:8]
	}
	app.Usage = usage
	vntIncludeFlag = includeFlag
	wasmCeptionFlag = wasmFlag
	return app
}

var (
	// flags that configure the node
	contractCodeFlag = cli.StringFlag{
		Name:  "code",
		Usage: "Specific a contract code path",
	}
	outputFlag = cmdutils.DirectoryFlag{
		Name:  "output",
		Usage: "Specific a output directory path",
	}
	includeFlag = cmdutils.DirectoryFlag{
		Name:  "include",
		Usage: "Specific the head file directory need by contract",
	}
	abiFlag = cli.StringFlag{
		Name:  "abi",
		Usage: "Specific a abi path need by contract",
	}
	wasmFlag = cli.StringFlag{
		Name:  "wasm",
		Usage: "Specific a wasm path",
	}
	compressFileFlag = cli.StringFlag{
		Name:  "file",
		Usage: "Specific a compress file path to decompress",
	}
	CompileCmd = cli.Command{
		Action:    compile,
		Name:      "compile",
		Usage:     "Compile contract code to wasm and compress",
		ArgsUsage: "<code file>",
		Category:  "COMPILE COMMANDS",
		Description: `
		bottle compile

Compile contract code to wasm and compress
		`,
		Flags: []cli.Flag{
			contractCodeFlag,
			outputFlag,
			includeFlag,
		},
	}
	CompressCmd = cli.Command{
		Action:    compress,
		Name:      "compress",
		Usage:     "Compress wasm and abi",
		ArgsUsage: "<code file> <abi file>",
		Category:  "COMPRESS COMMANDS",
		Description: `
		bottle compress

Compress wasm and abi
		`,
		Flags: []cli.Flag{
			wasmFlag,
			abiFlag,
			outputFlag,
		},
	}
	DecompressCmd = cli.Command{
		Action:    decompress,
		Name:      "decompress",
		Usage:     "Deompress file into wasm and abi",
		ArgsUsage: "<code file> <abi file>",
		Category:  "DECOMPRESS COMMANDS",
		Description: `
		bottle decompress

Deompress file into wasm and abi
		`,
		Flags: []cli.Flag{
			compressFileFlag,
			outputFlag,
		},
	}
	HintCmd = cli.Command{
		Action:    hint,
		Name:      "hint",
		Usage:     "Contract hint",
		ArgsUsage: "<code file> <abi file>",
		Category:  "HINT COMMANDS",
		Description: `
		bottle hint

Contract hint
		`,
		Flags: []cli.Flag{
			contractCodeFlag,
		},
	}
	InitCmd = cli.Command{
		Action:    initContract,
		Name:      "init",
		Usage:     "Initialize contract project",
		ArgsUsage: "<project name> <directory path>",
		Category:  "INIT COMMANDS",
		Description: `
		bottle init

Contract init
		`,
	}
	BuildCmd = cli.Command{
		Action:    buildContract,
		Name:      "build",
		Usage:     "Build contract",
		ArgsUsage: "<project name> <directory path>",
		Category:  "BUILD COMMANDS",
		Description: `
		bottle build

Contract build
		`,
	}
)

func compile(ctx *cli.Context) error {
	start := time.Now()

	if err := hint(ctx); err != nil {
		return err
	}

	codePath = ctx.String(contractCodeFlag.Name)
	includeDir = ctx.String(includeFlag.Name)
	outputDir = ctx.String(outputFlag.Name)
	if codePath == "" {
		fmt.Printf("Error:No Contract Code\n")
		os.Exit(-1)
	}
	mustCFile(codePath)
	if outputDir == "" {
		outputDir = path.Join(path.Dir(codePath), "output")
	}
	if includeDir == "" {
		includeDir = path.Dir(codePath)
	}

	code, err := ioutil.ReadFile(codePath)
	if err != nil {
		return err
	}
	cmd([]string{codePath})
	abigen := newAbiGen(code)
	abigen.removeComment()
	abigen.parseMethod()
	// abigen.parseKey()
	abigen.parseEvent()
	abigen.parseCall()
	abigen.parseConstructor()

	var pack []interface{}
	if abigen.abi.Constructor.Name != "" {
		pack = append(pack, abigen.abi.Constructor)
	}

	for _, v := range abigen.abi.Methods {
		pack = append(pack, v)
	}
	for _, v := range abigen.abi.Events {
		pack = append(pack, v)
	}
	for _, v := range abigen.abi.Calls {
		pack = append(pack, v)
	}
	res, err := json.Marshal(pack)
	if err != nil {
		return err
	}
	abijson := string(res)
	abires, err := abi.JSON(bytes.NewBuffer(res))
	if err != nil {
		return err
	}
	abipath := abires.Constructor.Name + ".abi"
	err = writeFile(path.Join(outputDir, abipath), res)
	if err != nil {
		return err
	}
	pre := abigen.insertRegistryCode()
	// pre = abigen.insertMutableCode(pre)
	codeOutput := path.Join(outputDir, abires.Constructor.Name+"_precompile.c")
	err = writeFile(codeOutput, pre)
	if err != nil {
		return err
	}
	// fmt.Printf("Precompile code path: %s\n", codeOutput)
	wasmOutput := path.Join(outputDir, abires.Constructor.Name+".wasm")
	SetEnvPath()
	BuildWasm(codeOutput, wasmOutput)
	// fmt.Printf("Wasm path: %s\n", wasmOutput)
	wasm, err := ioutil.ReadFile(wasmOutput)
	if err != nil {
		return err
	}
	cpsPath := path.Join(outputDir, abires.Constructor.Name+".compress")
	cpsRes := utils.CompressWasmAndAbi(res, wasm, nil)
	err = writeFile(cpsPath, cpsRes)
	if err != nil {
		return err
	}
	hexPath := path.Join(outputDir, abires.Constructor.Name+".hex")
	hexString := "0x" + hex.EncodeToString(cpsRes)
	err = writeFile(hexPath, []byte(hexString))
	if err != nil {
		return err
	}
	deployCodePath := path.Join(outputDir, abires.Constructor.Name+".js")
	err = writeFile(deployCodePath, []byte(deployText(string(res), hexString)))
	if err != nil {
		return err
	}

	contractJsonPath := path.Join(outputDir, abires.Constructor.Name+".json")
	contract := Contract{
		ContractName: abires.Constructor.Name,
		Abi:          abijson,
		Bytecode:     hexString,
	}
	contractJson, err := json.Marshal(contract)
	if err != nil {
		return err
	}
	err = writeFile(contractJsonPath, contractJson)
	if err != nil {
		return err
	}

	output := line.New(os.Stdout, "", "", line.WhiteColor)
	li := output.Prefix(">>>").Cyan()
	li.Printf("Compile finished. %s\n", time.Since(start).String())
	li.Printf("Input file\n")
	li = output.Prefix("   ").Cyan()
	li.Printf("Contract path :%s\n", codePath)
	li = output.Prefix(">>>").Cyan()
	li.Printf("Output file\n")
	li = output.Prefix("   ").Cyan()
	li.Printf("Abi path: %s\n", path.Join(outputDir, abipath))
	li.Printf("Precompile code path: %s\n", codeOutput)
	li.Printf("Wasm path: %s\n", wasmOutput)
	li.Printf("Compress Data path: %s\n", cpsPath)
	li.Printf("Compress Hex Data path: %s\n", hexPath)
	li.Printf("Deploy JS path: %s\n", deployCodePath)
	li.Printf("Contract JSON path: %s\n", contractJsonPath)
	li = output.Prefix(">>>").Cyan()
	li.Printf("Please use %s when you want to create a contract\n", abires.Constructor.Name+".compress")
	return nil
}

func compress(ctx *cli.Context) error {
	start := time.Now()
	wasmPath = ctx.String(wasmFlag.Name)
	abiPath = ctx.String(abiFlag.Name)
	outputDir = ctx.String(outputFlag.Name)
	if wasmPath == "" {
		fmt.Printf("Error:No Wasm File\n")
		os.Exit(-1)
	}
	if abiPath == "" {
		fmt.Printf("Error:No Abi File\n")
		os.Exit(-1)
	}

	if outputDir == "" {
		outputDir = path.Join(path.Dir(wasmPath), "output")
	}

	wasm, err := ioutil.ReadFile(wasmPath)
	if err != nil {
		return err
	}
	abijson, err := ioutil.ReadFile(abiPath)
	if err != nil {
		return err
	}

	abires, err := abi.JSON(bytes.NewBuffer(abijson))
	if err != nil {
		return err
	}
	cpsPath := path.Join(outputDir, abires.Constructor.Name+".compress")
	cpsRes := utils.CompressWasmAndAbi(abijson, wasm, nil)
	err = writeFile(cpsPath, cpsRes)
	if err != nil {
		return err
	}
	output := line.New(os.Stdout, "", "", line.WhiteColor)
	li := output.Prefix(">>>").Cyan()
	li.Printf("Compress finished. %s\n", time.Since(start).String())
	li.Printf("Input file\n")
	li = output.Prefix("   ").Cyan()
	li.Printf("Wasm path :%s\n", wasmPath)
	li.Printf("Abi path :%s\n", abiPath)
	li = output.Prefix(">>>").Cyan()
	li.Printf("Output file\n")
	li = output.Prefix("   ").Cyan()
	li.Printf("Compress Data path: %s\n", cpsPath)
	li = output.Prefix(">>>").Cyan()
	li.Printf("Please use %s when you want to create a contract\n", abires.Constructor.Name+".compress")
	return nil
}

func decompress(ctx *cli.Context) error {
	start := time.Now()
	compressPath = ctx.String(compressFileFlag.Name)
	outputDir = ctx.String(outputFlag.Name)
	if compressPath == "" {
		fmt.Printf("Error:No Compress File\n")
	}
	if outputDir == "" {
		outputDir = path.Join(path.Dir(compressPath), "output")
	}
	comContent, err := ioutil.ReadFile(compressPath)
	if err != nil {
		return err
	}
	wasmcode, _, err := utils.DecodeContractCode(comContent)
	if err != nil {
		return err
	}
	abires, err := abi.JSON(bytes.NewBuffer(wasmcode.Abi))
	if err != nil {
		return err
	}
	wasmoutputPath := path.Join(outputDir, abires.Constructor.Name+".wasm")
	err = writeFile(wasmoutputPath, wasmcode.Code)
	if err != nil {
		return err
	}
	abioutputPath := path.Join(outputDir, "abi.json")
	err = writeFile(abioutputPath, wasmcode.Abi)
	if err != nil {
		return err
	}
	output := line.New(os.Stdout, "", "", line.WhiteColor)
	li := output.Prefix(">>>").Cyan()
	li.Printf("Decompress finished. %s\n", time.Since(start).String())
	li.Printf("Input file\n")
	li = output.Prefix("   ").Cyan()
	li.Printf("Compress file path :%s\n", compressPath)
	li = output.Prefix(">>>").Cyan()
	li.Printf("Output file\n")
	li = output.Prefix("   ").Cyan()
	li.Printf("wasm path: %s\n", wasmoutputPath)
	li.Printf("abi path: %s\n", abioutputPath)
	return nil
}

func hint(ctx *cli.Context) error {
	codePath = ctx.String(contractCodeFlag.Name)
	// fileContent = readfile(codePath)
	var code []byte
	var err error
	code, err = ioutil.ReadFile(codePath)
	if err != nil {
		return err
	}
	code = removeComment(string(code))
	cmdErr := cmd([]string{codePath})
	if cmdErr != nil {
		return cmdErr
	}
	// jsonres, _ := json.Marshal(varLists)
	// fmt.Printf("vallist %s\n", jsonres)

	// structres, _ := json.Marshal(structLists)
	// fmt.Printf("structres %s\n", structres)
	hint := newHint(codePath, code)
	var msgs HintMessages
	msg, err := hint.contructorCheck()
	if err != nil {
		return err
	}
	msgs = append(msgs, msg...)
	msg, err = hint.keyCheck()
	if err != nil {
		return err
	}
	msgs = append(msgs, msg...)
	msg, err = hint.callCheck()
	if err != nil {
		return err
	}
	msgs = append(msgs, msg...)
	msg, err = hint.eventCheck()
	if err != nil {
		return err
	}
	msgs = append(msgs, msg...)
	msg, err = hint.payableCheck()
	if err != nil {
		return err
	}
	msgs = append(msgs, msg...)
	msg, err = hint.exportCheck()
	if err != nil {
		return err
	}
	msgs = append(msgs, msg...)
	msg, err = hint.checkUnmutableFunction()
	if err != nil {
		return err
	}
	msgs = append(msgs, msg...)
	if len(msgs) != 0 {
		return cli.NewExitError(msgs.ToString(), -1)
	} else {
		return nil
	}
}

func initContract(ctx *cli.Context) error {
	fmt.Printf("init\n")
	return nil
}

func buildContract(ctx *cli.Context) error {
	fmt.Printf("build\n")
	return nil
}
