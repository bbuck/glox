package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bbuck/glox/scanner"
	"github.com/bbuck/glox/tree/parser"
	"github.com/bbuck/glox/tree/printer"
)

func main() {
	prog := os.Args[0]
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Printf("USAGE: %s [script]\n", prog)
		os.Exit(64)
	} else if len(args) == 1 {
		if err := runFile(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Reading file: %s", err.Error())
			os.Exit(1)
		}
	} else {
		if err := runPrompt(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Running REPL: %s", err.Error())
			os.Exit(1)
		}
	}
}

func runFile(name string) error {
	bytes, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	if err = run(string(bytes)); err != nil {
		os.Exit(65)
	}

	return nil
}

func runPrompt() error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			run(scanner.Text())
		} else {
			return scanner.Err()
		}
	}
}

func run(contents string) error {
	scanner := scanner.New(contents)
	scanner.ScanTokens()
	p := parser.New(scanner.Tokens())
	ex := p.Parse()
	if ex == nil {
		return nil
	}

	fmt.Println(printer.Print(ex))

	return nil
}
