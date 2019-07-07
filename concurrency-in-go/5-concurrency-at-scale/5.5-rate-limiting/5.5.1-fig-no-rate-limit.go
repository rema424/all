package main

import (
	"context"
	"log"
	"os"
	"sync"
)

func execNoRateLimit() {
	defer log.Printf("Done.")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := Open1()
	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if err := apiConnection.ReadFile(context.Background()); err != nil {
				log.Printf("cannot ReadFile: %v", err)
			}
			log.Printf("FeadFile")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if err := apiConnection.ResolveAddress(context.Background()); err != nil {
				log.Printf("cannot ResolveAddress: %v", err)
			}
			log.Printf("ResolveAddress")
		}()
	}

	wg.Wait()
}

type APIConecction1 struct{}

func Open1() *APIConecction1 {
	return &APIConecction1{}
}

func (a *APIConecction1) ReadFile(ctx context.Context) error {
	// 何か処理をしたということにする
	return nil
}

func (a *APIConecction1) ResolveAddress(ctx context.Context) error {
	// 何か処理をしたということにする
	return nil
}
