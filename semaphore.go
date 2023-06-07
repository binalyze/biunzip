package main

type semphore chan struct{}

func newSemaphore(size int) semphore {
	return make(chan struct{}, size)
}

func (sem semphore) acquire() {
	sem <- struct{}{}
}

func (sem semphore) release() {
	<-sem
}

func (sem semphore) releaseAll() {
	for i := 0; i < cap(sem); i++ {
		sem.release()
	}
}

func (sem semphore) close() {
	close(sem)
}
