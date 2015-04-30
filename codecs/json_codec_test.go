package codecs

import (
	"github.com/nelsam/silverback/codecs/constants"
	"github.com/stretchr/testify/assert"
	"testing"
)

var codec JSON

func TestJSONInterface(t *testing.T) {

	assert.Implements(t, (*Codec)(nil), new(JSON), "JSON")

}

func TestJSONMarshal(t *testing.T) {

	obj := make(map[string]string)
	obj["name"] = "Mat"

	// normally Marshal returns a []byte
	jsonString, jsonError := codec.Marshal(obj)

	if jsonError != nil {
		t.Errorf("Shouldn't return error: %s", jsonError)
	}

	// ask sam about this
	assert.Equal(t, string(jsonString.([]byte)), `{"name":"Mat"}`)

}

func TestJSONUnmarshal(t *testing.T) {

	jsonString := `{"name":"Mat"}`
	var object map[string]interface{}

	err := codec.Unmarshal([]byte(jsonString), &object)

	if err != nil {
		t.Errorf("Shouldn't return error: %s", err)
	}

	assert.Equal(t, "Mat", object["name"])

}

func TestJSONResponseContentType(t *testing.T) {

	assert.Equal(t, codec.ContentType(), constants.ContentTypeJSON)

}

func TestJSONFileExtension(t *testing.T) {

	assert.Equal(t, constants.FileExtensionJSON, codec.FileExtension())

}

func TestJSONCanMarshalWithCallback(t *testing.T) {

	assert.False(t, codec.CanMarshalWithCallback())

}
