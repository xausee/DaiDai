package models

import (
	// gio "io"
	"log"

	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/io"
	// rio "github.com/qiniu/api/resumable/io"
	"github.com/qiniu/api/auth/digest"
	"github.com/qiniu/api/rs"
)

//将本地文件上传到七牛云存储，返回key。
func UploadToQiniu(filepath string) (string, error) {

	conf.ACCESS_KEY = QiNiuAccessKey
	conf.SECRET_KEY = QiNiuSecretKey

	//获取uptoken
	putPolicy := rs.PutPolicy{Scope: QiNiuBucket}
	uptoken := putPolicy.Token(nil)

	//上传
	var ret io.PutRet
	err := io.PutFileWithoutKey(nil, &ret, uptoken, filepath, nil)

	if err != nil {
		//上传产生错误
		log.Print("io.PutFileWithoutKey failed:", err)
	}

	return ret.Key, err
}

func DeleteFileOnQiNiu(key string) error {
	mac := digest.Mac{AccessKey: QiNiuAccessKey, SecretKey: []byte(QiNiuSecretKey)}
	reCli := rs.New(&mac)
	err := reCli.Delete(nil, QiNiuBucket, key)
	if err != nil {
		// 产生错误
		log.Println("rs.Delete failed:", err)
	}
	return err
}
