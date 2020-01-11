package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/adamjedlicka/go-blu/src/compiler"
	"github.com/adamjedlicka/go-blu/src/parser"
	"github.com/adamjedlicka/go-blu/src/vm"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		runRepl()
	} else {
		runFile(os.Args[1])
	}
}

func runFile(name string) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}

	runes := bytes.Runes(data)

	p := parser.NewParser(runes)
	c := compiler.NewCompiler(name, p)

	vm := vm.NewVM()

	result := vm.Interpret(c.Compile())

	fmt.Println(result)
}

func runRepl() {
	reader := bufio.NewReader(os.Stdin)

	vm := vm.NewVM()

	for true {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		result := vm.Exec(line)

		fmt.Println(result)
	}
}
