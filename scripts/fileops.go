package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

type CustomFileInfo struct {
	FullPath             string
	StandardizedFileName string
	OriginalFileName     string
}

func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Printf("Failed to close file: %v", err)
	}
}

func RecursiveReadDir(dir string, fileMap map[string]*CustomFileInfo) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			err := RecursiveReadDir(filepath.Join(dir, entry.Name()), fileMap)
			if err != nil {
				return err
			}
		} else {
			cfi := GetCustomFileInfo(dir, entry)
			if cfi != nil {
				fileMap[cfi.StandardizedFileName] = cfi
			}
		}
	}

	return nil
}

func GetCustomFileInfo(dir string, de os.DirEntry) *CustomFileInfo {
	if de.IsDir() {
		panic(fmt.Errorf("entry is a directory"))
	}

	sfn, err := Standard(de.Name())
	if err != nil {
		return nil
	}

	return &CustomFileInfo{
		FullPath:             filepath.Join(dir, de.Name()),
		StandardizedFileName: sfn,
		OriginalFileName:     de.Name(),
	}
}

const DatePrefix = `^\d{4}-\d{2}-\d{2}`
const DatePrefixNoDelimiters = `^(\d{8})\d*-pub-(.+)$`

func Standard(fileName string) (string, error) {
	if matched, _ := regexp.MatchString(DatePrefix, fileName); matched {
		return fileName, nil
	}

	re := regexp.MustCompile(DatePrefixNoDelimiters)
	matches := re.FindStringSubmatch(fileName)
	if matches == nil {
		return "", fmt.Errorf("filename does not match expected pattern")
	}

	datePart := matches[1] // First 8 digits
	rest := matches[2]     // The remaining part (title and extension)

	// Reformat date: "YYYYMMDD" -> "YYYY-MM-DD"
	formattedDate := fmt.Sprintf("%s-%s-%s", datePart[:4], datePart[4:6], datePart[6:8])

	// Return combined result
	return formattedDate + "-" + rest, nil
}

func SyncFiles(srcfi *CustomFileInfo, dstfi *CustomFileInfo) error {
	srcf, err := os.Open(srcfi.FullPath)
	if err != nil {
		return fmt.Errorf("failed to open src file: %v", err)
	}
	defer CloseFile(srcf)

	dstf, err := os.Open(dstfi.FullPath)
	if err != nil {
		return fmt.Errorf("failed to open dst file: %v", err)
	}
	defer CloseFile(dstf)

	srcb := new(bytes.Buffer)
	_, err = io.Copy(srcb, srcf)
	if err != nil {
		return fmt.Errorf("failed to read src file to buffer: %v", err)
	}

	dstb := new(bytes.Buffer)
	_, err = io.Copy(dstb, dstf)
	if err != nil {
		return fmt.Errorf("failed to read dst file to buffer: %v", err)
	}

	eq := bytes.Equal(srcb.Bytes(), dstb.Bytes())
	if !eq {
		dstfw, err := os.OpenFile(dstfi.FullPath, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("failed to open dst file for writing: %v", err)
		}
		defer CloseFile(dstfw)
		_, err = srcb.WriteTo(dstfw)
		if err != nil {
			return fmt.Errorf("failed to write file: %v", err)
		}
	}

	if eq {
		fmt.Println("Files are equal, no sync needed for", srcfi.OriginalFileName)
		return nil
	}
	fmt.Println("Synced file", srcfi.OriginalFileName)
	return nil
}

func CreatePost(srcfi *CustomFileInfo, dstDir string) error {
	srcf, err := os.Open(srcfi.FullPath)
	if err != nil {
		return fmt.Errorf("failed to open src file: %v", err)
	}
	defer CloseFile(srcf)

	srcb := new(bytes.Buffer)
	_, err = io.Copy(srcb, srcf)
	if err != nil {
		return fmt.Errorf("failed to read src file to buffer: %v", err)
	}

	dstFileName := filepath.Join(dstDir, srcfi.StandardizedFileName)
	dstfw, err := os.OpenFile(dstFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open dst file for writing: %v", err)
	}
	defer CloseFile(dstfw)

	_, err = srcb.WriteTo(dstfw)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	fmt.Println("Created new post file", srcfi.StandardizedFileName)
	return nil
}
