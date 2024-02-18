package cmd

import (
	"badrodem/localize"
	"badrodem/platform"
	"badrodem/run"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Verbose bool

	rootCmd = &cobra.Command{
		Use:   "badrodem <file>...",
		Short: "An UTF-8 BOM adder.",
		Long: `badrodem is an anagram of "bom-adder".
It adds an UTF-8 BOM (Byte Order Mark, 0xEF 0xBB 0xBF)
to the beginning of the text file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if Verbose {
				fmt.Printf("Language: %s\n", localize.I18n.Tag)
			}
			err := run.Run()
			return err
		},
	}
)

func Execute(versionString string) {
	rootCmd.Version = versionString

	err := rootCmd.Execute()
	if runtime.GOOS == "windows" && platform.IsDoubleClickRun() {
		// keep console open on exit
		fmt.Print(localize.T("press_any_key_to_continue"))
		os.Stdin.Read(make([]byte, 1))
	}
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.MousetrapHelpText = ""
	rootCmd.DisableFlagsInUseLine = true

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "", false, "verbose output")
}
