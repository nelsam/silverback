package silverback

import (
	"encoding/json"
	"net/http"
)

// A Codec contains methods for marshaling and unmarshaling data.
type Codec interface {
	// New takes a map of options passed in the MIME type, and returns
	// a codec which is set up to handle those options.  For example,
	// "application/json; encoding=utf-8" would pass
	// map[string]string{"encoding":"utf-8"} as the options argument.
	New(options map[string]string) Codec

	Match(*http.Request) bool
	Marshal(target interface{}) ([]byte, error)
	Unmarshal(raw []byte, targetAddr interface{}) error
}

// JSON is a codec that handles json marshalling and unmarshalling.
type JSON struct{}

// New returns j.  This is because the JSON codec currently has no
// context to alter, so there's no need to use a separate copy across
// threads.
func (j *JSON) New(map[string]string) Codec {
	return j
}

// Match returns whether or not the provided request is asking for the
// JSON codec.  This is usually from an entry in the Accept header.
func (j *JSON) Match(r *http.Request) bool {
	// Simplest cases
	switch r.Header.Get("Accept") {
	case "application/json", "text/json":
		return true
	}
	return false
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
