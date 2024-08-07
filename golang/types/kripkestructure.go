package cav

import (
	"strings"
)

type IKripkeStructure interface {
	NewLabel(name string) ILabel
	NewState(name string, label ...ILabel) IState
	GetStates() ISet[IState]
	Validate() bool
	MakeTrueFormula() IFormula
	MakeFalseFormula() IFormula
	MakeNotFormula(formula IFormula) IFormula
	MakeAndFormula(formula1 IFormula, formula2 IFormula) IFormula
	MakeOrFormula(formula1 IFormula, formula2 IFormula) IFormula
	MakeEXFormula(formula IFormula) IFormula
	MakeEGFormula(formula IFormula) IFormula
	MakeEFFormula(formula IFormula) IFormula
	MakeEUFormula(formula1 IFormula, formula2 IFormula) IFormula
	MakeERFormula(formula1 IFormula, formula2 IFormula) IFormula
	MakeAXFormula(formula IFormula) IFormula
	MakeAGFormula(formula IFormula) IFormula
	MakeAFFormula(formula IFormula) IFormula
	MakeAUFormula(formula1 IFormula, formula2 IFormula) IFormula
	MakeARFormula(formula1 IFormula, formula2 IFormula) IFormula
	DetailString() string
	String() string
}

type KripkeStructure struct {
	labels        ISet[ILabel]
	states        ISet[IState]
	initialStates ISet[IState]
}

func (ks *KripkeStructure) NewLabel(name string) ILabel {
	label := &Label{
		kripkeStructure: ks,
		name:            name,
	}
	ks.labels.Add(label)
	return label
}

func (ks *KripkeStructure) NewState(name string, label ...ILabel) IState {
	state := &State{
		kripkeStructure: ks,
		name:            name,
		labels:          MakeSet[ILabel](),
		children:        MakeSet[IState](),
		parents:         MakeSet[IState](),
	}
	for _, l := range label {
		state.AddLabel(l)
	}
	ks.states.Add(state)
	return state
}

func (ks *KripkeStructure) GetStates() ISet[IState] {
	return ks.states
}

func (ks *KripkeStructure) Validate() bool {
	result := true
	ks.states.ForEach(func(state IState) {
		state.GetLabels().ForEach(func(label ILabel) {
			if !ks.labels.Contains(label) {
				result = false
			}
		})
		state.GetChildren().ForEach(func(child IState) {
			if !ks.states.Contains(child) {
				result = false
			}
		})
	})
	return result
}

func (ks *KripkeStructure) MakeLabelFormula(label ILabel) IFormula {
	return &LabelFormula{ks, label}
}

func (ks *KripkeStructure) MakeTrueFormula() IFormula {
	return &TrueFormula{ks}
}

func (ks *KripkeStructure) MakeFalseFormula() IFormula {
	return &FalseFormula{ks}
}

func (ks *KripkeStructure) MakeNotFormula(formula IFormula) IFormula {
	return &NotFormula{ks, formula}
}

func (ks *KripkeStructure) MakeAndFormula(formula1 IFormula, formula2 IFormula) IFormula {
	return &AndFormula{ks, formula1, formula2}
}

func (ks *KripkeStructure) MakeOrFormula(formula1 IFormula, formula2 IFormula) IFormula {
	return &OrFormula{ks, formula1, formula2}
}

func (ks *KripkeStructure) MakeEXFormula(formula IFormula) IFormula {
	return &EXFormula{ks, formula}
}

func (ks *KripkeStructure) MakeEGFormula(formula IFormula) IFormula {
	return &EGFormula{ks, formula}
}

func (ks *KripkeStructure) MakeEFFormula(formula IFormula) IFormula {
	return &EFFormula{ks, formula, ks.MakeEUFormula(ks.MakeTrueFormula(), formula)}
}

func (ks *KripkeStructure) MakeEUFormula(formula1 IFormula, formula2 IFormula) IFormula {
	return &EUFormula{ks, formula1, formula2}
}

func (ks *KripkeStructure) MakeERFormula(formula1 IFormula, formula2 IFormula) IFormula {
	return &ERFormula{ks, formula1, formula2, ks.MakeNotFormula(ks.MakeAUFormula(ks.MakeNotFormula(formula1), ks.MakeNotFormula(formula2)))}
}

func (ks *KripkeStructure) MakeAXFormula(formula IFormula) IFormula {
	return &AXFormula{ks, formula, ks.MakeNotFormula(ks.MakeEXFormula(ks.MakeNotFormula(formula)))}
}

func (ks *KripkeStructure) MakeAGFormula(formula IFormula) IFormula {
	return &AGFormula{ks, formula, ks.MakeNotFormula(ks.MakeEFFormula(ks.MakeNotFormula(formula)))}
}

func (ks *KripkeStructure) MakeAFFormula(formula IFormula) IFormula {
	return &AFFormula{ks, formula, ks.MakeNotFormula(ks.MakeEGFormula(ks.MakeNotFormula(formula)))}
}

func (ks *KripkeStructure) MakeAUFormula(formula1 IFormula, formula2 IFormula) IFormula {
	return &AUFormula{ks, formula1, formula2, ks.MakeAndFormula(ks.MakeNotFormula(ks.MakeEUFormula(ks.MakeNotFormula(formula2), ks.MakeAndFormula(ks.MakeNotFormula(formula1), ks.MakeNotFormula(formula2)))), ks.MakeNotFormula(ks.MakeEGFormula(ks.MakeNotFormula(formula2))))}
}

func (ks *KripkeStructure) MakeARFormula(formula1 IFormula, formula2 IFormula) IFormula {
	return &ARFormula{ks, formula1, formula2, ks.MakeNotFormula(ks.MakeEUFormula(ks.MakeNotFormula(formula1), ks.MakeNotFormula(formula2)))}
}

func (ks *KripkeStructure) DetailString() string {
	result := "KripkeStructure:\n"
	result += "  Labels:\n"
	ks.labels.ForEach(func(label ILabel) {
		result += "    " + label.String() + "\n"
	})
	result += "  States:\n"
	ks.states.ForEach(func(state IState) {
		ss := strings.Split(state.DetailString(), "\n")
		for _, s := range ss {
			result += "    " + s + "\n"
		}
	})
	return result[:len(result)-1]
}

func (ks *KripkeStructure) String() string {
	return "States: " + ks.states.String() + ", Labels: " + ks.labels.String()
}

func MakeKripkeStructure() IKripkeStructure {
	return &KripkeStructure{
		labels:        MakeSet[ILabel](),
		states:        MakeSet[IState](),
		initialStates: MakeSet[IState](),
	}
}
