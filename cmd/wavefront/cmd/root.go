package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spaceapegames/go-wavefront/api"
	"github.com/spf13/cobra"
)

var address, token string
var rawResponse, debugMode bool
var client *wavefront.Client

var RootCmd = &cobra.Command{
	Use:   "wavefront",
	Short: "CLI Tool for accessing Wavefront API",
	Long: `CLI Tool for accessing Wavefront API
Ensure that Wavefront token and address are passed, either through the WAVEFRONT_TOKEN and WAVEFRONT_ADDRESS environment variables, or the command line.`,
}

func init() {
	//persistent flags attached to the root will be global flags
	RootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Wavefront access token")
	RootCmd.PersistentFlags().StringVarP(&address, "wavefront-address", "w", "", "Wavefront API address. e.g. example.wavefront.com")
	RootCmd.PersistentFlags().BoolVarP(&rawResponse, "raw", "j", false, "show raw JSON response")
	RootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "show HTTP request information")

	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		_ = RootCmd.Help()
		os.Exit(0)
	}

}

func createClient() {
	if token == "" {
		token = os.Getenv("WAVEFRONT_TOKEN")
	}
	if token == "" {
		log.Fatal("Please provide a Wavefront token either through the WAVEFRONT_TOKEN environment variable or the --token flag")
	}

	if address == "" {
		address = os.Getenv("WAVEFRONT_ADDRESS")
	}
	if address == "" {
		log.Fatal("Please provide a Wavefront address either through the WAVEFRONT_ADDRESS environment variable or the --address flag")
	}

	var err error
	client, err = wavefront.NewClient(&wavefront.Config{
		Address: address,
		Token:   token,
	})

	if err != nil {
		log.Fatal(err)
	}
	client.Debug(debugMode)
}

func Execute() {
	fmt.Println(RootCmd.PersistentFlags().Lookup("debug").Value)
	fmt.Println(rawResponse)
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func prettyPrint(raw *[]byte) {
	var output bytes.Buffer
	err := json.Indent(&output, *raw, "", "\t")
	if err != nil {
		log.Println("JSON error: ", err)
	}
	fmt.Println(string(output.Bytes()))
}
