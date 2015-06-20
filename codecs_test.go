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
		mime  silverback.MIMEType
	)

	Context("json", func() {
		BeforeEach(func() {
			codec = &silverback.JSON{}
			var err error
			req, err = http.NewRequest("GET", "foo", nil)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns a json codec when a new copy is requested", func() {
			newCodec := codec.New(map[string]string{})
			_, isJSON := newCodec.(*silverback.JSON)
			Expect(isJSON).To(BeTrue())
		})

		It("matches an application/json MIME type", func() {
			mime.Name = "application/json"
			Expect(codec.Match(mime)).To(BeTrue())
		})

		It("matches a text/json MIME type", func() {
			mime.Name = "text/json"
			Expect(codec.Match(mime)).To(BeTrue())
		})

		It("doesn't match a text/xml MIME type", func() {
			mime.Name = "text/xml"
			Expect(codec.Match(mime)).To(BeFalse())
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

		It("errors from improper json", func() {
			raw := []byte(`{"foo":"bar","baz":3.0,}`)
			m := map[string]string{}
			Expect(json.Unmarshal(raw, &m)).To(HaveOccurred())
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
