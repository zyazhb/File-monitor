package main

import (
	"crypto/sha256"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// GetAllFile 获取目录中所有文件
func GetAllFile(pathname string, level int) ([]string, error) {
	filenames := []string{}

	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return filenames, err
	}

	for _, fi := range rd {
		if fi.IsDir() && level > 1 {
			log.Printf("[+]Find dir: " + pathname + fi.Name())
			level--
			newfilenames, err := GetAllFile(pathname+fi.Name()+"/", level)
			if err != nil {
				log.Print("[-]Error: ", err)
			}
			filenames = append(filenames, newfilenames...)
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
