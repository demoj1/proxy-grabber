package grabber

type Grabber interface {
	Grab(ProxyType) (chan string, error)
}
