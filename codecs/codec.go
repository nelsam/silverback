package codecs

// Codec is the interface to which a codec must conform.
type Codec interface {

	// Marshals the target to output of this codec's type.
	Marshal(target interface{}) (output interface{}, err error)

	// Unmarshals input of this codec's type to target.
	Unmarshal(input, target interface{}) error

	// Takes a map of Content-Type options and returns a copy of this codec with
	// those options applied.
	Options(map[string]string) Codec

	// ContentType returns the content type for this codec.
	ContentType() string

	// CanMarshalWithCallback returns whether this codec is capable of
	// marshalling a response containing a callback.
	CanMarshalWithCallback() bool

	// FileExtention returns the file extension for this codec.
	FileExtension() string
}

type ContentTypeMatcherCodec interface {
	Codec

	// ContentTypeSupported returns true if the passed in content type can be
	// handled by this codec, false otherwise
	ContentTypeSupported(contentType string) bool
}
