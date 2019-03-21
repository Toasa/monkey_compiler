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
                code.Make(code.OpAdd),
                code.Make(code.OpPop),
            },
        },
        {
            input: "1 - 2",
            expectedConstants: []interface{}{1, 2},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpConst, 0),
                code.Make(code.OpConst, 1),
                code.Make(code.OpSub),
                code.Make(code.OpPop),
            },
        },
        {
            input: "1 * 2",
            expectedConstants: []interface{}{1, 2},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpConst, 0),
                code.Make(code.OpConst, 1),
                code.Make(code.OpMul),
                code.Make(code.OpPop),
            },
        },
        {
            input: "2 / 1",
            expectedConstants: []interface{}{2, 1},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpConst, 0),
                code.Make(code.OpConst, 1),
                code.Make(code.OpDiv),
                code.Make(code.OpPop),
            },
        },
        {
            input: "-1",
            expectedConstants: []interface{}{1},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpConst, 0),
                code.Make(code.OpMinus),
                code.Make(code.OpPop),
            },
        },
    }

    runCompilerTest(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
    tests := []compilerTestCase {
        {
            input: "true",
            expectedConstants: []interface{}{},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpTrue),
                code.Make(code.OpPop),
            },
        },
        {
            input: "false",
            expectedConstants: []interface{}{},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpFalse),
                code.Make(code.OpPop),
            },
        },
        {
            input: "1 > 2",
            expectedConstants: []interface{}{1, 2},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpConst, 0),
                code.Make(code.OpConst, 1),
                code.Make(code.OpGT),
                code.Make(code.OpPop),
            },
        },
        {
            input: "1 < 2",
            expectedConstants: []interface{}{2, 1},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpConst, 0),
                code.Make(code.OpConst, 1),
                code.Make(code.OpGT),
                code.Make(code.OpPop),
            },
        },
        {
            input: "1 == 2",
            expectedConstants: []interface{}{1, 2},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpConst, 0),
                code.Make(code.OpConst, 1),
                code.Make(code.OpEq),
                code.Make(code.OpPop),
            },
        },
        {
            input: "1 != 2",
            expectedConstants: []interface{}{1, 2},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpConst, 0),
                code.Make(code.OpConst, 1),
                code.Make(code.OpNE),
                code.Make(code.OpPop),
            },
        },
        {
            input: "true == false",
            expectedConstants: []interface{}{},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpTrue),
                code.Make(code.OpFalse),
                code.Make(code.OpEq),
                code.Make(code.OpPop),
            },
        },
        {
            input: "true != false",
            expectedConstants: []interface{}{},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpTrue),
                code.Make(code.OpFalse),
                code.Make(code.OpNE),
                code.Make(code.OpPop),
            },
        },
        {
            input: "!true",
            expectedConstants: []interface{}{},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpTrue),
                code.Make(code.OpBang),
                code.Make(code.OpPop),
            },
        },
    }

    runCompilerTest(t, tests)
}

func TestConditionals(t *testing.T) {
    tests := []compilerTestCase {
        {
            input : `
            if (true) { 10 }; 3333;
            `,
            expectedConstants: []interface{}{10, 3333},
            expectedInstructions: []code.Instructions {
                // 0000
                code.Make(code.OpTrue),
                // 0001
                code.Make(code.OpJumpNotTruthy, 10),
                // 0004
                code.Make(code.OpConst, 0),
                // 0007
                code.Make(code.OpJump, 11),
                // 0010
                code.Make(code.OpNull),
                // 0011
                code.Make(code.OpPop),
                // 0012
                code.Make(code.OpConst, 1),
                // 0015
                code.Make(code.OpPop),
            },
        },
        {
            input : `
            if (true) { 10 } else { 20 }; 3333;
            `,
            expectedConstants: []interface{}{10, 20, 3333},
            expectedInstructions: []code.Instructions {
                // 0000
                code.Make(code.OpTrue),
                // 0001
                code.Make(code.OpJumpNotTruthy, 10),
                // 0004
                code.Make(code.OpConst, 0),
                // 0007
                code.Make(code.OpJump, 13),
                // 0010
                code.Make(code.OpConst, 1),
                // 0013
                code.Make(code.OpPop),
                // 0014
                code.Make(code.OpConst, 2),
                // 0017
                code.Make(code.OpPop),
            },
        },
    }

    runCompilerTest(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
    tests := []compilerTestCase {
        {
            input: `
            let one = 1;
            let two = 2;
            `,
            expectedConstants: []interface{}{1, 2},
            expectedInstructions: []code.Instructions {
                code.Make(code.OpConst, 0),
                code.Make(code.OpSetGlobal, 0),
                code.Make(code.OpConst, 1),
                code.Make(code.OpSetGlobal, 1),
            },
        },
        {
            input: `
            let one = 1;
            one;
            `,
            expectedConstants: []interface{}{1},
            expectedInstructions: []code.Instructions {
                code.Make(code.OpConst, 0),
                code.Make(code.OpSetGlobal, 0),
                code.Make(code.OpGetGlobal, 0),
                code.Make(code.OpPop),
            },
        },
        {
            input: `
            let one = 1;
            let two = one;
            two;
            `,
            expectedConstants: []interface{}{1},
            expectedInstructions: []code.Instructions {
                code.Make(code.OpConst, 0),
                code.Make(code.OpSetGlobal, 0),
                code.Make(code.OpGetGlobal, 0),
                code.Make(code.OpSetGlobal, 1),
                code.Make(code.OpGetGlobal, 1),
                code.Make(code.OpPop),
            },
        },
    }

    runCompilerTest(t, tests)
}

func TestStringExpressions(t *testing.T) {
    tests := []compilerTestCase {
        {
            input: `"monkey"`,
            expectedConstants: []interface{}{"monkey"},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpConst, 0),
                code.Make(code.OpPop),
            },
        },
        {
            input: `"mon" + "key"`,
            expectedConstants: []interface{}{"mon", "key"},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpConst, 0),
                code.Make(code.OpConst, 1),
                code.Make(code.OpAdd),
                code.Make(code.OpPop),
            },
        },
    }

    runCompilerTest(t, tests)
}

func TestArrayLiterals(t *testing.T) {
    tests := []compilerTestCase {
        {
            input: "[]",
            expectedConstants: []interface{}{},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpArray, 0),
                code.Make(code.OpPop),
            },
        },
        {
            input: "[1, 2, 3]",
            expectedConstants: []interface{}{1, 2, 3},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpConst, 0),
                code.Make(code.OpConst, 1),
                code.Make(code.OpConst, 2),
                code.Make(code.OpArray, 3),
                code.Make(code.OpPop),
            },
        },
        {
            input: "[1 + 2, 3 - 4, 5 * 6]",
            expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
            expectedInstructions: []code.Instructions{
                code.Make(code.OpConst, 0),
                code.Make(code.OpConst, 1),
                code.Make(code.OpAdd),
                code.Make(code.OpConst, 2),
                code.Make(code.OpConst, 3),
                code.Make(code.OpSub),
                code.Make(code.OpConst, 4),
                code.Make(code.OpConst, 5),
                code.Make(code.OpMul),
                code.Make(code.OpArray, 3),
                code.Make(code.OpPop),
            },
        },
    }

    runCompilerTest(t, tests)
}

func TestHashLiterals(t *testing.T) {
    tests := []compilerTestCase {
        {
            input: "{}",
            expectedConstants: []interface{}{},
            expectedInstructions: []code.Instructions {
                code.Make(code.OpHash, 0),
                code.Make(code.OpPop),
            },
        },
        {
            input: "{1: 2, 3: 4, 5: 6}",
            expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
            expectedInstructions: []code.Instructions {
                code.Make(code.OpConst, 0),
                code.Make(code.OpConst, 1),
                code.Make(code.OpConst, 2),
                code.Make(code.OpConst, 3),
                code.Make(code.OpConst, 4),
                code.Make(code.OpConst, 5),
                code.Make(code.OpHash, 6),
                code.Make(code.OpPop),
            },
        },
        {
            input: "{1: 2 + 3, 4: 5 * 6}",
            expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
            expectedInstructions: []code.Instructions {
                code.Make(code.OpConst, 0),
                code.Make(code.OpConst, 1),
                code.Make(code.OpConst, 2),
                code.Make(code.OpAdd),
                code.Make(code.OpConst, 3),
                code.Make(code.OpConst, 4),
                code.Make(code.OpConst, 5),
                code.Make(code.OpMul),
                code.Make(code.OpHash, 4),
                code.Make(code.OpPop),
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
        return fmt.Errorf("wrong number of constants\n expected: %d, but got: %d", len(expected), len(actual))
    }

    for i, cons := range expected {
        switch cons := cons.(type) {
        case int:
            err := testIntegerObject(int64(cons), actual[i])
            if err != nil {
                return fmt.Errorf("incorrect value")
            }
        case string:
            err := testStringObject(cons, actual[i])
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

func testStringObject(expected string, actual object.Object) error {
    s, ok := actual.(*object.String)
    if !ok {
        return fmt.Errorf("type assertion error")
    }

    if expected != s.Value {
        return fmt.Errorf("incorrect value")
    }

    return nil
}
