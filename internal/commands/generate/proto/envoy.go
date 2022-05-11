package proto

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/Clink-n-Clank/Eitri/internal/handler"
	"github.com/Clink-n-Clank/Eitri/internal/log"
	"github.com/spf13/cobra"
)

const (
	envoyFlagSource     = "source"
	envoyFlagProto      = "proto"
	envoyFlagOutDesc    = "out"
	envoyFlagDefaultVal = ""
)

// GenerateEnvoyTranscodingCMD represents the new command.
var GenerateEnvoyTranscodingCMD = &cobra.Command{
	Use:     "envoy",
	Short:   "Build Envoy: transcoder: gRPC-JSON transcoder in Base 64",
	Long:    "Build Envoy: Generates transcode from proto file to use filter which allows a RESTful JSON API client to send requests to Envoy for gRPC",
	Example: "scanx envoy --source=/home/user/project/proto --proto=org/project/my-api/v1/my-api.proto --out=my-api.proto-description",
	Run:     runEnvoyTranscodingGen,
}

func init() {
	GenerateEnvoyTranscodingCMD.PersistentFlags().String(envoyFlagSource, envoyFlagDefaultVal, "Source a root folder where proto files are stored")
	GenerateEnvoyTranscodingCMD.PersistentFlags().String(envoyFlagProto, envoyFlagDefaultVal, "Proto file to convert to descriptor that must be used in Envoy")
	GenerateEnvoyTranscodingCMD.PersistentFlags().String(envoyFlagOutDesc, envoyFlagDefaultVal, "Where proto description must be created")
}

func runEnvoyTranscodingGen(c *cobra.Command, args []string) {
	input := []string{"-I"}
	if pErr := c.ParseFlags(args); pErr != nil {
		log.PrintError(fmt.Sprintf("failed to parse envoy command flags, error: %s", pErr.Error()))
		return
	}

	log.PrintInfo("Generating proto descriptor for envoy...")

	cPath, cPathHasErr := handler.FindCurrentFolder()
	if cPathHasErr {
		log.PrintError("unable to resolve current directory path...")
		return
	}

	protoDescriptionFileName := "proto.description"
	protoDescriptionFileName, input = envoyCmdParseInput(c, protoDescriptionFileName, input, cPath)

	fd := exec.Command("protoc", input...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = "."

	if err := fd.Run(); err != nil {
		log.PrintError(fmt.Sprintf("failed to run protoc, error: %s", err.Error()))
		return
	}

	if envoyCmdEncodeToBase64(cPath, protoDescriptionFileName) {
		return
	}

	log.PrintInfo("Done...")
}

func envoyCmdParseInput(c *cobra.Command, protoDescriptionFileName string, input []string, cPath string) (string, []string) {
	if v, _ := c.Flags().GetString(envoyFlagOutDesc); v != envoyFlagDefaultVal {
		protoDescriptionFileName = filepath.Base(v)
	}
	if v, _ := c.Flags().GetString(envoyFlagSource); v != envoyFlagDefaultVal {
		input = append(input, v)
	} else {
		input = append(input, cPath)
	}

	input = append(input, "--descriptor_set_out="+protoDescriptionFileName)

	if v, _ := c.Flags().GetString(envoyFlagProto); v != envoyFlagDefaultVal {
		input = append(input, "--include_imports")
		input = append(input, v)
	}

	return protoDescriptionFileName, input
}

func envoyCmdEncodeToBase64(cPath string, protoDescriptionFileName string) (hasErr bool) {
	log.PrintInfo("Encode to Base64...")

	protoFile, protoFileErr := os.Open(path.Join(cPath, protoDescriptionFileName))
	if protoFileErr != nil {
		log.PrintError(fmt.Sprintf("failed to open %s to read content, error: %s...", protoDescriptionFileName, protoFileErr.Error()))
		return true
	}

	protoDescriptionContent, protoDescriptionContentErr := ioutil.ReadAll(bufio.NewReader(protoFile))
	if protoDescriptionContentErr != nil {
		log.PrintError("failed to read proto description content...")
	}
	if cErr := handler.CreateFile(cPath, protoDescriptionFileName, base64.StdEncoding.EncodeToString(protoDescriptionContent)); cErr != nil {
		log.PrintError(fmt.Sprintf("unable to create proto description file (%s), err: %s...", protoDescriptionFileName, cErr.Error()))
		return true
	}

	return false
}
