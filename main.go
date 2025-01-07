package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"

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
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.CreateLexer(line)
		p := parser.CreateParser(l)
		program, err := p.ParseProgram()
		if len(err) != 0 {
			printParserErrors(out, err)
			continue
		}
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
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
