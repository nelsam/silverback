package codecs

import (
	"github.com/nelsam/silverback/codecs/constants"
	"labix.org/v2/mgo/bson"
)

// BSON converts objects to and from BSON.
type BSON struct{
	options map[string]string
}

// Marshal converts an object to BSON.
func (b *BSON) Marshal(target interface{}) (output interface{}, err error) {
	return bson.Marshal(target)
}

// Unmarshal converts BSON into an object.
func (b *BSON) Unmarshal(input, target interface{}) error {
	// take the input and convert it to target
	return bson.Unmarshal(input.([]byte), target)
}

// Options returns a copy of itself with the input options set.
func (b *BSON) Options(o map[string]string) Codec {
	bc := b
	bc.options = o
	return bc
}

// CanMarshalWithCallback returns whether this codec is capable of marshalling
// a response containing a callback.
func (b *BSON) CanMarshalWithCallback() bool {
	return false
}

// ContentType returns the content type for this codec.
func (b *BSON) ContentType() string {
	return constants.ContentTypeBSON
}

// FileExtension returns the file extension for this codec.
func (b *BSON) FileExtension() string {
	return constants.FileExtensionBSON
}
