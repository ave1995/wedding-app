package factory

import (
	"wedding-app/domain/event"
	"wedding-app/ws"
)

func (f *Factory) Hub() *ws.Hub {
	f.hubOnce.Do(func() {
		f.hub = ws.NewHub(f.Logger())
	})
	return f.hub
}

func (f *Factory) EventPublisher() event.EventPublisher {
	f.publisherOnce.Do(func() {
		f.publisher = ws.NewPublisher(f.Hub())
	})
	return f.publisher
}
