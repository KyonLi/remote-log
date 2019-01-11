package command

type Server struct {
	ServerName string
	Hostname   string
	Port       int
	User       string
	Password   string
	PrivateKey string
	TailFile   string
}

type Config struct {
	TailFile string
	Servers  map[string]Server
}
