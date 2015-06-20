package silverback_test

import (
	"encoding/json"
	"net/http"

	"github.com/nelsam/silverback"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Codecs", func() {
	var (
		codec silverback.Codec
		req   *http.Request
	)

	Context("json", func() {
		BeforeEach(func() {
			codec = &silverback.JSON{}
			var err error
			req, err = http.NewRequest("GET", "foo", nil)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns itself when a new copy is requested", func() {
			Expect(codec.New(map[string]string{})).To(Equal(codec))
		})

		It("matches a request for application/json", func() {
			req.Header.Set("Accept", "application/json")
			Expect(codec.Match(req)).To(BeTrue())
		})

		It("matches a request for text/json", func() {
			req.Header.Set("Accept", "text/json")
			Expect(codec.Match(req)).To(BeTrue())
		})

		It("marshals to proper json", func() {
			val := map[string]interface{}{
				"foo": "bar",
				"baz": 3.0,
			}
			expected, err := json.Marshal(val)
			Expect(err).ToNot(HaveOccurred())
			Expect(codec.Marshal(val)).To(MatchJSON(expected))
		})

		It("unmarshals from proper json", func() {
			raw := []byte(`{"foo":"bar","baz":3.0}`)
			var expected, actual map[string]interface{}
			Expect(json.Unmarshal(raw, &expected)).ToNot(HaveOccurred())
			Expect(codec.Unmarshal(raw, &actual)).ToNot(HaveOccurred())
			Expect(actual).To(BeEquivalentTo(expected))
		})
	})
})
