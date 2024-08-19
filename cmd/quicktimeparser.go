/*
Copyright © 2024 Krzysztof Heinke <Krzysztof.Heinke@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/KrzysztofHeinke/quicktime-movie-parser/internal/parser"
	"github.com/spf13/cobra"
)

// quicktimeparserCmd represents the quicktimeparser command
var quicktimeparserCmd = &cobra.Command{
	Use: "parse",
	Short: `Parse the 'moov' atom from a MOV/MP4 file to extract and analyze" 
			metadata such as track information, sample rates, and video dimensions.`,
	Long: `The 'moov' atom parser is a tool designed to extract and analyze the 'moov' atom
	 from MOV/MP4 files. The 'moov' atom contains critical metadata that describes the structure
	  of the video, including track information, sample rates, video dimensions, and more. 
	  This tool reads the 'moov' atom from the specified file, identifies and processes all child
	   	atoms, and extracts key details such as audio sample rates, video width, and height. 
	   It is essential for tasks such as media file analysis, editing, and metadata extraction.`,
	Run: func(cmd *cobra.Command, args []string) {
		info, err := os.Stat(args[0])
		if os.IsNotExist(err) {
			fmt.Printf("File %s do not exist!", args[0])
			return
		} else if info.IsDir() {
			fmt.Printf("%s that is a directory, not a file!", args[0])
			return
		}
		parser.Parse(args[0])
	},
}

func init() {
	rootCmd.AddCommand(quicktimeparserCmd)
}
