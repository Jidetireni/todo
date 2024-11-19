package filemanger

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type FileManager struct {
	inputFile  string
	outputFile string
}

func (fm *FileManager) ReadTasksToFIle(tasks interface{}) error {
	file, err := os.Open(fm.inputFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(tasks)
	if err != nil {
		return fmt.Errorf("error decoding tasks from file: %w", err)
	}
	return nil

}

func (fm *FileManager) WriteTaskToFile(tasks interface{}) error {
	file, err := os.OpenFile(fm.outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		return errors.New("failed to create file")
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(tasks)

	if err != nil {
		return errors.New("failed to convert data to json")
	}
	return nil

}

func New(inputFile, outputFile string) *FileManager {
	return &FileManager{
		inputFile:  inputFile,
		outputFile: outputFile,
	}
}
