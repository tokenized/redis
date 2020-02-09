package redis

import (
	"log"

	"github.com/gomodule/redigo/redis"
	redigo "github.com/gomodule/redigo/redis"
)

// Subscriber subscribes to a Redis key, and writes received messages to a
// channel.
type Subscriber struct {
	pool   *redis.Pool
	key    string
	output chan []byte
}

// NewSubscriber creates a Subscriber that will use the provided
// redis.Pool to subscribe to the given key.
//
// Any messages that are received by the Subscriber will be written to the
// channel.
func NewSubscriber(pool *redigo.Pool, key string, out chan []byte) Subscriber {
	return Subscriber{
		pool:   pool,
		key:    key,
		output: out,
	}
}

// Run receives pubsub messages from Redis after establishing a connection.
//
// When a valid message is received it is written to the channel.
func (s *Subscriber) Run() error {
	conn := s.pool.Get()
	defer conn.Close()

	psc := redigo.PubSubConn{conn}
	if err := psc.Subscribe(s.key); err != nil {
		return err
	}

	for {
		switch v := psc.Receive().(type) {
		case redigo.Message:
			// log.Printf("incoming [%v] message=%v",
			// 	v.Channel,
			// 	string(v.Data))

			s.output <- v.Data

		case redigo.Subscription:
			log.Printf("Redis subscription received : channel=%v : kind=%v : count=%v", v.Channel, v.Kind, v.Count)

		case error:
			log.Printf("Error while subscribed to Redis channel %s : %v",
				s.key,
				v)
		default:
			log.Println("Unknown Redis receive during subscription : v=%v", v)
		}
	}

	return nil
}
