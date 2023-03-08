package config

type Auth struct {
	AccessToken  string `mapstructure:"AccessToken"`
	RefreshToken string `mapstructure:"RefreshToken"`
}
