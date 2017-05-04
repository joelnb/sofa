package sofa

// Attachment represents files or other data which is attached to documents in the
// Database.
type Attachment struct {
	ContentType   string  `json:"content_type,omitempty"`
	Data          string  `json:"data,omitempty"`
	Digest        string  `json:"digest,omitempty"`
	EncodedLength float64 `json:"encoded_length,omitempty"`
	Encoding      string  `json:"encoding,omitempty"`
	Length        int64   `json:"length,omitempty"`
	RevPos        float64 `json:"revpos,omitempty"`
	Stub          bool    `json:"stub,omitempty"`
	Follows       bool    `json:"follows,omitempty"`
}
