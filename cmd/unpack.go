package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/MaximumTroubles/go-archiver/lib/compression"
	"github.com/MaximumTroubles/go-archiver/lib/compression/vlc"
	"github.com/MaximumTroubles/go-archiver/lib/compression/vlc/table/shannon_fann"
	"github.com/spf13/cobra"
)

const unpackedExtension = "txt"

var unpackCmd = &cobra.Command{
	Use:   "unpack",
	Short: "Decompres file command command",
	Run:   unpack,
}

func unpack(cmd *cobra.Command, args []string) {
	var decoder compression.Decode

	if len(args) == 0 || args[0] == "" {
		handleErr(ErrEmptyPath)
	}

	method := cmd.Flag("method").Value.String()

	switch method {
	case "vlc":
		decoder = vlc.New(shannon_fann.NewGenerator())
	default: 
		cmd.PrintErr("unknown method")
	}

	// here we recieve an argument from command line (path to file)
	filePath := args[0]

	// We got here file and open it.
	r, err := os.Open(filePath)
	if err != nil {
		handleErr(err)
	}

	// Close file always
	defer r.Close()

	// Here we read all data from file.
	data, err := io.ReadAll(r)
	if err != nil {
		handleErr(err)
	}

	// Here we provide data to file compresor
	packed := decoder.Decode(data)

	// Here we write down compresed data to new file with perm: 0644 which means that current user can read and write other only read
	err = os.WriteFile(unpackedFileName(filePath), []byte(packed), 0644)
	if err != nil {
		handleErr(err)
	}
}

// Here we form file's name
func unpackedFileName(path string) string {
	var parts []string

	// Get file name filename.txt
	fileName := filepath.Base(path)
	// Get file ext .txt
	ext := filepath.Ext(fileName)
	// Get base name without ext
	parts = append(parts, strings.TrimSuffix(fileName, ext), unpackedExtension) // filename

	return strings.Join(parts, ".")
}

func init() {
	rootCmd.AddCommand(unpackCmd)

	unpackCmd.Flags().StringP("method", "m", "", "Decompression method: vlc")

	if err := unpackCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
}
