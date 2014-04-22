package controllers

import (
	// "SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
	//"io/ioutil"
	//"log"
	"os"
	//"path/filepath"
	// "strconv"
	//"strings"
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

	for f := range pictures {
		fmt.Println(f)
	}

	fmt.Println(pictures[0])
	// fmt.Println(pictures[1])
	// fmt.Println(pictures[2])
	// fmt.Println(pictures[3])

	buf0 := make([]byte, 900000)
	fmt.Println(buf0)
	n0, err0 := pictures[0].Read(buf0)
	fmt.Println(buf0, n0, err0)

	// open output file
	fo, err := os.Create("output.jpg")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if erre := fo.Close(); erre != nil {
			panic(erre)
		}
	}()

	// write a chunk
	if _, erro := fo.Write(buf0[:n0]); erro != nil {
		panic(erro)
	}

	// buf1 := make([]byte, 20)
	// n1, err1 := pictures[0].Read(buf1)
	// fmt.Println(buf1, n1, err1)

	// buf2 := make([]byte, 20)
	// n2, err2 := pictures[0].Read(buf2)
	// fmt.Println(buf2, n2, err2)

	// buf3 := make([]byte, 20)
	// n3, err3 := pictures[0].Read(buf3)
	// fmt.Println(buf3, n3, err3)

	fmt.Println(len(pictures))
	// // 使用revel requst formfile获取文件数据
	// file, handler, err := this.Request.FormFile("pictures[]")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// // 读取所有数据
	// data, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// // 获取当前路径
	// dir, patherr := filepath.Abs(filepath.Dir(os.Args[0]))
	// if patherr != nil {
	// 	log.Fatal(patherr)
	// }

	// // 文件路径
	// dir = strings.Replace(dir, "bin", "", 1)
	// filePath := dir + "/" + "src/SanWenJia/public/img/" + handler.Filename

	// // 保存到文件
	// err = ioutil.WriteFile(filePath, data, 0777)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	return this.Redirect((*Admin).UploadPicture)
}
