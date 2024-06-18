// this file might seem redundant and maybe it is but I'm just practising
// seperation of concerns here man
package main

import "sort"

func initialCart() cart {
	return cart{
		make(map[string]float64),
	}
}

func (c *cart) addToCart(b book) {
	c.items[b.Title] = b.Price
}

func (c *cart) removeFromCart(b book) {
	delete(c.items, b.Title)
}

func (c *cart) allTitles() []string {
	// Extract keys
	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}

	// Sort keys
	sort.Strings(keys)

	return keys
}

func (c *cart) cartItemsToDisp(top int, spatials dimensions) []string {
    itemCount := spatials.innerH / 3
	// extract keys
	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}

	// sort keys
	sort.Strings(keys)

	// ensure top and bot are within the valid range
	if top < 0 {
		top = 0
	}

	return keys[top:len(keys) % itemCount]
}

func (c *cart) bookInCart(b book) bool {
	_, ok := c.items[b.Title]
	return ok
}

func (c *cart) booksTotal() float64 {
	var price float64
	for _, p := range c.items {
		price += p
	}

	return price
}
