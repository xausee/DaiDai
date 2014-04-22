package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
	"os"
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
	err := models.BackupPicture()
	if err != nil {
		fmt.Println("备份图片失败")
	}

	models.SavePicture(pictures)
	if err != nil {
		fmt.Println("保存新图片失败")
	}

	return this.Redirect((*Admin).UploadPicture)
}
