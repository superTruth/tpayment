package fileutils

import (
	"fmt"
	"testing"
)

func TestSeparateFilePath(t *testing.T) {
	dirPath, fileName, _ := SeparateFilePath("/12345/54321/asdfa.apk")

	fmt.Println("dirPath->", dirPath, ", fileName->", fileName)

}
