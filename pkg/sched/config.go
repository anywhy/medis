package sched

type Config struct {
	Master    string
	Name      string
	User      string
	Principal string
	Zk        string
}

func NewConfig(master string, name string, user string, principal string, zk string) *Config {
	return &Config{
		Master:    master,
		Name:      name,
		User:      user,
		Principal: principal,
		Zk:        zk,
	}
}

func (c *Config) GetMaster() string  {
	return c.Master
}

func (c *Config) GetName() string  {
	return c.Name
}

func (c *Config) GetUser() string  {
	return c.User
}

func (c *Config) GetPrincipal() string  {
	return c.Principal
}

func (c *Config) GetZk() string  {
	return c.Zk
}