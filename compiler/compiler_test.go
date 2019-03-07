package compiler

import (
    "testing"
    "fmt"
    "monkey_interpreter/ast"
    "monkey_interpreter/lexer"
    "monkey_interpreter/parser"
    "monkey_interpreter/object"
    "monkey_compiler/code"
)

type compilerTestCase struct {
    input string
    expectedConsts []interface{}
    expectedInsts []code.Insts
}

func TestIntegerArithmetic(t *testing.T) {
    tests := []compilerTestCase {
        {
            input: "1 + 2",
            expectedConsts: []interface{}{1, 2},
            expectedInsts: []code.Insts{
                code.Make(code.OpConst, 1),
                code.Make(code.OpConst, 2),
            },
        },
    }

    runCompilerTest(t, tests)
}

func runCompilerTest(t *testing.T, tests []compilerTestCase) {
    t.Helper()

    for _, test := range tests {
        program := parse(test.input)

        compiler := New()
        err := compiler.Compile(program)

        if err != nil {
            t.Fatalf("compiler error: %s", err)
        }

        bc := compiler.Bytecode()

        err = testInsts(test.expectedInsts, bc.Insts)
        if err != nil {
            t.Fatalf("testInsts failed: %s", err)
        }

        err = testConsts(t, test.expectedConsts, bc.Consts)
        if err != nil {
            t.Fatalf("testConsts failed: %s", err)
        }
    }
}

func parse(input string) *ast.Program {
    l := lexer.New(input)
    p := parser.New(l)
    return p.ParseProgram()
}

func testInsts(expected []code.Insts, actual code.Insts) error {
    concatted := concatInsts(expected)
    if len(actual) != len(concatted) {
        return fmt.Errorf("wrong instructions length.\nwant=%q\ngot=%q", concatted, actual)
    }

    for i, ins := range concatted {
        if actual[i] != ins {
            return fmt.Errorf("wrong instruction")
        }
    }

    return nil
}

func concatInsts(s []code.Insts) code.Insts {
    out := code.Insts{}

    for _, ins := range s {
        out = append(out, ins...)
    }

    return out
}

func testConsts(t *testing.T, expected []interface{}, actual []object.Object) error {
    if len(expected) != len(actual) {
        return fmt.Errorf("wrong number of constants")
    }

    for i, cons := range expected {
        switch cons := cons.(type) {
        case int:
            err := testIntegerObject(int64(cons), actual[i])
            if err != nil {
                return fmt.Errorf("incorrect value")
            }
        }
    }

    return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
    result, ok := actual.(*object.Integer)
    if !ok {
        return fmt.Errorf("object is not Integer")
    }

    if expected != result.Value {
        return fmt.Errorf("incorrect value")
    }

    return nil
}
