package jwt

import (
	"encoding/base64"
)

// Secret represents a byte slice used for storing secret data, such as cryptographic keys.
type Secret []byte

// Bytes returns the byte slice representation of the Secret.
func (r *Secret) Bytes() []byte {
	return *r
}

// MarshalText converts the Secret into a base64-encoded byte slice suitable for use in text-based encodings.
// It returns the encoded byte slice and any error that occurs during the encoding process.
func (r *Secret) MarshalText() (text []byte, err error) {
	base64.StdEncoding.Encode(text, *r)

	return
}

// String returns the base64-encoded string representation of the Secret.
func (r *Secret) String() string {
	return base64.StdEncoding.EncodeToString(*r)
}

// UnmarshalText decodes the provided Base64-encoded text into the Secret. Returns an error if decoding fails.
func (r *Secret) UnmarshalText(text []byte) (err error) {
	_, err = base64.StdEncoding.Decode(*r, text)

	return
}
