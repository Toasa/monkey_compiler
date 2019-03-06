package code

import (
    "encoding/binary"
    "fmt"
)

type Opcode byte

const (
    OpConst Opcode = iota
)

type Instructions []byte

type Definition struct {
    Name string
    OperandWidths []int
}

var defs = map[Opcode]*Definition {
    OpConst: {"OpConst", []int{2}},
}

func Lookup(op byte) (*Definition, error) {
    def, ok := defs[Opcode(op)]
    if !ok {
        return nil, fmt.Errorf("opcode %d undefined", op)
    }

    return def, nil
}

func Make(op Opcode, operands ...int) []byte {
    def, ok := defs[op]
    if !ok {
        return []byte{}
    }

    instLen := 1
    for _, w := range def.OperandWidths {
        instLen += w
    }

    // 命令バイト列の生成
    inst := make([]byte, instLen)
    inst[0] = byte(op)

    offset := 1
    for i, operand := range operands {
        w := def.OperandWidths[i]
        switch w {
        case 2:
            binary.BigEndian.PutUint16(inst[offset:], uint16(operand))
        }
        offset += w
    }

    return inst
}
