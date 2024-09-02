package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var object = &cobra.Command{
	Use:   "object",
	Short: "Does operation on the object",
	Long:  "can make the object, delete the object, modify the meta-data of the object, retrive the object, replace the object",

	Run: opObject,
}

func init() {
	rootCmd.AddCommand(object)

	object.Flags().BoolP("make", "m", false, "makes the bucket at NewBucket/object")
	object.Flags().StringP("modify", "M", "", "modify the meta-data -> bucket.object.key.new-value")
	object.Flags().StringP("replace", "R", "", "replace the object")
	object.Flags().StringP("delete", "d", "", "delete the object")
	object.Flags().StringP("retrive", "r", "", "delete the object")
}

func opObject(cmd *cobra.Command, args []string) {

	make, _ := cmd.Flags().GetBool("make")
	fmt.Println("Make:", make)
	// modify, _ := cmd.Flags().GetString("modify")
	// replace, _ := cmd.Flags().GetString("replace")
	// retrive, _ := cmd.Flags().GetString("retrive")
	// delete, _ := cmd.Flags().GetString("delete")

	if make {
		//filePath := cmd.Flag("make").Value.String()
		fmt.Println("Making the object")
		//fmt.Println(filePath)
		var filePath string
		fmt.Println("Enter the file path: ")
		fmt.Scanln(&filePath)
		filePath = filepath.Join("C:/Users/HEET B JHAVERI/OneDrive", "Desktop", "Heet_Jhaveri.pdf")
		jsonOutput, err := FileToJSON(filePath)
		if err != nil {
			log.Fatal(err)
		}

		jsonInput := jsonOutput

		// Parse JSON into FileData struct
		var fileData FileData
		err = json.Unmarshal([]byte(jsonInput), &fileData)
		if err != nil {
			log.Fatalf("Failed to parse JSON: %v", err)
		}

		// Specify output directory (current directory in this case)
		outputDir := "."

		// Save file from FileData
		err = SaveFile(fileData, outputDir)
		if err != nil {
			log.Fatalf("Failed to save file: %v", err)
		}
	}
}

type FileData struct {
	Data     string            `json:"data"`
	MetaData map[string]string `json:"meta-data"`
}

func ReadFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func ExtractMetaData(filePath string) (map[string]string, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	metaData := map[string]string{
		"file_name":       fileInfo.Name(),
		"file_size_bytes": fmt.Sprintf("%d", fileInfo.Size()),
		"modified_time":   fileInfo.ModTime().Format(time.RFC3339),
		"file_type":       mime.TypeByExtension(filepath.Ext(filePath)),
	}

	return metaData, nil
}

func FileToJSON(filePath string) ([]byte, error) {

	fmt.Println(filePath)
	data, err := ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	metaData, err := ExtractMetaData(filePath)
	if err != nil {
		return nil, err
	}

	fileData := FileData{
		Data:     string(data),
		MetaData: metaData,
	}

	jsonData, err := json.MarshalIndent(fileData, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func SaveFile(fileData FileData, outputDir string) error {
	// Extract file name from metadata
	fileName, exists := fileData.MetaData["file_name"]
	if !exists {
		return fmt.Errorf("file_name not found in metadata")
	}

	// Define output file path
	outputPath := filepath.Join(outputDir, fileName)

	// Write data to file
	err := ioutil.WriteFile(outputPath, []byte(fileData.Data), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	// Restore modification time if available
	if modTimeStr, exists := fileData.MetaData["modified_time"]; exists {
		modTime, err := time.Parse(time.RFC3339, modTimeStr)
		if err != nil {
			return fmt.Errorf("failed to parse modified_time: %v", err)
		}
		err = os.Chtimes(outputPath, time.Now(), modTime)
		if err != nil {
			return fmt.Errorf("failed to set modified_time: %v", err)
		}
	}

	// Restore permissions if available
	if permStr, exists := fileData.MetaData["permissions"]; exists {
		permInt, err := strconv.ParseUint(permStr, 8, 32)
		if err != nil {
			return fmt.Errorf("failed to parse permissions: %v", err)
		}
		err = os.Chmod(outputPath, os.FileMode(permInt))
		if err != nil {
			return fmt.Errorf("failed to set permissions: %v", err)
		}
	}

	fmt.Printf("File successfully reconstructed at: %s\n", outputPath)
	return nil
}
