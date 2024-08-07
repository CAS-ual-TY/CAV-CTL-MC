package parser

import (
	"bufio"
	"cav/golang/types"
	"fmt"
	"io"
	"os"
	"strings"
)

type IFileParser interface {
	ParseFile(path string) (cav.IKripkeStructure, []cav.IFormula, error)
}

type FileParser struct {
	scanner   *bufio.Scanner
	line      string
	lineNr    int
	ks        cav.IKripkeStructure
	formulas  []cav.IFormula
	statesMap map[string]cav.IState
	labelsMap map[string]cav.ILabel
}

func (p *FileParser) nextLine() error {
	for p.scanner.Scan() {
		p.lineNr++
		line := p.scanner.Text()
		line = strings.SplitN(line, "//", 2)[0]
		line = strings.Replace(line, "\t", " ", -1)
		line = strings.Trim(line, " ")
		for strings.Contains(line, "  ") {
			line = strings.Replace(line, "  ", " ", -1)
		}
		if len(line) > 0 {
			p.line = line
			return nil
		}
	}
	return io.EOF
}

func (p *FileParser) errorf(s string, ss ...any) error {
	if len(ss) <= 0 {
		return fmt.Errorf("%d: %s", p.lineNr, s)
	}
	return fmt.Errorf("%d: %s", p.lineNr, fmt.Sprintf(s, ss...))
}

func (p *FileParser) parseFormula(s string) (cav.IFormula, error) {
	s = strings.Trim(s, " ")

	var i int
	runes := []rune(s)
	pCounter := 0 // counts '('
	bCounter := 0 // counts '['

	for i = 0; i < len(runes); i++ {
		if runes[i] == '(' {
			pCounter++
		} else if runes[i] == ')' {
			pCounter--
		} else if runes[i] == '[' {
			bCounter++
		} else if runes[i] == ']' {
			bCounter--
		}
		if pCounter != 0 || bCounter != 0 {
			continue
		} else if i < len(runes)-3 && s[i:i+3] == "AND" {
			break
		} else if i < len(runes)-2 && s[i:i+2] == "OR" {
			break
		}
	}

	if i < len(runes)-3 && s[i:i+3] == "AND" {
		left := strings.Trim(s[:i], " ")
		leftFormula, errLeft := p.parseFormula(left)
		if errLeft != nil {
			return nil, errLeft
		}
		right := strings.Trim(s[i+3:], " ")
		rightFormula, errRight := p.parseFormula(right)
		if errRight != nil {
			return nil, errRight
		}
		return p.ks.MakeAndFormula(leftFormula, rightFormula), nil
	} else if i < len(runes)-2 && s[i:i+2] == "OR" {
		left := strings.Trim(s[:i], " ")
		leftFormula, errLeft := p.parseFormula(left)
		if errLeft != nil {
			return nil, errLeft
		}
		right := strings.Trim(s[i+2:], " ")
		rightFormula, errRight := p.parseFormula(right)
		if errRight != nil {
			return nil, errRight
		}
		return p.ks.MakeOrFormula(leftFormula, rightFormula), nil
	}

	if s == "true" {
		return p.ks.MakeTrueFormula(), nil
	} else if s == "false" {
		return p.ks.MakeFalseFormula(), nil
	} else if strings.HasPrefix(s, "NOT") {
		formula, err := p.parseFormula(s[3:])
		if err == nil {
			return p.ks.MakeNotFormula(formula), nil
		} else {
			return nil, err
		}
	} else if strings.HasPrefix(s, "E") {
		s = strings.Trim(s[1:], " ")
		if strings.HasPrefix(s, "X") {
			s = strings.Trim(s[1:], " ")
			formula, err := p.parseFormula(s)
			if err == nil {
				return p.ks.MakeEXFormula(formula), nil
			} else {
				return nil, err
			}
		} else if strings.HasPrefix(s, "G") {
			s = strings.Trim(s[1:], " ")
			formula, err := p.parseFormula(s)
			if err == nil {
				return p.ks.MakeEGFormula(formula), nil
			} else {
				return nil, err
			}
		} else if strings.HasPrefix(s, "F") {
			s = strings.Trim(s[1:], " ")
			formula, err := p.parseFormula(s)
			if err == nil {
				return p.ks.MakeEFFormula(formula), nil
			} else {
				return nil, err
			}
		} else if (strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")")) || (strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]")) {
			s = strings.Trim(s[1:len(s)-1], " ")
			runes = []rune(s)
			pCounter = 0
			bCounter = 0
			for i = 0; i < len(runes); i++ {
				if runes[i] == '(' {
					pCounter++
				} else if runes[i] == ')' {
					pCounter--
				} else if runes[i] == '[' {
					bCounter++
				} else if runes[i] == ']' {
					bCounter--
				}
				if pCounter != 0 || bCounter != 0 {
					continue
				} else if runes[i] == 'U' || runes[i] == 'R' {
					break
				}
			}

			left := strings.Trim(s[:i], " ")
			leftFormula, errLeft := p.parseFormula(left)
			if errLeft != nil {
				return nil, errLeft
			}
			right := strings.Trim(s[i+1:], " ")
			rightFormula, errRight := p.parseFormula(right)
			if errRight != nil {
				return nil, errRight
			}
			if runes[i] == 'U' {
				return p.ks.MakeEUFormula(leftFormula, rightFormula), nil
			} else if runes[i] == 'R' {
				return p.ks.MakeERFormula(leftFormula, rightFormula), nil
			}
		}
	} else if strings.HasPrefix(s, "A") {
		s = strings.Trim(s[1:], " ")
		if strings.HasPrefix(s, "X") {
			s = strings.Trim(s[1:], " ")
			formula, err := p.parseFormula(s)
			if err == nil {
				return p.ks.MakeAXFormula(formula), nil
			} else {
				return nil, err
			}
		} else if strings.HasPrefix(s, "G") {
			s = strings.Trim(s[1:], " ")
			formula, err := p.parseFormula(s)
			if err == nil {
				return p.ks.MakeAGFormula(formula), nil
			} else {
				return nil, err
			}
		} else if strings.HasPrefix(s, "F") {
			s = strings.Trim(s[1:], " ")
			formula, err := p.parseFormula(s)
			if err == nil {
				return p.ks.MakeAFFormula(formula), nil
			} else {
				return nil, err
			}
		} else if (strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")")) || (strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]")) {
			s = strings.Trim(s[1:len(s)-1], " ")
			runes = []rune(s)
			pCounter = 0
			bCounter = 0
			for i = 0; i < len(runes); i++ {
				if runes[i] == '(' {
					pCounter++
				} else if runes[i] == ')' {
					pCounter--
				} else if runes[i] == '[' {
					bCounter++
				} else if runes[i] == ']' {
					bCounter--
				}
				if pCounter != 0 || bCounter != 0 {
					continue
				} else if runes[i] == 'U' || runes[i] == 'R' {
					break
				}
			}

			left := strings.Trim(s[:i], " ")
			leftFormula, errLeft := p.parseFormula(left)
			if errLeft != nil {
				return nil, errLeft
			}
			right := strings.Trim(s[i+1:], " ")
			rightFormula, errRight := p.parseFormula(right)
			if errRight != nil {
				return nil, errRight
			}
			if runes[i] == 'U' {
				return p.ks.MakeAUFormula(leftFormula, rightFormula), nil
			} else if runes[i] == 'R' {
				return p.ks.MakeARFormula(leftFormula, rightFormula), nil
			}
		}
	} else if (strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")")) || (strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]")) {
		s = strings.Trim(s[1:len(s)-1], " ")
		return p.parseFormula(s)
	}

	s = strings.Split(s, " ")[0]
	label, ok := p.labelsMap[s]
	if !ok {
		return nil, p.errorf("unknown label in formula: %s", s)
	}
	return label.MakeLabelFormula(), nil
}

func (p *FileParser) parseEverything() error {
	if err := p.nextLine(); err != nil {
		return err
	}

	if p.line != "states" {
		return p.errorf("Expected \"states\", but got %s", p.line)
	}
	p.ks = cav.MakeKripkeStructure()
	p.statesMap = map[string]cav.IState{}
	p.labelsMap = map[string]cav.ILabel{}

	// -------------------------------------------
	// states
	// -------------------------------------------

	if err := p.nextLine(); err != nil {
		return err
	}

	for p.line != "transitions" {
		stateName := p.line

		if _, ok := p.statesMap[stateName]; ok {
			return p.errorf("duplicate state: %s", stateName)
		}

		p.statesMap[stateName] = p.ks.NewState(stateName)

		if err := p.nextLine(); err != nil {
			if err == io.EOF {
				return p.errorf("expected \"transitions\", but could not find it")
			}
			return err
		}
	}

	// -------------------------------------------
	// transitions
	// -------------------------------------------

	if err := p.nextLine(); err != nil {
		return err
	}

	for p.line != "labels" {
		parts := strings.Split(p.line, " ")

		var prevState cav.IState
		var right bool

		for i, part := range parts {
			if i%2 == 1 {
				if part == "->" {
					right = true
				} else if part == "<-" {
					right = false
				} else {
					return p.errorf("invalid transition, expected \"->\" or \"<-\", but got: %s", part)
				}
				continue
			} else {
				if len(part) <= 0 {
					return p.errorf("missing state for transition")
				}

				nextState, ok := p.statesMap[part]

				if !ok {
					return p.errorf("unknown state for transition: %s", part)
				}

				if i > 0 {
					if right {
						prevState.AddChildren(nextState)
					} else {
						nextState.AddChildren(prevState)
					}
				}

				prevState = nextState
			}
		}

		if err := p.nextLine(); err != nil {
			if err == io.EOF {
				return p.errorf("expected \"labels\", but could not find it")
			}
			return err
		}
	}

	// -------------------------------------------
	// labels
	// -------------------------------------------

	if err := p.nextLine(); err != nil {
		return err
	}

	for p.line != "formulas" {
		parts := strings.Split(p.line, ":")

		if len(parts) != 2 {
			return p.errorf("invalid label definition, expected exactly one ':' but got: %s", p.line)
		}

		labelName := parts[0]

		if _, ok := p.labelsMap[labelName]; ok {
			return p.errorf("duplicate label: %s", labelName)
		}

		label := p.ks.NewLabel(labelName)
		p.labelsMap[labelName] = label

		stateNames := strings.Split(parts[1], ",")

		for _, stateName := range stateNames {
			stateName = strings.Trim(stateName, " ")
			state, ok := p.statesMap[stateName]
			if !ok {
				return p.errorf("unknown state for label %s: %s", labelName, stateName)
			}
			state.AddLabel(label)
		}

		if err := p.nextLine(); err != nil {
			if err == io.EOF {
				return p.errorf("expected \"formulas\", but could not find it")
			}
			return err
		}
	}

	// -------------------------------------------
	// formulas
	// -------------------------------------------

	if err := p.nextLine(); err != nil {
		return err
	}

	p.formulas = make([]cav.IFormula, 0)
	for {
		formula, err := p.parseFormula(p.line)
		if err != nil {
			return err
		}
		p.formulas = append(p.formulas, formula)

		if err := p.nextLine(); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

func (p *FileParser) ParseFile(path string) (cav.IKripkeStructure, []cav.IFormula, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("file %s does not exist: %s", path, err.Error())
		} else {
			return nil, nil, err
		}
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	p.scanner = bufio.NewScanner(file)
	p.line = ""
	p.lineNr = 0

	err = p.parseEverything()

	if err == io.EOF {
		err = p.errorf("unexpected end of file")
	}

	return p.ks, p.formulas, err
}

var PARSER IFileParser = &FileParser{}
