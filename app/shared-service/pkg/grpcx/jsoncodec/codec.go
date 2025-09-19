// Package jsoncodec регистрирует json-кодек для gRPC
package jsoncodec

import (
	"encoding/json"

	"google.golang.org/grpc/encoding"
)

// Codec реализует gRPC encoding.Codec с json
// имя кодека "json" => content-subtype: application/grpc+json
type Codec struct{}

func (Codec) Name() string                               { return "json" }
func (Codec) Marshal(v interface{}) ([]byte, error)      { return json.Marshal(v) }
func (Codec) Unmarshal(data []byte, v interface{}) error { return json.Unmarshal(data, v) }

func init() {
	encoding.RegisterCodec(Codec{})
}
