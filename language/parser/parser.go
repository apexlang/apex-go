package parser

import (
	"fmt"

	"github.com/wapc/widl-go/errors"
	"github.com/wapc/widl-go/language/ast"
	"github.com/wapc/widl-go/language/lexer"
	"github.com/wapc/widl-go/language/source"
)

type parseFn func(parser *Parser) (interface{}, error)

// parse operation, fragment, typeSystem{schema, type..., extension, directives} definition
type parseDefinitionFn func(parser *Parser) (ast.Node, error)

var tokenDefinitionFn map[string]parseDefinitionFn

func init() {
	tokenDefinitionFn = make(map[string]parseDefinitionFn)
	{
		// for sign
		tokenDefinitionFn[lexer.GetTokenKindDesc(lexer.STRING)] = parseTypeSystemDefinition
		tokenDefinitionFn[lexer.GetTokenKindDesc(lexer.BLOCK_STRING)] = parseTypeSystemDefinition
		tokenDefinitionFn[lexer.GetTokenKindDesc(lexer.NAME)] = parseTypeSystemDefinition
		// for NAME
		tokenDefinitionFn[lexer.NAMESPACE] = parseNamespaceDefinition
		tokenDefinitionFn[lexer.SCALAR] = parseScalarTypeDefinition
		tokenDefinitionFn[lexer.TYPE] = parseObjectTypeDefinition
		tokenDefinitionFn[lexer.INTERFACE] = parseInterfaceTypeDefinition
		tokenDefinitionFn[lexer.UNION] = parseUnionTypeDefinition
		tokenDefinitionFn[lexer.ENUM] = parseEnumTypeDefinition
		tokenDefinitionFn[lexer.DIRECTIVE] = parseAnnotationDefinition
	}
}

type ParseOptions struct {
	NoLocation bool
	NoSource   bool
}

type ParseParams struct {
	Source  interface{}
	Options ParseOptions
}

type Parser struct {
	LexToken lexer.Lexer
	Source   *source.Source
	Options  ParseOptions
	PrevEnd  int
	Token    lexer.Token
}

func Parse(p ParseParams) (*ast.Document, error) {
	var sourceObj *source.Source
	switch src := p.Source.(type) {
	case *source.Source:
		sourceObj = src
	default:
		body, _ := p.Source.(string)
		sourceObj = source.NewSource(&source.Source{Body: []byte(body)})
	}
	parser, err := makeParser(sourceObj, p.Options)
	if err != nil {
		return nil, err
	}
	doc, err := parseDocument(parser)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// TODO: test and expose parseValue as a public
func parseValue(p ParseParams) (ast.Value, error) {
	var value ast.Value
	var sourceObj *source.Source
	switch src := p.Source.(type) {
	case *source.Source:
		sourceObj = src
	default:
		body, _ := p.Source.(string)
		sourceObj = source.NewSource(&source.Source{Body: []byte(body)})
	}
	parser, err := makeParser(sourceObj, p.Options)
	if err != nil {
		return value, err
	}
	value, err = parseValueLiteral(parser, false)
	if err != nil {
		return value, err
	}
	return value, nil
}

// Converts a name lex token into a name parse node.
func parseName(parser *Parser) (*ast.Name, error) {
	token, err := expect(parser, lexer.TokenKind[lexer.NAME])
	if err != nil {
		return nil, err
	}
	return ast.NewName(&ast.Name{
		Value: token.Value,
		Loc:   loc(parser, token.Start),
	}), nil
}

func makeParser(s *source.Source, opts ParseOptions) (*Parser, error) {
	lexToken := lexer.Lex(s)
	token, err := lexToken(0)
	if err != nil {
		return &Parser{}, err
	}
	return &Parser{
		LexToken: lexToken,
		Source:   s,
		Options:  opts,
		PrevEnd:  0,
		Token:    token,
	}, nil
}

/* Implements the parsing rules in the Document section. */

func parseDocument(parser *Parser) (*ast.Document, error) {
	var (
		nodes []ast.Node
		node  ast.Node
		item  parseDefinitionFn
		err   error
	)
	start := parser.Token.Start
	for {
		if skp, err := skip(parser, lexer.TokenKind[lexer.EOF]); err != nil {
			return nil, err
		} else if skp {
			break
		}
		switch parser.Token.Kind {
		case lexer.TokenKind[lexer.BRACE_L]:
			item = tokenDefinitionFn[lexer.GetTokenKindDesc(lexer.TokenKind[lexer.BRACE_L])]
		case lexer.TokenKind[lexer.NAME]:
			item = tokenDefinitionFn[lexer.GetTokenKindDesc(lexer.TokenKind[lexer.NAME])]
		case lexer.TokenKind[lexer.STRING]:
			item = tokenDefinitionFn[lexer.GetTokenKindDesc(lexer.TokenKind[lexer.STRING])]
		case lexer.TokenKind[lexer.BLOCK_STRING]:
			item = tokenDefinitionFn[lexer.GetTokenKindDesc(lexer.TokenKind[lexer.BLOCK_STRING])]
		default:
			return nil, unexpected(parser, lexer.Token{})
		}
		if node, err = item(parser); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return ast.NewDocument(&ast.Document{
		Loc:         loc(parser, start),
		Definitions: nodes,
	}), nil
}

/* Implements the parsing rules in the Operations section. */

/**
 * Arguments : ( Argument+ )
 */
func parseArguments(parser *Parser) ([]*ast.Argument, error) {
	arguments := []*ast.Argument{}
	if peek(parser, lexer.TokenKind[lexer.PAREN_L]) {
		if iArguments, err := reverse(parser,
			lexer.TokenKind[lexer.PAREN_L], parseArgument, lexer.TokenKind[lexer.PAREN_R],
			true,
		); err != nil {
			return arguments, err
		} else {
			for _, iArgument := range iArguments {
				arguments = append(arguments, iArgument.(*ast.Argument))
			}
		}
	}
	return arguments, nil
}

/**
 * Argument : Name : Value
 */
func parseArgument(parser *Parser) (interface{}, error) {
	var (
		err   error
		name  *ast.Name
		value ast.Value
	)
	start := parser.Token.Start
	if name, err = parseName(parser); err != nil {
		return nil, err
	}
	if _, err = expect(parser, lexer.TokenKind[lexer.COLON]); err != nil {
		return nil, err
	}
	if value, err = parseValueLiteral(parser, false); err != nil {
		return nil, err
	}
	return ast.NewArgument(&ast.Argument{
		Name:  name,
		Value: value,
		Loc:   loc(parser, start),
	}), nil
}

/* Implements the parsing rules in the Values section. */

/**
 * Value[Const] :
 *   - [~Const] Variable
 *   - IntValue
 *   - FloatValue
 *   - StringValue
 *   - BooleanValue
 *   - EnumValue
 *   - ListValue[?Const]
 *   - ObjectValue[?Const]
 *
 * BooleanValue : one of `true` `false`
 *
 * EnumValue : Name but not `true`, `false` or `null`
 */
func parseValueLiteral(parser *Parser, isConst bool) (ast.Value, error) {
	token := parser.Token
	switch token.Kind {
	case lexer.TokenKind[lexer.BRACE_L]:
		return parseMap(parser, isConst)
	case lexer.TokenKind[lexer.BRACKET_L]:
		return parseList(parser, isConst)
	case lexer.TokenKind[lexer.BRACE_L]:
		return parseObject(parser, isConst)
	case lexer.TokenKind[lexer.INT]:
		if err := advance(parser); err != nil {
			return nil, err
		}
		return ast.NewIntValue(&ast.IntValue{
			Value: token.Value,
			Loc:   loc(parser, token.Start),
		}), nil
	case lexer.TokenKind[lexer.FLOAT]:
		if err := advance(parser); err != nil {
			return nil, err
		}
		return ast.NewFloatValue(&ast.FloatValue{
			Value: token.Value,
			Loc:   loc(parser, token.Start),
		}), nil
	case lexer.TokenKind[lexer.BLOCK_STRING], lexer.TokenKind[lexer.STRING]:
		return parseStringLiteral(parser)
	case lexer.TokenKind[lexer.NAME]:
		if token.Value == "true" || token.Value == "false" {
			if err := advance(parser); err != nil {
				return nil, err
			}
			value := true
			if token.Value == "false" {
				value = false
			}
			return ast.NewBooleanValue(&ast.BooleanValue{
				Value: value,
				Loc:   loc(parser, token.Start),
			}), nil
		} else if token.Value != "null" {
			if err := advance(parser); err != nil {
				return nil, err
			}
			return ast.NewEnumValue(&ast.EnumValue{
				Value: token.Value,
				Loc:   loc(parser, token.Start),
			}), nil
		}
	}

	return nil, unexpected(parser, lexer.Token{})
}

func parseConstValue(parser *Parser) (interface{}, error) {
	value, err := parseValueLiteral(parser, true)
	if err != nil {
		return value, err
	}
	return value, nil
}

func parseValueValue(parser *Parser) (interface{}, error) {
	return parseValueLiteral(parser, false)
}

/**
 * ListValue[Const] :
 *   - [ ]
 *   - [ Value[?Const]+ ]
 */
func parseList(parser *Parser, isConst bool) (*ast.ListValue, error) {
	start := parser.Token.Start
	var item parseFn = parseValueValue
	if isConst {
		item = parseConstValue
	}
	values := []ast.Value{}
	if iValues, err := reverse(parser,
		lexer.TokenKind[lexer.BRACKET_L], item, lexer.TokenKind[lexer.BRACKET_R],
		false,
	); err != nil {
		return nil, err
	} else {
		for _, iValue := range iValues {
			values = append(values, iValue.(ast.Value))
		}
	}
	return ast.NewListValue(&ast.ListValue{
		Values: values,
		Loc:    loc(parser, start),
	}), nil
}

/**
 * MapValue[Const] :
 *   - [ ]
 *   - [ Value[?Const]+ ]
 */
func parseMap(parser *Parser, isConst bool) (*ast.MapValue, error) {
	start := parser.Token.Start
	var item parseFn = parseValueValue
	if isConst {
		item = parseConstValue
	}
	values := []ast.Value{}
	if iValues, err := reverse(parser,
		lexer.TokenKind[lexer.BRACE_L], item, lexer.TokenKind[lexer.BRACE_R],
		false,
	); err != nil {
		return nil, err
	} else {
		for _, iValue := range iValues {
			values = append(values, iValue.(ast.Value))
		}
	}
	return ast.NewMapValue(&ast.MapValue{
		Values: values,
		Loc:    loc(parser, start),
	}), nil
}

/**
 * ObjectValue[Const] :
 *   - { }
 *   - { ObjectField[?Const]+ }
 */
func parseObject(parser *Parser, isConst bool) (*ast.ObjectValue, error) {
	start := parser.Token.Start
	if _, err := expect(parser, lexer.TokenKind[lexer.BRACE_L]); err != nil {
		return nil, err
	}
	fields := []*ast.ObjectField{}
	for {
		if skp, err := skip(parser, lexer.TokenKind[lexer.BRACE_R]); err != nil {
			return nil, err
		} else if skp {
			break
		}
		if field, err := parseObjectField(parser, isConst); err != nil {
			return nil, err
		} else {
			fields = append(fields, field)
		}
	}
	return ast.NewObjectValue(&ast.ObjectValue{
		Fields: fields,
		Loc:    loc(parser, start),
	}), nil
}

/**
 * ObjectField[Const] : Name : Value[?Const]
 */
func parseObjectField(parser *Parser, isConst bool) (*ast.ObjectField, error) {
	var (
		name  *ast.Name
		value ast.Value
		err   error
	)
	start := parser.Token.Start
	if name, err = parseName(parser); err != nil {
		return nil, err
	}
	if _, err = expect(parser, lexer.TokenKind[lexer.COLON]); err != nil {
		return nil, err
	}
	if value, err = parseValueLiteral(parser, isConst); err != nil {
		return nil, err
	}
	return ast.NewObjectField(&ast.ObjectField{
		Name:  name,
		Value: value,
		Loc:   loc(parser, start),
	}), nil
}

/* Implements the parsing rules in the Annotations section. */

/**
 * Annotations : Annotation+
 */
func parseAnnotations(parser *Parser) ([]*ast.Annotation, error) {
	directives := []*ast.Annotation{}
	for peek(parser, lexer.TokenKind[lexer.AT]) {
		if directive, err := parseAnnotation(parser); err != nil {
			return directives, err
		} else {
			directives = append(directives, directive)
		}
	}
	return directives, nil
}

/**
 * Annotation : @ Name Arguments?
 */
func parseAnnotation(parser *Parser) (*ast.Annotation, error) {
	var (
		err  error
		name *ast.Name
		args []*ast.Argument
	)
	start := parser.Token.Start
	if _, err = expect(parser, lexer.TokenKind[lexer.AT]); err != nil {
		return nil, err
	}
	if name, err = parseName(parser); err != nil {
		return nil, err
	}
	if args, err = parseArguments(parser); err != nil {
		return nil, err
	}
	return ast.NewAnnotation(&ast.Annotation{
		Name:      name,
		Arguments: args,
		Loc:       loc(parser, start),
	}), nil
}

/* Implements the parsing rules in the Types section. */

/**
 * Type :
 *   - NamedType
 *   - ListType
 *   - NonNullType
 */
func parseType(parser *Parser) (ttype ast.Type, err error) {
	token := parser.Token
	var keyType, valueType ast.Type
	// [ String! ]!
	switch token.Kind {
	case lexer.TokenKind[lexer.BRACKET_L]:
		if err = advance(parser); err != nil {
			return nil, err
		}
		if ttype, err = parseType(parser); err != nil {
			return nil, err
		}
		fallthrough
	case lexer.TokenKind[lexer.BRACKET_R]:
		if err = advance(parser); err != nil {
			return nil, err
		}
		ttype = ast.NewList(&ast.List{
			Type: ttype,
			Loc:  loc(parser, token.Start),
		})
	case lexer.TokenKind[lexer.BRACE_L]:
		if err = advance(parser); err != nil {
			return nil, err
		}
		if keyType, err = parseType(parser); err != nil {
			return nil, err
		}
		if _, err = expect(parser, lexer.TokenKind[lexer.COLON]); err != nil {
			return nil, err
		}
		if valueType, err = parseType(parser); err != nil {
			return nil, err
		}
		fallthrough
	case lexer.TokenKind[lexer.BRACE_R]:
		if err = advance(parser); err != nil {
			return nil, err
		}
		ttype = ast.NewMap(&ast.Map{
			KeyType:   keyType,
			ValueType: valueType,
			Loc:       loc(parser, token.Start),
		})
	case lexer.TokenKind[lexer.NAME]:
		if ttype, err = parseNamed(parser); err != nil {
			return nil, err
		}
	}

	// QUESTION must be executed
	if skp, err := skip(parser, lexer.TokenKind[lexer.QUESTION]); err != nil {
		return nil, err
	} else if skp {
		ttype = ast.NewOptional(&ast.Optional{
			Type: ttype,
			Loc:  loc(parser, token.Start),
		})
	}
	return ttype, nil
}

/**
 * NamedType : Name
 */
func parseNamed(parser *Parser) (*ast.Named, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewNamed(&ast.Named{
		Name: name,
		Loc:  loc(parser, start),
	}), nil
}

/* Implements the parsing rules in the Type Definition section. */

/**
 * TypeSystemDefinition :
 *   - SchemaDefinition
 *   - TypeDefinition
 *   - TypeExtension
 *   - AnnotationDefinition
 *
 * TypeDefinition :
 *   - NamespaceDefinition
 *   - ObjectTypeDefinition
 *   - InterfaceTypeDefinition
 *   - UnionTypeDefinition
 *   - EnumTypeDefinition
 *   - InputObjectTypeDefinition
 */
func parseTypeSystemDefinition(parser *Parser) (ast.Node, error) {
	var (
		item parseDefinitionFn
		err  error
	)
	// Many definitions begin with a description and require a lookahead.
	keywordToken := parser.Token
	if peekDescription(parser) {
		if keywordToken, err = lookahead(parser); err != nil {
			return nil, err
		}
	}

	if keywordToken.Kind != lexer.NAME {
		return nil, unexpected(parser, keywordToken)
	}
	var ok bool
	if item, ok = tokenDefinitionFn[keywordToken.Value]; !ok {
		return nil, unexpected(parser, keywordToken)
	}
	return item(parser)
}

/**
 * ScalarTypeDefinition : Description? scalar Name Annotations?
 */
func parseNamespaceDefinition(parser *Parser) (ast.Node, error) {
	start := parser.Token.Start
	description, err := parseDescription(parser)
	if err != nil {
		return nil, err
	}
	_, err = expectKeyWord(parser, lexer.NAMESPACE)
	if err != nil {
		return nil, err
	}
	ns := parser.Token
	if ns.Kind == lexer.NS || ns.Kind == lexer.NAME || ns.Kind == lexer.STRING {
		advance(parser)
	} else {
		return nil, unexpected(parser, ns)
	}
	name := ast.NewName(&ast.Name{
		Value: ns.Value,
		Loc:   loc(parser, ns.Start),
	})
	directives, err := parseAnnotations(parser)
	if err != nil {
		return nil, err
	}
	def := ast.NewNamespaceDefinition(&ast.NamespaceDefinition{
		Name:        name,
		Description: description,
		Annotations: directives,
		Loc:         loc(parser, start),
	})
	return def, nil
}

/**
 * ScalarTypeDefinition : Description? scalar Name Directives?
 */
func parseScalarTypeDefinition(parser *Parser) (ast.Node, error) {
	start := parser.Token.Start
	description, err := parseDescription(parser)
	if err != nil {
		return nil, err
	}
	_, err = expectKeyWord(parser, lexer.SCALAR)
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	_, err = expect(parser, lexer.TokenKind[lexer.COLON])
	if err != nil {
		return nil, err
	}
	ttype, err := parseType(parser)
	if err != nil {
		return nil, err
	}
	annotations, err := parseAnnotations(parser)
	if err != nil {
		return nil, err
	}
	def := ast.NewScalarDefinition(&ast.ScalarDefinition{
		Name:        name,
		Description: description,
		Type:        ttype,
		Annotations: annotations,
		Loc:         loc(parser, start),
	})
	return def, nil
}

/**
 * ObjectTypeDefinition :
 *   Description?
 *   type Name ImplementsInterfaces? Annotations? { FieldDefinition+ }
 */
func parseObjectTypeDefinition(parser *Parser) (ast.Node, error) {
	start := parser.Token.Start
	description, err := parseDescription(parser)
	if err != nil {
		return nil, err
	}
	_, err = expectKeyWord(parser, lexer.TYPE)
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	interfaces, err := parseImplementsInterfaces(parser)
	if err != nil {
		return nil, err
	}
	directives, err := parseAnnotations(parser)
	if err != nil {
		return nil, err
	}
	iFields, err := reverse(parser,
		lexer.TokenKind[lexer.BRACE_L], parseFieldDefinition, lexer.TokenKind[lexer.BRACE_R],
		false,
	)
	if err != nil {
		return nil, err
	}
	fields := []*ast.FieldDefinition{}
	for _, iField := range iFields {
		if iField != nil {
			fields = append(fields, iField.(*ast.FieldDefinition))
		}
	}
	return ast.NewObjectDefinition(&ast.ObjectDefinition{
		Name:        name,
		Description: description,
		Loc:         loc(parser, start),
		Interfaces:  interfaces,
		Annotations: directives,
		Fields:      fields,
	}), nil
}

/**
 * ImplementsInterfaces :
 *   - implements `&`? NamedType
 *   - ImplementsInterfaces & NamedType
 */
func parseImplementsInterfaces(parser *Parser) ([]*ast.Named, error) {
	types := []*ast.Named{}
	if parser.Token.Value == "implements" {
		if err := advance(parser); err != nil {
			return nil, err
		}
		// optional leading ampersand
		skip(parser, lexer.TokenKind[lexer.AMP])
		for {
			ttype, err := parseNamed(parser)
			if err != nil {
				return types, err
			}
			types = append(types, ttype)
			if skipped, err := skip(parser, lexer.TokenKind[lexer.AMP]); !skipped {
				break
			} else if err != nil {
				return types, err
			}
		}
	}
	return types, nil
}

/**
 * OperationDefinition : Description? Name ArgumentsDefinition? : Type Annotations?
 */
func parseOperationDefinition(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	description, err := parseDescription(parser)
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	args, input, err := parseArgumentDefs(parser, true)
	if err != nil {
		return nil, err
	}
	_, err = expect(parser, lexer.TokenKind[lexer.COLON])
	if err != nil {
		return nil, err
	}
	ttype, err := parseType(parser)
	if err != nil {
		return nil, err
	}
	annotations, err := parseAnnotations(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewOperationDefinition(&ast.OperationDefinition{
		Name:        name,
		Description: description,
		Arguments:   args,
		Input:       input,
		Type:        ttype,
		Annotations: annotations,
		Loc:         loc(parser, start),
	}), nil
}

/**
 * FieldDefinition : Description? Name ArgumentsDefinition? : Type Annotations?
 */
func parseFieldDefinition(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	description, err := parseDescription(parser)
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	_, err = expect(parser, lexer.TokenKind[lexer.COLON])
	if err != nil {
		return nil, err
	}
	ttype, err := parseType(parser)
	if err != nil {
		return nil, err
	}
	var defaultValue ast.Value
	if skp, err := skip(parser, lexer.TokenKind[lexer.EQUALS]); err != nil {
		return nil, err
	} else if skp {
		if defaultValue, err = parseValueLiteral(parser, true); err != nil {
			return nil, err
		}
	}
	directives, err := parseAnnotations(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewFieldDefinition(&ast.FieldDefinition{
		Name:        name,
		Description: description,
		Type:        ttype,
		Default:     defaultValue,
		Annotations: directives,
		Loc:         loc(parser, start),
	}), nil
}

/**
 * ArgumentsDefinition : ( InputValueDefinition+ )
 */
func parseArgumentDefs(parser *Parser, unary bool) ([]*ast.InputValueDefinition, *ast.InputValueDefinition, error) {
	if peek(parser, lexer.TokenKind[lexer.PAREN_L]) {
		// Arguments operation
		iInputValueDefinitions, err := reverse(parser,
			lexer.TokenKind[lexer.PAREN_L], parseInputValueDef, lexer.TokenKind[lexer.PAREN_R],
			true,
		)
		if err != nil {
			return nil, nil, err
		}

		inputValueDefinitions := make([]*ast.InputValueDefinition, 0, len(iInputValueDefinitions))
		for _, iInputValueDefinition := range iInputValueDefinitions {
			if iInputValueDefinition != nil {
				inputValueDefinitions = append(inputValueDefinitions, iInputValueDefinition.(*ast.InputValueDefinition))
			}
		}

		return inputValueDefinitions, nil, nil
	} else if unary && peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
		// Unary operation
		if err := advance(parser); err != nil {
			return nil, nil, err
		}
		iInputValueDef, err := parseInputValueDef(parser)
		if err != nil {
			return nil, nil, err
		}

		if _, err := expect(parser, lexer.TokenKind[lexer.BRACE_R]); err != nil {
			return nil, nil, err
		}

		return nil, iInputValueDef.(*ast.InputValueDefinition), nil
	}

	return nil, nil, unexpected(parser, parser.Token)
}

/**
 * InputValueDefinition : Description? Name : Type DefaultValue? Annotations?
 */
func parseInputValueDef(parser *Parser) (interface{}, error) {
	var (
		description *ast.StringValue
		name        *ast.Name
		ttype       ast.Type
		directives  []*ast.Annotation
		err         error
	)
	start := parser.Token.Start
	if description, err = parseDescription(parser); err != nil {
		return nil, err
	}
	if name, err = parseName(parser); err != nil {
		return nil, err
	}
	if _, err = expect(parser, lexer.TokenKind[lexer.COLON]); err != nil {
		return nil, err
	}
	if ttype, err = parseType(parser); err != nil {
		return nil, err
	}
	var defaultValue ast.Value
	if skp, err := skip(parser, lexer.TokenKind[lexer.EQUALS]); err != nil {
		return nil, err
	} else if skp {
		val, err := parseConstValue(parser)
		if err != nil {
			return nil, err
		}
		if val, ok := val.(ast.Value); ok {
			defaultValue = val
		}
	}
	if directives, err = parseAnnotations(parser); err != nil {
		return nil, err
	}
	return ast.NewInputValueDefinition(&ast.InputValueDefinition{
		Name:         name,
		Description:  description,
		Type:         ttype,
		DefaultValue: defaultValue,
		Annotations:  directives,
		Loc:          loc(parser, start),
	}), nil
}

/**
 * InterfaceTypeDefinition :
 *   Description?
 *   interface Name Annotations? { FieldDefinition+ }
 */
func parseInterfaceTypeDefinition(parser *Parser) (ast.Node, error) {
	start := parser.Token.Start
	description, err := parseDescription(parser)
	if err != nil {
		return nil, err
	}
	_, err = expectKeyWord(parser, lexer.INTERFACE)
	if err != nil {
		return nil, err
	}
	// name, err := parseName(parser)
	// if err != nil {
	// 	return nil, err
	// }
	directives, err := parseAnnotations(parser)
	if err != nil {
		return nil, err
	}
	iOperations, err := reverse(parser,
		lexer.TokenKind[lexer.BRACE_L], parseOperationDefinition, lexer.TokenKind[lexer.BRACE_R],
		false,
	)
	if err != nil {
		return nil, err
	}
	operations := []*ast.OperationDefinition{}
	for _, iOperation := range iOperations {
		if iOperation != nil {
			operations = append(operations, iOperation.(*ast.OperationDefinition))
		}
	}
	return ast.NewInterfaceDefinition(&ast.InterfaceDefinition{
		//Name:        name,
		Description: description,
		Annotations: directives,
		Loc:         loc(parser, start),
		Operations:  operations,
	}), nil
}

/**
 * UnionTypeDefinition : Description? union Name Annotations? = UnionMembers
 */
func parseUnionTypeDefinition(parser *Parser) (ast.Node, error) {
	start := parser.Token.Start
	description, err := parseDescription(parser)
	if err != nil {
		return nil, err
	}
	_, err = expectKeyWord(parser, lexer.UNION)
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	directives, err := parseAnnotations(parser)
	if err != nil {
		return nil, err
	}
	_, err = expect(parser, lexer.TokenKind[lexer.EQUALS])
	if err != nil {
		return nil, err
	}
	types, err := parseUnionMembers(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewUnionDefinition(&ast.UnionDefinition{
		Name:        name,
		Description: description,
		Annotations: directives,
		Loc:         loc(parser, start),
		Types:       types,
	}), nil
}

/**
 * UnionMembers :
 *   - NamedType
 *   - UnionMembers | NamedType
 */
func parseUnionMembers(parser *Parser) ([]*ast.Named, error) {
	members := []*ast.Named{}
	for {
		member, err := parseNamed(parser)
		if err != nil {
			return members, err
		}
		members = append(members, member)
		if skp, err := skip(parser, lexer.TokenKind[lexer.PIPE]); err != nil {
			return nil, err
		} else if !skp {
			break
		}
	}
	return members, nil
}

/**
 * EnumTypeDefinition : Description? enum Name Annotations? { EnumValueDefinition+ }
 */
func parseEnumTypeDefinition(parser *Parser) (ast.Node, error) {
	start := parser.Token.Start
	description, err := parseDescription(parser)
	if err != nil {
		return nil, err
	}
	_, err = expectKeyWord(parser, lexer.ENUM)
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	directives, err := parseAnnotations(parser)
	if err != nil {
		return nil, err
	}
	iEnumValueDefs, err := reverse(parser,
		lexer.TokenKind[lexer.BRACE_L], parseEnumValueDefinition, lexer.TokenKind[lexer.BRACE_R],
		false,
	)
	if err != nil {
		return nil, err
	}
	values := []*ast.EnumValueDefinition{}
	for _, iEnumValueDef := range iEnumValueDefs {
		if iEnumValueDef != nil {
			values = append(values, iEnumValueDef.(*ast.EnumValueDefinition))
		}
	}
	return ast.NewEnumDefinition(&ast.EnumDefinition{
		Name:        name,
		Description: description,
		Annotations: directives,
		Loc:         loc(parser, start),
		Values:      values,
	}), nil
}

/**
 * EnumValueDefinition : Description? EnumValue Annotations?
 *
 * EnumValue : Name
 */
func parseEnumValueDefinition(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	description, err := parseDescription(parser)
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	directives, err := parseAnnotations(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewEnumValueDefinition(&ast.EnumValueDefinition{
		Name:        name,
		Description: description,
		Annotations: directives,
		Loc:         loc(parser, start),
	}), nil
}

/**
 * AnnotationDefinition :
 *   - directive @ Name ArgumentsDefinition? on AnnotationLocations
 */
func parseAnnotationDefinition(parser *Parser) (ast.Node, error) {
	var (
		err         error
		description *ast.StringValue
		name        *ast.Name
		args        []*ast.InputValueDefinition
		locations   []*ast.Name
	)
	start := parser.Token.Start
	if description, err = parseDescription(parser); err != nil {
		return nil, err
	}
	if _, err = expectKeyWord(parser, lexer.DIRECTIVE); err != nil {
		return nil, err
	}
	if _, err = expect(parser, lexer.TokenKind[lexer.AT]); err != nil {
		return nil, err
	}
	if name, err = parseName(parser); err != nil {
		return nil, err
	}
	if args, _, err = parseArgumentDefs(parser, false); err != nil {
		return nil, err
	}
	if _, err = expectKeyWord(parser, "on"); err != nil {
		return nil, err
	}
	if locations, err = parseAnnotationLocations(parser); err != nil {
		return nil, err
	}

	return ast.NewAnnotationDefinition(&ast.AnnotationDefinition{
		Loc:         loc(parser, start),
		Name:        name,
		Description: description,
		Arguments:   args,
		Locations:   locations,
	}), nil
}

/**
 * AnnotationLocations :
 *   - Name
 *   - AnnotationLocations | Name
 */
func parseAnnotationLocations(parser *Parser) ([]*ast.Name, error) {
	locations := []*ast.Name{}
	for {
		if name, err := parseName(parser); err != nil {
			return locations, err
		} else {
			locations = append(locations, name)
		}

		if hasPipe, err := skip(parser, lexer.TokenKind[lexer.PIPE]); err != nil {
			return locations, err
		} else if !hasPipe {
			break
		}
	}
	return locations, nil
}

func parseStringLiteral(parser *Parser) (*ast.StringValue, error) {
	token := parser.Token
	if err := advance(parser); err != nil {
		return nil, err
	}
	return ast.NewStringValue(&ast.StringValue{
		Value: token.Value,
		Loc:   loc(parser, token.Start),
	}), nil
}

/**
 * Description : StringValue
 */
func parseDescription(parser *Parser) (*ast.StringValue, error) {
	if peekDescription(parser) {
		return parseStringLiteral(parser)
	}
	return nil, nil
}

/* Core parsing utility functions */

// Returns a location object, used to identify the place in
// the source that created a given parsed object.
func loc(parser *Parser, start int) *ast.Location {
	if parser.Options.NoLocation {
		return nil
	}
	if parser.Options.NoSource {
		return ast.NewLocation(&ast.Location{
			Start: start,
			End:   parser.PrevEnd,
		})
	}
	return ast.NewLocation(&ast.Location{
		Start:  start,
		End:    parser.PrevEnd,
		Source: parser.Source,
	})
}

// Moves the internal parser object to the next lexed token.
func advance(parser *Parser) error {
	parser.PrevEnd = parser.Token.End
	token, err := parser.LexToken(parser.PrevEnd)
	if err != nil {
		return err
	}
	parser.Token = token
	return nil
}

// lookahead retrieves the next token
func lookahead(parser *Parser) (lexer.Token, error) {
	return parser.LexToken(parser.Token.End)
}

// Determines if the next token is of a given kind
func peek(parser *Parser, Kind int) bool {
	return parser.Token.Kind == Kind
}

// peekDescription determines if the next token is a string value
func peekDescription(parser *Parser) bool {
	return peek(parser, lexer.STRING) || peek(parser, lexer.BLOCK_STRING)
}

// If the next token is of the given kind, return true after advancing
// the parser. Otherwise, do not change the parser state and return false.
func skip(parser *Parser, Kind int) (bool, error) {
	if parser.Token.Kind == Kind {
		return true, advance(parser)
	}
	return false, nil
}

// If the next token is of the given kind, return that token after advancing
// the parser. Otherwise, do not change the parser state and return error.
func expect(parser *Parser, kind int) (lexer.Token, error) {
	token := parser.Token
	if token.Kind == kind {
		return token, advance(parser)
	}
	descp := fmt.Sprintf("Expected %s, found %s", lexer.GetTokenKindDesc(kind), lexer.GetTokenDesc(token))
	return token, errors.NewSyntaxError(parser.Source, token.Start, descp)
}

// If the next token is a keyword with the given value, return that token after
// advancing the parser. Otherwise, do not change the parser state and return false.
func expectKeyWord(parser *Parser, value string) (lexer.Token, error) {
	token := parser.Token
	if token.Kind == lexer.TokenKind[lexer.NAME] && token.Value == value {
		return token, advance(parser)
	}
	descp := fmt.Sprintf("Expected \"%s\", found %s", value, lexer.GetTokenDesc(token))
	return token, errors.NewSyntaxError(parser.Source, token.Start, descp)
}

// Helper function for creating an error when an unexpected lexed token
// is encountered.
func unexpected(parser *Parser, atToken lexer.Token) error {
	var token = atToken
	if (atToken == lexer.Token{}) {
		token = parser.Token
	}
	description := fmt.Sprintf("Unexpected %v", lexer.GetTokenDesc(token))
	return errors.NewSyntaxError(parser.Source, token.Start, description)
}

func unexpectedEmpty(parser *Parser, beginLoc int, openKind, closeKind int) error {
	description := fmt.Sprintf("Unexpected empty IN %s%s",
		lexer.GetTokenKindDesc(openKind),
		lexer.GetTokenKindDesc(closeKind),
	)
	return errors.NewSyntaxError(parser.Source, beginLoc, description)
}

//  Returns list of parse nodes, determined by
// the parseFn. This list begins with a lex token of openKind
// and ends with a lex token of closeKind. Advances the parser
// to the next lex token after the closing token.
// if zinteger is true, len(nodes) > 0
func reverse(parser *Parser, openKind int, parseFn parseFn, closeKind int, zinteger bool) ([]interface{}, error) {
	token, err := expect(parser, openKind)
	if err != nil {
		return nil, err
	}
	var nodes []interface{}
	for {
		if skp, err := skip(parser, closeKind); err != nil {
			return nil, err
		} else if skp {
			break
		}
		node, err := parseFn(parser)
		if err != nil {
			return nodes, err
		}
		nodes = append(nodes, node)
	}
	if zinteger && len(nodes) == 0 {
		return nodes, unexpectedEmpty(parser, token.Start, openKind, closeKind)
	}
	return nodes, nil
}
