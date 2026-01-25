package alerts

type Config struct {
	SMTPHost string
	SMTPPort string
	Username string
	Password string
	From     string
	To       string
}

var cfg Config

func Init(c Config) {
	cfg = c
}
