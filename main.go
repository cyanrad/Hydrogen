package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"

	"main/evaluator"
	"main/lexer"
	"main/object"
	"main/parser"
)

// to do:
// - support unicode
// - add file name and line
// - support float, hex, oct, bin numbers
// - postfix in parser

const PROMPT = ">> "

func main() {
	var filepath string
	flag.StringVar(&filepath, "file", "", "Specify entry point")
	flag.Parse()

	if filepath != "" {
		interpretFile(filepath)
	} else {
		repl()
	}

}

func interpretFile(filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	// tokenizing
	l := lexer.CreateLexer(string(data))

	// parsing
	p := parser.CreateParser(l)
	program, errs := p.ParseProgram()
	if len(errs) != 0 {
		for _, e := range errs {
			fmt.Println(e.Error())
		}
		return
	}

	// interpreting
	env := evaluator.NewEnvironment()
	evaluator.Eval(program, env)
}

func repl() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	StartRepl(os.Stdin, os.Stdout)
}

func StartRepl(in io.Reader, out io.Writer) {
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
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
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
