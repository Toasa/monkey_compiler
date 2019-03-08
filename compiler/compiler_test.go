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
    expectedConstants []interface{}
    expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
    tests := []compilerTestCase {
        {
            input: "1 + 2",
            expectedConstants: []interface{}{1, 2},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpConst, 0),
                code.Make(code.OpConst, 1),
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

        err = testInstructions(test.expectedInstructions, bc.Instructions)
        if err != nil {
            t.Fatalf("testInstructions failed: %s", err)
        }

        err = testConstants(t, test.expectedConstants, bc.Constants)
        if err != nil {
            t.Fatalf("testConstants failed: %s", err)
        }
    }
}

func parse(input string) *ast.Program {
    l := lexer.New(input)
    p := parser.New(l)
    return p.ParseProgram()
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
    concatted := concatInstructions(expected)
    if len(actual) != len(concatted) {
        return fmt.Errorf("wrong instructions length.\nwant=%q\ngot=%q", concatted, actual)
    }

    for i, ins := range concatted {
        if actual[i] != ins {
            return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot=%q", i, concatted, actual)
        }
    }

    return nil
}

func concatInstructions(s []code.Instructions) code.Instructions {
    out := code.Instructions{}

    for _, ins := range s {
        out = append(out, ins...)
    }

    return out
}

func testConstants(t *testing.T, expected []interface{}, actual []object.Object) error {
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
