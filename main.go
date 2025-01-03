package main

import (
	"fmt"
	"github.com/perepeli/ipv4trie/bitset"
	"github.com/perepeli/ipv4trie/trie"
	"runtime"
)

/*
Logical Reasoning for Design:

1. Initial Approach - Trie Structure**:
   - The initial idea was to use a trie to "reuse shared prefixes" of IPv4 addresses.
   - But for the full range of IPv4 addresses (0.0.0.0 to 255.255.255.255), every possible address would need to be stored, leading to the worst-case scenario of 2^33 nodes.
   - Each node in the trie would require at least 3 pointers (node reference + two child references for '0' and '1' bits).
   - Assuming 64-bit systems where pointers are 64 bits (8 bytes), the memory cost for each node is 3 x 8 = 24 bytes.
   - Total memory required for 8,589,934,592 nodes: 8,589,934,592 * 24 = 192 GB
   - This is an impractical amount of RAM, and clearly way too much.

P.S, the above ^ trie solution can be optimized by storing not objects, but representing them as a complete binary tree as a bitset array, (similar to binary heap),
where the last depth nodes would represent the presence or absence of an ip, and all depths above would be '0' and required to navigate to the last level,
this solution would require 2^33 * 1 bit = 1 GB, but we waste a lot of space for nodes above the last level just for the sake of navigation...

2. Next Approach - Naive Array of Bytes:
   - Instead of a complex trie, I decided to start with using a simple array where each byte represents a boolean (whether an IP has been seen).
   - To store 2^32 IPs, I would need: 2^32 * 1 byte (boolean) = 4 GB
   - This is significantly better than the 192 GB required by the trie, but it's still not the theoretical minimum.
   - Question is can we shrink usage of boolean (byte) as little as single bit?
   - Is it even possible? - yes and for that I had to research a bit about a new for me medata structure - Bitset

3. Optimal Approach with Bitset:
   - A boolean technically only requires 1 bit, not 1 byte.
   - Instead of using an array of booleans (where each boolean takes 8 bits), we pack multiple booleans into a single integer.
   - Since a 32-bit integer holds 32 bits, we can track 32 IP addresses in 1 integer.
   - To track 2^32 IPs, we need:  134,217,728 integers

   - Since each integer is 4 bytes, the total space required is: 34,217,728 * 4 = 512 MB
   - This is obviously much smaller than the 192 GB for trie and also much smaller than the 4 GB for byte (boolean) array.

   - To set a bit, we compute two indices:
     - Array index: index / 32
     - Bit position: index % 32
   - This allows us to track IPs efficiently and achieve optimal memory usage with as little as 512 MB for worse case of complete IPv4 address range.
*/

func main() {
	fmt.Println("Mem stats before:")
	printMemStats() // Alloc = 0 MB,   TotalAlloc = 0 MB

	bitSet := bitset.NewIPv4BitSet()

	bitSet.Insert("192.168.0.1")
	bitSet.Insert("192.168.0.1") // duplicate
	bitSet.Insert("192.168.0.2")
	bitSet.Insert("10.0.0.1")
	bitSet.Insert("8.8.8.8")
	bitSet.Insert("8.8.8.86")
	bitSet.Insert("8.8.8.81")
	bitSet.Insert("8.8.8.82")
	bitSet.Insert("8.8.8.83")
	bitSet.Insert("255.255.255.255")

	// print the number of unique addresses
	fmt.Println("Unique IP count:", bitSet.UniqueCount()) // 9

	// search if ip was previously inserted
	fmt.Println("Search for 192.168.0.1:", bitSet.Search("192.168.0.1"))         // true
	fmt.Println("Search for 8.8.8.1:", bitSet.Search("8.8.8.1"))                 // false
	fmt.Println("Search for 8.8.8.8:", bitSet.Search("8.8.8.8"))                 // true
	fmt.Println("Search for 255.255.255.255:", bitSet.Search("255.255.255.255")) // true

	var value uint32
	for value = 0; value <= 0xFFFFFFFF; value++ {
		bitSet.InsertInternal(value)
		if value%100000000 == 0 || value == 0xFFFFFFFF {
			fmt.Printf("Current value: %032b (%d)\n", value, value)
		}

		if value == 0xFFFFFFFF {
			fmt.Println("Loop is finished, full ipv4 range (0.0.0.0 - 255.255.255.255) is inserted.")
			break
		}
	}

	fmt.Println("Unique IP count:", bitSet.UniqueCount()) // 2^32 or 4294967296

	fmt.Println("Mem stats after:")
	printMemStats() // Alloc = 512 MB, TotalAlloc = 512 MB
}

func mainOld() { // initial naive approach with trie <<
	trie := trie.NewIPv4Trie()

	// insert ip addresses
	trie.Insert("192.168.0.1")
	trie.Insert("192.168.0.1") // duplicate
	trie.Insert("192.168.0.2")
	trie.Insert("10.0.0.1")
	trie.Insert("8.8.8.8")
	trie.Insert("8.8.8.86")
	trie.Insert("8.8.8.81")
	trie.Insert("8.8.8.82")
	trie.Insert("8.8.8.83")
	trie.Insert("8.8.8.8")

	// print the number of unique addresses
	fmt.Println("Unique IP count:", trie.UniqueCount())

	// search if ip was previously inserted (why not?)
	fmt.Println("Search for 192.168.0.1:", trie.Search("192.168.0.1")) // true
	fmt.Println("Search for 8.8.8.1:", trie.Search("8.8.8.1"))         // false
}

func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("\tAlloc = %v MB,", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MB\n", bToMb(m.TotalAlloc))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
