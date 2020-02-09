package redis

import (
	"net/url"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

func NewPool(uri string) (*redigo.Pool, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	var password string
	if u.User != nil {
		password, _ = u.User.Password()
	}

	return &redigo.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", u.Host)
			if err != nil {
				return nil, err
			}

			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}, nil
}
