package code

import (
    "encoding/binary"
    "fmt"
    "bytes"
)

type Opcode byte

const (
    OpConst Opcode = iota
    OpGetGlobal
    OpSetGlobal
    OpAdd
    OpSub
    OpMul
    OpDiv
    OpPop
    OpTrue
    OpFalse
    OpEq
    OpNE
    OpGT
    OpMinus
    OpBang
    OpJumpNotTruthy
    OpJump
    OpArray
    OpNull
)

type Instructions []byte

type Definition struct {
    Name string
    OperandWidths []int
}

var definitions = map[Opcode]*Definition {
    OpConst: {"OpConst", []int{2}},
    OpGetGlobal: {"OpGetGlobal", []int{2}},
    OpSetGlobal: {"OpSetGlobal", []int{2}},
    OpAdd: {"OpAdd", []int{}},
    OpSub: {"OpSub", []int{}},
    OpMul: {"OpMul", []int{}},
    OpDiv: {"OpDiv", []int{}},
    OpPop: {"OpPop", []int{}},
    OpTrue: {"OpTrue", []int{}},
    OpFalse: {"OpFalse", []int{}},
    OpEq: {"OpEq", []int{}},
    OpNE: {"OpNE", []int{}},
    OpGT: {"OpGT", []int{}},
    OpMinus: {"OpMinus", []int{}},
    OpBang: {"OpBang", []int{}},
    OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
    OpJump: {"OpJump", []int{2}},
    OpArray: {"OpArray", []int{2}},
    OpNull: {"OpNull", []int{}},
}

func Lookup(op byte) (*Definition, error) {
    def, ok := definitions[Opcode(op)]
    if !ok {
        return nil, fmt.Errorf("opcode %d undefined", op)
    }

    return def, nil
}

// operatorとoperandをbytecodeの命令列へencodeする
func Make(op Opcode, operands ...int) []byte {
    def, ok := definitions[op]
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

//　bytecodeの命令列からoperandのスライスへdecodeする
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
    ow := def.OperandWidths
    operands := make([]int, len(ow))

    offset := 0
    for i, width := range def.OperandWidths {
        switch width {
        case 2:
            operands[i] = int(ReadUint16(ins[offset:]))
        }
        offset += width
    }

    return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
    return binary.BigEndian.Uint16(ins)
}

func (ins Instructions)String() string {
    var out bytes.Buffer

    i := 0
    for i < len(ins) {
        def, err := Lookup(ins[i])
        if err != nil {
            fmt.Fprintf(&out, "ERROR: %s\n", err)
            continue
        }

        // 1byte分のoperatorは飛ばす
        operands, read_n := ReadOperands(def, ins[i+1:])

        fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

        i += 1 + read_n
    }

    return out.String()
}

func (ins Instructions)fmtInstruction(def *Definition, operands []int) string {
    operandCount := len(def.OperandWidths)

    if len(operands) != operandCount {
        return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
            len(operands), operandCount)
    }

    switch operandCount {
    case 0:
        return def.Name
    case 1:
        return fmt.Sprintf("%s %d", def.Name, operands[0])
    }

    return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}
