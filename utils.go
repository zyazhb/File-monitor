package main

import (
	"io/ioutil"
	"log"
)

// GetAllFile 获取目录中所有文件
func GetAllFile(pathname string) ([]string, error) {
	filenames := []string{}

	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return filenames, err
	}

	for _, fi := range rd {
		if fi.IsDir() {
			log.Printf("[%s]\n", pathname+"\\"+fi.Name())
			GetAllFile(pathname + fi.Name() + "\\")
		} else {
			log.Printf("[+]Find file: " + pathname + fi.Name())
			filenames = append(filenames, pathname+fi.Name())
		}
	}
	return filenames, nil
}

// func printlog(info string, data string, err error) {
// 	errexist := value.(err)
// 	if err {

// 	} else {
// 		log.Println(info, data)
// 	}
// }
