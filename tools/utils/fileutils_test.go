package toolsutils

import (
	"fmt"
	"testing"
)

func TestListDir(t *testing.T) {
	files, err := ListDir("e:\\goproject\\src\\github.com\\colefan\\gsgo\\gameprotocol", ".xml")
	if err != nil {
		fmt.Println("ListDir err,", err)
	}
	fmt.Println("ListDir,files:", files)
}

func TestWalkDir(t *testing.T) {
	files, err := WalkDir("e:\\goproject\\src\\github.com\\colefan\\gsgo\\gameprotocol", ".go")
	if err != nil {
		fmt.Println("WalkDir,error,", err)
	}
	fmt.Println("WorkDir,files", files)
}

func TestMkDir(t *testing.T) {
	err1 := MakeFile("e:\\goproject\\src\\github.com\\colefan\\gsgo\\gameprotocol\\login", "login2.go", "nihaoatest1")
	if err1 != nil {
		fmt.Println("TestMkdir error,", err1)
	}
	err2 := MakeFile("e:\\goproject\\src\\github.com\\colefan\\gsgo\\gameprotocol\\login3", "login3.go", "nihaoatest2")
	if err2 != nil {
		fmt.Println("TestMkdir error,", err2)
	}
}
