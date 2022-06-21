package vlc

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type BinaryChuncks []BinaryChunck

type BinaryChunck string

type HexChunks []HexChunk

type HexChunk string

type encodingTable map[rune]string

const chunkSize = 8

const hexChunkSeparator = " "

func NewHexChunks(str string) HexChunks {
	parts := strings.Split(str, hexChunkSeparator)

	res := make(HexChunks, 0, len(parts))

	for _, part := range parts {
		res = append(res, HexChunk(part))
	}

	return res
}

func (bcs BinaryChuncks) ToHex() HexChunks {
	res := make(HexChunks, 0, len(bcs))

	for _, chunk := range bcs {
		//Here we use method ToHex on 1 byte string = "00100100" and convert it into Hexadecimal notation
		res = append(res, chunk.ToHex())
	}

	return res
}


// Join joins chucks into one line and return as string 
func (bcs BinaryChuncks) Join() string {
	var buf strings.Builder

	for _, bc := range bcs {
		buf.WriteString(string(bc))
	}

	return buf.String()
}
 
func (bc BinaryChunck) ToHex() HexChunk {
	num, err := strconv.ParseUint(string(bc), 2, chunkSize) // 2f 3a 9c
	if err != nil {
		panic("can't parse binary chunk: " + err.Error())
	}

	// here we have to make chunk ToUpper Case
	res := strings.ToUpper(fmt.Sprintf("%x", num)) // 2F 3A 9C

	// occassionaly convertation can return single result i.g:  "1 2F 3 9C", so we could prefix 0 to hex chunk
	// ex. "1 2F 3 9C" => "01 2F 03 9C"
	if len(res) == 1 {
		res = "0" + res
	}

	return HexChunk(res)
}

func (hcs HexChunks) hexToString() string {

	switch len(hcs) {
	case 0:
		return ""
	case 1:
		return string(hcs[0])
	}

	var buf strings.Builder

	buf.WriteString(string(hcs[0]))

	for _, hc := range hcs[1:] {
		buf.WriteString(hexChunkSeparator)
		buf.WriteString(string(hc))
	}

	return buf.String()
}

func splitByChunks(bStr string, chunkSize int) BinaryChuncks {
	strLen := utf8.RuneCountInString(bStr)

	chunksCount := strLen / chunkSize // 57 / 8 ==> go round it to interger = 7

	//Remainder
	if strLen%chunkSize != 0 {
		chunksCount++
	}

	res := make(BinaryChuncks, 0, chunksCount)

	var buf strings.Builder

	for i, ch := range bStr {
		buf.WriteString(string(ch))

		if (i+1)%chunkSize == 0 {
			res = append(res, BinaryChunck(buf.String()))
			buf.Reset()
		}
	}

	if buf.Len() != 0 {
		lastChunk := buf.String()

		lastChunk += strings.Repeat("0", chunkSize-len(lastChunk))

		res = append(res, BinaryChunck(lastChunk))
	}

	return res
}

func (hcs HexChunks) ToBinary() BinaryChuncks {
	res := make(BinaryChuncks, 0, len(hcs))

	for _, chunk := range hcs {
		res = append(res, chunk.ToBinary())
	}

	return res
}

func (hc HexChunk) ToBinary() BinaryChunck {
	num, err := strconv.ParseUint(string(hc), 16, chunkSize)
	if err != nil {
		panic("can't parse hex chunk" + err.Error())
	}

	res := fmt.Sprintf("%08b", num)

	return BinaryChunck(res)
}
