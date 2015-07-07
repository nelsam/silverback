package codecs

import (
	"github.com/nelsam/silverback/codecs/constants"
	jsonEncoding "encoding/json"
)

// JSON converts objects to and from JSON.
type JSON struct{
	options map[string]string
}

// Marshal converts an object to JSON.
func (j *JSON) Marshal(target interface{}) (output interface{}, err error) {
	return jsonEncoding.Marshal(target)
}

// Unmarshal converts JSON into an object.
func (j *JSON) Unmarshal(input, target interface{}) error {
	// take the input and convert it to target
	return jsonEncoding.Unmarshal(input.([]byte), target)
}

// Options returns a copy of itself with the input options set.
func (j *JSON) Options(o map[string]string) Codec {
	jc := j
	jc.options = o
	return jc
}

// CanMarshalWithCallback returns whether this codec is capable of marshalling
// a response containing a callback.
func (j *JSON) CanMarshalWithCallback() bool {
	return false
}

// ContentType returns the content type for this codec.
func (j *JSON) ContentType() string {
	return constants.ContentTypeJSON
}

// FileExtension returns the file extension for this codec.
func (j *JSON) FileExtension() string {
	return constants.FileExtensionJSON
}
