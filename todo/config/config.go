package config

// Config struct def
type Config struct {
	Database string
	TodoColl string
}

// New returns a new config
func New() Config {
	return Config{
		Database: "TODO",
		TodoColl: "todo",
	}
}
