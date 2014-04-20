package controllers

import (
	// "SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	// "strconv"
	"strings"
)

type Admin struct {
	*revel.Controller
}

func (this *Admin) Index() revel.Result {
	return this.Render()
}

func (this *Admin) UploadPicture() revel.Result {
	return this.Render()
}

func (this *Admin) PostUploadPicture(pictures []*os.File) revel.Result {
	// 使用revel requst formfile获取文件数据
	file, handler, err := this.Request.FormFile("pictures")
	if err != nil {
		fmt.Println(err)
	}
	// 读取所有数据
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	// 获取当前路径
	dir, patherr := filepath.Abs(filepath.Dir(os.Args[0]))
	if patherr != nil {
		log.Fatal(patherr)
	}

	// 文件路径
	dir = strings.Replace(dir, "bin", "", 1)
	filePath := dir + "/" + "src/SanWenJia/public/img/" + handler.Filename

	// 保存到文件
	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		fmt.Println(err)
	}

	return this.Redirect((*Admin).UploadPicture)
}
