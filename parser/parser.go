package parser

import (
	"fmt"
	"strconv"

	"github.com/wapc/widl-go/ast"
	"github.com/wapc/widl-go/errors"
	"github.com/wapc/widl-go/lexer"
	"github.com/wapc/widl-go/source"
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
		tokenDefinitionFn[lexer.IMPORT] = parseImportDefinition
		tokenDefinitionFn[lexer.ALIAS] = parseAliasDefinition
		tokenDefinitionFn[lexer.TYPE] = parseTypeDefinition
		tokenDefinitionFn[lexer.INTERFACE] = parseInterfaceDefinition
		tokenDefinitionFn[lexer.ROLE] = parseRoleDefinition
		tokenDefinitionFn[lexer.UNION] = parseUnionDefinition
		tokenDefinitionFn[lexer.ENUM] = parseEnumDefinition
		tokenDefinitionFn[lexer.DIRECTIVE] = parseDirectiveDefinition
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
	PrevEnd  uint
	Token    lexer.Token
}

func Parse(p ParseParams) (*ast.Document, error) {
	var sourceObj *source.Source
	switch src := p.Source.(type) {
	case *source.Source:
		sourceObj = src
	default:
		body, _ := p.Source.(string)
		sourceObj = source.NewSource("", []byte(body))
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
		sourceObj = source.NewSource("", []byte(body))
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
	return ast.NewName(
		loc(parser, token.Start),
		token.Value,
	), nil
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
	return ast.NewDocument(
		loc(parser, start),
		nodes,
	), nil
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
	if !peek(parser, lexer.TokenKind[lexer.NAME]) {
		name = ast.NewName(nil, "value")
	} else {
		if name, err = parseName(parser); err != nil {
			return nil, err
		}
		if _, err = expect(parser, lexer.TokenKind[lexer.COLON]); err != nil {
			return nil, err
		}
	}
	if value, err = parseValueLiteral(parser, false); err != nil {
		return nil, err
	}
	return ast.NewArgument(
		loc(parser, start),
		name,
		value,
	), nil
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
	// case lexer.TokenKind[lexer.BRACE_L]:
	// 	return parseMap(parser, isConst)
	case lexer.TokenKind[lexer.BRACKET_L]:
		return parseList(parser, isConst)
	case lexer.TokenKind[lexer.BRACE_L]:
		return parseObject(parser, isConst)
	case lexer.TokenKind[lexer.INT]:
		if err := advance(parser); err != nil {
			return nil, err
		}
		intVal, err := strconv.Atoi(token.Value)
		if err != nil {
			return nil, err
		}
		return ast.NewIntValue(
			loc(parser, token.Start),
			intVal,
		), nil
	case lexer.TokenKind[lexer.FLOAT]:
		if err := advance(parser); err != nil {
			return nil, err
		}
		floatVal, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			return nil, err
		}
		return ast.NewFloatValue(
			loc(parser, token.Start),
			floatVal,
		), nil
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
			return ast.NewBooleanValue(
				loc(parser, token.Start),
				value,
			), nil
		} else if token.Value != "null" {
			if err := advance(parser); err != nil {
				return nil, err
			}
			return ast.NewEnumValue(
				loc(parser, token.Start),
				token.Value,
			), nil
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
	return ast.NewListValue(
		loc(parser, start),
		values,
	), nil
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
	return ast.NewObjectValue(
		loc(parser, start),
		fields,
	), nil
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
	if parser.Token.Kind == lexer.TokenKind[lexer.NS] ||
		parser.Token.Kind == lexer.TokenKind[lexer.NAME] ||
		parser.Token.Kind == lexer.TokenKind[lexer.STRING] {
		name = ast.NewName(
			loc(parser, parser.Token.Start),
			parser.Token.Value,
		)
		advance(parser)
	} else {
		return nil, unexpected(parser, parser.Token)
	}
	if _, err = expect(parser, lexer.TokenKind[lexer.COLON]); err != nil {
		return nil, err
	}
	if value, err = parseValueLiteral(parser, isConst); err != nil {
		return nil, err
	}
	return ast.NewObjectField(
		loc(parser, start),
		name,
		value,
	), nil
}

/* Implements the parsing rules in the Annotations section. */

/**
 * Annotations : Annotation+
 */
func parseAnnotations(parser *Parser) ([]*ast.Annotation, error) {
	annotations := []*ast.Annotation{}
	for peek(parser, lexer.TokenKind[lexer.AT]) {
		if annotation, err := parseAnnotation(parser); err != nil {
			return annotations, err
		} else {
			annotations = append(annotations, annotation)
		}
	}
	return annotations, nil
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
	return ast.NewAnnotation(
		loc(parser, start),
		name,
		args,
	), nil
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
		ttype = ast.NewListType(
			loc(parser, token.Start),
			ttype,
		)
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
		ttype = ast.NewMapType(
			loc(parser, token.Start),
			keyType,
			valueType,
		)
	case lexer.TokenKind[lexer.NAME]:
		if ttype, err = parseNamed(parser); err != nil {
			return nil, err
		}
	}

	// QUESTION must be executed
	if skp, err := skip(parser, lexer.TokenKind[lexer.QUESTION]); err != nil {
		return nil, err
	} else if skp {
		ttype = ast.NewOptional(
			loc(parser, token.Start),
			ttype,
		)
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
	return ast.NewNamed(
		loc(parser, start),
		name,
	), nil
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
	name := ast.NewName(
		loc(parser, ns.Start),
		ns.Value,
	)
	annotations, err := parseAnnotations(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewNamespaceDefinition(
		loc(parser, start),
		name,
		description,
		annotations,
	), nil
}

/**
 * ImportDefinition : Description? import Name from StringValue
 */
func parseImportDefinition(parser *Parser) (ast.Node, error) {
	start := parser.Token.Start
	all := false
	importNames := []*ast.ImportName{}
	description, err := parseDescription(parser)
	if err != nil {
		return nil, err
	}
	if _, err = expectKeyWord(parser, lexer.IMPORT); err != nil {
		return nil, err
	}

	if peek(parser, lexer.TokenKind[lexer.STAR]) {
		all = true
		advance(parser)
	} else if peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
		// Parameters operation
		iImportNames, err := reverse(parser,
			lexer.TokenKind[lexer.BRACE_L], parseImportName, lexer.TokenKind[lexer.BRACE_R],
			true,
		)
		if err != nil {
			return nil, err
		}

		importNames = make([]*ast.ImportName, 0, len(iImportNames))
		for _, iImportName := range iImportNames {
			if iImportName != nil {
				importNames = append(importNames, iImportName.(*ast.ImportName))
			}
		}
	} else {
		return nil, unexpected(parser, parser.Token)
	}

	if _, err = expectKeyWord(parser, "from"); err != nil {
		return nil, err
	}

	from, err := parseStringLiteral(parser)
	if err != nil {
		return nil, err
	}
	annotations, err := parseAnnotations(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewImportDefinition(
		loc(parser, start),
		description,
		all,
		importNames,
		from,
		annotations,
	), nil
}

func parseImportName(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	var alias *ast.Name
	if _, ok, err := optionalKeyWord(parser, "as"); err != nil {
		return nil, err
	} else if ok {
		if alias, err = parseName(parser); err != nil {
			return nil, err
		}
	}
	return ast.NewImportName(
		loc(parser, start),
		name,
		alias,
	), nil
}

/**
 * AliasDefinition : Description? alias Name Directives?
 */
func parseAliasDefinition(parser *Parser) (ast.Node, error) {
	start := parser.Token.Start
	description, err := parseDescription(parser)
	if err != nil {
		return nil, err
	}
	_, err = expectKeyWord(parser, lexer.ALIAS)
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	_, err = expect(parser, lexer.TokenKind[lexer.EQUALS])
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
	return ast.NewAliasDefinition(
		loc(parser, start),
		name,
		description,
		ttype,
		annotations,
	), nil
}

/**
 * TypeDefinition :
 *   Description?
 *   type Name ImplementsInterfaces? Annotations? { FieldDefinition+ }
 */
func parseTypeDefinition(parser *Parser) (ast.Node, error) {
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
	annotations, err := parseAnnotations(parser)
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
	return ast.NewTypeDefinition(
		loc(parser, start),
		name,
		description,
		interfaces,
		annotations,
		fields,
	), nil
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
	params, unary, err := parseParameterDefs(parser, true)
	if err != nil {
		return nil, err
	}
	_, colon, err := optional(parser, lexer.TokenKind[lexer.COLON])
	if err != nil {
		return nil, err
	}
	var ttype ast.Type = ast.NewNamed(nil, ast.NewName(nil, "void"))
	if colon {
		ttype, err = parseType(parser)
		if err != nil {
			return nil, err
		}
	}
	annotations, err := parseAnnotations(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewOperationDefinition(
		loc(parser, start),
		name,
		description,
		ttype,
		annotations,
		unary,
		params,
	), nil
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
	annotations, err := parseAnnotations(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewFieldDefinition(
		loc(parser, start),
		name,
		description,
		ttype,
		defaultValue,
		annotations,
	), nil
}

/**
 * ParametersDefinition : ( ParameterDefinition+ )
 */
func parseParameterDefs(parser *Parser, unary bool) ([]*ast.ParameterDefinition, bool, error) {
	if peek(parser, lexer.TokenKind[lexer.PAREN_L]) {
		// Parameters operation
		iParameterDefinitions, err := reverse(parser,
			lexer.TokenKind[lexer.PAREN_L], parseParameterDef, lexer.TokenKind[lexer.PAREN_R],
			true,
		)
		if err != nil {
			return nil, false, err
		}

		parameterDefinitions := make([]*ast.ParameterDefinition, 0, len(iParameterDefinitions))
		for _, iParameterDefinition := range iParameterDefinitions {
			if iParameterDefinition != nil {
				parameterDefinitions = append(parameterDefinitions, iParameterDefinition.(*ast.ParameterDefinition))
			}
		}

		return parameterDefinitions, false, nil
	} else if unary && peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
		// Unary operation
		if err := advance(parser); err != nil {
			return nil, true, err
		}
		iInputValueDef, err := parseParameterDef(parser)
		if err != nil {
			return nil, true, err
		}

		if _, err := expect(parser, lexer.TokenKind[lexer.BRACE_R]); err != nil {
			return nil, true, err
		}

		return []*ast.ParameterDefinition{iInputValueDef.(*ast.ParameterDefinition)}, true, nil
	}

	return nil, false, unexpected(parser, parser.Token)
}

/**
 * ParameterDefinition : Description? Name : Type DefaultValue? Annotations?
 */
func parseParameterDef(parser *Parser) (interface{}, error) {
	var (
		description *ast.StringValue
		name        *ast.Name
		ttype       ast.Type
		annotations []*ast.Annotation
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
	if annotations, err = parseAnnotations(parser); err != nil {
		return nil, err
	}
	return ast.NewParameterDefinition(
		loc(parser, start),
		name,
		description,
		ttype,
		defaultValue,
		annotations,
	), nil
}

/**
 * InterfaceDefinition :
 *   Description?
 *   interface Annotations? { FieldDefinition+ }
 */
func parseInterfaceDefinition(parser *Parser) (ast.Node, error) {
	start := parser.Token.Start
	description, err := parseDescription(parser)
	if err != nil {
		return nil, err
	}
	_, err = expectKeyWord(parser, lexer.INTERFACE)
	if err != nil {
		return nil, err
	}
	annotations, err := parseAnnotations(parser)
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
	return ast.NewInterfaceDefinition(
		loc(parser, start),
		description,
		annotations,
		operations,
	), nil
}

/**
 * RoleDefinition :
 *   Description?
 *   role Name Annotations? { FieldDefinition+ }
 */
func parseRoleDefinition(parser *Parser) (ast.Node, error) {
	start := parser.Token.Start
	description, err := parseDescription(parser)
	if err != nil {
		return nil, err
	}
	_, err = expectKeyWord(parser, lexer.ROLE)
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	annotations, err := parseAnnotations(parser)
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
	return ast.NewRoleDefinition(
		loc(parser, start),
		name,
		description,
		annotations,
		operations,
	), nil
}

/**
 * UnionDefinition : Description? union Name Annotations? = UnionMembers
 */
func parseUnionDefinition(parser *Parser) (ast.Node, error) {
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
	annotations, err := parseAnnotations(parser)
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
	return ast.NewUnionDefinition(
		loc(parser, start),
		name,
		description,
		annotations,
		types,
	), nil
}

/**
 * UnionMembers :
 *   - NamedType
 *   - UnionMembers | NamedType
 */
func parseUnionMembers(parser *Parser) ([]ast.Type, error) {
	members := []ast.Type{}
	for {
		member, err := parseType(parser)
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
 * EnumDefinition : Description? enum Name Annotations? { EnumValueDefinition+ }
 */
func parseEnumDefinition(parser *Parser) (ast.Node, error) {
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
	return ast.NewEnumDefinition(
		loc(parser, start),
		name,
		description,
		directives,
		values,
	), nil
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
	_, err = expect(parser, lexer.TokenKind[lexer.EQUALS])
	if err != nil {
		return nil, err
	}
	token, err := expect(parser, lexer.TokenKind[lexer.INT])
	if err != nil {
		return nil, err
	}
	indexValue, err := strconv.Atoi(token.Value)
	index := ast.NewIntValue(loc(parser, token.Start), indexValue)
	annotations, err := parseAnnotations(parser)
	if err != nil {
		return nil, err
	}
	var display *ast.StringValue
	if _, ok, err := optionalKeyWord(parser, "as"); err != nil {
		return nil, err
	} else if ok {
		if display, err = parseStringLiteral(parser); err != nil {
			return nil, err
		}
	}
	return ast.NewEnumValueDefinition(
		loc(parser, start),
		name,
		description,
		index,
		display,
		annotations,
	), nil
}

/**
 * AnnotationDefinition :
 *   - directive @ Name ArgumentsDefinition? on AnnotationLocations
 */
func parseDirectiveDefinition(parser *Parser) (ast.Node, error) {
	var (
		err         error
		description *ast.StringValue
		name        *ast.Name
		params      []*ast.ParameterDefinition
		locations   []*ast.Name
		requires    []*ast.DirectiveRequire
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
	if params, _, err = parseParameterDefs(parser, false); err != nil {
		return nil, err
	}
	if _, err = expectKeyWord(parser, "on"); err != nil {
		return nil, err
	}
	if locations, err = parseDirectiveLocations(parser); err != nil {
		return nil, err
	}
	if _, ok, err := optionalKeyWord(parser, "require"); err != nil {
		return nil, err
	} else if ok {
		if requires, err = parseDirectiveRequires(parser); err != nil {
			return nil, err
		}
	}
	return ast.NewDirectiveDefinition(
		loc(parser, start),
		name,
		description,
		params,
		locations,
		requires,
	), nil
}

/**
 * DirectiveLocations :
 *   - Name
 *   - AnnotationLocations | Name
 */
func parseDirectiveLocations(parser *Parser) ([]*ast.Name, error) {
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

func parseDirectiveRequires(parser *Parser) ([]*ast.DirectiveRequire, error) {
	requires := []*ast.DirectiveRequire{}
	for {
		token, err := expect(parser, lexer.TokenKind[lexer.AT])
		if err != nil {
			return requires, err
		}
		name, err := parseName(parser)
		if err != nil {
			return requires, err
		}

		locations, err := parseDirectiveLocations(parser)
		if err != nil {
			return requires, err
		}

		requires = append(requires, ast.NewDirectiveRequire(
			loc(parser, token.Start),
			name,
			locations,
		))

		if hasPipe, err := skip(parser, lexer.TokenKind[lexer.PIPE]); err != nil {
			return requires, err
		} else if !hasPipe {
			break
		}
	}
	return requires, nil
}

func parseStringLiteral(parser *Parser) (*ast.StringValue, error) {
	token := parser.Token
	if err := advance(parser); err != nil {
		return nil, err
	}
	return ast.NewStringValue(
		loc(parser, token.Start),
		token.Value,
	), nil
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
func loc(parser *Parser, start uint) *ast.Location {
	if parser.Options.NoLocation {
		return nil
	}
	if parser.Options.NoSource {
		return ast.NewLocation(
			start,
			parser.PrevEnd,
			nil,
		)
	}
	return ast.NewLocation(
		start,
		parser.PrevEnd,
		parser.Source,
	)
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

func optional(parser *Parser, kind int) (lexer.Token, bool, error) {
	token := parser.Token
	if token.Kind == kind {
		return token, true, advance(parser)
	}
	return token, false, nil
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

func optionalKeyWord(parser *Parser, value string) (lexer.Token, bool, error) {
	token := parser.Token
	if token.Kind == lexer.TokenKind[lexer.NAME] && token.Value == value {
		return token, true, advance(parser)
	}
	return token, false, nil
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

func unexpectedEmpty(parser *Parser, beginLoc uint, openKind, closeKind int) error {
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
	_, err := expect(parser, openKind)
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
	// if zinteger && len(nodes) == 0 {
	// 	return nodes, unexpectedEmpty(parser, token.Start, openKind, closeKind)
	// }
	return nodes, nil
}
