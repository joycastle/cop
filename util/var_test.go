package util

import (
	"testing"
)

type VarTest struct {
	Name string
}

func TestCase_IsVarPointor(t *testing.T) {
	if IsVarPointor(nil) != true {
		t.Fatal()
	}

	var a int
	if IsVarPointor(a) == false {
		t.Fatal()
	}

	if IsVarPointor(&a) != true {
		t.Fatal()
	}

	var vt VarTest
	if IsVarPointor(vt) == false {
		t.Fatal()
	}

	if IsVarPointor(&vt) != true {
		t.Fatal()
	}

	var mp map[string]string
	if IsVarPointor(mp) == false {
		t.Fatal()
	}

	var sl []int

	if IsVarPointor(sl) == false {
		t.Fatal()
	}
}
