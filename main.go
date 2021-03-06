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
	"sort"

	"github.com/vntchain/bottle/core"

	cli "gopkg.in/urfave/cli.v1"
)

// get wasmCeptionFlag and vntIncludeFlag from go install -ldflags "-X env.wasmCeptionFlag=xxx -X env.vntIncludeFlag=xxx"
var wasmCeptionFlag string
var vntIncludeFlag string
var nodeFlag string
var gitCommit string

func main() {
	app := core.NewApp(gitCommit, "the bottle command line interface", vntIncludeFlag, wasmCeptionFlag, nodeFlag)
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2018-2019 The bottle Authors"
	app.Commands = []cli.Command{
		core.CompileCmd,
		core.CompressCmd,
		core.DecompressCmd,
		core.HintCmd,
		core.InitCmd,
		core.BuildCmd,
		core.MigrateCmd,
		// core.ServerCmd,
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
