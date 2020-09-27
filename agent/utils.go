package main

import (
	"crypto/sha256"
	"io/ioutil"
	"io"
	"os"
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

// calcHash 计算文件sha256hash
func calcHash(filename string) []byte {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Println("文件读取失败！")
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatalln(err)
	}
	sum := hash.Sum(nil)

	return sum
}