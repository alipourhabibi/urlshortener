package config

type Posgres struct {
	Host     string `mapstructure:"Host"`
	Port     string `mapstructure:"Port"`
	DBName   string `mapstructure:"DBName"`
	User     string `mapstructure:"User"`
	Password string `mapstructure:"Password"`
}
