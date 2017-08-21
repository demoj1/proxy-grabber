package grabber

type ProxyType int

func (p ProxyType) String() string {
	switch p {
	case HTTP:
		return "HTTP"
	case HTTPS:
		return "HTTPS"
	}

	return "NONE"
}

const (
	HTTP ProxyType = iota
	HTTPS
)
