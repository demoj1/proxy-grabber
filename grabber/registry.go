package grabber

import "proxy_grabber/grabber/sites"

var Registry *registry = initRegistry()

type registry struct {
	grabbers map[string]Grabber
}

func initRegistry() *registry {
	return registry{}.Add(
		"fresh", sites.NewFreshProxy(),
	).Add(
		"hidemy", sites.NewHidemy(),
	).Add(
		"multiproxy", sites.NewMultiProxy(),
	).Add(
		"primespeed", sites.NewPrimeSpeed(),
	).Add(
		"therealist", sites.NewThereAList(),
	)
}

func (r *registry) Add(name string, grabber Grabber) *registry {
	r.grabbers[name] = grabber
	return r
}

func (r *registry) Delete(name string) *registry {
	delete(r.grabbers, name)
	return r
}

func (r *registry) Grab(proxyType ProxyType) (chan string, error) {
	var chains []chan string

	for _, grabber := range r.grabbers {
		channel, err := grabber.Grab(proxyType)
		if err != nil {
			return nil, err
		}

		chains = append(chains, channel)
	}

	return Merge(chains...), nil
}
