// this file might seem redundant and maybe it is but I'm just practising
// seperation of concerns here man
package main

func initialCart() cart {
	return cart{
		make(map[string]float64),
	}
}

func (c cart) addToCart(b book) {
	c.items[b.Title] = b.Price
}

func (c cart) removeFromCart(b book) {
	delete(c.items, b.Title)
}

func (c cart) allTitles() []string {
	var titles = make([]string, 0)

	for k := range c.items {
		titles = append(titles, k)
	}

	return titles
}

func (c cart) allBooksInCart() map[string]float64 {
	return c.items
}

func (c cart) bookInCart(b book) bool {
	_, ok := c.items[b.Title]
	return ok
}

func (c cart) booksTotal() float64 {
	var price float64
	for _, p := range c.items {
		price += p
	}

	return price
}
