package cav

import "fmt"

type IFormula interface {
	Check() ISet[IState]
	GetKripkeStructure() IKripkeStructure
	String() string
}

type emptyFormula struct {
	kripkeStructure IKripkeStructure
}

type subFormula struct {
	kripkeStructure IKripkeStructure
	formula         IFormula
}

type biSubFormula struct {
	kripkeStructure IKripkeStructure
	formula1        IFormula
	formula2        IFormula
}

type equivalencyFormula struct {
	kripkeStructure    IKripkeStructure
	formula            IFormula
	equivalenceFormula IFormula
}

type biEquivalencyFormula struct {
	kripkeStructure    IKripkeStructure
	formula1           IFormula
	formula2           IFormula
	equivalenceFormula IFormula
}

type LabelFormula struct {
	kripkeStructure IKripkeStructure
	label           ILabel
}

func (f *LabelFormula) Check() ISet[IState] {
	result := MakeSet[IState]()
	f.kripkeStructure.GetStates().ForEach(func(state IState) {
		if state.HasLabel(f.label) {
			result.Add(state)
		}
	})
	return result
}

func (f *LabelFormula) String() string {
	return f.label.String()
}

func (f *LabelFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}

type TrueFormula emptyFormula

func (f *TrueFormula) Check() ISet[IState] {
	return f.kripkeStructure.GetStates()
}

func (f *TrueFormula) String() string {
	return "true"
}

func (f *TrueFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}

type FalseFormula emptyFormula

func (f *FalseFormula) Check() ISet[IState] {
	return MakeSet[IState]()
}

func (f *FalseFormula) String() string {
	return "false"
}

func (f *FalseFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}

type NotFormula subFormula

func (f *NotFormula) Check() ISet[IState] {
	return f.kripkeStructure.GetStates().Minus(f.formula.Check())
}

func (f *NotFormula) String() string {
	return fmt.Sprintf("(NOT %s)", f.formula.String())
}

func (f *NotFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}

type AndFormula biSubFormula // obviously this can also be done by using doubleEquivalencyFormula, containing De-Morgan

func (f *AndFormula) Check() ISet[IState] {
	//Alternative, by using: NOT[(NOT f1) OR (NOT f2)]:
	return f.formula1.Check().Intersect(f.formula2.Check())
}

func (f *AndFormula) String() string {
	return fmt.Sprintf("(%s AND %s)", f.formula1.String(), f.formula2.String())
}

func (f *AndFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}

type OrFormula biSubFormula

func (f *OrFormula) Check() ISet[IState] {
	return f.formula1.Check().Union(f.formula2.Check())
}

func (f *OrFormula) String() string {
	return fmt.Sprintf("(%s OR %s)", f.formula1.String(), f.formula2.String())
}

func (f *OrFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}

type EXFormula subFormula

func (f *EXFormula) Check() ISet[IState] {
	check := f.formula.Check()
	result := MakeSet[IState]()
	f.kripkeStructure.GetStates().ForEach(func(state IState) {
		check.ForEach(func(nextState IState) {
			if state.HasChild(nextState) {
				result.Add(state)
			}
		})
	})
	return result
}

func (f *EXFormula) String() string {
	return fmt.Sprintf("EX%s", f.formula.String())
}

func (f *EXFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}

type EGFormula subFormula

func (f *EGFormula) Check() ISet[IState] {
	p := f.formula.Check()

	var prevZ ISet[IState]
	var nextZ ISet[IState] = f.kripkeStructure.GetStates()

	for !nextZ.Equals(prevZ) {
		prevZ = nextZ

		exz := MakeSet[IState]()
		f.kripkeStructure.GetStates().ForEach(func(state IState) {
			prevZ.ForEach(func(nextState IState) {
				if state.HasChild(nextState) {
					exz.Add(state)
				}
			})
		})

		nextZ = p.Intersect(exz)
	}
	return prevZ
}

func (f *EGFormula) String() string {
	return fmt.Sprintf("EG%s", f.formula.String())
}

func (f *EGFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}

type EFFormula equivalencyFormula

func (f *EFFormula) Check() ISet[IState] {
	return f.equivalenceFormula.Check()
}

func (f *EFFormula) String() string {
	return fmt.Sprintf("EF%s", f.formula.String())
}

func (f *EFFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}

type EUFormula biSubFormula

func (f *EUFormula) Check() ISet[IState] {
	p := f.formula1.Check()
	q := f.formula2.Check()

	var prevZ ISet[IState]
	var nextZ ISet[IState] = MakeSet[IState]()

	for !nextZ.Equals(prevZ) {
		prevZ = nextZ

		exz := MakeSet[IState]()
		f.kripkeStructure.GetStates().ForEach(func(state IState) {
			prevZ.ForEach(func(nextState IState) {
				if state.HasChild(nextState) {
					exz.Add(state)
				}
			})
		})

		nextZ = q.Union(p.Intersect(exz))
	}
	return prevZ
}

func (f *EUFormula) String() string {
	return fmt.Sprintf("E[%s U %s]", f.formula1.String(), f.formula2.String())
}

func (f *EUFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}

type ERFormula biEquivalencyFormula

func (f *ERFormula) Check() ISet[IState] {
	return f.equivalenceFormula.Check()
}

func (f *ERFormula) String() string {
	return fmt.Sprintf("E[%s R %s]", f.formula1.String(), f.formula2.String())
}

func (f *ERFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}

type AXFormula equivalencyFormula

func (f *AXFormula) Check() ISet[IState] {
	return f.equivalenceFormula.Check()
}

func (f *AXFormula) String() string {
	return fmt.Sprintf("AX%s", f.formula.String())
}

func (f *AXFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}

type AGFormula equivalencyFormula

func (f *AGFormula) Check() ISet[IState] {
	return f.equivalenceFormula.Check()
}

func (f *AGFormula) String() string {
	return fmt.Sprintf("AG%s", f.formula.String())
}

func (f *AGFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}

type AFFormula equivalencyFormula

func (f *AFFormula) Check() ISet[IState] {
	return f.equivalenceFormula.Check()
}

func (f *AFFormula) String() string {
	return fmt.Sprintf("AF%s", f.formula.String())
}

func (f *AFFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}

type AUFormula biEquivalencyFormula

func (f *AUFormula) Check() ISet[IState] {
	return f.equivalenceFormula.Check()
}

func (f *AUFormula) String() string {
	return fmt.Sprintf("A[%s U %s]", f.formula1.String(), f.formula2.String())
}

func (f *AUFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}

type ARFormula biEquivalencyFormula

func (f *ARFormula) Check() ISet[IState] {
	return f.equivalenceFormula.Check()
}

func (f *ARFormula) String() string {
	return fmt.Sprintf("A[%s R %s]", f.formula1.String(), f.formula2.String())
}

func (f *ARFormula) GetKripkeStructure() IKripkeStructure {
	return f.kripkeStructure
}
