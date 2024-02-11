package cmd

import (
	"badrodem/run"
	"os"

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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.MousetrapHelpText = ""
	rootCmd.DisableFlagsInUseLine = true
}
