package code

import (
	"bytes"
	"fmt"
)

type Instructions []byte

func (instructions Instructions) String() string {
	var out bytes.Buffer

	i := 0

	for i < len(instructions) {
		definition, err := Lookup(instructions[i])

		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, bytesRead := ReadOperands(definition, instructions[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, instructions.fmtInstruction(definition, operands))

		i += bytesRead
	}

	return out.String()
}

func (instructions Instructions) fmtInstruction(definition *Definition, operands []int) string {
	operandCount := len(definition.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand length %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 1:
		return fmt.Sprintf("%s %d", definition.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", definition.Name)
}
