package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func createBaseFolder(baseFolder string) {
	err := os.MkdirAll(baseFolder, os.ModePerm)
	if err != nil {
		log.Printf("Error creating base directory %s: %v", baseFolder, err)
	}
}

func createFolderStructure(baseFolder string, minFiles int) {
	subFolders := []string{"documents", "pictures", "videos", "music", "archives"}
	for _, folder := range subFolders {
		createSubFolder(baseFolder, folder, minFiles)
	}
}

func createSubFolder(baseFolder, folderName string, minFiles int) error {
	path := filepath.Join(baseFolder, folderName)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Printf("Error creating directory %s: %v", path, err)
		return err
	}

	for i := 0; i < minFiles; i++ {
		_ = createFile(path, fmt.Sprintf("file%d.txt", i))
	}

	return nil
}

func createFile(folderPath, fileName string) error {
	file, err := os.Create(filepath.Join(folderPath, fileName))
	if err != nil {
		log.Printf("Error creating file %s: %v", filepath.Join(folderPath, fileName), err)
		return err
	}

	_, err = file.WriteString(strings.Repeat(fmt.Sprintf("Data in %s %s\n", folderPath, fileName), 100))
	if err != nil {
		log.Printf("Error writing to file %s: %v", filepath.Join(folderPath, fileName), err)
		return err
	}

	return file.Close()
}

func performRandomOperations(baseFolder string) {
	filepath.Walk(baseFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing a path %s: %v", path, err)
			return nil // continue to the next item
		}
		if path != baseFolder && info.IsDir() {
			if err := performOperationOnDirectory(path); err != nil {
				log.Printf("Error performing operation on directory %s: %v", path, err)
			}
		}
		return nil
	})
}

func performOperationOnDirectory(path string) error {
	switch seededRand.Intn(4) {
	case 0:
		return modifyFilesInDirectory(path)
	case 1:
		return createSubFolder(path, fmt.Sprintf("folder_new_%s", randomString(5)), seededRand.Intn(5)+1)
	case 2:
		return removeDataFromFilesInDirectory(path)
	case 3:
		return os.RemoveAll(path)
	}
	return nil
}

func modifyFilesInDirectory(directoryPath string) error {
	var lastError error
	filepath.Walk(directoryPath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if fileInfo.IsDir() {
			return nil
		}
		if err := modifyFile(filePath); err != nil {
			log.Printf("Error modifying file %s: %v", filePath, err)
			lastError = err
		}
		return nil
	})
	return lastError
}

func modifyFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading the file %s: %v", filePath, err)
	}
	lines := strings.Split(string(content), "\n")
	if len(lines) > 1 {
		lines[1] = randomString(20) // modify the second line
	}
	err = os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0)
	if err != nil {
		return fmt.Errorf("error writing to the file %s: %v", filePath, err)
	}
	return nil
}

func removeDataFromFilesInDirectory(directoryPath string) error {
	var lastError error
	filepath.Walk(directoryPath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if fileInfo.IsDir() {
			return nil
		}
		if err := removeDataFromFile(filePath); err != nil {
			log.Printf("Error removing data from file %s: %v", filePath, err)
			lastError = err
		}
		return nil
	})
	return lastError
}

func removeDataFromFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading the file %s: %v", filePath, err)
	}
	lines := strings.Split(string(content), "\n")
	if len(lines) > 1 {
		lines = lines[:len(lines)-2] // remove the last two lines
	}
	err = os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0)
	if err != nil {
		return fmt.Errorf("error writing to the file %s: %v", filePath, err)
	}
	return nil
}

func main() {
	baseFolder := flag.String("base", "testData", "Base folder for the test data")
	minFiles := flag.Int("minFiles", 5, "Minimum number of files to be created in each folder")
	randomOps := flag.Bool("randomOps", true, "Perform random operations on the data")

	flag.Parse()

	createBaseFolder(*baseFolder)
	createFolderStructure(*baseFolder, *minFiles)

	if *randomOps {
		performRandomOperations(*baseFolder)
	}

	fmt.Println("Data generation and modification completed successfully.")
}
