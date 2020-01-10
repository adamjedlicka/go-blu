package main

import (
	"bufio"
	"fmt"
	"github.com/adamjedlicka/go-blue/src/vm"
	"os"
)

func main() {
	runRepl()
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
