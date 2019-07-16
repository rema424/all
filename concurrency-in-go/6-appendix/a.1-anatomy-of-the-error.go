package main

func anatomy() {
	waitForever := make(chan interface{})
	go func() {
		panic("test panic")
	}()
	<-waitForever
}
