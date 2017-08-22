package grabber

var Registry *registry = &registry{
	make(map[string]Grabber)}

type registry struct {
	grabbers map[string]Grabber
}

func (r *registry) Add(name string, grabber Grabber) *registry {
	r.grabbers[name] = grabber
	return r
}

func (r *registry) Delete(name string) *registry {
	delete(r.grabbers, name)
	return r
}

func (r *registry) Grab(proxyType ProxyType) error {
	for _, grabber := range r.grabbers {
		err := grabber.Grab(proxyType)
		if err != nil {
			return err
		}
	}

	return nil
}
