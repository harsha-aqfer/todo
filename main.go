package main

type Config struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Database string `json:"database"`
	Host     string `json:"host"`
}

func NewConfig() *Config {
	return &Config{}
}

type Service struct {
}

func main() {

}
