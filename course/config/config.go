package config

type DBConfig struct {
	Dbname string
	Course string
}

func New() DBConfig {
	return DBConfig{
		Dbname: "UCAM",
		Course: "course",
	}

}
