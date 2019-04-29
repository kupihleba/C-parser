package cs_parser

import (
	"fmt"
	"parser/cs_lexer"
	"strconv"
)

type State int

const (
	GO State = iota
	BACKTRACK
	SUCCESS
	FAILURE
)

const DEBUG = true

func (s State) String() string {
	switch s {
	case GO:
		return "Q"
	case BACKTRACK:
		return "B"
	case SUCCESS:
		return "SUCCESS"
	case FAILURE:
		return "FAILURE"
	default:
		return "???"
	}
}

var state State
var pos int // Next token pos
var l1 []Instance
var l2 []Action
var Tokens []int
var Grammar []Rule
var ExpToRule map[int]*Rule
var Expressions []Expression

func init() {
	if DEBUG {
		initSampleGrammar()
	} else {
		initGrammar()
	}
	initExpressions()
}

func reset() {
	state = GO
	pos = 0
	Tokens = Tokens[:0]
	l1 = l1[:0]
	l2 = l2[:0]
}

type Rule struct {
	Name    string
	ExprSet []Expression
}

func (r Rule) String() string {
	return r.Name
}

type Expression []Instance

type Command int

const (
	READ_TOKEN Command = iota
	RULE_MATCHED
)

type Action struct {
	value int
	cmd   Command
}

func (a Action) String() string {
	if a.cmd == READ_TOKEN {
		return "-"
	} else {
		return strconv.Itoa(a.value)
	}
}

func (this *Expression) matchSuffix(pattern []Instance) bool {
	if len(pattern) < len(*this) {
		return false
	}
	suffix := pattern[len(pattern)-len(*this):] // get suffix of expression length
	for i, s := range *this {
		if !equals(suffix[i], s) {
			//fmt.Printf("%v != %v\n", suffix[i], s)
			return false
		} else {
			//fmt.Printf("%v == %v\n", suffix[i], s)
		}
	}
	return true
}

func equals(a Instance, b Instance) bool {
	if a.t == b.t {
		if a.t == TERMINAL {
			return a.id_or_ref.(int) == b.id_or_ref.(int)
		} else if a.t == NONTERMINAL {
			return a.id_or_ref.(*Rule).Name == b.id_or_ref.(*Rule).Name
		}
	}
	return false
}

func initGrammar() {
	var PROGRAM, CLASS_DECL, CLASS_BODY, USING, VAR_DECL, IDENT_CHAIN, ASSIGN, EXPRESSION, METHOD, PREFIX, TYPE,
		NAMESPACE, EXPRESSION_PROP_CHAIN, STATEMENT, EXPR_CHAIN, METHOD_BODY Rule
	CLASS_BODY.Name = "CLASS BODY"
	CLASS_DECL.Name = "CLASS DECLARATION"
	USING.Name = "USING KEYWORD"
	VAR_DECL.Name = "VARIABLE DECLARATION"
	IDENT_CHAIN.Name = "CHAIN OF IDENTIFICATORS"
	ASSIGN.Name = "ASSIGN EXPRESSION"
	EXPRESSION.Name = "EXPRESSION"
	PROGRAM.Name = "PROGRAM"
	METHOD.Name = "METHOD DECLARATION"
	PREFIX.Name = "PREFIX KEYWORD"
	NAMESPACE.Name = "NAMESPACE"
	EXPRESSION_PROP_CHAIN.Name = "EXPRESSION WITH PROPERTY OPERATIONS"
	STATEMENT.Name = "STATEMENT"
	EXPR_CHAIN.Name = "EXPRESSION CHAIN"
	METHOD_BODY.Name = "METHOD BODY"

	PROGRAM.ExprSet = []Expression{
		{ref(&USING)},
		{ref(&CLASS_DECL)},
		{ref(&METHOD)},
		{ref(&NAMESPACE)},
		{ref(&VAR_DECL)},
		{ref(&PROGRAM), ref(&USING)},
		{ref(&PROGRAM), ref(&CLASS_DECL)},
		{ref(&PROGRAM), ref(&METHOD)},
		{ref(&PROGRAM), ref(&NAMESPACE)},
		{ref(&PROGRAM), ref(&VAR_DECL)},
	}
	NAMESPACE.ExprSet = []Expression{
		{tok("namespace"), tok("IDENTIFIER"), tok("{"), tok("}")},
		{tok("namespace"), tok("IDENTIFIER"), tok("{"), ref(&PROGRAM), tok("}")},
	}
	USING.ExprSet = []Expression{
		{tok("using"), tok("IDENTIFIER"), tok(";")},
	}
	CLASS_DECL.ExprSet = []Expression{
		{tok("class"), tok("IDENTIFIER"), tok("{"), ref(&CLASS_BODY), tok("}")},
		{tok("class"), tok("IDENTIFIER"), tok("{"), tok("}")},
	}

	EXPR_CHAIN.ExprSet = []Expression{
		{ref(&EXPRESSION)},
		{ref(&EXPRESSION), tok(","), ref(&EXPR_CHAIN)},
		{ref(&ASSIGN)},
		{ref(&ASSIGN), tok(","), ref(&EXPR_CHAIN)},
	}

	ASSIGN.ExprSet = []Expression{
		{tok("IDENTIFIER"), tok("="), ref(&EXPRESSION)},
	}
	METHOD_BODY.ExprSet = []Expression{
		{ref(&STATEMENT)},
		{ref(&METHOD_BODY), ref(&STATEMENT)},
	}

	EXPRESSION.ExprSet = []Expression{
		{tok("IDENTIFIER")},
		{tok("NUMBER")},
		{tok("STRING")},
		{tok("("), ref(&EXPRESSION), tok(")")},
		{tok("IDENTIFIER"), ref(&EXPRESSION_PROP_CHAIN)},
		{tok("NUMBER"), ref(&EXPRESSION_PROP_CHAIN)},
		{tok("("), ref(&EXPRESSION), ref(&EXPRESSION_PROP_CHAIN), tok(")")},
	}
	STATEMENT.ExprSet = []Expression{
		{ref(&EXPRESSION), tok(";")},
		{ref(&ASSIGN), tok(";")},
	}
	PREFIX.ExprSet = []Expression{
		{tok("static")},
	}
	CLASS_BODY.ExprSet = []Expression{
		{ref(&VAR_DECL)},
		{ref(&METHOD)},
		{ref(&METHOD), ref(&CLASS_BODY)},
		{ref(&VAR_DECL), ref(&CLASS_BODY)},
	}
	METHOD.ExprSet = []Expression{
		{ref(&PREFIX), ref(&TYPE), tok("IDENTIFIER"), tok("("), ref(&VAR_DECL), tok(")"), tok("{"), tok("}")},
		{ref(&PREFIX), ref(&TYPE), tok("IDENTIFIER"), tok("("), ref(&VAR_DECL), tok(")"), tok("{"), ref(&METHOD_BODY), tok("}")},

		{ref(&PREFIX), ref(&TYPE), tok("IDENTIFIER"), tok("("), ref(&VAR_DECL), tok(")"), tok(";")},

		{ref(&PREFIX), ref(&TYPE), tok("IDENTIFIER"), tok("("), tok(")"), tok("{"), tok("}")},
		{ref(&PREFIX), ref(&TYPE), tok("IDENTIFIER"), tok("("), tok(")"), tok("{"), ref(&METHOD_BODY), tok("}")},

		{ref(&PREFIX), ref(&TYPE), tok("IDENTIFIER"), tok("("), tok(")"), tok(";")},

		{ref(&TYPE), tok("IDENTIFIER"), tok("("), ref(&VAR_DECL), tok(")"), tok("{"), tok("}")},
		{ref(&TYPE), tok("IDENTIFIER"), tok("("), ref(&VAR_DECL), tok(")"), tok("{"), ref(&METHOD_BODY), tok("}")},
		{ref(&TYPE), tok("IDENTIFIER"), tok("("), ref(&VAR_DECL), tok(")"), tok(";")},

		{ref(&TYPE), tok("IDENTIFIER"), tok("("), tok(")"), tok("{"), tok("}")},
		{ref(&TYPE), tok("IDENTIFIER"), tok("("), tok(")"), tok("{"), ref(&METHOD_BODY), tok("}")},
		{ref(&TYPE), tok("IDENTIFIER"), tok("("), tok(")"), tok(";")},
	}
	TYPE.ExprSet = []Expression{
		{tok("string")},
		{tok("int")},
		{tok("void")},
		{tok("IDENTIFIER")},
	}

	EXPRESSION_PROP_CHAIN.ExprSet = []Expression{
		{tok("."), tok("IDENTIFIER")},
		{tok("("), tok(")")},
		{tok("("), ref(&EXPR_CHAIN), tok(")")},

		{ref(&EXPRESSION_PROP_CHAIN), tok("."), tok("IDENTIFIER")},
		{ref(&EXPRESSION_PROP_CHAIN), tok("["), ref(&EXPRESSION), tok("]")},
		{ref(&EXPRESSION_PROP_CHAIN), tok("("), ref(&EXPR_CHAIN), tok(")")},
		{ref(&EXPRESSION_PROP_CHAIN), tok("("), tok(")")},

		{tok("["), ref(&EXPRESSION), tok("]")},
		{ref(&EXPRESSION_PROP_CHAIN), tok("["), ref(&EXPRESSION), tok("]")},
	}
	VAR_DECL.ExprSet = []Expression{
		{ref(&TYPE), ref(&IDENT_CHAIN), tok(";")},
	}
	IDENT_CHAIN.ExprSet = []Expression{
		{tok("IDENTIFIER")},
		{tok("IDENTIFIER"), tok(","), ref(&IDENT_CHAIN)},
		{ref(&ASSIGN)},
		{ref(&ASSIGN), tok(","), ref(&IDENT_CHAIN)},
	}
	fmt.Printf("%v", Grammar)
	Grammar = append(Grammar, PROGRAM, CLASS_DECL, CLASS_BODY, USING, VAR_DECL, IDENT_CHAIN, ASSIGN, EXPRESSION, METHOD,
		PREFIX, TYPE, NAMESPACE, EXPRESSION_PROP_CHAIN, STATEMENT, EXPR_CHAIN, METHOD_BODY)
}

/**
 * Grammar for DEBUG mode is here:
 */
func initSampleGrammar() {
	var S, B Rule
	S.Name = "S"
	S.ExprSet = []Expression{ // S -> "a" | B
		{tok("a"), ref(&B)},
		//{ref(&C)},
	}
	B.Name = "B"
	B.ExprSet = []Expression{ // B -> "a" B | "b"
		{tok("a"), ref(&B)},
		{tok("b")},
	}
	//C.ExprSet = []Expression{
	//	{tok("b")},
	//}

	Grammar = append(Grammar, S, B) // ADD ALL RULES HERE!!
}

func initExpressions() {
	ExpToRule = make(map[int]*Rule)
	c := 0
	for i := 0; i < len(Grammar); i++ {
		Expressions = append(Expressions, Grammar[i].ExprSet...)
		for j := 0; j < len(Grammar[i].ExprSet); j++ {
			ExpToRule[c] = &Grammar[i]
			fmt.Printf("%s(%d): %s -> %d\n", ExpToRule[c].Name, j, ExpToRule[c].ExprSet, c)
			c++
		}
	}
	fmt.Printf("%d expressions loaded\n", c)
}

func ref(ref *Rule) Instance {
	var i Instance
	i.t = NONTERMINAL
	i.id_or_ref = ref
	return i
}

func tok(token_type string) Instance {
	return tokID(cs_lexer.TokenIdentifiers[token_type])
}

func tokID(id int) Instance {
	var i Instance
	i.t = TERMINAL
	i.id_or_ref = id
	return i
}

type Type int

const (
	TERMINAL Type = iota
	NONTERMINAL
)

type Instance struct {
	id_or_ref interface{}
	t         Type
}

func (i Instance) String() string {
	if i.NonTerminal() {
		return i.id_or_ref.(*Rule).Name
	}
	return cs_lexer.TokToStr(i.id_or_ref.(int))
}

func (this *Instance) NonTerminal() bool {
	return this.t == NONTERMINAL
}
func (this *Instance) IsTerminal() bool {
	return this.t == TERMINAL
}

func reachedTheEnd() bool {
	return pos == len(Tokens)
}

func findSuitableExpression() (int, error) {
	if len(l1) == 0 {
		return -1, fmt.Errorf("No expressions supplied")
	}

	for i, expr := range Expressions {
		if expr.matchSuffix(l1) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("No matchSuffix")
}

func findSuitableExpressionAfter(index int) (int, error) {
	if len(l1) == 0 {
		return -1, fmt.Errorf("No expressions supplied")
	}

	for i := index + 1; i < len(Expressions); i++ {
		if Expressions[i].matchSuffix(l1) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("No matchSuffix")
}

func printState(step int) {
	fmt.Printf("%d:\t(%s, %d, %v, %v)\n", step, state, pos, l1, l2)
}

func readToken() {
	l1 = append(l1, tokID(Tokens[pos]))
	l2 = append(l2, Action{pos, READ_TOKEN})
	pos++
}
func throwToken() {
	if pos <= 0 || l2[len(l2)-1].cmd != READ_TOKEN || l1[len(l1)-1].t != TERMINAL {
		panic("throw token failure")
	}
	l1 = l1[:len(l1)-1]
	l2 = l2[:len(l2)-1]
	pos--
}

func goExpression(expr_n int) {
	l1 = l1[:len(l1)-len(Expressions[expr_n])] // pop last elements, that matched rule
	l1 = append(l1, ref(ExpToRule[expr_n]))
	l2 = append(l2, Action{expr_n, RULE_MATCHED}) // l2 push
}

func rollback() Action {
	last_reduce := l2[len(l2)-1] // get previous reduce expression
	if last_reduce.cmd != RULE_MATCHED {
		panic("call for rollback on invalid rule")
	}

	l1 = l1[:len(l1)-1]                                // pop l1 rule
	l1 = append(l1, Expressions[last_reduce.value]...) // push old expressions back

	l2 = l2[:len(l2)-1]

	return last_reduce
}

func step() {
	switch state {
	case GO: // Q state
		if reachedTheEnd() &&
			len(l1) == 1 && // and non terminal is on top of the stack
			l1[0].NonTerminal() { // we may want to compare non terminal with S
			state = SUCCESS
			return
		}

		expr, err := findSuitableExpression() // Trying to find suitable expr
		if err != nil {                       // if there was no matchSuffix
			//println("No matchSuffix!")
			if pos < len(Tokens) { // if we can read more tokens
				readToken()
			} else {
				state = BACKTRACK
			}
		} else { // if we found a suitable expr
			//println("MATCHED RULE:", ExpToRule[expr].Name)
			goExpression(expr)
		}

		break
	case BACKTRACK: // B state
		if len(l1) == 0 {
			state = FAILURE
			return
		}

		if l1[len(l1)-1].IsTerminal() { // if last element is a terminal
			throwToken() // throw it away
			//println("Throw token!")
			return
		}

		// We reach this point in case
		// the last element is non terminal

		lastReduce := rollback() // we rollback changes
		expr, err := findSuitableExpressionAfter(lastReduce.value)

		if err == nil { // if alternative rule exists
			goExpression(expr)
			state = GO
		} else if pos < len(Tokens) { // ROLLBACK TEMP CHANGES
			readToken()
			state = GO
		}

	case SUCCESS:
		println("FILE HAS BEEN SUCCESSFULLY PARSED")
		return

	case FAILURE:
		println("ERRORS WERE ENCOUNTERED WHILE PARSING")
		return

	}
}

const PARSE_LIMIT = 1000000

func Parse(tokens []int) {
	reset()
	Tokens = tokens
	var i int
	for i = 1; state != SUCCESS && state != FAILURE; i++ {
		step()
		if DEBUG || i%1000000 == 0 {
			printState(i)
		}
	}
	if !DEBUG {
		printState(i)
	}
}
