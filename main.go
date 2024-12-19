package main

import (
	"fmt"
	"github.com/perepeli/ipv4trie/trie"
)

func main() {
	t := trie.NewIPv4Trie()

	// insert ip addresses
	t.Insert("192.168.0.1")
	t.Insert("192.168.0.1") // duplicate
	t.Insert("192.168.0.2")
	t.Insert("10.0.0.1")
	t.Insert("8.8.8.8")

	// print the number of unique addresses
	fmt.Println("Unique IP count:", t.UniqueCount())

	// search if ip was previously inserted (why not?)
	fmt.Println("Search for 192.168.0.1:", t.Search("192.168.0.1")) // true
	fmt.Println("Search for 8.8.8.1:", t.Search("8.8.8.1"))         // false
}
