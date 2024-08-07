package test

import (
	cav2 "cav/golang/parser"
	"cav/golang/types"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testFormula(t *testing.T, fla cav.IFormula, expected cav.ISet[cav.IState]) {
	flaResult := fla.Check()
	t.Log(fla.String() + ": " + flaResult.String())
	if !flaResult.Equals(expected) {
		t.Errorf("%s is invalid: Expected %s but got %s", fla.String(), expected.String(), flaResult.String())
	}
}

func TestKripkeStructure(t *testing.T) {
	/*
	 * p      p      p      r
	 * s1 --> s2 <-- s3 <-- s4
	 *       ^  \            ^
	 *      /     \          |
	 *    /         \        |
	 *  /            v       |
	 * s5 <-- s6 <-- s7     s8
	 * q      p      p      p
	 *
	 */

	ks := cav.MakeKripkeStructure()

	p := ks.NewLabel("p")
	q := ks.NewLabel("q")
	r := ks.NewLabel("r")

	s1 := ks.NewState("s1", p)
	s2 := ks.NewState("s2", p)
	s3 := ks.NewState("s3", p)
	s4 := ks.NewState("s4", r)
	s5 := ks.NewState("s5", q)
	s6 := ks.NewState("s6", p)
	s7 := ks.NewState("s7", p)
	s8 := ks.NewState("s8", p)

	s1.AddChildren(s2)
	s2.AddChildren(s7)
	s3.AddChildren(s2)
	s4.AddChildren(s3)
	s5.AddChildren(s2)
	s6.AddChildren(s5)
	s7.AddChildren(s6)
	s8.AddChildren(s4)

	if !s1.HasLabel(p) || !s4.HasLabel(r) || s4.HasLabel(p) || !ks.Validate() {
		t.Errorf("Kripke structure is invalid")
	}
	t.Log(ks.String())
	t.Log(ks.DetailString())

	fla_p := p.MakeLabelFormula()
	fla_q := q.MakeLabelFormula()
	fla_r := r.MakeLabelFormula()

	testFormula(t, fla_p, cav.MakeSetOf(s1, s2, s3, s6, s7, s8))
	testFormula(t, fla_q, cav.MakeSetOf(s5))
	testFormula(t, fla_r, cav.MakeSetOf(s4))

	testFormula(t, ks.MakeTrueFormula(), ks.GetStates())
	testFormula(t, ks.MakeFalseFormula(), cav.MakeSet[cav.IState]())
	testFormula(t, ks.MakeNotFormula(fla_p), cav.MakeSetOf(s4, s5))
	testFormula(t, ks.MakeAndFormula(fla_p, fla_q), cav.MakeSet[cav.IState]())
	testFormula(t, ks.MakeOrFormula(fla_p, fla_q), cav.MakeSetOf(s1, s2, s3, s5, s6, s7, s8))
	testFormula(t, ks.MakeEXFormula(fla_p), cav.MakeSetOf(s1, s2, s3, s4, s5, s7))
	testFormula(t, ks.MakeEGFormula(fla_p), cav.MakeSet[cav.IState]())
	testFormula(t, ks.MakeEFFormula(fla_p), ks.GetStates())
	testFormula(t, ks.MakeEUFormula(fla_p, fla_q), cav.MakeSetOf(s1, s2, s3, s5, s6, s7))
	testFormula(t, ks.MakeERFormula(fla_p, fla_q), cav.MakeSet[cav.IState]())
	testFormula(t, ks.MakeAXFormula(fla_p), cav.MakeSetOf(s1, s2, s3, s4, s5, s7))
	testFormula(t, ks.MakeAGFormula(fla_p), cav.MakeSet[cav.IState]())
	testFormula(t, ks.MakeAFFormula(fla_p), ks.GetStates())
	testFormula(t, ks.MakeAUFormula(fla_p, fla_q), cav.MakeSetOf(s1, s2, s3, s5, s6, s7))
	testFormula(t, ks.MakeARFormula(fla_p, fla_q), cav.MakeSet[cav.IState]())

	testFormula(t, ks.MakeEXFormula(fla_p), cav.MakeSetOf(s1, s2, s3, s4, s5, s7))
	testFormula(t, ks.MakeEGFormula(fla_p), cav.MakeSet[cav.IState]())
	testFormula(t, ks.MakeEUFormula(fla_p, fla_q), cav.MakeSetOf(s1, s2, s3, s5, s6, s7))

	fla1 := ks.MakeEXFormula(fla_p)
	fla2 := ks.MakeEGFormula(fla_p)
	fla3 := ks.MakeEUFormula(fla_p, fla_q)

	fmt.Println(fla1.Check())
	fmt.Println(fla2.Check())
	fmt.Println(fla3.Check())
}

func TestParser(t *testing.T) {
	// Testing package always sets the working directory to the package directory
	// so we go up to project dir
	wd, _ := os.Getwd()
	path := ""
	for !strings.HasSuffix(wd, "CAV") {
		wd = filepath.Dir(wd)
		path = path + "../"
	}

	t.Log("Working directory: " + wd)

	ks, flas, err := cav2.PARSER.ParseFile(fmt.Sprintf("%s/kripkestructure_exam.txt", wd))
	if err != nil {
		t.Errorf("Failed to parse file:")
		t.Error(err)
		return
	}

	t.Log(ks.String())
	t.Log(ks.DetailString())
	for _, fla := range flas {
		t.Log(fla.String() + ":")
		t.Log(fla.Check().String())
	}
}
