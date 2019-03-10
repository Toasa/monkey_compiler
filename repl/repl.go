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

        comp := compiler.New()
        err := comp.Compile(program)
        if err != nil {
            fmt.Fprintf(out, "Compilation failed")
            continue
        }

        machine := vm.New(comp.Bytecode())
        err = machine.Run()
        if err != nil {
            fmt.Fprintf(out, "Execution bytecode failed")
            continue
        }

        stackTop = machine.StackTpp()
        io.WriteString(out, stackTop.Inspect())
        io.WriteString(out, "\n")
    }
}
