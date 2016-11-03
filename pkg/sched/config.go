package sched

type Config struct {
	master    string
	Name      string
	User      string
	Principal string
	Checkpoint bool
	WebuiUrl string
	Role string
	FailoverTimeout float64
	Secret string
	AuthProvider string
	Address string
}

func NewConfig(master string) *Config {
	return &Config{
		master:    master,
	}
}

func (c *Config) GetMaster() string  {
	return c.master
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

func (c *Config) GetCheckpoint() bool  {
	return c.Checkpoint
}

func (c *Config) GetWebuiUrl() string  {
	return c.WebuiUrl
}

func (c *Config) GetRole() string  {
	return c.Role
}

func (c *Config) GetFailoverTimeout() float64  {
	return c.FailoverTimeout
}

func (c *Config) GetSecret() string  {
	return c.Secret
}

func (c *Config) GetAuthProvider() string  {
	return c.AuthProvider
}

func (c *Config) GetAddress() string  {
	return c.Address
}