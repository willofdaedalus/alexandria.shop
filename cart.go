// this file might seem redundant and maybe it is but I'm just practising
// seperation of concerns here man
package main

import (
	"fmt"
	"slices"
)

func initialCart() cart {
	return cart{
		make([]cartItem, 0),
	}
}

func (c *cart) addToCart(b book) {
	item := cartItem{b.Title, b.Price}
	if slices.Contains(c.items, item) {
		return
	}
	c.items = append(c.items, item)
}

func (c *cart) removeFromCartStr(t string) {
	idx := -1
	for i, ci := range c.items {
		if ci.title == t {
			idx = i
			break
		}
	}
	if idx == -1 {
		return
	}
	c.items = append(c.items[:idx], c.items[idx+1:]...)
}

func (c *cart) removeFromCart(b book) {
	item := cartItem{b.Title, b.Price}
	itemIdx := slices.Index(c.items, item)
	// https://pkg.go.dev/slices#Index
	if itemIdx == -1 {
		return
	}
	c.items = append(c.items[:itemIdx], c.items[itemIdx+1:]...)
}

func (c *cart) cartItemsToDisp(top int, spatials dimensions) []string {
	// var limit int
	itemCount := spatials.innerH / 3
	// extract keys
	keys := make([]string, 0, len(c.items))
	for _, k := range c.items {
		keys = append(keys, k.title)
	}

	// ensure top and bot are within the valid range
	if top < 0 {
		top = 0
	}

	limit := len(keys)
	if limit > itemCount {
		limit = top + (itemCount - 1)
		if limit > len(keys) {
			limit = len(keys)
		}
	}
	logToFile(fmt.Sprintf("limit is: %d", limit))

	retval := make([]string, limit)
	copy(retval, keys[top:limit])

	return retval
}

func (c *cart) bookInCart(b book) bool {
	item := cartItem{b.Title, b.Price}
	return slices.Contains(c.items, item)
}

func (c *cart) booksTotal() float64 {
	var price float64
	for _, p := range c.items {
		price += p.price
	}

	return price
}
