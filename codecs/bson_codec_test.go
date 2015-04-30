package codecs

import (
	"github.com/nelsam/silverback/codecs/constants"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBSONInterface(t *testing.T) {

	assert.Implements(t, (*Codec)(nil), new(BSON))

}

func TestBSONMarshal(t *testing.T) {

	codec := new(BSON)

	obj := make(map[string]string)
	obj["name"] = "Tyler"
	expectedResult := []byte{0x15, 0x0, 0x0, 0x0, 0x2, 0x6e, 0x61, 0x6d, 0x65, 0x0, 0x6, 0x0, 0x0, 0x0, 0x54, 0x79, 0x6c, 0x65, 0x72, 0x0, 0x0}

	bsonData, bsonError := codec.Marshal(obj)

	if bsonError != nil {
		t.Errorf("Shouldn't return error: %s", bsonError)
	}

	assert.Equal(t, bsonData, expectedResult)

}

func TestBSONUnmarshal(t *testing.T) {

	codec := new(BSON)
	bsonData := []byte{0x15, 0x0, 0x0, 0x0, 0x2, 0x6e, 0x61, 0x6d, 0x65, 0x0, 0x6, 0x0, 0x0, 0x0, 0x54, 0x79, 0x6c, 0x65, 0x72, 0x0, 0x0}
	var object map[string]interface{}

	err := codec.Unmarshal(bsonData, &object)

	if assert.Nil(t, err) {
		assert.Equal(t, "Tyler", object["name"])
	}

}

func TestBSONResponseContentType(t *testing.T) {

	codec := new(BSON)
	assert.Equal(t, codec.ContentType(), constants.ContentTypeBSON)
}

func TestBSONFileExtension(t *testing.T) {

	codec := new(BSON)
	assert.Equal(t, constants.FileExtensionBSON, codec.FileExtension())

}

func TestBSONCanMarshalWithCallback(t *testing.T) {

	codec := new(BSON)
	assert.False(t, codec.CanMarshalWithCallback())

}
