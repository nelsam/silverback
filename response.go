package silverback

import "net/http"

// A Response is a container for typical response fields.  We use it
// mainly so that the handler methods don't need to worry about
// writing (or, for that matter, marshalling) data.
type Response struct {
	Status  int
	Headers http.Header
	Body    interface{}

	// codec is not exported because we don't want to give the
	// impression that the user should normally be setting this.
	// Normally, it will be assigned automatically from codecs, using
	// the Accept header of request.  SetCodec can be used to override
	// that behavior.
	codec Codec

	// codecs is a slice of codecs available for this Response to use
	// for formatting data.
	codecs []Codec

	request *http.Request
}

// NewResponse returns a *Response set up to correspond to r.  r will
// not be processed until it is needed (for example, when Codec() is
// called, but SetCodec() has not been called).
func NewResponse(r *http.Request, codecs []Codec) *Response {
	return &Response{
		request: r,
		codecs:  codecs,
	}
}

// Codec returns the codec that will be used for this response.
func (r *Response) Codec() Codec {
	if r.codec == nil {
		accept := ParseAcceptHeader(r.request.Header.Get("Accept"))
		r.codec = accept.Codec(r.codecs)
	}
	return r.codec
}

// SetCodec sets the codec to be used for this response.
func (r *Response) SetCodec(codec Codec) {
	r.codec = codec
}
