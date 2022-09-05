package sdk

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

// 七牛 SDK_go       相关文档 https://developer.qiniu.com/kodo/1238/go#1

type QiniuPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

func QiNiuToken(accessKey, secretKey, bucket string) string {
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	return putPolicy.UploadToken(mac)
}

func QiNiuOverWriteToken(accessKey, secretKey, bucket, overWriteFile string) string {
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", bucket, overWriteFile),
	}
	mac := qbox.NewMac(accessKey, secretKey)
	return putPolicy.UploadToken(mac)
}

func QiNiuUploadFile(accessKey, secretKey, bucket string, file string, data []byte) (url string, err error) {
	buff := bytes.NewBuffer(data)
	token := QiNiuToken(accessKey, secretKey, bucket)
	cfg := storage.Config{Zone: &storage.ZoneHuanan}
	formUploader := storage.NewFormUploader(&cfg)
	putExtra := &storage.PutExtra{
		Params:   map[string]string{},
		MimeType: "application/octet-stream",
	}
	ret := QiniuPutRet{}
	err = formUploader.Put(context.Background(), &ret, token, file, buff, int64(len(data)), putExtra)
	if err == nil {
		url = ret.Key
	}
	return
}

func QiNiuReplaceUploadFile(accessKey, secretKey, bucket string, file string, data []byte) (url string, err error) {
	buff := bytes.NewBuffer(data)
	overWriteToken := QiNiuOverWriteToken(accessKey, secretKey, bucket, file)
	cfg := storage.Config{Zone: &storage.ZoneHuanan}
	formUploader := storage.NewFormUploader(&cfg)
	putExtra := &storage.PutExtra{
		Params:   map[string]string{},
		MimeType: "application/octet-stream",
	}
	ret := QiniuPutRet{}
	err = formUploader.Put(context.Background(), &ret, overWriteToken, file, buff, int64(len(data)), putExtra)
	if err == nil {
		url = ret.Key
	}
	return
}

func QiNiuGetFile(qiniuUrl string, file string) (bs []byte, err error) {
	qiniuUrl = checkPath(qiniuUrl)
	resp, err := http.Get(fmt.Sprintf("%s%s", qiniuUrl, file))
	if err != nil {
		return nil, fmt.Errorf("file[%s]下载失败,  err: %v", file, err)
	}
	defer resp.Body.Close()
	bs, err = ioutil.ReadAll(resp.Body)
	return bs, err
}
