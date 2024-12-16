package main

import (
	"bufio"
	"fmt"
	"io"
	"main/token"
	"os"
	"os/user"

	"main/lexer"
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
		for tok := l.GetNextToken(); tok.Type != token.EOF; tok = l.GetNextToken() {
			fmt.Printf("%+v\n", tok)
		}
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
