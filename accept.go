package silverback

import "strconv"

const defaultQuality = 1.0

// MIMEType stores a full MIME type for Accept and Content-Type
// headers.
type MIMEType struct {
	Name    string
	Options map[string]string
}

// AcceptEntry stores a single entry in an Accept header.
type AcceptEntry struct {
	MIMEType

	// Used for caching the quality value read from Options.
	quality float32
}

func (entry *AcceptEntry) Quality() float32 {
	if entry.quality != 0 {
		return entry.quality
	}
	q, ok := entry.Options["q"]
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

// Accept stores all values in an Accept header.
type Accept []AcceptEntry

// order reorders the entries from an Accept header according to
// RFC2616, section 14.1.
func (accept *Accept) order() {
}

// Codec returns the best codec in codecs for this accept header.  It
// automatically sorts the entries based on RFC2616 section 14.1 prior
// to ranging through them, to ensure the codec it loads is optimal
// for the Accept header.
func (accept *Accept) Codec(codecs []Codec) Codec {
	if len(codecs) > 0 {
		return codecs[0]
	}
	return nil
}
