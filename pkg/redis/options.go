package redis

// Option -.
type Option func(client *Client)

// Addr -.
func Addr(addr string) Option {
	return func(c *Client) {
		c.opts.Addr = addr
	}
}

// Db -.
func Db(db int) Option {
	return func(c *Client) {
		c.opts.DB = db
	}
}
