package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var bucket = &cobra.Command{
	Use:   "bucket",
	Short: "Does operation on the bucket",
	Long:  "can make the bucket, delete the bucket, modify the meta-data of the bucket, list the objects in the bucket",

	Run: opBucket,
}

func init() {
	rootCmd.AddCommand(bucket)

	bucket.Flags().StringP("make", "m", "", "makes the bucket at Folder/NewBucket")
	bucket.Flags().StringP("modify", "M", "", "modify the meta-data -> bucket.key.new-value")
	bucket.Flags().StringP("list objects", "l", "", "list the objects in the bucket bucketName")
	bucket.Flags().StringP("delete", "d", "", "delete the bucket")
}

func opBucket(cmd *cobra.Command, args []string) {

	make, _ := cmd.Flags().GetString("make")
	modify, _ := cmd.Flags().GetString("modify")
	list, _ := cmd.Flags().GetString("list objects")
	delete, _ := cmd.Flags().GetString("delete")

	if make != "" {
		fmt.Println("Making the bucket")
		fmt.Println("Bucket name: ", cmd.Flag("make").Value)
		path := cmd.Flag("make").Value.String()
		fmt.Println("Path:", path)
		makeBucket(path)
	}
	if modify != "" {
		fmt.Println("Modifying the bucket")
		bucketName := cmd.Flag("modify").Value.String()
		fmt.Println("Bucket name:", bucketName)
		modifyMetaData(bucketName)
	}
	if list != "" {
		fmt.Println("Listing the objects in the bucket")
		bucketName := cmd.Flag("list objects").Value.String()
		listObjects(bucketName)
	}
	if delete != "" {
		fmt.Println("Deleting the bucket")
		fmt.Println("Are you sure? [N/Y]")
		var choice string
		fmt.Scanln(&choice)
		if choice == "N" {
			return
		}
		bucketName := cmd.Flag("delete").Value.String()
		deleteBucket(bucketName)
	}
}

func listObjects(bucketName string) {
	fmt.Println("Listing objects in the bucket", bucketName)
	//get the name of objects from database
}

func makeBucket(path string) {
	folder := strings.Split(path, "/")[0]
	file := strings.Split(path, "/")[1]

	directory := filepath.Join("C:/Users/HEET B JHAVERI/OneDrive", "Desktop", folder)
	filePath := filepath.Join(directory, file+".json")

	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	createdBy := ""
	createdAt := ""
	LastModifiedBy := ""
	LastModifiedAt := ""

	fmt.Print("Enter your name: ")
	fmt.Scanln(&createdBy)
	LastModifiedBy = createdBy
	createdAt = time.Now().String()
	LastModifiedAt = createdAt

	objects := make(map[string]interface{})
	data := map[string]interface{}{
		"objects": objects,
		"metaData": map[string]interface{}{
			"createdBy":      createdBy,
			"createdAt":      createdAt,
			"LastModifiedBy": LastModifiedBy,
			"LastModifiedAt": LastModifiedAt,
		},
	}

	newFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	defer newFile.Close()

	encoder := json.NewEncoder(newFile)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("JSON file created at", filePath)
}

func modifyMetaData(bucketName string) {

	key := strings.Split(bucketName, ".")[1]
	newValue := strings.Split(bucketName, ".")[2]
	bucketName = strings.Split(bucketName, ".")[0]

	filePath := "C:/Users/HEET B JHAVERI/OneDrive/Desktop/hi/" + bucketName + ".json"

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(file, &jsonData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	metaData, ok := jsonData["metaData"].(map[string]interface{})
	if !ok {
		fmt.Println("metaData not found or is not a valid map")
		return
	}

	metaData["LastModifiedAt"]=time.Now().String()
	metaData[key] = newValue

	modifiedJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = ioutil.WriteFile(filePath, modifiedJSON, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("JSON file successfully modified")
}

func deleteBucket(bucketName string) {
	filePath := "C:/Users/HEET B JHAVERI/OneDrive/Desktop/hi/" + bucketName + ".json"
	err := os.Remove(filePath)
	if err != nil {
		fmt.Println("Error deleting file:", err)
		return
	}
}
