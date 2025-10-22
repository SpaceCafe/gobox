package types

// Metadata contains metadata information about an application.
type Metadata struct {

	// AppName is the name of the application.
	AppName string `json:"app_name" yaml:"app_name"`

	// Slug is a unique identifier for the application.
	Slug string `json:"slug" yaml:"slug"`

	// Version represents the version of the application.
	Version string `json:"version" yaml:"version"`

	// Author is the name of the person or entity who created or owns the application.
	Author string `json:"author" yaml:"author"`

	// Description is a brief summary or overview of the application.
	Description string `json:"description" yaml:"description"`

	// License is the license type of the application, if any.
	License string `json:"license" yaml:"license"`

	// URL is a link associated with the application. It could reference the homepage, repository, or documentation.
	URL string `json:"url" yaml:"url"`

	// Keywords are used to categorize or tag the application for searchability.
	Keywords []string `json:"keywords" yaml:"keywords"`
}

// MetadataProvider provides access to metadata information about an application.
type MetadataProvider interface {
	Metadata() *Metadata
}
