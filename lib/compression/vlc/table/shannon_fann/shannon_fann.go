package shannon_fann

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/MaximumTroubles/go-archiver/lib/compression/vlc/table"
)

type Generator struct{}

func NewGenerator() Generator {
	return Generator{}
}

type encodingTable map[rune]code

type code struct {
	Char     rune
	Quantity int
	Bits     uint32
	Size     int
}

// we specified type for char stats as key = char, value = binary code
type charStat map[rune]int

func (g Generator) NewTable(text string) table.EncodingTable {
	// we should count characters stat how often they meet in text
	charStat := newCharStat(text)

	// based on char stats we build encoding table
	t := build(charStat)

	// and return outside in specific format
	return t.Export()
}

func (et encodingTable) Export() map[rune]string {
	res := make(map[rune]string)

	for ch, code := range et {
		byteStr := fmt.Sprintf("%b", code.Bits)

		// for ex. code.Size = 0011 aterf fmt.Sprinf() we lossing two 00 infront of bits
		if lenDiff := code.Size - len(byteStr); lenDiff > 0 {
			byteStr = strings.Repeat("0", lenDiff) + byteStr
		}

		res[ch] = byteStr
	}

	return res
}

func build(stat charStat) encodingTable {
	// here we init slice of type code sturct
	// we make slice here cuz we must sort list of char by decreasining the frequency meet ups
	codes := make([]code, 0, len(stat))

	for ch, qty := range stat {
		codes = append(codes, code{
			Char:     ch,
			Quantity: qty,
		})
	}

	// sorting
	sort.Slice(codes, func(i, j int) bool {
		if codes[i].Quantity != codes[j].Quantity {
			return codes[i].Quantity > codes[j].Quantity
		}

		return codes[i].Char < codes[j].Char
	})

	assignCodes(codes)

	res := make(encodingTable)

	for _, code := range codes {
		res[code.Char] = code
	}

	return res
}

func assignCodes(codes []code) {
	if len(codes) < 2 {
		return
	}

	// here we have sorted list of code type structs with char and quantity of it in decs 20 -> 18 -> 17
	// we have to divede it by 2 equeal groups or almost equeal
	divider := bestDividerPosition(codes)

	// in loop
	for i := 0; i < len(codes); i++ {
		// we just added to field Bits 00000010. this operator means shift bit 1 to left
		// for all left side bits add 0 at the end. and 1 for all right bits
		codes[i].Bits <<= 1
		codes[i].Size++

		// here is check for verify right side of slpit groups
		if i >= divider {
			codes[i].Bits |= 1
		}
	}

	assignCodes(codes[:divider])
	assignCodes(codes[divider:])
}

func bestDividerPosition(codes []code) int {

	// count total chars quantity
	total := 0
	for _, code := range codes {
		total += code.Quantity
	}

	// it's var with max possible num
	prevDiff := math.MaxInt
	left := 0
	bestPosition := 0

	for i := 0; i < len(codes)-1; i++ {
		// a | b c f d -- first we take first index the biggest quantity number and compare with sum of other side's nums
		left += codes[0].Quantity
		right := total - left

		// calc diff ex. 150 - 15, and abs func helps when 40 - 50 = -10 and we got negative num, so it invert it
		diff := abs(right - left)
		// here we need atleast one iteration so we make diff smaller than prevDiff on first interation
		if diff >= prevDiff {
			break
		}

		prevDiff = diff
		bestPosition = i + 1
	}

	return bestPosition
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func newCharStat(text string) charStat {
	res := make(charStat)

	for _, ch := range text {
		res[ch]++
	}

	return res
}
