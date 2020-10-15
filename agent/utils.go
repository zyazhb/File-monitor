package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"

	"github.com/google/logger"
)

// GetAllFile 获取目录中所有文件
func GetAllFile(pathname string, level int) ([]string, error) {
	filenames := []string{}

	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return filenames, err
	}
	filenames = append(filenames, pathname)
	for _, fi := range rd {
		if fi.IsDir() && level > 1 {
			logger.Info("\033[1;32m [+]Find dir: " + pathname + fi.Name() + " \033[0m")
			level--
			newfilenames, err := GetAllFile(pathname+fi.Name()+"/", level)
			if err != nil {
				logger.Error("\033[1;31m [-]Error: ", err, " \033[0m")
			}
			filenames = append(filenames, newfilenames...)
		} else {
			logger.Info("\033[1;32m [+]Find file: " + pathname + fi.Name() + " \033[0m")
			filenames = append(filenames, pathname+fi.Name())
		}
	}
	return filenames, nil
}

// calcHash 计算文件sha256hash
func calcHash(filename string) string {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		logger.Error("\033[1;31m [-]Can't read the file! \033[0m")
		return ""
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		logger.Error("\033[1;31m ", err, "\033[0m")
		return "directory"
	}
	sum := hash.Sum(nil)

	return string(hex.EncodeToString(sum))
}
