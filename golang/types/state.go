package cav

type IState interface {
	GetKripkeStructure() IKripkeStructure
	GetName() string
	AddLabel(label ILabel)
	HasLabel(label ILabel) bool
	GetLabels() ISet[ILabel]
	AddChildren(child ...IState)
	HasChild(child IState) bool
	GetChildren() ISet[IState]
	DetailString() string
	String() string
}

type State struct {
	kripkeStructure IKripkeStructure
	name            string
	labels          ISet[ILabel]
	children        ISet[IState]
	parents         ISet[IState]
}

func (s *State) GetKripkeStructure() IKripkeStructure {
	return s.kripkeStructure
}

func (s *State) GetName() string {
	return s.name
}

func (s *State) AddLabel(label ILabel) {
	s.labels.Add(label)
}

func (s *State) HasLabel(label ILabel) bool {
	return s.labels.Contains(label)
}

func (s *State) GetLabels() ISet[ILabel] {
	return s.labels
}

func (s *State) AddChildren(children ...IState) {
	for _, child := range children {
		s.children.Add(child)

		n2, ok := child.(*State)
		if ok {
			n2.AddParent(s)
		}
	}
}

func (s *State) HasChild(child IState) bool {
	return s.children.Contains(child)
}

func (s *State) GetChildren() ISet[IState] {
	return s.children
}

func (s *State) AddParent(parent IState) {
	s.parents.Add(parent)
}

func (s *State) DetailString() string {
	result := "State \"" + s.name + "\"\n"
	result += "  Labels:\n"
	s.labels.ForEach(func(label ILabel) {
		result += "    " + label.String() + "\n"
	})
	result += "  Children:\n"
	s.children.ForEach(func(child IState) {
		result += "    " + child.GetName() + "\n"
	})
	return result[:len(result)-1]

}

func (s *State) String() string {
	return s.name
}
