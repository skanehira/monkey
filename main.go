package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/skanehira/monkey/compiler"
	"github.com/skanehira/monkey/lexer"
	"github.com/skanehira/monkey/object"
	"github.com/skanehira/monkey/parser"
	"github.com/skanehira/monkey/repl"
	"github.com/skanehira/monkey/vm"
	"golang.org/x/term"
)

func onExit(msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func main() {
	run := func(input string) {
		constants := []object.Object{}
		globals := make([]object.Object, vm.GlobalsSize)
		symbolTable := compiler.NewSymbolTable()
		for i, v := range object.Builtins {
			symbolTable.DefineBuiltin(i, v.Name)
		}

		l := lexer.New(string(input))
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			onExit(strings.Join(p.Errors(), "\n"))
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			onExit(err)
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalsStore(code, globals)
		err = machine.Run()
		if err != nil {
			onExit(err)
		}

		stackTop := machine.LastPoppedStackElem()
		fmt.Println(stackTop.Inspect())

	}
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		input, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			onExit(err)
		}
		run(string(input))
	} else if len(os.Args) > 1 {
		input, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			onExit(err)
		}
		run(string(input))
	} else {
		repl.Start(os.Stdin, os.Stdout)
	}
}
