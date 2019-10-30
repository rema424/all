package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("./service_account_key.json"))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Buckets:")
	it := client.Buckets(ctx, "michael-sandbox7")
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(battrs.Name)
	}

}
