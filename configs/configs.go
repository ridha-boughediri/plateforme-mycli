package configs

type AliasConfig struct {
	Alias    string `json:"alias"`
	Remote   string `json:"remote"`
	Username string `json:"username"`
	Password string `json:"password"`
}
