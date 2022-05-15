package bot

// NewEventCollector returns a channel in which the events of type T gets sent which pass the passed filter and a function which can be used to stop the event collector.
// The close function needs to be called to stop the event collector.
func NewEventCollector[E Event](disgo Client, filterFunc func(e E) bool) (<-chan E, func()) {
	ch := make(chan E)

	coll := &collector[E]{
		FilterFunc: filterFunc,
		Chan:       ch,
	}
	disgo.EventManager().AddEventListeners(coll)

	return ch, func() {
		disgo.EventManager().RemoveEventListeners(coll)
		close(ch)
	}
}

type collector[E Event] struct {
	FilterFunc func(e E) bool
	Chan       chan<- E
}

func (c *collector[E]) OnEvent(e Event) {
	if event, ok := e.(E); ok {
		if !c.FilterFunc(event) {
			return
		}
		c.Chan <- event
	}
}