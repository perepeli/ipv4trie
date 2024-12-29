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
	numOfInts := (uint64(1<<32) + 31) / 32
	return &BitSet{bitArray: make([]uint32, numOfInts)}
}

func (bitSet *BitSet) Insert(ip string) {
	index := ipToInt(ip)
	bitSet.InsertInternal(index)
}

func (bitSet *BitSet) InsertInternal(ip uint32) {
	index := ip
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
	return bitSet.SearchInternal(index)
}

func (bitSet *BitSet) SearchInternal(ip uint32) bool {
	index := ip
	arrayIndex := index / 32  // which int to use
	bitPosition := index % 32 // which bit in that int
	return (bitSet.bitArray[arrayIndex] & (1 << bitPosition)) != 0
}

func (bitSet *BitSet) UniqueCount() uint64 {
	return bitSet.uniqueCount
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
