package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/briandowns/spinner"
	"github.com/cavaliercoder/grab"
)

var (
	tempPath     = path.Join("/Users/weisaizhang/Documents/go/src/github.com/vntchain/bottle/template")
	tempFileName = "bottle-contract-template-master"
	tempUrl      = "https://github.com/ooozws/bottle-contract-template/archive/master.zip"
)

type contractInit struct {
	dst string
}

func newContractInit(dst string) contractInit {

	return contractInit{
		dst: dst,
	}
}

func (init contractInit) download() error {
	f, err := init.createTempDir()
	if err != nil {
		return err
	}
	client := grab.NewClient()
	req, err := grab.NewRequest(f, tempUrl)
	if err != nil {
		return err
	}

	fmt.Print("✔ Preparing to download...\n")
	resp := client.Do(req)

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	s.Suffix = "  :Downloading"
Loop:
	for {
		select {
		case <-resp.Done:
			// download is complete
			time.Sleep(1 * time.Second)
			s.Stop()
			break Loop
		}
	}

	if err := resp.Err(); err != nil {
		return err
	}
	fmt.Print("✔ Download completed...\n")
	if err := unpackZip(init.dst, path.Join(f, tempFileName+".zip"), 1); err != nil {
		return err
	}
	fmt.Print("✔ Unzip file...\n")
	if err := os.RemoveAll(f); err != nil {
		return err
	}
	fmt.Print("✔ Cleaning up temporary files\n")
	fmt.Print("✔ Init successful\n")
	return nil
}

func (init contractInit) createTempDir() (string, error) {
	f, err := ioutil.TempDir(".", ".tmp")
	if err != nil {
		return "", err
	}
	return f, nil
}
