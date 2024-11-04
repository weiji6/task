package savePaper

import (
	"fmt"
	"github.com/spf13/afero" //需要使用go.mod进行管理
)

// SavePaper 函数将指定的文本保存到给定路径的文件中，使用 afero 提供的文件系统
func SavePaper(path string, text string) error {
	// 使用本地文件系统
	fs := afero.NewOsFs()
	// 创建或打开指定路径的文件
	file, err := fs.Create(path)
	if err != nil {
		// 如果发生错误，打印错误并返回
		fmt.Println("Failed to create file:", err)
		return nil
	}
	defer file.Close()

	// 写入文本内容到文件
	_, err = file.WriteString(text)
	if err != nil {
		// 如果写入发生错误，打印错误信息
		fmt.Println("Failed to write to file:", err)
		return nil
	}

	// 文件保存成功的信息
	fmt.Println("File saved successfully at", path)
	return nil
}
