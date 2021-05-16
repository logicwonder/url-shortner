package msgpack

import (
	"fmt"

	"github.com/logicwonder/url-shortner/shortner"
	"github.com/vmihailenco/msgpack"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*shortner.Redirect, error) {
	redirect := &shortner.Redirect{}
	if err := msgpack.Unmarshal(input, redirect); err != nil {
		return nil, fmt.Errorf("serializer.Redirect.Decode: %w", err)
	}
	return redirect, nil
}

func (r *Redirect) Encode(input *shortner.Redirect) ([]byte, error) {
	rawMsg, err := msgpack.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("serializer.Redirect.Encode: %w", err)
	}
	return rawMsg, nil
}
