package codecs

import (
	"encoding/json"

	"github.com/nelsam/silverback"
)

// JSON is a codec that handles json marshalling and unmarshalling.
type JSON struct{}

// New returns j.  This is because the JSON codec currently has no
// context to alter, so there's no need to use a separate copy across
// threads.
func (j *JSON) New(silverback.MIMEType) silverback.Codec {
	return j
}

// Types returns the MIME types that this codec is capable of handling.
func (j *JSON) Types() []silverback.MIMEType {
	return []silverback.MIMEType{
		{
			Type:    "application",
			SubType: "json",
		},
		{
			Type:    "text",
			SubType: "json",
		},
	}
}

// Marshal marshals target to a JSON string, returning the bytes and
// any errors encountered.
func (j *JSON) Marshal(target interface{}) ([]byte, error) {
	return json.Marshal(target)
}

// Unmarshal unmarshals a JSON string to the value that is pointed to
// by targetAddr, which must be a pointer.  It returns any errors
// encountered.
func (j *JSON) Unmarshal(raw []byte, targetAddr interface{}) error {
	return json.Unmarshal(raw, targetAddr)
}
