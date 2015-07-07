package services

import (
	"github.com/nelsam/silverback/codecs"
)

// contentTypeCodecWrapper is a wrapper for a Codec.  It is used to
// return any given Codec value, but with an overridden ContentType()
// value, usually for the purposes of returning the ContentType that
// was requested in an Accept header.
type contentTypeCodecWrapper struct {
	codec       codecs.Codec
	contentType string
	options     map[string]string
}

// wrapCodecWithContentType takes a codecs.Codec and a mime type
// string, returning a codecs.Codec that acts exactly like the passed
// in codecs.Codec except that the ContentType() method will return
// the passed in mime type instead of the underlying codec's default
// mime type.
func wrapCodecWithContentType(c codecs.Codec, typeString string) codecs.Codec {
	return &contentTypeCodecWrapper{
		codec:       c,
		contentType: typeString,
	}
}

func (c *contentTypeCodecWrapper) Marshal(object interface{}) (interface{}, error) {
	// Pass the matched content type as a codec option
	if c.options == nil {
		c.options = make(map[string]string)
	}
	c.options["matched_type"] = c.contentType
	c.codec = c.codec.Options(c.options)
	return c.codec.Marshal(object.([]byte))
}

func (c *contentTypeCodecWrapper) Unmarshal(input, target interface{}) error {
	return c.codec.Unmarshal(input, target)
}

func (c *contentTypeCodecWrapper) ContentType() string {
	return c.contentType
}

func (c *contentTypeCodecWrapper) FileExtension() string {
	return c.codec.FileExtension()
}

func (c *contentTypeCodecWrapper) CanMarshalWithCallback() bool {
	return c.codec.CanMarshalWithCallback()
}

func (c *contentTypeCodecWrapper) Options(o map[string]string) codecs.Codec {
	c.options = o
	return c
}
