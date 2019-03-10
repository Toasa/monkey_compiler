package vm

import (
    "fmt"
    "testing"
    "monkey_interpreter/ast"
    "monkey_interpreter/lexer"
    "monkey_interpreter/parser"
    "monkey_interpreter/object"
    "monkey_compiler/compiler"
)

type vmTestCase struct {
    input string
    expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
    tests := []vmTestCase {
        {"1", 1},
        {"2", 2},
        {"1 + 2", 3}, // incorrect
    }

    runVmTest(t, tests)
}

func parse(input string) *ast.Program {
    l := lexer.New(input)
    p := parser.New(l)
    return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
    result, ok := actual.(*object.Integer)
    if !ok {
        return fmt.Errorf("type assertion error")
    }

    if result.Value != expected {
        return fmt.Errorf("incorrect value")
    }

    return nil
}

func runVmTest(t *testing.T, tests []vmTestCase) {
    t.Helper()

    for _, test := range tests {
        program := parse(test.input)

        comp := compiler.New()
        err := comp.Compile(program)
        if err != nil {
            t.Fatalf("compiler err: %s", err)
        }

        vm := New(comp.Bytecode())
        err = vm.Run()
        if err != nil {
            t.Fatalf("vm err: %s", err)
        }

        stackElem := vm.StackTop()

        testExpectedObject(t, test.expected, stackElem)
    }
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
    t.Helper()

    switch expected := expected.(type) {
    case int:
        err := testIntegerObject(int64(expected), actual)
        if err != nil {
            t.Errorf("testIntegerObject failed")
        }
    }
}
