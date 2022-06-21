package vlc

import (
	"strings"
	"unicode"
)

func Encode(str string) string {
	// algorithm

	// prepare text: M -> !m
	str = prepareText(str)

	// encode to binary: some text -> 10010101
	bStr := encodeBinary(str)

	//split binary by chunks (8): bits to bytes -> '10010101 10010101 10010101'
	chunks := splitByChunks(bStr, chunkSize)

	// bytes to hex -> '20 30 3C'
	return chunks.ToHex().hexToString()
}

func Decode(encodedText string) string {
	// hex chunks -> splited by " " and putted in HexChunks slice
	hChunks := NewHexChunks(string(encodedText))

	// Hex chunks convert to binary chunks // "100110", "1011001" ect.
	bChunks := hChunks.ToBinary()

	// bChunks -> join in one string "11010101001011010010101"
	bStr := bChunks.Join()

	// build decoding tree
	dTree := getEncodingTable().DecodingTree()

	return exportText(dTree.Decode(bStr)) // My name is Ted -> !my name is !ted
}

// prepareText prepares text to be fit for encode:
// changes upper case letters to: ! + lower case latter
// i.g: My name is Ted -> !my name is !ted
func prepareText(str string) string {
	var buf strings.Builder

	for _, ch := range str {
		if unicode.IsUpper(ch) {
			buf.WriteRune('!')
			buf.WriteRune(unicode.ToLower(ch))
		} else {
			buf.WriteRune(ch)
		}
	}

	return buf.String()
}

// exportText is opposite to prepareText, it prepares decoded text to export:
// it chages: ! + <lower case latter> ->to upper case latter.
// i.g.: !my name is !ted -> My name is Ted.
func exportText(str string) string {
	var buf strings.Builder

	var isCapital bool

	for _, ch := range str {
		if isCapital {
			buf.WriteRune(unicode.ToUpper(ch))
			isCapital = false

			continue
		}

		if ch == '!' {
			isCapital = true

			continue
		} else {
			buf.WriteRune(ch)
		}
	}

	return buf.String()
}

// encodeBinary encodes str into binary codes string without spaces
// i.g: !my name is !ted -> "001000000011000000111100000110000111011101001010111001000100110100101"
func encodeBinary(str string) string {
	var buf strings.Builder

	for _, ch := range str {
		buf.WriteString(bin(ch))
	}

	return buf.String()
}

func bin(ch rune) string {
	table := getEncodingTable()

	res, ok := table[ch]
	if !ok {
		panic("unkown character:" + string(ch))
	}

	return res
}

func getEncodingTable() encodingTable {
	return encodingTable{
		' ': "11",
		't': "1001",
		'n': "10000",
		's': "0101",
		'r': "01000",
		'd': "00101",
		'!': "001000",
		'c': "000101",
		'm': "000011",
		'g': "0000100",
		'b': "0000010",
		'v': "00000001",
		'k': "0000000001",
		'q': "000000000001",
		'e': "101",
		'o': "10001",
		'a': "011",
		'i': "01001",
		'h': "0011",
		'l': "001001",
		'u': "00011",
		'f': "000100",
		'p': "0000101",
		'w': "0000011",
		'y': "0000001",
		'j': "000000001",
		'x': "00000000001",
		'z': "000000000000",
	}
}
