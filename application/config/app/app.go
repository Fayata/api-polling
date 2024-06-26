package app

var Load config

type config struct {
	WebServer webServer
	Database  database
}

type webServer struct {
	Host string
	Port string
}

type database struct {
	AutoMigrate bool
	Host        string
	Port        string
	User        string
	Password    string
	Name        string
}

func init() {
	Load = config{
		WebServer: webServer{
			Host: "0.0.0.0",
			Port: "8080",
		},
		Database: database{
			Host:     "172.19.0.1",
			Port:     "3307",
			User:     "root",
			Password: "root",
			Name:     "dev_interactive_poll",
		},
	}
}
