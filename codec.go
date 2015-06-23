package silverback

// A Codec contains methods for marshaling and unmarshaling data.
type Codec interface {
	// New takes a map of options passed in the MIME type, and returns
	// a codec which is set up to handle those options.  For example,
	// "application/json; encoding=utf-8" would pass
	// map[string]string{"encoding":"utf-8"} as the options argument.
	New(options map[string]string) Codec

	// Types returns a slice of MIME types that this codec can handle.
	Types() []MIMEType

	Marshal(target interface{}) ([]byte, error)
	Unmarshal(raw []byte, targetAddr interface{}) error
}
