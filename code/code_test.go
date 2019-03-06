package code

import "testing"

func TestMake(t *testing.T) {
    tests := []struct {
        op Opcode
        operands []int
        expected []byte
    }{
        {OpConst, []int{65534}, []byte{byte(OpConst), 255, 254}},
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
