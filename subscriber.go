package redis

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	redigo "github.com/gomodule/redigo/redis"
)

// Subscriber subscribes to a Redis key, and writes receive messages to a channel.
type Subscriber struct {
	pool   *redis.Pool
	key    string
	output chan []byte
	Logger Logger
}

// NewSubscriber creates a Subscriber that will use the provided redis.Pool to subscribe to the
// given key.
//
// Any messages that are received by the Subscriber will be written to the channel.
func NewSubscriber(logger Logger, pool *redigo.Pool, key string, out chan []byte) Subscriber {
	return Subscriber{
		pool:   pool,
		key:    key,
		output: out,
		Logger: logger,
	}
}

// Run receives pubsub messages from Redis after establishing a connection.
//
// When a valid message is received it is written to the channel.
func (s *Subscriber) Run(ctx context.Context) error {
	conn := s.pool.Get()
	defer conn.Close()

	psc := redigo.PubSubConn{conn}
	if err := psc.Subscribe(s.key); err != nil {
		return err
	}

	for {
		switch v := psc.Receive().(type) {
		case redigo.Message:
			s.output <- v.Data

		case redigo.Subscription:
			s.Logger.Info(ctx, "Subscription channel=%v kind=%v count=%v", v.Channel, v.Kind, v.Count)

		case error:
			s.Logger.Error(ctx, "Error while subscribed to channel=%v : %v", s.key, v)

			// back off a little
			time.Sleep(2 * time.Second)

		default:
			s.Logger.Error(ctx, "Unknown Redis message : %v", v)
		}
	}

	return nil
}
