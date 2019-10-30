package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0510/users/61b9fb41fab22e8dd0e5441cc2acb80024595b55/i-img690x530-1570514714uroqjb128.jpg"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(byteArray)
}
