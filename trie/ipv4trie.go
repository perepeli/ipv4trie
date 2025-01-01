// key ideas:
// - we use trie (prefix tree) to store ip addresses memory efficiently
// - if I was asked to store all english words I would use trie same way (key idea is to reuse prefixes to save space vs naive set approach)
// - each bit of 32 bit ip address representation is stored only once and if "shared" with other address is reused
// - a bit worse approach would use octet strings as trie nodes e.g "192" node, but after some brainstorming 32-bit int representation looked much more optimal
// - each IP address is converted to a 32-bit integer.
// - each bit of the integer (from 31 to 0) is used to navigate the trie using bitwise operations

package trie

import (
	"strconv"
	"strings"
)

type TrieNode struct {
	children [2]*TrieNode
}

type IPv4Trie struct {
	root        *TrieNode
	uniqueCount uint64
}

func NewIPv4Trie() *IPv4Trie {
	return &IPv4Trie{root: &TrieNode{}}
}

func (t *IPv4Trie) Insert(ipAddress string) {
	ip := ipToInt(ipAddress)
	node := t.root
	for i := 31; i >= 0; i-- {
		bit := (ip >> i) & 1
		if node.children[bit] == nil {
			node.children[bit] = &TrieNode{}
			if i == 0 {
				t.uniqueCount++
			}
		}
		node = node.children[bit]
	}
}

func (t *IPv4Trie) Search(ipAddress string) bool {
	ip := ipToInt(ipAddress)
	node := t.root

	for i := 31; i >= 0; i-- {
		bit := (ip >> i) & 1
		if node.children[bit] == nil {
			return false
		}
		node = node.children[bit]
	}
	return true
}

func (t *IPv4Trie) UniqueCount() uint64 {
	return t.uniqueCount
}

func ipToInt(ipAddress string) uint32 {
	parts := strings.Split(ipAddress, ".")
	if len(parts) != 4 {
		panic("Invalid IP address format")
	}

	var result uint32
	for _, part := range parts {
		val, err := strconv.Atoi(part)
		if err != nil || val < 0 || val > 255 {
			panic("Invalid IP octet: " + part)
		}
		result = (result << 8) | uint32(val)
	}
	return result
}
