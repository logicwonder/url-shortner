package json

import (
	"encoding/json"
	"fmt"

	"github.com/logicwonder/url-shortner/shortner"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*shortner.Redirect, error) {
	redirect := &shortner.Redirect{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, fmt.Errorf("serializer.Redirect.Decode: %w", err)
	}
	return redirect, nil
}

func (r *Redirect) Encode(input *shortner.Redirect) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("serializer.Redirect.Encode: %w", err)
	}
	return rawMsg, nil
}
