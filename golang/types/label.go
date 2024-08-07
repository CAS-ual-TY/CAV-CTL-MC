package cav

type ILabel interface {
	GetKripkeStructure() IKripkeStructure
	MakeLabelFormula() IFormula
	String() string
}

type Label struct {
	kripkeStructure IKripkeStructure
	name            string
}

func (l *Label) GetKripkeStructure() IKripkeStructure {
	return l.kripkeStructure
}

func (l *Label) MakeLabelFormula() IFormula {
	return &LabelFormula{
		kripkeStructure: l.kripkeStructure,
		label:           l,
	}
}

func (l *Label) String() string {
	return l.name
}
