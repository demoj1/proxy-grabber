package grabber

type Grabber interface {
	Grab() (error, []Proxy)
}
