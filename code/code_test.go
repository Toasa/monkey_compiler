package code

import "testing"

func TestMake(t *testing.T) {
    tests := []struct {
        op Opcode
        operands []int
        expected []byte
    }{
        {OpConst, []int{65534}, []byte{byte(OpConst), 255, 254}},
        {OpAdd, []int{}, []byte{byte(OpAdd)}},
    }

    for _, test := range tests {
        inst := Make(test.op, test.operands...)

        if len(inst) != len(test.expected) {
            t.Errorf("Instruction has wrong length")
        }

        for i, b := range test.expected {
            if inst[i] != test.expected[i] {
                t.Errorf("wrong byte at pos %d. want=%d, but got=%d", i, b, inst[i])
            }
        }
    }
}

func TestInstsString(t *testing.T) {
    insts := []Instructions {
        Make(OpAdd),
        Make(OpConst, 2),
        Make(OpConst, 65534),
    }

    expected := `0000 OpAdd
0001 OpConst 2
0004 OpConst 65534
`

    concatted := Instructions{}
    for _, ins := range insts {
        concatted = append(concatted, ins...)
    }

    if concatted.String() != expected {
        t.Errorf("instructions wrongly formatted.\nwant=%q\ngot=%q", expected, concatted.String())
    }
}

func TestReadOperands(t *testing.T) {
    tests := []struct {
        op Opcode
        operands []int
        bytesRead int
    }{
        {OpConst, []int{65535}, 2},
    }

    for _, test := range tests {
        inst := Make(test.op, test.operands...)

        def, err := Lookup(byte(test.op))
        if err != nil {
            t.Fatalf("defnition not found")
        }

        operandsRead, n := ReadOperands(def, inst[1:])
        if n != test.bytesRead {
            t.Fatalf("n wrong")
        }

        for i, want := range test.operands {
            if operandsRead[i] != want {
                t.Errorf("operand wrong")
            }
        }
    }
}
