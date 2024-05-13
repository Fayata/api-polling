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
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func init() {
	Load = config{
		WebServer: webServer{
			Host: "0.0.0.0",
			Port: "9000",
		},
		Database: database{
			Host:     "127.19.0.1",
			Port:     "3307",
			User:     "root",
			Password: "",
			Name:     "db_api",
		},
	}
}
