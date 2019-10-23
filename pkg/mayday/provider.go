package mayday

// ClientProvider - s
type ClientProvider interface {
	Set(Client)
	Get() Client
}

type clientProvider struct {
	client Client
}

// NewClientProvider -
func NewClientProvider() ClientProvider {
	return &clientProvider{}
}

func (c *clientProvider) Set(cc Client) {
	c.client = cc
}
func (c *clientProvider) Get() Client {
	return c.client
}
