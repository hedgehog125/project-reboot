package common

import (
	"bytes"
	"encoding/json"

	"github.com/xeipuuv/gojsonschema"
)

type PublicJSONSchema struct {
	*gojsonschema.Schema
	PublicSchema json.RawMessage
}

func NewPublicJSONSchema(schema *gojsonschema.Schema, rawSchema []byte) (*PublicJSONSchema, WrappedError) {
	buf := bytes.NewBuffer([]byte{})
	stdErr := json.Compact(buf, rawSchema)
	if stdErr != nil {
		return nil, ErrWrapperNewPublicJSONSchema.Wrap(stdErr)
	}

	return &PublicJSONSchema{
		Schema:       schema,
		PublicSchema: json.RawMessage(buf.Bytes()),
	}, nil
}
