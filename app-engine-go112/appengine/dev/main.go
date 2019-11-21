package main

import "fmt"

func main() {
	type Item struct {
		ID    int
		Name  string
		Price int
	}
	fn := func(v interface{}) {
		i, ok := v.(*Item)
		fmt.Println(i, ok)
	}
	var i *Item
	fn(i)
	// type Inventory struct {
	// 	ID string
	// 	Item
	// }

	// i := Inventory{
	// 	ID: "XXXX-XXXX",
	// 	Item: Item{
	// 		ID:    1,
	// 		Name:  "商品A",
	// 		Price: 1000,
	// 	},
	// }

	// fmt.Println(i.ID, i.ID, i.Item.ID, i.Name, i.Price)
}
