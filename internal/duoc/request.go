package duoc

type Request struct {
	client *Client
	auth   Credentials
}

func NewRequest(c *Client, auth Credentials) *Request {
	return &Request{
		client: c,
		auth:   auth,
	}
}
