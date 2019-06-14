package tests

import (
	"path"
	"testing"
)

func TestStruct(t *testing.T) {
	tests := []struct {
		contract string
		wanted   string
	}{
		{
			contract: "test_struct.c",
		},
	}
	for _, v := range tests {
		execBottle(path.Join(contractsPath, v.contract))
	}

}
