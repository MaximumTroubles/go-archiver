package vlc

import (
	"reflect"
	"testing"
)

func Test_splitByChunks(t *testing.T) {
	type args struct {
		bStr      string
		chunkSize int
	}
	tests := []struct {
		name string
		args args
		want BinaryChuncks
	}{
		{
			name: "splitByChunks func test",
			args: args{
				bStr:      "001000100110100101",
				chunkSize: 8,
			},
			want: BinaryChuncks{"00100010", "01101001", "01000000"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitByChunks(tt.args.bStr, tt.args.chunkSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitByChunks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinaryChuncks_ToHex(t *testing.T) {
	tests := []struct {
		name string
		bcs  BinaryChuncks
		want HexChunks
	}{
		{
			name: "base test",
			bcs:  BinaryChuncks{"01011111", "10000000"},
			want: HexChunks{"5F", "80"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bcs.ToHex(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BinaryChuncks.ToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHexChunks(t *testing.T) {

	tests := []struct {
		name string
		str  string
		want HexChunks
	}{
		{
			name: "base name",
			str:  "20 30 3C",
			want: HexChunks{"20", "30", "3C"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHexChunks(tt.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHexChunks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHexChunks_ToBinary(t *testing.T) {
	tests := []struct {
		name string
		hcs  HexChunks
		want BinaryChuncks
	}{
		{
			name: "base test",
			hcs:  HexChunks{"20", "30", "3C"},
			want: BinaryChuncks{"00100000", "00110000", "00111100"},
		},

		{
			name: "base test case 2",
			hcs:  HexChunks{"20", "30", "3C", "00"},
			want: BinaryChuncks{"00100000", "00110000", "00111100", "00000000"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hcs.ToBinary(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HexChunks.ToBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHexChunk_ToBinary(t *testing.T) {
	tests := []struct {
		name string
		hc   HexChunk
		want BinaryChunck
	}{
		{
			name: "base test 1 case",
			hc:   HexChunk("20"),
			want: BinaryChunck("00100000"),
		},
		{
			name: "base test 2 case",
			hc:   HexChunk("3C"),
			want: BinaryChunck("00111100"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hc.ToBinary(); got != tt.want {
				t.Errorf("HexChunk.ToBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinaryChuncks_Join(t *testing.T) {
	tests := []struct {
		name string
		bcs  BinaryChuncks
		want string
	}{
		{
			name: "base test",
			bcs:  BinaryChuncks{"00100000", "00110000", "00111100", "00000000"},
			want: "00100000001100000011110000000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bcs.Join(); got != tt.want {
				t.Errorf("BinaryChuncks.Join() = %v, want %v", got, tt.want)
			}
		})
	}
}
