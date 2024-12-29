package bitset

import (
	"strconv"
	"strings"
)

type BitSet struct {
	bitArray    []uint32
	uniqueCount uint64
}

func NewIPv4BitSet() *BitSet {
	numOfInts := (1<<32 + 31) / 2
	return &BitSet{bitArray: make([]uint32, numOfInts)}
}

func (bitSet *BitSet) Insert(ip string) {
	index := ipToInt(ip)
	arrIndex := index / 32
	bitPosition := index % 32
	mask := uint32(1 << bitPosition)
	if (bitSet.bitArray[arrIndex] & mask) == 0 {
		bitSet.uniqueCount++
	}
	bitSet.bitArray[arrIndex] |= mask
}

func (bitSet *BitSet) Search(ip string) bool {
	index := ipToInt(ip)
	arrayIndex := index / 32  // Which int to use
	bitPosition := index % 32 // Which bit in that int
	return (bitSet.bitArray[arrayIndex] & (1 << bitPosition)) != 0
}

func (bitSet *BitSet) UniqueCount() uint64 {
	return bitSet.uniqueCount
}

func ipToInt(ipAddress string) int {
	parts := strings.Split(ipAddress, ".")
	if len(parts) != 4 {
		panic("Invalid IP address format")
	}

	var result int
	for _, part := range parts {
		val, err := strconv.Atoi(part)
		if err != nil || val < 0 || val > 255 {
			panic("Invalid IP octet: " + part)
		}
		result = (result << 8) | val
	}
	return result
}
