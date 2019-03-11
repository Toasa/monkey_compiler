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
        {"1 + 2", 3},
        {"1 - 2", -1},
        {"1 * 2", 2},
        {"4 / 2", 2},
        {"50 / 2 * 2 + 10 - 5", 55},
        {"5 * (2 + 10)", 60},
        {"(1 + 2) * (3 + 4)", 21},
        {"1 + 2 * 3 + 4", 11},
    }

    runVmTest(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
    tests := []vmTestCase {
        {"true", true},
        {"false", false},
        {"1 < 2", true},
        {"1 > 2", false},
        {"1 < 1", false},
        {"1 > 1", false},
        {"1 == 1", true},
        {"1 != 1", false},
        {"1 == 2", false},
        {"1 != 2", true},
        {"true == true", true},
        {"false == false", true},
        {"true == false", false},
        {"true != false", true},
        {"false != true", true},
        {"(1 < 2) == true", true},
        {"(1 < 2) == false", false},
        {"(1 > 2) == true", false},
        {"(1 > 2) == false", true},
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
        return fmt.Errorf("expected %d, but got %d", expected, result.Value)
    }

    return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
    b, ok := actual.(*object.Boolean)
    if !ok {
        return fmt.Errorf("type assertion error")
    }

    if b.Value != expected {
        return fmt.Errorf("expected value faild")
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

        stackElem := vm.LastPoppedStackElem()

        testExpectedObject(t, test.expected, stackElem)
    }
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
    t.Helper()

    switch expected := expected.(type) {
    case int:
        err := testIntegerObject(int64(expected), actual)
        if err != nil {
            t.Errorf("testIntegerObject failed, %s", err)
        }
    case bool:
        err := testBooleanObject(expected, actual)
        if err != nil {
            t.Errorf("testBooleanObject faild, %s", err)
        }
    }
}
