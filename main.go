package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ExtensionReplacement struct {
	OldExt string
	NewExt string
}

type TextReplacement struct {
	Extension string
	OldStr    string
	NewStr    string
}

type FileConverter struct {
	SourceDirectory      string
	DestinationDirectory string
	Extensions           []ExtensionReplacement
	TextReplacements     []TextReplacement
}

func (fc *FileConverter) CopyDirectory() error {
	err := filepath.Walk(fc.SourceDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(fc.SourceDirectory, path)
		if err != nil {
			return err
		}

		destinationPath := filepath.Join(fc.DestinationDirectory, relativePath)

		if info.IsDir() {
			err := os.MkdirAll(destinationPath, info.Mode())
			if err != nil {
				return err
			}
		} else {
			err := fc.CopyFile(path, destinationPath)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (fc *FileConverter) CopyFile(sourcePath, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func (fc *FileConverter) RenameFiles() error {
	err := filepath.Walk(fc.DestinationDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		extension := filepath.Ext(path)
		for _, replacement := range fc.Extensions {
			if extension == replacement.OldExt {
				newPath := strings.TrimSuffix(path, replacement.OldExt) + replacement.NewExt
				err := os.Rename(path, newPath)
				if err != nil {
					return err
				}
				break
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("File renaming completed successfully.")
	return nil
}

func (fc *FileConverter) ReplaceTextInFiles() error {
	err := filepath.Walk(fc.DestinationDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		extension := filepath.Ext(path)
		for _, replacement := range fc.TextReplacements {
			if extension == replacement.Extension {
				content, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				newContent := strings.ReplaceAll(string(content), replacement.OldStr, replacement.NewStr)

				err = ioutil.WriteFile(path, []byte(newContent), info.Mode())
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("Text replacement completed successfully.")
	return nil
}

func main() {
	// 从命令行参数中获取目录路径
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <source_directory>")
		return
	}
	sourceDirectory := os.Args[1]
	destinationDirectory := filepath.Join(filepath.Dir(sourceDirectory), filepath.Base(sourceDirectory)+"_tt")

	// 检查源目录是否存在
	if _, err := os.Stat(sourceDirectory); os.IsNotExist(err) {
		fmt.Println("Source directory does not exist.")
		return
	}

	// 创建 FileConverter 对象
	converter := FileConverter{
		SourceDirectory:      sourceDirectory,
		DestinationDirectory: destinationDirectory,
		Extensions: []ExtensionReplacement{
			{OldExt: ".wxml", NewExt: ".ttml"},
			{OldExt: ".wxss", NewExt: ".ttss"},
		},
		TextReplacements: []TextReplacement{
			{Extension: ".ttml", OldStr: "wx:", NewStr: "tt:"},
			{Extension: ".ttss", OldStr: ".wxss", NewStr: ".ttss"},
			{Extension: ".js", OldStr: "wx.", NewStr: "tt."},
			{Extension: ".ttml", OldStr: ".wxml", NewStr: ".ttml"},
			{Extension: ".ts", OldStr: "wx.", NewStr: "tt."},
		},
	}

	// 复制源文件夹到目标文件夹
	err := converter.CopyDirectory()
	if err != nil {
		fmt.Println("Error copying directory:", err)
		return
	}

	// 重命名文件
	err = converter.RenameFiles()
	if err != nil {
		fmt.Println("Error renaming files:", err)
		return
	}

	// 替换文本内容
	err = converter.ReplaceTextInFiles()
	if err != nil {
		fmt.Println("Error replacing content:", err)
		return
	}

	fmt.Println("File conversion completed successfully.")
}
