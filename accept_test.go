package silverback_test

import (
	"net/http"
	"sort"

	"github.com/nelsam/silverback"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Accept", func() {
	Context("MIMEType", func() {
		var (
			mimeString    string
			mime          silverback.MIMEType
			acceptOptions map[string]string
		)

		JustBeforeEach(func() {
			mime, acceptOptions = silverback.ParseMIMEType(mimeString)
		})

		Context("No Params", func() {
			BeforeEach(func() {
				mimeString = "application/json"
			})

			It("parses a MIME type without params", func() {
				Expect(mime.Type).To(Equal("application"))
				Expect(mime.SubType).To(Equal("json"))
				Expect(mime.Options).To(BeEmpty())
				Expect(acceptOptions).To(BeEmpty())
			})
		})

		Context("MIME Params", func() {
			var expectedOptions = map[string]string{
				"foo": "bar",
				"baz": "",
			}
			BeforeEach(func() {
				mimeString = "application/json; foo=bar; baz"
			})

			It("parses a MIME type with params", func() {
				Expect(mime.Type).To(Equal("application"))
				Expect(mime.SubType).To(Equal("json"))
				Expect(mime.Options).To(BeEquivalentTo(expectedOptions))
				Expect(acceptOptions).To(BeEmpty())
			})
		})

		Context("Quality Param", func() {
			BeforeEach(func() {
				mimeString = "application/json; q=0.1"
			})

			It("parses a quality param without any other params", func() {
				Expect(mime.Type).To(Equal("application"))
				Expect(mime.SubType).To(Equal("json"))
				Expect(mime.Options).To(BeEmpty())
				Expect(acceptOptions).To(HaveKeyWithValue("q", "0.1"))
			})
		})

		Context("MIME and Accept params", func() {
			var (
				expectedMIMEOptions = map[string]string{
					"foo": "bar",
				}
				expectedAcceptOptions = map[string]string{
					"q":     "0.2",
					"bacon": "eggs",
				}
			)
			BeforeEach(func() {
				mimeString = "application/json; foo=bar; q=0.2; bacon=eggs"
			})

			It("parses MIME and Accept params separately", func() {
				Expect(mime.Type).To(Equal("application"))
				Expect(mime.SubType).To(Equal("json"))
				Expect(mime.Options).To(BeEquivalentTo(expectedMIMEOptions))
				Expect(acceptOptions).To(BeEquivalentTo(expectedAcceptOptions))
			})
		})
	})

	Context("AcceptEntry Parsing", func() {
		It("parses accept-params", func() {
			entryString := "application/json; foo=bar; q=0.1; bacon=eggs"
			expectedOptions := map[string]string{"foo": "bar"}
			expectedAcceptOptions := map[string]string{
				"q":     "0.1",
				"bacon": "eggs",
			}

			entry := silverback.ParseAcceptEntry(entryString)
			Expect(entry.Type).To(Equal("application"))
			Expect(entry.SubType).To(Equal("json"))
			Expect(entry.Options).To(BeEquivalentTo(expectedOptions))
			Expect(entry.AcceptOptions).To(BeEquivalentTo(expectedAcceptOptions))
			Expect(entry.Quality()).To(Equal(float32(0.1)))
		})
	})

	Context("AcceptEntry Quality", func() {
		var (
			options map[string]string
			entry   *silverback.AcceptEntry
		)

		JustBeforeEach(func() {
			mimeType := silverback.MIMEType{
				Type: "foo",
			}
			entry = &silverback.AcceptEntry{
				MIMEType:      mimeType,
				AcceptOptions: options,
			}
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

	Context("Accept Header Parsing", func() {
		var (
			accept  silverback.Accept
			headers http.Header
		)

		JustBeforeEach(func() {
			accept = silverback.ParseAcceptHeader(headers)
		})

		Context("Empty Header", func() {
			BeforeEach(func() {
				headers = http.Header{}
			})

			It("returns an empty slice", func() {
				Expect(accept).To(BeEmpty())
			})
		})

		Context("Single Header Entry", func() {
			BeforeEach(func() {
				headers = header("Accept", "application/json")
			})

			It("returns a single entry", func() {
				Expect(accept).To(HaveLen(1))
				Expect(accept[0].Type).To(Equal("application"))
				Expect(accept[0].SubType).To(Equal("json"))
			})
		})

		Context("Multiple Simple Entries", func() {
			var expected silverback.Accept

			BeforeEach(func() {
				headers = header("Accept", "application/json, text/xml")
				jsonEntry := silverback.ParseAcceptEntry("application/json")
				xmlEntry := silverback.ParseAcceptEntry("text/xml")
				expected = silverback.Accept{jsonEntry, xmlEntry}
				// sort them to ensure all operations done during
				// sorting are done on the expected entries, too.
				sort.Sort(expected)
			})

			It("returns multiple entries", func() {
				Expect(accept).To(HaveLen(2))
				Expect(accept).To(BeEquivalentTo(expected))
			})
		})

		Context("Multiple Headers", func() {
			var expected silverback.Accept

			BeforeEach(func() {
				headers = header("Accept", "application/json, text/xml", "text/html")
				jsonEntry := silverback.ParseAcceptEntry("application/json")
				xmlEntry := silverback.ParseAcceptEntry("text/xml")
				htmlEntry := silverback.ParseAcceptEntry("text/html")
				expected = silverback.Accept{jsonEntry, xmlEntry, htmlEntry}
				// again, we need to ensure all values that are
				// assigned during sorting are assigned here.
				sort.Sort(expected)
			})

			It("handles multiple header entries", func() {
				Expect(accept).To(HaveLen(3))
				Expect(accept).To(BeEquivalentTo(expected))
			})
		})

		Context("Sorting", func() {
			var orderedNames = []string{
				"application/json", // default quality, 0.1
				"text/xml",         // quality 0.9
				"text/*",           // quality 0.9, less specific
				"text/html",        // quality 0.8
				"image/jpeg",       // quality 0.3
				"image/*",          // quality 0.3, less specific
				"*/*",              // quality 0.3, least specific
			}

			BeforeEach(func() {
				// Avoid any ambiguous sorting - nothing at the same
				// quality with the same specificity.
				headers = header("Accept", "image/*; q=0.3, text/html; q=0.8, image/jpeg; q=0.3, application/json, text/xml; q=0.9, text/*; q=0.9, */*; q=0.3")
			})

			It("sorts the accept header according to RFC 2616", func() {
				Expect(accept).To(HaveLen(len(orderedNames)))
				for index, name := range orderedNames {
					expectedMIME, _ := silverback.ParseMIMEType(name)
					actual := accept[index]
					Expect(actual.Type).To(Equal(expectedMIME.Type))
					Expect(actual.SubType).To(Equal(expectedMIME.SubType))
				}
			})
		})
	})

	Context("Accept Codec", func() {
		var (
			accept silverback.Accept
			codecs []silverback.Codec
		)

		Context("Empty Codecs", func() {
			BeforeEach(func() {
				accept = silverback.Accept{&silverback.AcceptEntry{}}
				codecs = []silverback.Codec{}
			})

			It("returns a nil codec if codecs is empty", func() {
				Expect(accept.Codec(codecs)).To(BeNil())
			})
		})

		Context("Empty Accept Header", func() {
			BeforeEach(func() {
				accept = silverback.Accept{}
				codecs = []silverback.Codec{&mockCodec{}, &mockCodec{}}
			})

			It("returns nil if the accept header is empty", func() {
				Expect(accept.Codec(codecs)).To(BeNil())
			})
		})

		Context("Direct Match", func() {
			var (
				mockJSON      = makeCodec("application", "json")
				mockXML       = makeCodec("text", "xml")
				mockGibberish = makeCodec("gibberish", "doomsday")
			)

			BeforeEach(func() {
				accept = silverback.Accept{
					silverback.ParseAcceptEntry("application/json"),
					silverback.ParseAcceptEntry("text/xml"),
				}
			})

			It("returns the first matching codec", func() {
				codecs = []silverback.Codec{mockXML, mockJSON}
				Expect(accept.Codec(codecs)).To(Equal(mockJSON))

				codecs = codecs[:1]
				Expect(accept.Codec(codecs)).To(Equal(mockXML))
			})

			It("returns nil if there is no matching codec", func() {
				codecs = []silverback.Codec{mockGibberish}
				Expect(accept.Codec(codecs)).To(BeNil())
			})
		})

		Context("Wildcard Match", func() {
			var (
				mockJSON = makeCodec("application", "json")
			)

			BeforeEach(func() {
				accept = silverback.Accept{
					silverback.ParseAcceptEntry("application/*"),
				}
			})

			It("returns codecs matching wildcard values", func() {
				Expect(accept.Codec([]silverback.Codec{mockJSON})).To(Equal(mockJSON))
			})
		})
	})
})

func header(key string, values ...string) http.Header {
	return http.Header{
		key: values,
	}
}
