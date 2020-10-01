package models

type MimeFormat struct {
	Source       string   `json:"source,omitempty"`
	Compressible bool     `json:"compressible,omitempty"`
	Extensions   []string `json:"extensions"`
	Charset      string   `json:"charset,omitempty"`
}
