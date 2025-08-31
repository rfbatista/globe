package globe

type Config struct {
	Name string
}

var config Config

func init() {
	config.Name = ""
}
