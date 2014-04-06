package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/robfig/revel"
	"io"
	"io/ioutil"
	// "os"
)

type FileUploader struct {
	*revel.Controller
}

func (file *FileUploader) Upload(files io.Reader, localFilePath string) revel.Result {
	fmt.Println(file.Params.Files, files)
	// os.OpenFile(name, flag, perm)
	//files.Write()
	// for p, v := range file.Params {
	// 	fmt.Println(p, v)
	// }
	// file, handler, err := file.Params.FormFile("file")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	data, err := ioutil.ReadAll(files)
	fmt.Println(":$$$$$$$$$$$$$$", files, err)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("handler.Filename.png", data, 0777)
	if err != nil {
		fmt.Println(err)
	}

	localFile := "E:/GitHubRepos/gocode/src/SanWenJia/public/img/3.jpg"

	fmt.Println(localFilePath)
	key, e := models.UploadToQiniu(localFile)

	fmt.Println(key, e)
	return file.Render()
}
