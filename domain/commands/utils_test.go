package commands

import "sync"

func waitGroupToChan(wg *sync.WaitGroup) chan (bool) {
	ch := make(chan bool)
	go func() {
		defer close(ch)
		wg.Wait()
		ch <- true
	}()
	return ch
}
