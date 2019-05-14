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

	"github.com/vntchain/bottle/js"

	"github.com/dollarshaveclub/line"
	"github.com/vntchain/go-vnt/accounts/abi"
	cmdutils "github.com/vntchain/go-vnt/cmd/utils"
	"github.com/vntchain/go-vnt/core/wavm/utils"
	"github.com/vntchain/go-vnt/rpc"
	cli "gopkg.in/urfave/cli.v1"
)

var wasmCeptionFlag string
var vntIncludeFlag string
var nodePathFlag string

const (
	VersionMajor = 0      // Major version component of the current release
	VersionMinor = 6      // Minor version component of the current release
	VersionPatch = 0      // Patch version component of the current release
	VersionMeta  = "beta" // Version metadata to append to the version string
)

var JS = deps.MustAsset("bottle.js")

// Version holds the textual version string.
var Version = func() string {
	v := fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	if VersionMeta != "" {
		v += "-" + VersionMeta
	}
	return v
}()

// NewApp creates an app with sane defaults.
func NewApp(gitCommit, usage, includeFlag, wasmFlag, nodeFlag string) *cli.App {
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
	nodePathFlag = nodeFlag
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
		Usage: "Specific a abi path needed by contract",
	}
	wasmFlag = cli.StringFlag{
		Name:  "wasm",
		Usage: "Specific a wasm path",
	}
	compressFileFlag = cli.StringFlag{
		Name:  "file",
		Usage: "Specific a compress file path to decompress",
	}
	resetFlag = cli.BoolFlag{
		Name:  "reset",
		Usage: "Run all migrations from the beginning, instead of running from the last completed migration",
	}
	fromFlag = cli.IntFlag{
		Name:  "f",
		Usage: " Run contracts from a specific migration. The number refers to the prefix of the migration file",
	}
	toFlag = cli.IntFlag{
		Name:  "t",
		Usage: "Run contracts to a specific migration. The number refers to the prefix of the migration file",
	}
	networkFlag = cli.StringFlag{
		Name:  "network",
		Usage: "Specify the network to use, saving artifacts specific to that network. Network name must exist in the configuration",
	}
	verboseRpcFlag = cli.BoolFlag{
		Name:  "verbose-rpc",
		Usage: "Log communication between bottle and the VNTChain client",
	}
	CompileCmd = cli.Command{
		Action:   compile,
		Name:     "compile",
		Usage:    "Compile contract source file",
		Category: "COMPILE COMMANDS",
		Description: `
bottle compile [-code <code_path>] [-output <output_dir_path>] [-include <include_dir_path>]

Compile contract source file
		`,
		Flags: []cli.Flag{
			contractCodeFlag,
			outputFlag,
			includeFlag,
		},
	}
	CompressCmd = cli.Command{
		Action:   compress,
		Name:     "compress",
		Usage:    "Compress wasm and abi file",
		Category: "COMPRESS COMMANDS",
		Description: `
bottle compress [-wasm <wasm_path>] [-abi <abi_path>] [-output <output_dir_path>]

Compress wasm and abi file
		`,
		Flags: []cli.Flag{
			wasmFlag,
			abiFlag,
			outputFlag,
		},
	}
	DecompressCmd = cli.Command{
		Action:   decompress,
		Name:     "decompress",
		Usage:    "Deompress file into wasm and abi file",
		Category: "DECOMPRESS COMMANDS",
		Description: `
bottle decompress [-file <compress_file_path>] [-output <output_dir_path>]

Deompress file into wasm and abi file
		`,
		Flags: []cli.Flag{
			compressFileFlag,
			outputFlag,
		},
	}
	HintCmd = cli.Command{
		Action:   hint,
		Name:     "hint",
		Usage:    "Contract hint",
		Category: "HINT COMMANDS",
		Description: `
bottle hint [-code <contract_path>]

Contract hint
		`,
		Flags: []cli.Flag{
			contractCodeFlag,
		},
	}
	InitCmd = cli.Command{
		Action:   initContract,
		Name:     "init",
		Usage:    "Initialize dapp project",
		Category: "INIT COMMANDS",
		Description: `
bottle init

Initialize dapp project
		`,
	}
	BuildCmd = cli.Command{
		Action:   buildContract,
		Name:     "build",
		Usage:    "Build contracts",
		Category: "BUILD COMMANDS",
		Description: `
bottle build

Build contracts in dapp directory
		`,
	}
	MigrateCmd = cli.Command{
		Action:   migrateContract,
		Name:     "migrate",
		Usage:    "Run migrations to deploy contracts",
		Category: "MIGRATE COMMANDS",
		Description: `
bottle migrate [-reset] [-f <from_number>] [-t <to_number>] [-network <network>] [verbose-rpc]

Run migrations to deploy contracts
		`,
		Flags: []cli.Flag{
			resetFlag,
			fromFlag,
			toFlag,
			networkFlag,
			verboseRpcFlag,
		},
	}
	ServerCmd = cli.Command{
		Action:   startServer,
		Name:     "server",
		Usage:    "Run rpc and ipc server",
		Category: "SERVER COMMANDS",
		Description: `
bottle server

Contract server
		`,
		Flags: []cli.Flag{},
	}
)

func compile(ctx *cli.Context) error {

	_codePath := ctx.String(contractCodeFlag.Name)
	_includeDir := ctx.String(includeFlag.Name)
	_outputDir := ctx.String(outputFlag.Name)
	return compileWith(ctx, _codePath, _includeDir, _outputDir)

}

func compileWith(ctx *cli.Context, _codePath, _includeDir, _outputDir string) error {
	start := time.Now()
	codePath = _codePath
	includeDir = _includeDir
	outputDir = _outputDir
	if err := hintWith(ctx, codePath); err != nil {
		return err
	}
	if codePath == "" {
		return cli.NewExitError("Error:No Contract Code", -1)
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
	if err := writeFile(path.Join(outputDir, abipath), res); err != nil {
		return err
	}
	pre := abigen.insertRegistryCode()
	// pre = abigen.insertMutableCode(pre)
	codeOutput := path.Join(outputDir, abires.Constructor.Name+"_precompile.c")
	if err := writeFile(codeOutput, pre); err != nil {
		return err
	}
	// fmt.Printf("Precompile code path: %s\n", codeOutput)
	wasmOutput := path.Join(outputDir, abires.Constructor.Name+".wasm")
	SetEnvPath()
	BuildWasm(codeOutput, wasmOutput)
	wasm, err := ioutil.ReadFile(wasmOutput)
	if err != nil {
		return err
	}
	if err := ValidWasm(wasm, abires); err != nil {
		return err
	}

	cpsPath := path.Join(outputDir, abires.Constructor.Name+".compress")
	cpsRes := utils.CompressWasmAndAbi(res, wasm, nil)
	if err := writeFile(cpsPath, cpsRes); err != nil {
		return err
	}
	hexPath := path.Join(outputDir, abires.Constructor.Name+".hex")
	hexString := "0x" + hex.EncodeToString(cpsRes)
	if err := writeFile(hexPath, []byte(hexString)); err != nil {
		return err
	}
	deployCodePath := path.Join(outputDir, abires.Constructor.Name+".js")
	if err := writeFile(deployCodePath, []byte(deployText(string(res), hexString))); err != nil {
		return err
	}

	contractJsonPath := path.Join(outputDir, abires.Constructor.Name+".json")
	abs, err := filepath.Abs(codePath)
	if err != nil {
		return err
	}
	contract := Contract{
		ContractName: abires.Constructor.Name,
		Abi:          abijson,
		Bytecode:     hexString,
		SourcePath:   abs,
		UpdatedAt:    time.Now().Format("2006-01-02T15:04:05.999Z"),
	}
	file, err := ioutil.ReadFile(contractJsonPath)
	if err != nil {
		contractJson, err := json.Marshal(contract)
		if err != nil {
			return err
		}
		err = writeFile(contractJsonPath, contractJson)
		if err != nil {
			return err
		}
	} else {
		var originContract Contract
		err := json.Unmarshal(file, &originContract)
		if err != nil {
			return err
		}
		originContract = merge(originContract, contract)
		contractJson, err := json.Marshal(originContract)
		if err != nil {
			return err
		}
		err = writeFile(contractJsonPath, contractJson)
		if err != nil {
			return err
		}
	}

	output := line.New(os.Stdout, "", "", line.WhiteColor)
	PrintfHeader(output, "Compile contract: "+abires.Constructor.Name+"\n")
	PrintfHeader(output, "Input file\n")
	PrintfBody(output, "Contract path:", codePath)
	PrintfHeader(output, "Output file\n")
	PrintfBody(output, "Abi path:", path.Join(outputDir, abipath))
	PrintfBody(output, "Precompile code path:", codeOutput)
	PrintfBody(output, "Wasm path:", wasmOutput)
	PrintfBody(output, "Compress data path:", cpsPath)
	PrintfBody(output, "Compress hex Data path:", hexPath)
	PrintfBody(output, "Deploy js path:", deployCodePath)
	PrintfBody(output, "Contract json path:", contractJsonPath)
	PrintfBody(output, fmt.Sprintf("Please use %s when you want to create a contract", abires.Constructor.Name+".compress"), "")
	PrintfHeader(output, "Compile finished. %s\n", time.Since(start).String())
	return nil
}

func compress(ctx *cli.Context) error {
	start := time.Now()
	wasmPath = ctx.String(wasmFlag.Name)
	abiPath = ctx.String(abiFlag.Name)
	outputDir = ctx.String(outputFlag.Name)
	if wasmPath == "" {
		return cli.NewExitError("Error:No Wasm File", -1)
	}
	if abiPath == "" {
		return cli.NewExitError("Error:No Abi File", -1)
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

	if err := ValidWasm(wasm, abires); err != nil {
		return err
	}

	cpsPath := path.Join(outputDir, abires.Constructor.Name+".compress")
	cpsRes := utils.CompressWasmAndAbi(abijson, wasm, nil)
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

	contractJsonPath := path.Join(outputDir, abires.Constructor.Name+".json")
	abs, err := filepath.Abs(codePath)
	if err != nil {
		return err
	}
	contract := Contract{
		ContractName: abires.Constructor.Name,
		Abi:          string(abijson),
		Bytecode:     hexString,
		SourcePath:   abs,
		UpdatedAt:    time.Now().Format("2019-04-24T07:49:32.705Z"),
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
	PrintfHeader(output, "Compress finished. %s\n", time.Since(start).String())
	PrintfBody(output, "Wasm path:", wasmPath)
	PrintfBody(output, "Abi path:", abiPath)
	PrintfHeader(output, "Output file\n")
	PrintfBody(output, "Compress data path:", cpsPath)
	PrintfBody(output, "Contract json path:", contractJsonPath)
	PrintfBody(output, fmt.Sprintf("Please use %s when you want to create a contract", abires.Constructor.Name+".compress"), "")
	return nil
}

func decompress(ctx *cli.Context) error {
	start := time.Now()
	compressPath = ctx.String(compressFileFlag.Name)
	outputDir = ctx.String(outputFlag.Name)
	if compressPath == "" {
		return cli.NewExitError("Error:No Compress File", -1)
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

	if err := ValidWasm(wasmcode.Code, abires); err != nil {
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
	PrintfHeader(output, "Decompress finished. %s\n", time.Since(start).String())
	PrintfHeader(output, "Input file\n")
	PrintfBody(output, "Compress file path:", compressPath)
	PrintfHeader(output, "Output file\n")
	PrintfBody(output, "wasm path:", wasmoutputPath)
	PrintfBody(output, "abi path:", abioutputPath)
	return nil
}

func hint(ctx *cli.Context) error {
	_codePath := ctx.String(contractCodeFlag.Name)
	return hintWith(ctx, _codePath)
}

func hintWith(ctx *cli.Context, _codePath string) error {
	codePath = _codePath
	if codePath == "" {
		return cli.NewExitError("Error:No Contract Code", -1)
	}
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
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	if !isEmpty(dir) {
		return cli.NewExitError("Warning: directory is't empty, can't create dapp project", -1)
	}
	init := newContractInit(dir)
	// init.copy(dir)
	err = init.download()
	return err
}

func buildContract(ctx *cli.Context) error {
	mig := NewMigrate()
	mig.FindCmd(ctx)
	if err := mig.Compile(ctx); err != nil {
		return err
	}
	return nil
}

func migrateContract(ctx *cli.Context) error {
	mig := NewMigrate()
	mig.FindCmd(ctx)
	if err := mig.Compile(ctx); err != nil {
		return err
	}
	jsPath := path.Join(os.TempDir(), "bottle.js")
	if err := writeFile(jsPath, JS); err != nil {
		return err
	}
	args := []string{jsPath, "migrate"}
	if ctx.IsSet("reset") {
		args = append(args, []string{"--reset"}...)
	}
	if ctx.IsSet("f") {
		args = append(args, []string{"--f", fmt.Sprintf("%d", ctx.Int("f"))}...)
	}
	if ctx.IsSet("t") {
		args = append(args, []string{"--t", fmt.Sprintf("%d", ctx.Int("t"))}...)
	}
	if ctx.IsSet("network") {
		args = append(args, []string{"--network", ctx.String("network")}...)
	}
	if ctx.IsSet("compile-all") {
		args = append(args, []string{"--compile-all"}...)
	}
	if ctx.IsSet("verbose-rpc") {
		args = append(args, []string{"--verbose-rpc"}...)
	}
	if err := mig.Run(ctx, args); err != nil {
		return err
	}
	return nil
}

func startServer(ctx *cli.Context) error {
	s := newServer("bottle.ipc")
	apis := []rpc.API{}
	if err := s.startIPC(apis); err != nil {
		fmt.Printf("err %s\n", err.Error())
		return err
	}
	s.wait()
	return nil
}
