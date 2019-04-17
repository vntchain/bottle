package core

import (
	"reflect"
	"testing"
)

func Test_RemoveSpaceAndParen(t *testing.T) {
	tests := []struct {
		signatrue   string
		leftWanted  []string
		rightWanted []string
	}{
		{
			signatrue:   "void ()",
			leftWanted:  []string{"void"},
			rightWanted: []string{},
		},
		{
			signatrue:   "void (int)",
			leftWanted:  []string{"void"},
			rightWanted: []string{"int"},
		},
		{
			signatrue:   "void (int, int)",
			leftWanted:  []string{"void"},
			rightWanted: []string{"int", "int"},
		},
		{
			signatrue:   "void *()",
			leftWanted:  []string{"void*"},
			rightWanted: []string{},
		},
		{
			signatrue:   "void *(int *)",
			leftWanted:  []string{"void*"},
			rightWanted: []string{"int*"},
		},
		{
			signatrue:   "void *(int *, int *)",
			leftWanted:  []string{"void*"},
			rightWanted: []string{"int*", "int*"},
		},
	}
	for _, v := range tests {
		left, right := removeSpaceAndParen(v.signatrue)
		if !reflect.DeepEqual(left, v.leftWanted) {
			t.Fatalf("left  mismatch, got %v, want %v", left, v.leftWanted)
		}
		if !reflect.DeepEqual(right, v.rightWanted) {
			t.Fatalf("right  mismatch, got %v, want %v", right, v.rightWanted)
		}
	}
}

func Test_RemoveComment(t *testing.T) {
	tests := []struct {
		code   []byte
		wanted []byte
	}{
		{
			code: []byte(`
// a
// b
// c
// d
`),
			wanted: []byte{
				byte(10),
				byte(32), byte(32), byte(32), byte(32), byte(10),
				byte(32), byte(32), byte(32), byte(32), byte(10),
				byte(32), byte(32), byte(32), byte(32), byte(10),
				byte(32), byte(32), byte(32), byte(32), byte(10),
			},
		},
		{
			code: []byte(`
/*a
b
c
d*/
`),
			wanted: []byte{
				byte(10),
				byte(32), byte(32), byte(32), byte(10),
				byte(32), byte(10),
				byte(32), byte(10),
				byte(32), byte(32), byte(32), byte(10),
			},
		},
		{
			code: []byte(`
/*a
//b
//c
d*/
`),
			wanted: []byte{
				byte(10),
				byte(32), byte(32), byte(32), byte(10),
				byte(32), byte(32), byte(32), byte(10),
				byte(32), byte(32), byte(32), byte(10),
				byte(32), byte(32), byte(32), byte(10),
			},
		},
		{
			code: []byte(`
///*a
//b
//c
//d*/
`),
			wanted: []byte{
				byte(10),
				byte(32), byte(32), byte(32), byte(32), byte(32), byte(10),
				byte(32), byte(32), byte(32), byte(10),
				byte(32), byte(32), byte(32), byte(10),
				byte(32), byte(32), byte(32), byte(32), byte(32), byte(10),
			},
		},
		{
			code: []byte(`
aa/*a
//b
//c
d*/aa
`),
			wanted: []byte{
				byte(10),
				byte('a'), byte('a'), byte(32), byte(32), byte(32), byte(10),
				byte(32), byte(32), byte(32), byte(10),
				byte(32), byte(32), byte(32), byte(10),
				byte(32), byte(32), byte(32), byte('a'), byte('a'), byte(10),
			},
		},
		{
			code: []byte(`
/**a
a
*/
aa
/**a
a
*/
bb
`),
			wanted: []byte{
				byte(10),
				byte(32), byte(32), byte(32), byte(32), byte(10),
				byte(32), byte(10),
				byte(32), byte(32), byte(10),
				byte('a'), byte('a'), byte(10),
				byte(32), byte(32), byte(32), byte(32), byte(10),
				byte(32), byte(10),
				byte(32), byte(32), byte(10),
				byte('b'), byte('b'), byte(10),
			},
		},
	}

	for _, v := range tests {
		res := removeComment(string(v.code))
		if string(v.wanted) != string(res) {
			t.Fatalf("mismatch, got %v, want %v", res, v.wanted)
		}
	}

}
