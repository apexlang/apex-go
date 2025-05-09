package wapc

import (
	"context"

	"github.com/apexlang/apex-go/model"
	msgpack "github.com/wapc/tinygo-msgpack"
	"github.com/wapc/wapc-guest-tinygo"
)

type ResolverImpl struct {
	binding string
}

func NewResolver(binding ...string) *ResolverImpl {
	var bindingName string
	if len(binding) > 0 {
		bindingName = binding[0]
	}
	return &ResolverImpl{
		binding: bindingName,
	}
}

func (h *ResolverImpl) Resolve(ctx context.Context, location string, from string) (string, error) {
	inputArgs := model.ResolverResolveArgs{
		Location: location,
		From:     from,
	}
	inputBytes, err := msgpack.ToBytes(&inputArgs)
	if err != nil {
		return "", err
	}
	payload, err := wapc.HostCall(
		h.binding,
		"apexlang.v1.Resolver",
		"resolve",
		inputBytes,
	)
	if err != nil {
		return "", err
	}
	decoder := msgpack.NewDecoder(payload)
	return decoder.ReadString()
}

func RegisterParser(svc model.Parser) {
	wapc.RegisterFunction("apexlang.v1.Parser/parse", parserParseWrapper(svc))
}

func parserParseWrapper(svc model.Parser) wapc.Function {
	return func(payload []byte) ([]byte, error) {
		ctx := context.Background()
		decoder := msgpack.NewDecoder(payload)
		var inputArgs model.ParserParseArgs
		inputArgs.Decode(&decoder)
		response, err := svc.Parse(ctx, inputArgs.Source)
		if err != nil {
			return nil, err
		}
		return msgpack.ToBytes(response)
	}
}
