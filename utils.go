package main

import (
	"fmt"
	"io/ioutil"
)

// GetAllFile 1
func GetAllFile(pathname string) ([]string, error) {
	filenames := []string{}

	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return filenames, err
	}

	for _, fi := range rd {
		if fi.IsDir() {
			fmt.Printf("[%s]\n", pathname+"\\"+fi.Name())
			GetAllFile(pathname + fi.Name() + "\\")
		} else {
			filenames = append(filenames, fi.Name())
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
