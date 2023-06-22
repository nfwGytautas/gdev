package file

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

// Check if file exists
func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// Create file if not exists
func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if fileInfo.IsDir() {
		return true, nil
	} else {
		return false, nil
	}
}

/*
Get all directories inside the root directory
*/
func GetDirectories(root string) ([]string, error) {
	var files []string

	f, err := os.Open(root)

	if err != nil {
		return files, err
	}

	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			files = append(files, file.Name())
		}
	}

	return files, nil
}

// Write a template to file
func WriteTemplate(path, templateString string, data any) error {
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	t, err := template.New("template").Parse(templateString)
	if err != nil {
		return err
	}

	out := &bytes.Buffer{}
	err = t.Execute(out, data)
	if err != nil {
		return err
	}

	_, err = file.Write(out.Bytes())
	if err != nil {
		return err
	}
	file.Sync()

	return nil
}

// Append a template to file
func AppendTemplate(path, templateString string, data any) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	t, err := template.New("template").Parse(templateString)
	if err != nil {
		return err
	}

	out := &bytes.Buffer{}
	err = t.Execute(out, data)
	if err != nil {
		return err
	}

	_, err = file.Write(out.Bytes())
	if err != nil {
		return err
	}
	file.Sync()

	return nil
}

// Append content to file
func Append(file, content string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err := f.WriteString(content); err != nil {
		return err
	}

	err = f.Close()
	return err
}

// PRIVATE FUNCTIONS
// ========================================================================
