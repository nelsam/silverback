package silverback_test

import "github.com/nelsam/silverback"

type mockCodec struct {
	silverback.Codec
	types []silverback.MIMEType
}

func makeCodec(mimeType, mimeSubType string) *mockCodec {
	return &mockCodec{
		types: []silverback.MIMEType{
			{
				Type:    mimeType,
				SubType: mimeSubType,
			},
		},
	}
}

func (m *mockCodec) Types() []silverback.MIMEType {
	return m.types
}
