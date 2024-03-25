package app

type Config struct {
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	} `yaml:"postgres"`
	HttpServer struct {
		Port string `yaml:"port"`
	} `yaml:"httpServer"`
	Telegram struct {
		BotToken string `yaml:"botToken"`
	} `yaml:"telegram"`
}
