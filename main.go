package main

import (
	"fmt"
	"os"
	"sort"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {

	// app.Action = gvnt
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2018-2019 The go-vnt Authors"
	app.Commands = []cli.Command{
		compileCmd,
		compressCmd,
		decompressCmd,
		hintCmd,
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
