package calculation

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type participleParser struct {
	parser        *participle.Parser[Expression]
	lexingSymbols map[lexer.TokenType]string
}

func (parser *participleParser) Parse(input string) (*Expression, error) {
	if len(strings.TrimSpace(input)) == 0 {
		log.Printf("Received empty input")
		return nil, errors.New(EMPTY_INPUT_ERROR_ID)
	}

	tokens, lexerError := parser.parser.Lex("", strings.NewReader(input))

	if lexerError != nil {
		log.Println("Lexing error thrown:", lexerError.Error())
		return nil, errors.New(LEXING_ERROR_ID)

	}

	logLexingResult(input, tokens, parser.lexingSymbols)

	expression, err := parser.parser.ParseString("", input) // This actually lexes again. We accept this..

	if err != nil {
		log.Println("Parsing error thrown:", err.Error())
		return nil, errors.New(PARSING_ERROR_ID)
	}

	log.Printf("Parsing succesful. Has parsed to %s", expression)
	return expression, nil
}

func newParticipleParser() participleParser {
	var myLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "UnsignedInteger", Pattern: `\d+`},
		{Name: "OpAdd", Pattern: `\+`},
		{Name: "OpSub", Pattern: `-`},
		{Name: "OpMul", Pattern: `\*`},
		{Name: "OpQuo", Pattern: `/`},
		{Name: "StartParen", Pattern: `\(`},
		{Name: "EndParen", Pattern: `\)`},
		{Name: "WhiteSpace", Pattern: `[\s]*`},
	})

	parser := participle.MustBuild[Expression](
		participle.Lexer(myLexer),
		participle.Elide("WhiteSpace"), // The parser should ignore any whitespace

	)

	return participleParser{parser: parser, lexingSymbols: lexer.SymbolsByRune(myLexer)}
}

func logLexingResult(input string, tokens []lexer.Token, lexingSymbols map[lexer.TokenType]string) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Lexed input \"%s\". Tokens: ", input))

	for _, token := range tokens {
		sb.WriteString(fmt.Sprintf("%s(%s)", lexingSymbols[token.Type], token.Value))
	}

	log.Println(sb.String())
}
