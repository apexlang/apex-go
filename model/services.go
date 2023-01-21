package model

import (
	"context"

	"github.com/apexlang/apex-go/errors"
	"github.com/apexlang/apex-go/location"
	"github.com/apexlang/apex-go/parser"
	"github.com/apexlang/apex-go/rules"
)

type parserImpl struct {
	resolver Resolver
}

func NewParser(resolver Resolver) Parser {
	return &parserImpl{
		resolver: resolver,
	}
}

func (p *parserImpl) Parse(ctx context.Context, source string) (*ParserResult, error) {
	doc, err := parser.Parse(parser.ParseParams{
		Source: source,
		Options: parser.ParseOptions{
			NoSource: true,
			Resolver: func(location, from string) (string, error) {
				return p.resolver.Resolve(ctx, location, from)
			},
		},
	})
	if err != nil {
		return nil, err
	}

	errs := rules.Validate(doc, rules.Rules...)
	if len(errs) > 0 {
		return &ParserResult{
			Errors: convertErrors(errs),
		}, nil
	}

	ns, errs := Convert(doc)
	if len(errs) > 0 {
		return &ParserResult{
			Errors: convertErrors(errs),
		}, nil
	}

	return &ParserResult{
		Namespace: ns,
	}, nil
}

func convertErrors(errs []error) []Error {
	e := make([]Error, len(errs))
	for i, err := range errs {
		switch v := err.(type) {
		case *errors.Error:
			e[i] = Error{
				Message:   v.Message,
				Positions: convertAny(v.Positions, func(p uint) uint32 { return uint32(p) }),
				Locations: convertAny(v.Locations, func(l location.SourceLocation) Location {
					return Location{
						Line:   uint32(l.Line),
						Column: uint32(l.Column),
					}
				}),
			}
		default:
			e[i] = Error{
				Message: err.Error(),
			}
		}

	}
	return e
}

func convertAny[S, D any](source []S, fn func(S) D) []D {
	dest := make([]D, len(source))
	for i, value := range source {
		dest[i] = fn(value)
	}
	return dest
}
