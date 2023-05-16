package file

import (
	"bytes"
	"html/template"
	"io/fs"
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
	_, err := os.Stat(path)
	return err == nil
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

// Copy source to target
func CopyFile(source, target string) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	err = os.WriteFile(target, data, fs.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// PRIVATE FUNCTIONS
// ========================================================================
