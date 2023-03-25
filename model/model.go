// Code generated by @apexlang/codegen. DO NOT EDIT.

package model

import (
	"context"
	"encoding/json"
	"errors"
)

type Parser interface {
	Parse(ctx context.Context, source string) (*ParserResult, error)
}

type Resolver interface {
	Resolve(ctx context.Context, location string, from string) (string, error)
}

type ParserResult struct {
	Namespace *Namespace `json:"namespace,omitempty" yaml:"namespace,omitempty" msgpack:"namespace,omitempty"`
	Errors    []Error    `json:"errors,omitempty" yaml:"errors,omitempty" msgpack:"errors,omitempty"`
}

// DefaultParserResult returns a `ParserResult` struct populated with its default
// values.
func DefaultParserResult() ParserResult {
	return ParserResult{}
}

type Error struct {
	Message   string     `json:"message" yaml:"message" msgpack:"message"`
	Positions []uint32   `json:"positions" yaml:"positions" msgpack:"positions"`
	Locations []Location `json:"locations" yaml:"locations" msgpack:"locations"`
}

// DefaultError returns a `Error` struct populated with its default values.
func DefaultError() Error {
	return Error{}
}

type Location struct {
	Line   uint32 `json:"line" yaml:"line" msgpack:"line"`
	Column uint32 `json:"column" yaml:"column" msgpack:"column"`
}

// DefaultLocation returns a `Location` struct populated with its default values.
func DefaultLocation() Location {
	return Location{}
}

// Namespace encapsulates is used to identify and refer to elements contained in
// the Apex specification.
type Namespace struct {
	Name        string       `json:"name" yaml:"name" msgpack:"name"`
	Description *string      `json:"description,omitempty" yaml:"description,omitempty" msgpack:"description,omitempty"`
	Annotations []Annotation `json:"annotations,omitempty" yaml:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Imports     []Import     `json:"imports,omitempty" yaml:"imports,omitempty" msgpack:"imports,omitempty"`
	Directives  []Directive  `json:"directives,omitempty" yaml:"directives,omitempty" msgpack:"directives,omitempty"`
	Aliases     []Alias      `json:"aliases,omitempty" yaml:"aliases,omitempty" msgpack:"aliases,omitempty"`
	Functions   []Operation  `json:"functions,omitempty" yaml:"functions,omitempty" msgpack:"functions,omitempty"`
	Interfaces  []Interface  `json:"interfaces,omitempty" yaml:"interfaces,omitempty" msgpack:"interfaces,omitempty"`
	Types       []Type       `json:"types,omitempty" yaml:"types,omitempty" msgpack:"types,omitempty"`
	Unions      []Union      `json:"unions,omitempty" yaml:"unions,omitempty" msgpack:"unions,omitempty"`
}

// DefaultNamespace returns a `Namespace` struct populated with its default values.
func DefaultNamespace() Namespace {
	return Namespace{}
}

// Apex can integrate external definitions using the import keyword.
type Import struct {
	Description *string      `json:"description,omitempty" yaml:"description,omitempty" msgpack:"description,omitempty"`
	All         bool         `json:"all" yaml:"all" msgpack:"all"`
	Names       []ImportRef  `json:"names,omitempty" yaml:"names,omitempty" msgpack:"names,omitempty"`
	From        string       `json:"from" yaml:"from" msgpack:"from"`
	Annotations []Annotation `json:"annotations,omitempty" yaml:"annotations,omitempty" msgpack:"annotations,omitempty"`
}

// DefaultImport returns a `Import` struct populated with its default values.
func DefaultImport() Import {
	return Import{}
}

type ImportRef struct {
	Name string  `json:"name" yaml:"name" msgpack:"name"`
	As   *string `json:"as,omitempty" yaml:"as,omitempty" msgpack:"as,omitempty"`
}

// DefaultImportRef returns a `ImportRef` struct populated with its default values.
func DefaultImportRef() ImportRef {
	return ImportRef{}
}

// Types are the most basic component of an Apex specification. They represent data
// structures with fields. Types are defined in a language-agnostic way. This means
// that complex features like nested structures, inheritance, and
// generics/templates are omitted by design.
type Type struct {
	Name        string       `json:"name" yaml:"name" msgpack:"name"`
	Description *string      `json:"description,omitempty" yaml:"description,omitempty" msgpack:"description,omitempty"`
	Fields      []Field      `json:"fields" yaml:"fields" msgpack:"fields"`
	Annotations []Annotation `json:"annotations,omitempty" yaml:"annotations,omitempty" msgpack:"annotations,omitempty"`
}

// DefaultType returns a `Type` struct populated with its default values.
func DefaultType() Type {
	return Type{}
}

// Interfaces are conceptual groups of operations that allow the developer to
// divide communication into multiple components. Typically, interfaces are named
// according to their purpose.
type Interface struct {
	Name        string       `json:"name" yaml:"name" msgpack:"name"`
	Description *string      `json:"description,omitempty" yaml:"description,omitempty" msgpack:"description,omitempty"`
	Operations  []Operation  `json:"operations" yaml:"operations" msgpack:"operations"`
	Annotations []Annotation `json:"annotations,omitempty" yaml:"annotations,omitempty" msgpack:"annotations,omitempty"`
}

// DefaultInterface returns a `Interface` struct populated with its default values.
func DefaultInterface() Interface {
	return Interface{}
}

// Alias types are used for cases when scalar types (like string) should be parsed
// our treated like a different data type in the generated code.
type Alias struct {
	Name        string       `json:"name" yaml:"name" msgpack:"name"`
	Description *string      `json:"description,omitempty" yaml:"description,omitempty" msgpack:"description,omitempty"`
	Type        TypeRef      `json:"type" yaml:"type" msgpack:"type"`
	Annotations []Annotation `json:"annotations,omitempty" yaml:"annotations,omitempty" msgpack:"annotations,omitempty"`
}

// DefaultAlias returns a `Alias` struct populated with its default values.
func DefaultAlias() Alias {
	return Alias{}
}

type Operation struct {
	Name        string       `json:"name" yaml:"name" msgpack:"name"`
	Description *string      `json:"description,omitempty" yaml:"description,omitempty" msgpack:"description,omitempty"`
	Parameters  []Parameter  `json:"parameters,omitempty" yaml:"parameters,omitempty" msgpack:"parameters,omitempty"`
	Unary       *Parameter   `json:"unary,omitempty" yaml:"unary,omitempty" msgpack:"unary,omitempty"`
	Returns     *TypeRef     `json:"returns,omitempty" yaml:"returns,omitempty" msgpack:"returns,omitempty"`
	Annotations []Annotation `json:"annotations,omitempty" yaml:"annotations,omitempty" msgpack:"annotations,omitempty"`
}

// DefaultOperation returns a `Operation` struct populated with its default values.
func DefaultOperation() Operation {
	return Operation{}
}

type Parameter struct {
	Name         string       `json:"name" yaml:"name" msgpack:"name"`
	Description  *string      `json:"description,omitempty" yaml:"description,omitempty" msgpack:"description,omitempty"`
	Type         TypeRef      `json:"type" yaml:"type" msgpack:"type"`
	DefaultValue *Value       `json:"defaultValue,omitempty" yaml:"defaultValue,omitempty" msgpack:"defaultValue,omitempty"`
	Annotations  []Annotation `json:"annotations,omitempty" yaml:"annotations,omitempty" msgpack:"annotations,omitempty"`
}

// DefaultParameter returns a `Parameter` struct populated with its default values.
func DefaultParameter() Parameter {
	return Parameter{}
}

type Field struct {
	Name         string       `json:"name" yaml:"name" msgpack:"name"`
	Description  *string      `json:"description,omitempty" yaml:"description,omitempty" msgpack:"description,omitempty"`
	Type         TypeRef      `json:"type" yaml:"type" msgpack:"type"`
	DefaultValue *Value       `json:"defaultValue,omitempty" yaml:"defaultValue,omitempty" msgpack:"defaultValue,omitempty"`
	Annotations  []Annotation `json:"annotations,omitempty" yaml:"annotations,omitempty" msgpack:"annotations,omitempty"`
}

// DefaultField returns a `Field` struct populated with its default values.
func DefaultField() Field {
	return Field{}
}

// Unions types denote that a type can have one of several representations.
type Union struct {
	Name        string       `json:"name" yaml:"name" msgpack:"name"`
	Description *string      `json:"description,omitempty" yaml:"description,omitempty" msgpack:"description,omitempty"`
	Types       []TypeRef    `json:"types" yaml:"types" msgpack:"types"`
	Annotations []Annotation `json:"annotations,omitempty" yaml:"annotations,omitempty" msgpack:"annotations,omitempty"`
}

// DefaultUnion returns a `Union` struct populated with its default values.
func DefaultUnion() Union {
	return Union{}
}

// Enumerations (or enums) are a type that is constrained to a finite set of
// allowed values.
type Enum struct {
	Name        string       `json:"name" yaml:"name" msgpack:"name"`
	Description *string      `json:"description,omitempty" yaml:"description,omitempty" msgpack:"description,omitempty"`
	Values      []EnumValue  `json:"values" yaml:"values" msgpack:"values"`
	Annotations []Annotation `json:"annotations,omitempty" yaml:"annotations,omitempty" msgpack:"annotations,omitempty"`
}

// DefaultEnum returns a `Enum` struct populated with its default values.
func DefaultEnum() Enum {
	return Enum{}
}

type EnumValue struct {
	Name        string       `json:"name" yaml:"name" msgpack:"name"`
	Description *string      `json:"description,omitempty" yaml:"description,omitempty" msgpack:"description,omitempty"`
	Index       uint64       `json:"index" yaml:"index" msgpack:"index"`
	Display     *string      `json:"display,omitempty" yaml:"display,omitempty" msgpack:"display,omitempty"`
	Annotations []Annotation `json:"annotations,omitempty" yaml:"annotations,omitempty" msgpack:"annotations,omitempty"`
}

// DefaultEnumValue returns a `EnumValue` struct populated with its default values.
func DefaultEnumValue() EnumValue {
	return EnumValue{}
}

// Directives are used to ensure that an annotation's arguments match an expected
// format.
type Directive struct {
	Name        string              `json:"name" yaml:"name" msgpack:"name"`
	Description *string             `json:"description,omitempty" yaml:"description,omitempty" msgpack:"description,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty" yaml:"parameters,omitempty" msgpack:"parameters,omitempty"`
	Locations   []DirectiveLocation `json:"locations" yaml:"locations" msgpack:"locations"`
	Require     []DirectiveRequire  `json:"require" yaml:"require" msgpack:"require"`
}

// DefaultDirective returns a `Directive` struct populated with its default values.
func DefaultDirective() Directive {
	return Directive{}
}

type DirectiveRequire struct {
	Directive string              `json:"directive" yaml:"directive" msgpack:"directive"`
	Locations []DirectiveLocation `json:"locations" yaml:"locations" msgpack:"locations"`
}

// DefaultDirectiveRequire returns a `DirectiveRequire` struct populated with its
// default values.
func DefaultDirectiveRequire() DirectiveRequire {
	return DirectiveRequire{}
}

// Annotations attach additional metadata to elements. These can be used in the
// code generation tool to implement custom functionality for your use case.
// Annotations have a name and zero or many arguments.
type Annotation struct {
	Name      string     `json:"name" yaml:"name" msgpack:"name"`
	Arguments []Argument `json:"arguments,omitempty" yaml:"arguments,omitempty" msgpack:"arguments,omitempty"`
}

// DefaultAnnotation returns a `Annotation` struct populated with its default
// values.
func DefaultAnnotation() Annotation {
	return Annotation{}
}

type Argument struct {
	Name  string `json:"name" yaml:"name" msgpack:"name"`
	Value Value  `json:"value" yaml:"value" msgpack:"value"`
}

// DefaultArgument returns a `Argument` struct populated with its default values.
func DefaultArgument() Argument {
	return Argument{}
}

type Named struct {
	Kind Kind   `json:"kind" yaml:"kind" msgpack:"kind"`
	Name string `json:"name" yaml:"name" msgpack:"name"`
}

// DefaultNamed returns a `Named` struct populated with its default values.
func DefaultNamed() Named {
	return Named{}
}

type List struct {
	Type TypeRef `json:"type" yaml:"type" msgpack:"type"`
}

// DefaultList returns a `List` struct populated with its default values.
func DefaultList() List {
	return List{}
}

type Map struct {
	KeyType   TypeRef `json:"keyType" yaml:"keyType" msgpack:"keyType"`
	ValueType TypeRef `json:"valueType" yaml:"valueType" msgpack:"valueType"`
}

// DefaultMap returns a `Map` struct populated with its default values.
func DefaultMap() Map {
	return Map{}
}

type Stream struct {
	Type TypeRef `json:"type" yaml:"type" msgpack:"type"`
}

// DefaultStream returns a `Stream` struct populated with its default values.
func DefaultStream() Stream {
	return Stream{}
}

type Optional struct {
	Type TypeRef `json:"type" yaml:"type" msgpack:"type"`
}

// DefaultOptional returns a `Optional` struct populated with its default values.
func DefaultOptional() Optional {
	return Optional{}
}

type Reference struct {
	Name string `json:"name" yaml:"name" msgpack:"name"`
}

// DefaultReference returns a `Reference` struct populated with its default values.
func DefaultReference() Reference {
	return Reference{}
}

type ListValue struct {
	Values []Value `json:"values" yaml:"values" msgpack:"values"`
}

// DefaultListValue returns a `ListValue` struct populated with its default values.
func DefaultListValue() ListValue {
	return ListValue{}
}

type ObjectValue struct {
	Fields []ObjectField `json:"fields" yaml:"fields" msgpack:"fields"`
}

// DefaultObjectValue returns a `ObjectValue` struct populated with its default
// values.
func DefaultObjectValue() ObjectValue {
	return ObjectValue{}
}

type ObjectField struct {
	Name  string `json:"name" yaml:"name" msgpack:"name"`
	Value Value  `json:"value" yaml:"value" msgpack:"value"`
}

// DefaultObjectField returns a `ObjectField` struct populated with its default
// values.
func DefaultObjectField() ObjectField {
	return ObjectField{}
}

type TypeRef struct {
	Scalar   *Scalar   `json:"Scalar,omitempty" yaml:"Scalar,omitempty" msgpack:"Scalar,omitempty"`
	Named    *Named    `json:"Named,omitempty" yaml:"Named,omitempty" msgpack:"Named,omitempty"`
	List     *List     `json:"List,omitempty" yaml:"List,omitempty" msgpack:"List,omitempty"`
	Map      *Map      `json:"Map,omitempty" yaml:"Map,omitempty" msgpack:"Map,omitempty"`
	Stream   *Stream   `json:"Stream,omitempty" yaml:"Stream,omitempty" msgpack:"Stream,omitempty"`
	Optional *Optional `json:"Optional,omitempty" yaml:"Optional,omitempty" msgpack:"Optional,omitempty"`
}

type Value struct {
	Bool        *bool        `json:"bool,omitempty" yaml:"bool,omitempty" msgpack:"bool,omitempty"`
	String      *string      `json:"string,omitempty" yaml:"string,omitempty" msgpack:"string,omitempty"`
	I64         *int64       `json:"i64,omitempty" yaml:"i64,omitempty" msgpack:"i64,omitempty"`
	F64         *float64     `json:"f64,omitempty" yaml:"f64,omitempty" msgpack:"f64,omitempty"`
	Reference   *Reference   `json:"Reference,omitempty" yaml:"Reference,omitempty" msgpack:"Reference,omitempty"`
	ListValue   *ListValue   `json:"ListValue,omitempty" yaml:"ListValue,omitempty" msgpack:"ListValue,omitempty"`
	ObjectValue *ObjectValue `json:"ObjectValue,omitempty" yaml:"ObjectValue,omitempty" msgpack:"ObjectValue,omitempty"`
}

type DirectiveLocation int32

const (
	DirectiveLocationNamespace DirectiveLocation = 0
	DirectiveLocationAlias     DirectiveLocation = 1
	DirectiveLocationUnion     DirectiveLocation = 2
	DirectiveLocationEnum      DirectiveLocation = 3
	DirectiveLocationEnumValue DirectiveLocation = 4
	DirectiveLocationType      DirectiveLocation = 5
	DirectiveLocationField     DirectiveLocation = 6
	DirectiveLocationInterface DirectiveLocation = 7
	DirectiveLocationOperation DirectiveLocation = 8
	DirectiveLocationParameter DirectiveLocation = 9
)

var toStringDirectiveLocation = map[DirectiveLocation]string{
	DirectiveLocationNamespace: "NAMESPACE",
	DirectiveLocationAlias:     "ALIAS",
	DirectiveLocationUnion:     "UNION",
	DirectiveLocationEnum:      "ENUM",
	DirectiveLocationEnumValue: "ENUM_VALUE",
	DirectiveLocationType:      "TYPE",
	DirectiveLocationField:     "FIELD",
	DirectiveLocationInterface: "INTERFACE",
	DirectiveLocationOperation: "OPERATION",
	DirectiveLocationParameter: "PARAMETER",
}

var toIDDirectiveLocation = map[string]DirectiveLocation{
	"NAMESPACE":  DirectiveLocationNamespace,
	"ALIAS":      DirectiveLocationAlias,
	"UNION":      DirectiveLocationUnion,
	"ENUM":       DirectiveLocationEnum,
	"ENUM_VALUE": DirectiveLocationEnumValue,
	"TYPE":       DirectiveLocationType,
	"FIELD":      DirectiveLocationField,
	"INTERFACE":  DirectiveLocationInterface,
	"OPERATION":  DirectiveLocationOperation,
	"PARAMETER":  DirectiveLocationParameter,
}

func (e DirectiveLocation) String() string {
	str, ok := toStringDirectiveLocation[e]
	if !ok {
		return "unknown"
	}
	return str
}

func (e *DirectiveLocation) FromString(str string) error {
	var ok bool
	*e, ok = toIDDirectiveLocation[str]
	if !ok {
		return errors.New("unknown value \"" + str + "\" for DirectiveLocation")
	}
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (e DirectiveLocation) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (e *DirectiveLocation) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	return e.FromString(str)
}

type Scalar int32

const (
	ScalarString   Scalar = 1
	ScalarBool     Scalar = 2
	ScalarI8       Scalar = 3
	ScalarI16      Scalar = 4
	ScalarI32      Scalar = 5
	ScalarI64      Scalar = 6
	ScalarU8       Scalar = 7
	ScalarU16      Scalar = 8
	ScalarU32      Scalar = 9
	ScalarU64      Scalar = 10
	ScalarF32      Scalar = 11
	ScalarF64      Scalar = 12
	ScalarBytes    Scalar = 13
	ScalarDatetime Scalar = 14
	ScalarAny      Scalar = 15
	ScalarRaw      Scalar = 16
)

var toStringScalar = map[Scalar]string{
	ScalarString:   "STRING",
	ScalarBool:     "BOOL",
	ScalarI8:       "I8",
	ScalarI16:      "I16",
	ScalarI32:      "I32",
	ScalarI64:      "I64",
	ScalarU8:       "U8",
	ScalarU16:      "U16",
	ScalarU32:      "U32",
	ScalarU64:      "U64",
	ScalarF32:      "F32",
	ScalarF64:      "F64",
	ScalarBytes:    "BYTES",
	ScalarDatetime: "DATETIME",
	ScalarAny:      "ANY",
	ScalarRaw:      "RAW",
}

var toIDScalar = map[string]Scalar{
	"STRING":   ScalarString,
	"BOOL":     ScalarBool,
	"I8":       ScalarI8,
	"I16":      ScalarI16,
	"I32":      ScalarI32,
	"I64":      ScalarI64,
	"U8":       ScalarU8,
	"U16":      ScalarU16,
	"U32":      ScalarU32,
	"U64":      ScalarU64,
	"F32":      ScalarF32,
	"F64":      ScalarF64,
	"BYTES":    ScalarBytes,
	"DATETIME": ScalarDatetime,
	"ANY":      ScalarAny,
	"RAW":      ScalarRaw,
}

func (e Scalar) String() string {
	str, ok := toStringScalar[e]
	if !ok {
		return "unknown"
	}
	return str
}

func (e *Scalar) FromString(str string) error {
	var ok bool
	*e, ok = toIDScalar[str]
	if !ok {
		return errors.New("unknown value \"" + str + "\" for Scalar")
	}
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (e Scalar) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (e *Scalar) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	return e.FromString(str)
}

type Kind int32

const (
	KindType      Kind = 1
	KindFunc      Kind = 2
	KindInterface Kind = 3
	KindAlias     Kind = 4
	KindUnion     Kind = 5
	KindEnum      Kind = 6
)

var toStringKind = map[Kind]string{
	KindType:      "TYPE",
	KindFunc:      "FUNC",
	KindInterface: "INTERFACE",
	KindAlias:     "ALIAS",
	KindUnion:     "UNION",
	KindEnum:      "ENUM",
}

var toIDKind = map[string]Kind{
	"TYPE":      KindType,
	"FUNC":      KindFunc,
	"INTERFACE": KindInterface,
	"ALIAS":     KindAlias,
	"UNION":     KindUnion,
	"ENUM":      KindEnum,
}

func (e Kind) String() string {
	str, ok := toStringKind[e]
	if !ok {
		return "unknown"
	}
	return str
}

func (e *Kind) FromString(str string) error {
	var ok bool
	*e, ok = toIDKind[str]
	if !ok {
		return errors.New("unknown value \"" + str + "\" for Kind")
	}
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (e Kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (e *Kind) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	return e.FromString(str)
}
