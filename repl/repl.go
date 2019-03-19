import (
    "bufio"
    "fmt"
    "io"
    "monkey_interpreter/lexer"
    "monkey_interpreter/parser"
    "monkey_compiler/compiler"
    "monkey_compiler/vm"
)

func Start(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)

    constants := []object.Object{}
    globals := make([]object.Object, vm.GlobalsSize)
    symbolTable := compiler.NewSymbolTable()

    for {
        fmt.Printf(PROMPT)
        scanned := scanner.Scan()
        if !scanned {
            return
        }

        line := scanner.Text()
        l := lexer.New(line)
        p := parser.New(l)

        program := p.ParseProgram()
        if len(p.Errors()) != {
            printParseErrors(out, p.Errors())
            continue
        }

        comp := compiler.NewWithState(symbolTable, constants)
        err := comp.Compile(program)
        if err != nil {
            fmt.Fprintf(out, "Compilation failed")
            continue
        }

        code := comp.Bytecode()
        constants := code.Constants

        machine := vm.NewWithGlobalsStore(code, constants)
        err = machine.Run()
        if err != nil {
            fmt.Fprintf(out, "Execution bytecode failed")
            continue
        }

        stackTop = machine.LastPoppedStackElem()
        io.WriteString(out, stackTop.Inspect())
        io.WriteString(out, "\n")
    }
}
