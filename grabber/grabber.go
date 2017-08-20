package grabber

type ADDR_TYPE int

const (
	HTTP ADDR_TYPE = iota
	HTTPS
)

type Grabber interface {
	Grab(ADDR_TYPE) (chan string, error)
}
