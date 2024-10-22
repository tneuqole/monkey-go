package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/tneuqole/monkey-go/compiler"
	"github.com/tneuqole/monkey-go/object"

	// "github.com/tneuqole/monkey-go/evaluator"
	"github.com/tneuqole/monkey-go/lexer"
	// "github.com/tneuqole/monkey-go/object"
	"github.com/tneuqole/monkey-go/parser"
	"github.com/tneuqole/monkey-go/vm"
)

const PROMPT = ">> "

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	// env := object.NewEnvironment()
	// macroEnv := object.NewEnvironment()

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

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
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		// evaluator.DefineMacros(program, macroEnv)
		// expanded := evaluator.ExpandMacros(program, macroEnv)
		//
		// evaluated := evaluator.Eval(expanded, env)
		// if evaluated != nil {
		// 	io.WriteString(out, evaluated.Inspect()+"\n")
		// }

		c := compiler.NewWithState(symbolTable, constants)
		err := c.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "compilation failed: %s", err)
		}

		bytecode := c.Bytecode()
		constants = bytecode.Constants
		machine := vm.NewWithGlobals(bytecode, globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "vm failed: %s", err)
		}

		result := machine.LastPoppedStackElem()
		io.WriteString(out, result.Inspect()+"\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "oopsy whoopsy\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
