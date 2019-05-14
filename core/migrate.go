package core

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

type Migrate struct {
}

func NewMigrate() Migrate {
	return Migrate{}
}

func (m Migrate) Compile(ctx *cli.Context) error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	dirs, err := ioutil.ReadDir(path.Join(dir, "contracts"))
	if err != nil {
		return err
	}
	outputDir := path.Join(dir, "build", "contracts")
	for _, v := range dirs {
		if !v.IsDir() {
			res := path.Ext(v.Name())
			if strings.Compare(res, ".c") == 0 {
				err := compileWith(ctx, path.Join(dir, "contracts", v.Name()), "", outputDir)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m Migrate) Run(ctx *cli.Context, cmdArgs []string) error {
	cmdPath := path.Join(nodePathFlag, "bin", "node")
	cmd := exec.Command(cmdPath, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (m Migrate) FindCmd(ctx *cli.Context) {

}
