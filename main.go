package main

import (
	"fmt"
	lex "github.com/timtadh/lexmachine"
	"io/ioutil"
	"log"
	"os"
	"parser/cs_lexer"
	"parser/cs_parser"
)

const DEBUG = true

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func tokenize(file string) ([]int, error) {
	data, err := ioutil.ReadFile(file)
	check(err)
	fmt.Println(string(data))

	s, err := cs_lexer.Lexer.Scanner(data)

	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("Lexer error!")
	}

	var tokens []int
	fmt.Println("Type    | Lexeme     | Position")
	fmt.Println("--------+------------+------------")
	for tok, err, eof := s.Next(); !eof; tok, err, eof = s.Next() {
		if err != nil {
			log.Fatal(err)
		}
		token := tok.(*lex.Token)
		tokens = append(tokens, token.Type)
		fmt.Printf("%-7v | %-10v | %v:%v-%v:%v\n",
			cs_lexer.Tokens[token.Type],
			string(token.Lexeme),
			token.StartLine,
			token.StartColumn,
			token.EndLine,
			token.EndColumn)
	}

	fmt.Printf("%v\n", tokens)
	return tokens, nil
}

func main() {

	path := os.Args[1:]
	if len(path) != 1 {
		panic("No argument file supplied!")
	}
	if DEBUG {
		tokens, err := tokenize("sample.txt")
		check(err)
		cs_parser.Parse(tokens)
	} else {

		files, err := ioutil.ReadDir(path[0])
		check(err)

		for _, file := range files {
			if !file.IsDir() {
				tokens, err := tokenize(path[0] + "/" + file.Name())
				check(err)
				cs_parser.Parse(tokens)
			}
		}
	}
}
