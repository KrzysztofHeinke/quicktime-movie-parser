/*
Copyright © 2024 Krzysztof Heinke <Krzysztof.Heinke@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "quicktime-movie-parser",
	Short: "A command-line tool for processing and analyzing media files, offering functionalities like parsing, metadata extraction, and format conversion.",
	Long: `A command-line tool for processing and analyzing media files, offering functionalities like parsing, metadata extraction, and format conversion.
	
	To run the parsing type:
	./quicktime-movie-parser parse <path_to_mov_file>
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Set the default log level and format
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	rootCmd.PersistentFlags().StringP("loglevel", "l", "info", "Set the logging level (debug, info, warn, error)")
	cobra.OnInitialize(initLogger)
}

func initLogger() {
	logLevel, err := rootCmd.Flags().GetString("loglevel")
	if err != nil {
		logrus.Fatalf("Could not read loglevel flag: %v", err)
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatalf("Invalid log level: %v", err)
	}
	logrus.SetLevel(level)

	logrus.Infof("Log level set to %s", logLevel)
}
