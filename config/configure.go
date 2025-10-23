package config

// Configure is an interface that defines methods for setting default values and validating configuration.
type Configure interface {
	SetDefaults()
	Validate() error
}
