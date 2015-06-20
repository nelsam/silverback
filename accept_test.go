package silverback_test

import (
	"github.com/nelsam/silverback"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Accept", func() {
	Context("AcceptEntry", func() {
		var (
			options map[string]string
			entry   *silverback.AcceptEntry
		)

		JustBeforeEach(func() {
			mimeType := silverback.MIMEType{
				Name:    "foo",
				Options: options,
			}
			entry = &silverback.AcceptEntry{MIMEType: mimeType}
		})

		Context("Defaults", func() {
			BeforeEach(func() {
				options = map[string]string{}
			})

			It("defaults to a quality of 1.0", func() {
				Expect(entry.Quality()).To(BeEquivalentTo(1.0))
			})
		})

		Context("Invalid", func() {
			BeforeEach(func() {
				options = map[string]string{"q": "invalid"}
			})

			It("falls back to the default quality", func() {
				Expect(entry.Quality()).To(BeEquivalentTo(1.0))
			})
		})

		Context("Options", func() {
			BeforeEach(func() {
				options = map[string]string{"q": "0.8"}
			})

			It("parses the q option for quality", func() {
				// BeEquivalentTo fails here because a float32(0.8)
				// ends up with a value *slightly* larger than 0.8.
				// So instead, we use the actual float32 value.
				Expect(entry.Quality()).To(Equal(float32(0.8)))
			})
		})
	})

	Context("Accept", func() {
		var (
			accept *silverback.Accept
			codecs []silverback.Codec
		)

		Context("Empty Codecs", func() {
			BeforeEach(func() {
				accept = &silverback.Accept{silverback.AcceptEntry{}}
				codecs = []silverback.Codec{}
			})

			It("returns a nil codec if codecs is empty", func() {
				Expect(accept.Codec(codecs)).To(BeNil())
			})
		})

		Context("Empty Accept Header", func() {
			BeforeEach(func() {
				accept = &silverback.Accept{}
				codecs = []silverback.Codec{&mockCodec{}, &mockCodec{}}
			})

			It("returns a codec if the accept header is empty, but codecs is not", func() {
				Expect(accept.Codec(codecs)).ToNot(BeNil())
			})
		})

		PContext("Non-Empty", func() {
			BeforeEach(func() {

			})
		})
	})
})
