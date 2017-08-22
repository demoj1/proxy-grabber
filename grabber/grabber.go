package grabber

type Grabber interface {
	Grab(ProxyType) error
}
