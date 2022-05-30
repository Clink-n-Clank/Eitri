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
	envoyFlagIsEncoded  = "isBase64Encoded"
	envoyFlagDefaultVal = ""
)

// GenerateEnvoyTranscodingCMD represents the new command.
var GenerateEnvoyTranscodingCMD = &cobra.Command{
	Use:     "envoy",
	Short:   "Build Envoy: transcoder: gRPC-JSON transcoder in Base 64",
	Long:    "Build Envoy: Generates transcode from proto file to use filter which allows a RESTful JSON API client to send requests to Envoy for gRPC",
	Example: "Eitri envoy --out envoy/descriptors/virtual_pos_definition.pb --proto workspace/services-interfaces/coopnorge/scanandpay/vpos/v1beta1/virtual_pos_api.proto --source workspace/services-interfaces --isBase64Encoded",
	Run:     runEnvoyTranscodingGen,
}

func init() {
	GenerateEnvoyTranscodingCMD.PersistentFlags().String(envoyFlagSource, envoyFlagDefaultVal, "Source a root folder where proto files are stored")
	GenerateEnvoyTranscodingCMD.PersistentFlags().String(envoyFlagProto, envoyFlagDefaultVal, "Proto file to convert to descriptor that must be used in Envoy")
	GenerateEnvoyTranscodingCMD.PersistentFlags().String(envoyFlagOutDesc, envoyFlagDefaultVal, "Where proto description must be created")
	GenerateEnvoyTranscodingCMD.PersistentFlags().Bool(envoyFlagIsEncoded, false, "Encode Envoy description to base 64 or not")
}

func runEnvoyTranscodingGen(c *cobra.Command, args []string) {
	input := []string{"-I"}
	if pErr := c.ParseFlags(args); pErr != nil {
		log.PrintError(fmt.Sprintf("failed to parse envoy command flags, error: %s", pErr.Error()))
		return
	}

	log.PrintInfo("Generating proto descriptor for envoy...")

	currentFolder, cPathHasErr := handler.FindCurrentFolder()
	if cPathHasErr {
		log.PrintError("unable to resolve current directory path...")
		return
	}

	protoDescriptionFileName := "proto.description"
	protoDescriptionFileName, input = envoyCmdParseInput(c, protoDescriptionFileName, input, currentFolder)

	fd := exec.Command("protoc", input...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = "."

	if err := fd.Run(); err != nil {
		log.PrintError(fmt.Sprintf("failed to run protoc, error: %s", err.Error()))
		return
	}

	isEncodedToBase64, _ := c.Flags().GetBool(envoyFlagIsEncoded)
	if envoySaveToFile(c, currentFolder, protoDescriptionFileName, isEncodedToBase64) {
		return
	}

	log.PrintInfo("Done...")
}

func envoyCmdParseInput(c *cobra.Command, protoDescriptionFileName string, input []string, currentFolder string) (string, []string) {
	if v, _ := c.Flags().GetString(envoyFlagOutDesc); v != envoyFlagDefaultVal {
		protoDescriptionFileName = filepath.Base(v)
	}
	if v, _ := c.Flags().GetString(envoyFlagSource); v != envoyFlagDefaultVal {
		input = append(input, v)
	} else {
		input = append(input, currentFolder)
	}

	input = append(input, "--descriptor_set_out="+protoDescriptionFileName)

	if v, _ := c.Flags().GetString(envoyFlagProto); v != envoyFlagDefaultVal {
		input = append(input, "--include_imports")
		input = append(input, v)
	}

	return protoDescriptionFileName, input
}

func envoySaveToFile(c *cobra.Command, pathToStoreFile, protoDescriptionFileName string, isEncodedToBase64 bool) (hasErr bool) {
	log.PrintInfo("Saving envoy proto description...")

	protoFile, protoFileErr := os.Open(path.Join(pathToStoreFile, protoDescriptionFileName))
	if protoFileErr != nil {
		log.PrintError(fmt.Sprintf("failed to open %s to read content, error: %s...", protoDescriptionFileName, protoFileErr.Error()))
		return true
	}

	protoDescriptionContent, protoDescriptionContentErr := ioutil.ReadAll(bufio.NewReader(protoFile))
	if protoDescriptionContentErr != nil {
		log.PrintError("failed to read proto description content...")
	}

	fileContent := string(protoDescriptionContent)
	if isEncodedToBase64 {
		log.PrintInfo("Encoding to Base64...")
		fileContent = base64.StdEncoding.EncodeToString(protoDescriptionContent)
	}

	if v, _ := c.Flags().GetString(envoyFlagOutDesc); v != envoyFlagDefaultVal {
		pathToStoreFile = v
	}

	fileInfo, fileInfoErr := os.Stat(pathToStoreFile)
	if fileInfoErr != nil {
		log.PrintError(fmt.Sprintf(
			"unable to resolve path to update or create file of envoy proto description, error: %s",
			fileInfoErr.Error(),
		))
		return
	}

	if fileInfo.IsDir() {
		if cErr := handler.CreateFile(pathToStoreFile, protoDescriptionFileName, fileContent); cErr != nil {
			log.PrintError(fmt.Sprintf("unable to create proto description file (%s), err: %s...", protoDescriptionFileName, cErr.Error()))
			return true
		}
	} else {
		if wErr := handler.OverwriteToFile(pathToStoreFile, fileContent); wErr != nil {
			log.PrintError(fmt.Sprintf("unable to overwirte proto description file (%s), err: %s...", protoDescriptionFileName, wErr.Error()))
			return true
		}
	}

	log.PrintInfo(fmt.Sprintf("Envoy proto description saved: %s", pathToStoreFile))

	return false
}
