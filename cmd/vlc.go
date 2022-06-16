package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var vlcCmd = &cobra.Command{
	Use:   "vlc",
	Short: "Compres file using variable-length code",
	Run:   pack,
}

var ErrEmptyPath = errors.New("path to file is not specified\n")

const packedExtension = "vlc"

func pack(_ *cobra.Command, args []string) {
	if len(args) == 0 || args[0] == "" {
		handleErr(ErrEmptyPath)
	}

	// here we recieve an argument from command line (path to file)
	filePath := args[0]

	// We got here file and open it.
	r, err := os.Open(filePath)
	if err != nil {
		handleErr(err)
	}

	// Here we read all data from file.
	data, err := io.ReadAll(r)
	if err != nil {
		handleErr(err)
	}

	// Here we provide data to file compresor
	// packed := Encode(data)
	packed := ""
	fmt.Println(string(data))

	// Here we write down compresed data to new file with perm: 0644 which means that current user can read and write other only read
	err = os.WriteFile(packedFileName(filePath), []byte(packed), 0644)
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
	packCmd.AddCommand(vlcCmd)
}
