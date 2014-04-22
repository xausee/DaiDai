package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func CopyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}

	return cerr
}

// 备份将要替换掉的图片
func BackupPicture() error {
	// 获取当前路径, 修改路径字串
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir := strings.Replace(path, "bin", "", 1)
	dir = dir + "/" + "src/SanWenJia/public/img/"

	// 改变当前路径，创建时间戳文件夹
	time_string := time.Now().Format("2006-01-02 15:04:05")
	time_string = strings.Replace(time_string, "-", "", -1)
	time_string = strings.Replace(time_string, " ", "", -1)
	time_string = strings.Replace(time_string, ":", "", -1)
	os.Chdir(dir)
	os.Mkdir(time_string, os.ModeDir)

	// 更新轮播的4个图片文件
	for i := 1; i < 5; i++ {
		src := dir + strconv.Itoa(i) + ".jpg"
		dst := dir + time_string + "/" + strconv.Itoa(i) + ".jpg"

		//保存到备份文件夹
		err = CopyFile(dst, src)
		if err != nil {
			fmt.Println("拷贝文件失败:", err)
		}

		// 删除原文件
		err = os.Remove(src)
		if err != nil {
			fmt.Println("删除原文件失败:", err)
		}
	}

	// 返回原始路径
	os.Chdir(path)

	return err
}

//保存上传的图片
func SavePicture(pictures []*os.File) (err error) {
	// 获取当前路径, 修改路径字串
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir := strings.Replace(path, "bin", "", 1)
	dir = dir + "/" + "src/SanWenJia/public/img/"

	for i, f := range pictures {
		// 获取文件信息
		fi, e := f.Stat()
		if e != nil {
			fmt.Println("获取文件信息失败")
		}

		// 根据文件大小分配空间
		buffer := make([]byte, fi.Size())
		n, e := f.Read(buffer)

		if e != nil {
			fmt.Println("读取文件失败")
		}
		// 创建文件
		file := dir + strconv.Itoa(i+1) + ".jpg"
		fo, e := os.Create(file)

		if e != nil {
			panic(e)
		}
		defer fo.Close()

		// 写入内容
		if _, e = fo.Write(buffer[:n]); e != nil {
			panic(e)
		}
		// 传递错误
		err = e
	}
	return err
}
