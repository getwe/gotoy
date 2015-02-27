//用于同步mac上的iTunes音乐到sd卡上
package main

import (
	"bufio"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	NoDestPath = errors.New("no such file or directory")
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage : %s srcPath destPath\n", os.Args[0])
		os.Exit(1)
	}

	srcPath := os.Args[1]
	destPath := os.Args[2]

	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		fmt.Println(NoDestPath)
		os.Exit(1)
	}

	log.Printf("开始计算目标文件夹md5签名")
	destSign, err := destPathCreateSign(destPath)
	if err != nil {
		fmt.Println(err)
	}

	srcFileList := make([]string, 0, 10240)
	err = filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		if !strings.HasSuffix(info.Name(), ".mp3") {
			return nil
		}

		srcFileList = append(srcFileList, path)

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	Shuffle(srcFileList)

	newSign, err := os.Create(filepath.Join(destPath, "md5.txt"))
	defer newSign.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	newSignIO := bufio.NewWriter(newSign)
	defer newSignIO.Flush()

	num := 0
	for _, file := range srcFileList {
		srcSign, err := signOneFile(file)
		if err != nil {
			continue
		}

		newName, err := getFileName(file, num)
		if err != nil {
			continue
		}
		num++

		destFile := filepath.Join(destPath, newName)

		newSignIO.WriteString(fmt.Sprintf("%s\t%s\n", srcSign, newName))

		srcFileName, ok := destSign[srcSign]
		if !ok {
			// cp
			copyFile(file, destFile)
			log.Printf("copy file [%s]\n", file)
		} else {
			// mv
			srcFile := filepath.Join(destPath, srcFileName)
			os.Rename(srcFile, destFile)
			log.Printf("mv file [%s]\n", srcFile)
		}
		delete(destSign, srcSign)
	}

	log.Printf("文件同步完成,开始删除操作")

	for _, v := range destSign {
		log.Printf("rm file [%s]\n", v)
		os.Remove(v)
	}
}

func destPathCreateSign(destPath string) (map[string]string, error) {
	// md5 -> fileName
	sign := make(map[string]string)

	md5Txt := filepath.Join(destPath, "md5.txt")
	if _, err := os.Stat(md5Txt); err == nil {

		log.Printf("找到缓存签名文件[%s],直接使用", md5Txt)

		md5fh, err := os.Open(md5Txt)
		defer md5fh.Close()
		if err == nil {
			s := bufio.NewScanner(md5fh)
			for s.Scan() {
				arr := strings.Split(s.Text(), "\t")
				sign[arr[0]] = arr[1]
			}
			return sign, nil
		}
	}

	log.Printf("没有签名缓存文件[%s],对整个文件夹重新计算", md5Txt)

	err := filepath.Walk(destPath, func(path string, info os.FileInfo, err error) error {
		//fmt.Println(path,info,err)
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}

		if !strings.HasSuffix(info.Name(), ".mp3") {
			return nil
		}

		signkey, err := signOneFile(path)
		if err != nil {
			return nil
		}

		sign[signkey] = filepath.Base(path)
		return nil
	})

	if err != nil {
		return sign, err
	}
	return sign, nil
}

func signOneFile(path string) (string, error) {

	log.Printf("计算文件[%s]md5", path)

	fh, err := os.Open(path)
	defer fh.Close()
	if err != nil {
		return "", err
	}

	mh := md5.New()
	io.Copy(mh, fh)

	sign := mh.Sum(nil)

	return fmt.Sprintf("%x", sign), nil
}

// 随机打乱数组
func Shuffle(arr []string) {
	rand.Seed(time.Now().UnixNano())

	N := len(arr)
	for i := 0; i < N; i++ {
		idx := rand.Intn(i + 1)
		arr[i], arr[idx] = arr[idx], arr[i]
	}
}

func getFileName(path string, n int) (string, error) {
	f := filepath.Base(path)

	name := strings.TrimLeft(f, " 0123456789-")
	return fmt.Sprintf("%04d %s", n, name), nil
}

func copyFile(from, to string) error {
	f, err := os.Open(from)
	defer f.Close()
	if err != nil {
		return err
	}

	t, err := os.Create(to)
	defer t.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(t, f)
	if err != nil {
		return err
	}

	return nil
}
