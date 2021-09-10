package redis

import (
	"math"
	"net/url"
	"strconv"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

type BorrowFunc func(c redigo.Conn, t time.Time) error

func NewPool(uri string) (*redigo.Pool, error) {
	return NewPoolWithBorrowFunc(uri, PingOnBorrow)
}

func NewPoolWithBorrowFunc(uri string, f BorrowFunc) (*redigo.Pool, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	var password string
	if u.User != nil {
		password, _ = u.User.Password()
	}

	// db ?
	db, _ := strconv.ParseInt(u.Query().Get("db"), 10, 64)

	return &redigo.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", u.Host)
			if err != nil {
				return nil, err
			}

			if db > 0 {
				if _, err := c.Do("SELECT", db); err != nil {
					return nil, err
				}
			}

			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: f,
	}, nil
}

func NoopOnBorrow(c redigo.Conn, t time.Time) error {
	return nil
}

func PingOnBorrow(c redigo.Conn, t time.Time) error {
	_, err := c.Do("PING")
	return err
}

func NewSamplingBorrow(mod float64) BorrowFunc {
	return func(c redigo.Conn, t time.Time) error {
		if math.Mod(float64(t.Unix()), mod) != 0 {
			return nil
		}

		return PingOnBorrow(c, t)
	}
}
