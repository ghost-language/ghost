package code

import (
	"encoding/binary"
	"fmt"
)

type Opcode byte
type Definition struct {
	Name          string
	OperandWidths []int
}

const (
	OpConstant Opcode = iota
)

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

func Lookup(op byte) (*Definition, error) {
	definition, ok := definitions[Opcode(op)]

	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return definition, nil
}

func Make(op Opcode, operands ...int) []byte {
	definition, ok := definitions[op]

	if !ok {
		return []byte{}
	}

	instructionLength := 1

	for _, w := range definition.OperandWidths {
		instructionLength += w
	}

	instruction := make([]byte, instructionLength)
	instruction[0] = byte(op)

	offset := 1

	for i, operand := range operands {
		width := definition.OperandWidths[i]

		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		}

		offset += width
	}

	return instruction
}

func ReadOperands(definition *Definition, instructions Instructions) ([]int, int) {
	operands := make([]int, len(definition.OperandWidths))
	offset := 0

	for i, width := range definition.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(instructions[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUint16(instructions Instructions) uint16 {
	return binary.BigEndian.Uint16(instructions)
}
