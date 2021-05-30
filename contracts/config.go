package contracts

// Environment is the data-type fo runtime environments.
type Environment int

// List of runtime environments.
const (
	Development Environment = iota
	Staging
	Production
)

// IConfigService provides access to configuration parameters.
type IConfigService interface {

	// IsDebug returns true if app is running in debug-mode.
	IsDebug() bool

	// IsProduction returns true if app is running in production environment.
	IsProduction() bool

	// Environment returns the current runtime environment.
	Environment() Environment

	// BaseURL returns the API base URL.
	BaseURL() string

	// WebURL returns the Web base URL.
	WebURL() string

	// StorageDir returns the directory path which files are stored.
	StorageDir() string

	// StaticDir returns the directory path which static files exist
	StaticDir() string

	// Int returns the configuration parameter key as an int.
	Int(key string) int

	// Bool returns the configuration parameter key as bool.
	Bool(key string) bool

	// String returns the configuration parameter key as string.
	String(key string) string

	// SetString changes the configuration parameter with a string value.
	SetString(key, value string)

	// Dump retrieves a collection of all configuration parameters.
	Dump() map[string]string
}
