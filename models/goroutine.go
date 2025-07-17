package models

type Goroutine struct {
	ID    int64
	Input chan any
	Stop  chan struct{}
}

func StartGoroutine(id int64, handler func(any)) *Goroutine {
	g := &Goroutine{
		ID:    id,
		Input: make(chan any),
		Stop:  make(chan struct{}),
	}

	go func() {
		for {
			select {
			case data := <-g.Input:
				handler(data)
			case <-g.Stop:
				return
			}
		}
	}()

	return g
}
