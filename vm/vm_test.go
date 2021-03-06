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
        {"-5", -5},
        {"-5 + 10 - 5", 0},
        {"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
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
        {"!true", false},
        {"!false", true},
        {"!5", false},
        {"!!true", true},
        {"!!false", false},
        {"!!5", true},
        {"!(if (false) { 5; })", true},
    }

    runVmTest(t, tests)
}

func TestConditionals(t *testing.T) {
    tests := []vmTestCase {
        {"if (true) { 10 }", 10},
        {"if (true) { 10 } else { 20 }", 10},
        {"if (false) { 10 } else { 20 }", 20},
        {"if (1) { 10 }", 10},
        {"if (1 < 2) { 10 }", 10},
        {"if (1 < 2) { 10 } else { 20 }", 10},
        {"if (1 > 2) { 10 } else { 20 }", 20},
        {"if (1 > 2) { 10 }", Null},
        {"if (false) { 10 }", Null},
        {"if (if (false) { 10 }) { 10 } else { 20 }", 20},
    }

    runVmTest(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
    tests := []vmTestCase {
        {"let one = 1; one", 1},
        {"let one = 1; let two = 2; one + two", 3},
        {"let one = 1; let two = one + one; one + two", 3},
    }

    runVmTest(t, tests)
}

func TestStringExpressions(t *testing.T) {
    tests := []vmTestCase {
        {`"monkey"`, "monkey"},
        {`"mon" + "key"`, "monkey"},
        {`"mon" + "key" + "ship"`, "monkeyship"},
    }

    runVmTest(t, tests)
}

func TestArrayLiterals(t *testing.T) {
    tests := []vmTestCase {
        {"[]", []int{}},
        {"[1, 2, 3]", []int{1, 2, 3}},
        {"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
    }

    runVmTest(t, tests)
}

func TestHashLiterals(t *testing.T) {
    tests := []vmTestCase {
        {
            "{}",
            map[object.HashKey]int64{},
        },
        {
            "{1: 2, 2: 3}",
            map[object.HashKey]int64{
                (&object.Integer{Value: 1}).HashKey(): 2,
                (&object.Integer{Value: 2}).HashKey(): 3,
            },
        },
        {
            "{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
            map[object.HashKey]int64{
                (&object.Integer{Value: 2}).HashKey(): 4,
                (&object.Integer{Value: 6}).HashKey(): 16,
            },
        },
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
        return fmt.Errorf("expected value failed")
    }

    return nil
}

func testStringObject(expected string, actual object.Object) error {
    s, ok := actual.(*object.String)
    if !ok {
        return fmt.Errorf("type assertion error")
    }

    if s.Value != expected {
        return fmt.Errorf("wrong string value")
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
            t.Errorf("testBooleanObject failed, %s", err)
        }
    case string:
        err := testStringObject(expected, actual)
        if err != nil {
            t.Errorf("testStringObject failed, %s", err)
        }
    case []int:
        arr, ok := actual.(*object.Array)
        if !ok {
            t.Errorf("object is not Array")
            return
        }

        if len(arr.Elems) != len(expected) {
            t.Errorf("wrong array length")
            return
        }

        for i, num := range expected {
            err := testIntegerObject(int64(num), arr.Elems[i])
            if err != nil {
                t.Errorf("testIntegerObject failed")
            }
        }
    case map[object.HashKey]int64:
        hash, ok := actual.(*object.Hash)
        if !ok {
            t.Errorf("object is not Hash")
            return
        }

        if len(hash.Pairs) != len(expected) {
            t.Errorf("wrong hash length")
            return
        }

        for expectedKey, expectedVal := range expected {
            pair, ok := hash.Pairs[expectedKey]

            if !ok {
                t.Errorf("not found a pair")
            }

            err := testIntegerObject(int64(expectedVal), pair.Value)
            if err != nil {
                t.Errorf("testIntegerObject() failed")
            }
        }

    case *object.Null:
        if actual != Null {
            t.Errorf("object is not Null")
        }
    }
}
