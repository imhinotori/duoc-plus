package duoc

type Credentials interface {
	JWT() string
	Request() *Request
}

type BearerCredentials struct {
	jwt    string
	client *Client
}

func NewBearerCredentials(c *Client, jwt string) *BearerCredentials {
	return &BearerCredentials{
		jwt:    jwt,
		client: c,
	}
}

func (c *BearerCredentials) JWT() string {
	return c.jwt
}

func (c *BearerCredentials) Request() *Request {
	return NewRequest(c.client, c)
}
