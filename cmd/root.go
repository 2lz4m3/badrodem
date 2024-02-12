package cmd

import (
	"badrodem/platform"
	"badrodem/run"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "badrodem <file>...",
	Short: "An UTF-8 BOM adder.",
	Long: `badrodem is an anagram of "bom-adder".
It adds an UTF-8 BOM (Byte Order Mark, 0xEF 0xBB 0xBF)
to the beginning of the text file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := run.Run()
		return err
	},
}

func Execute(versionString string) {
	rootCmd.Version = versionString

	err := rootCmd.Execute()
	if runtime.GOOS == "windows" && platform.IsDoubleClickRun() {
		// keep console open on exit
		fmt.Print("Press any key to continue . . .")
		os.Stdin.Read(make([]byte, 1))
	}
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.MousetrapHelpText = ""
	rootCmd.DisableFlagsInUseLine = true
}
