package vlc

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"strings"

	"github.com/MaximumTroubles/go-archiver/lib/compression/vlc/table"
)

type EncoderDecoder struct {
	tblGenerator table.Generator
}

func New(tblGenerator table.Generator) EncoderDecoder {
	return EncoderDecoder{
		tblGenerator: tblGenerator,
	}
}

func (ed EncoderDecoder) Encode(str string) []byte {
	// create empty encoding table
	tbl := ed.tblGenerator.NewTable(str)

	encoded := encodeBinary(str, tbl)

	return buildEncodedFile(tbl, encoded)
}

// encodeBinary encodes str into binary codes string without spaces
// i.g: !my name is !ted -> "001000000011000000111100000110000111011101001010111001000100110100101"
func encodeBinary(str string, table table.EncodingTable) string {
	var buf strings.Builder

	for _, ch := range str {
		buf.WriteString(bin(ch, table))
	}

	return buf.String()
}

func bin(ch rune, table table.EncodingTable) string {

	res, ok := table[ch]
	if !ok {
		panic("unkown character:" + string(ch))
	}

	return res
}

func buildEncodedFile(tbl table.EncodingTable, data string) []byte {
	// here we should encode table as well and put to the file. its makes for further decoding bites to text and we need this table
	// encodedTbl = subsequence of bites in form of a table
	encodedTbl := encodeTable(tbl)

	var buf bytes.Buffer

	// here we do like single sting of bytes where we put everything together so we need to know size of each elemens
	// as encoding table or data it self

	// Size of table converted to bites
	buf.Write(encodeInt(len(encodedTbl)))
	// Size of data converted to bites
	buf.Write(encodeInt(len(data)))
	// table itself
	buf.Write(encodedTbl)
	// data itself
	buf.Write((splitByChunks(data, chunkSize).Bytes()))

	return buf.Bytes()
}

func encodeTable(tbl table.EncodingTable) []byte {
	// here we make compress on our encoding table
	var tableBuf bytes.Buffer

	if err := gob.NewEncoder(&tableBuf).Encode(tbl); err != nil {
		log.Fatal("can't serialize table:", err)
	}

	return tableBuf.Bytes()
}

func encodeInt(num int) []byte {
	// BigEndian means that we put bites from left to right
	// 32 bits = 4 byte that why we knew size of bytes slice
	res := make([]byte, 4)
	binary.BigEndian.PutUint32(res, uint32(num))

	return res
}

//
// Decode
//

func (ed EncoderDecoder) Decode(encodedData []byte) string {
	// First parse and read data from file
	tbl, data := parseFile(encodedData)

	return tbl.Decode(data)
}

func parseFile(data []byte) (table.EncodingTable, string) {
	const (
		tableSizeBytesCount = 4
		dataSizeBytesCount  = 4
	)

	// tableSize = 0-3, data = 4...
	tableSizeBinary, data := data[:tableSizeBytesCount], data[tableSizeBytesCount:]
	// data size = 4-7, data 8...
	dataSizeBinary, data := data[:dataSizeBytesCount], data[dataSizeBytesCount:]

	tableSize := binary.BigEndian.Uint32(tableSizeBinary)
	dataSize := binary.BigEndian.Uint32(dataSizeBinary)

	tblBinary, data := data[:tableSize], data[tableSize:]

	table := decodeTable(tblBinary)

	body := NewBinChunks(data).Join()

	return table, body[:dataSize]
}

func decodeTable(tblBinary []byte) table.EncodingTable {
	var tbl table.EncodingTable

	r := bytes.NewReader(tblBinary)
	if err := gob.NewDecoder(r).Decode(&tbl); err != nil {
		log.Fatal("can't decode table: ", err)
	}

	return tbl
}
