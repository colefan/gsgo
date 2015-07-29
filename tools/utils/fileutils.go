package toolsutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//遍历指定目录下，所有的文件，可用后缀进行过滤
//
func ListDir(dirpath string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirpath)

	if err != nil {
		return nil, err
	}
	pathSep := string(os.PathSeparator)
	dirpathhassep := strings.HasSuffix(dirpath, pathSep)
	suffix = strings.ToLower(suffix)
	for _, f := range dir {
		if f.IsDir() {
			continue
		}

		if strings.HasSuffix(strings.ToLower(f.Name()), suffix) {
			if dirpathhassep {
				files = append(files, dirpath+f.Name())
			} else {
				files = append(files, dirpath+pathSep+f.Name())
			}

		}
	}
	return files, nil
}

//遍历当前目录以及所有子目录，可以通过后缀过滤
func WalkDir(dirpath, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToLower(suffix)
	filepath.Walk(dirpath, func(filename string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		if strings.HasSuffix(strings.ToLower(f.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

func MakeFile(dir string, filename string, content string) error {
	fmt.Println("dir = ", dir, "filename = ", filename)
	if err := os.Chdir(dir); err != nil {
		err2 := os.Mkdir(dir, os.ModePerm)
		if err2 != nil {
			return err2
		}
	}
	newFilePath := dir
	if strings.LastIndex(dir, string(os.PathSeparator)) < len(dir)-1 {
		newFilePath += string(os.PathSeparator) + filename
	} else {
		newFilePath += filename
	}
	fmt.Println("==>", newFilePath)
	f, err := os.Create(newFilePath)

	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

func ChFileExt(preFilePath string, newPostExt string) string {
	index := strings.LastIndex(preFilePath, ".")
	if index > 0 {
		newFilePath := preFilePath[0:index]
		newFilePath += newPostExt
		return newFilePath
	}
	return ""
}

func GetFileName(fullfilepath string) string {
	index := strings.LastIndex(fullfilepath, string(os.PathSeparator))
	if index >= 0 {
		fullfilepath = fullfilepath[index+1:]
	}

	return fullfilepath

}
