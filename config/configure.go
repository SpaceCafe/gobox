package config

// Configure is an interface that defines methods for setting default values and validating configuration.
type Configure interface {

	// SetDefaults initializes the default values for the relevant fields in the struct.
	SetDefaults()

	// Validate ensures the all necessary configurations are filled and within valid confines.
	Validate() error
}
