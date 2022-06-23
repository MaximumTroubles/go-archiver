package cmd

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/MaximumTroubles/go-archiver/lib/compression"
	"github.com/MaximumTroubles/go-archiver/lib/compression/vlc"
	"github.com/MaximumTroubles/go-archiver/lib/compression/vlc/table/shannon_fann"
	"github.com/spf13/cobra"
)

const packedExtension = "vlc"

var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "Compres file command",
	Run:   pack,
}

var ErrEmptyPath = errors.New("path to file is not specified")

func pack(cmd *cobra.Command, args []string) {
	var encoder compression.Encode

	if len(args) == 0 || args[0] == "" {
		handleErr(ErrEmptyPath)
	}

	method := cmd.Flag("method").Value.String()

	switch method {
	case "vlc":
		encoder = vlc.New(shannon_fann.NewGenerator())
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
	packed := encoder.Encode(string(data))

	// Here we write down compresed data to new file with perm: 0644 which means that current user can read and write other only read
	err = os.WriteFile(packedFileName(filePath), packed, 0644)
	if err != nil {
		handleErr(err)
	}
}

// Here we form file's name
func packedFileName(path string) string {
	var parts []string

	// Get file name filename.txt
	fileName := filepath.Base(path)
	// Get file ext .txt
	ext := filepath.Ext(fileName)
	// Get base name without ext
	parts = append(parts, strings.TrimSuffix(fileName, ext), packedExtension) // filename

	return strings.Join(parts, ".")
}

func init() {
	rootCmd.AddCommand(packCmd)

	packCmd.Flags().StringP("method", "m", "", "compression method: vlc")

	if err := packCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
}
