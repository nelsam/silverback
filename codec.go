package silverback

// A Codec contains methods for marshaling and unmarshaling data.
type Codec interface {
	// New takes a MIMEType that was matched against this codec and
	// returns a codec which is set up to handle that MIMEType.  This
	// can be useful for codecs that match multiple MIMETypes, or for
	// reading matched.Options to see which codec options should be
	// applied to the returned codec.
	New(matched MIMEType) Codec

	// Types returns a slice of MIME types that this codec can handle.
	Types() []MIMEType

	Marshal(target interface{}) ([]byte, error)
	Unmarshal(raw []byte, targetAddr interface{}) error
}
