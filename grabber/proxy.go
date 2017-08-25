package grabber

const (
	HTTP ProxyType = iota
	HTTPS
)

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

func StrToProxyType(s string) ProxyType {
	switch s {
	case "HTTP":
		return HTTP
	case "HTTPS":
		return HTTPS
	}

	return HTTP
}

type Proxy struct {
	Type    ProxyType `json:"type"`
	Address string    `json:"address"`
	Alive   bool      `json:"alive"`
}
