package redis

// Sender is an interface that is compatible with the redis.Send(..) function.
type Sender interface {
	Send(string, ...interface{}) error
}
