package silverback

import (
	"sort"
	"strconv"
	"strings"
)

const defaultQuality = 1.0

type Options map[string]string

func (o Options) Add(key, value string) {
	o[strings.TrimSpace(key)] = strings.TrimSpace(value)
}

// MIMEType stores a full MIME type for Accept and Content-Type
// headers.
type MIMEType struct {
	Type    string
	SubType string
	Options Options
}

// ParseMIMEType parses a MIME type entry, such as from a Content-Type
// or Accept header.
//
// The acceptOptions value is not from any MIME type spec, but from
// RFC 2616 section 14.1, on the Accept header.  It states that the
// option "q" is not related to the MIME type, but is instead an
// Accept param; and that any options following that param should be
// considered Accept extensions.
//
// If this function is being used to parse MIME types that are not
// coming from the Accept header, the acceptOptions value can be
// ignored if there is no valid "q" option for any supported MIME
// types; otherwise, just add the acceptOptions key/value pairs to
// mime.Options.
func ParseMIMEType(value string) (mime MIMEType, acceptOptions Options) {
	typeEnd := strings.IndexRune(value, '/')
	if typeEnd == -1 {
		return MIMEType{}, nil
	}
	mime.Type = value[:typeEnd]
	value = value[typeEnd+1:]
	subTypeEnd := strings.IndexRune(value, ';')
	if subTypeEnd == -1 {
		mime.SubType = value
		return mime, nil
	}
	mime.SubType = value[:subTypeEnd]
	mime.Options = make(Options)
	start := subTypeEnd + 1
	options := strings.FieldsFunc(value[start:], isOptionSplit)
	qFound := false
	for _, option := range options {
		nameEnd := strings.IndexRune(option, '=')
		if nameEnd == -1 {
			mime.Options.Add(option, "")
			continue
		}
		// Note: the 'q' option is required to have a value, so we
		// don't need to check for it until we know that this option
		// has an '='.
		name, value := option[:nameEnd], option[nameEnd+1:]
		if !qFound && strings.TrimSpace(name) == "q" {
			acceptOptions = make(Options)
			qFound = true
		}
		if qFound {
			acceptOptions.Add(name, value)
			continue
		}
		mime.Options.Add(option[:nameEnd], option[nameEnd+1:])
	}
	return mime, acceptOptions
}

// AcceptEntry stores a single entry in an Accept header.
type AcceptEntry struct {
	MIMEType
	AcceptOptions map[string]string

	// Used for caching the quality value read from AcceptOptions,
	// since parsing the float value from a string every time isn't
	// cheap.
	quality float32
}

func ParseAcceptEntry(value string) *AcceptEntry {
	mime, acceptOptions := ParseMIMEType(value)
	return &AcceptEntry{
		MIMEType:      mime,
		AcceptOptions: acceptOptions,
	}
}

func (entry *AcceptEntry) Quality() float32 {
	if entry.quality != 0 {
		return entry.quality
	}
	q, ok := entry.AcceptOptions["q"]
	if !ok {
		entry.quality = defaultQuality
		return entry.quality
	}
	quality, err := strconv.ParseFloat(q, 32)
	if err != nil {
		entry.quality = defaultQuality
		return entry.quality
	}
	entry.quality = float32(quality)
	return entry.quality
}

func (entry *AcceptEntry) Wildcards() int {
	if entry.Type == "*" {
		// Must be */*, since */sub-type is not valid.
		return 2
	}
	if entry.SubType == "*" {
		return 1
	}
	return 0
}

func (entry *AcceptEntry) match(codec Codec) bool {
	if entry.MIMEType.Type == "*" {
		return true
	}
	codecTypes := codec.Types()
	for _, supported := range codecTypes {
		if supported.Type != entry.MIMEType.Type {
			continue
		}
		if entry.MIMEType.SubType == "*" || entry.MIMEType.SubType == supported.SubType {
			return true
		}
	}
	return false
}

func (entry *AcceptEntry) bestCodec(codecs []Codec) Codec {
	for _, codec := range codecs {
		if entry.match(codec) {
			return codec
		}
	}
	return nil
}

// Accept stores all values in an Accept header.
type Accept []*AcceptEntry

// Len returns the length of accept.
func (accept Accept) Len() int {
	return len(accept)
}

// Less returns whether or not accept[i] should be earlier in a sorted
// list than accept[j] (i.e. accept[i] should be preferred over
// accept[j]), according to RFC 2616 section 14.1.
func (accept Accept) Less(i, j int) bool {
	if accept[i].Quality() != accept[j].Quality() {
		return accept[i].Quality() > accept[j].Quality()
	}
	return accept[i].Wildcards() < accept[j].Wildcards()
}

// Swap swaps accept[i] and accept[j].
func (accept Accept) Swap(i, j int) {
	accept[i], accept[j] = accept[j], accept[i]
}

// runeSplit performs strings.Split, but with a rune instead of a
// string.
func runeSplit(value string, split rune) []string {
	values := make([]string, 0, 10)
	start := 0

	for end := strings.IndexRune(value, split); end < len(value); end = start + strings.IndexRune(value[start:], split) {
		if end == -1 {
			end = len(value)
		}
		values = append(values, value[start:end])
		start = end + 1
	}
	return values
}

// ParseAcceptHeader parses an Accept header value into a *Accept.
func ParseAcceptHeader(acceptHeader string) Accept {
	if acceptHeader == "" {
		return nil
	}
	entries := strings.FieldsFunc(acceptHeader, isAcceptSplit)
	accept := make(Accept, 0, len(entries))
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		accept = append(accept, ParseAcceptEntry(entry))
	}
	sort.Sort(accept)
	return accept
}

// Codec returns the best codec in codecs for this accept header.  It
// automatically sorts the entries based on RFC2616 section 14.1 prior
// to ranging through them, to ensure the codec it loads is optimal
// for the Accept header.
func (accept Accept) Codec(codecs []Codec) Codec {
	if len(codecs) > 0 {
		return accept.bestCodec(codecs)
	}
	return nil
}

func (accept Accept) bestCodec(codecs []Codec) Codec {
	for _, entry := range accept {
		codec := entry.bestCodec(codecs)
		if codec != nil {
			return codec
		}
	}
	return nil
}

func isOptionSplit(r rune) bool {
	return r == ';'
}

func isAcceptSplit(r rune) bool {
	return r == ','
}
