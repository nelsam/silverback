package codecs_test

import (
	"encoding/json"
	"net/http"

	"github.com/nelsam/silverback"
	"github.com/nelsam/silverback/codecs"

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
			codec = &codecs.JSON{}
			var err error
			req, err = http.NewRequest("GET", "foo", nil)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns a json codec when a new copy is requested", func() {
			newCodec := codec.New(map[string]string{})
			_, isJSON := newCodec.(*codecs.JSON)
			Expect(isJSON).To(BeTrue())
		})

		It("supports application/json and text/json MIME types", func() {
			appJSON, _ := silverback.ParseMIMEType("application/json")
			textJSON, _ := silverback.ParseMIMEType("text/json")
			Expect(codec.Types()).To(ConsistOf(appJSON, textJSON))
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
