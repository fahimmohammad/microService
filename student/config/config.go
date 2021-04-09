package config

type DBConfig struct {
	Dbname    string
	Students  string
	Enrolment string
}

func New() DBConfig {
	return DBConfig{
		Dbname:    "UCAM",
		Students:  "student",
		Enrolment: "enrolment",
	}

}
