package silverback_test

import (
	"net/http"

	"github.com/nelsam/silverback"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Response", func() {
	var (
		resp     *silverback.Response
		req      *http.Request
		mockJSON = makeCodec("application", "json")
	)

	BeforeEach(func() {
		var err error
		req, err = http.NewRequest("GET", "foo", nil)
		req.Header.Set("Accept", "application/json")
		Expect(err).ToNot(HaveOccurred())
		codecs := []silverback.Codec{
			makeCodec("text", "xml"),
			mockJSON,
		}
		resp = silverback.NewResponseForCodecs(req, codecs)
	})

	It("doesn't overwrite codecs that are set manually", func() {
		c := new(mockCodec)
		resp.SetCodec(c)
		Expect(resp.Codec()).To(Equal(c))
	})

	It("loads the matching codec if it hasn't been set", func() {
		Expect(resp.Codec()).To(Equal(mockJSON))
	})
})
