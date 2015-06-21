package silverback_test

import "github.com/nelsam/silverback"

type mockCodec struct {
	silverback.Codec
	matcher func(mime silverback.MIMEType) bool
}

func (m mockCodec) Match(mime silverback.MIMEType) bool {
	return m.matcher(mime)
}
