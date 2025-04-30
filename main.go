package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"

	"main/evaluator"
	"main/lexer"
	"main/parser"
)

// to do:
// - support unicode
// - add file name and line
// - support float, hex, oct, bin numbers
// - postfix in parser

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := evaluator.NewEnvironment()
	for {
		// prompt user for input
		fmt.Print(PROMPT)

		// reading user input
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()

		// lexing
		l := lexer.CreateLexer(line)

		// parsing
		p := parser.CreateParser(l)
		program, err := p.ParseProgram()
		if len(err) != 0 {
			printParserErrors(out, err)
			continue
		}

		// interpreting
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []error) {
	for _, msg := range errors {
		str := fmt.Sprintf("%v", msg)
		io.WriteString(out, "\t"+str+"\n")
	}
}

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	Start(os.Stdin, os.Stdout)
}
