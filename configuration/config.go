package configuration

type Configuration struct {
	Database `json:"Database"`
}

type Database struct {
	MongoDB `json:"MongoDB"`
}

type MongoDB struct {
	LocalPort string `json:"LocalPort"`
	LocalURL  string `json:"LocalURL"`
	Name      string `json:"Name"`
}
