package sdk

import (
	"encoding/json"
	"fmt"
	"testing"
)

type QiNiuTestData struct {
	Id              int64  `db:"id"`
	AuthorId        int64  `db:"author_id"`
	Name            string `db:"name"`
	ProgramPb       []byte `db:"program_pb"`
	ThumbnailBase64 string `db:"thumbnail_base64"`
	TemplateType    int32  `db:"template_type"`
	IsNew           bool   `db:"is_new"`
}

var (
	// please set your QiNiu config
	qiniuUrl       = ""
	qiniuAccessKey = ""
	qiniuSecretKey = ""
	qiniuBucket    = ""
)

func Test_QiNiu(t *testing.T) {

	t.Run("create file", func(t *testing.T) {
		info := QiNiuTestData{
			Id:              333333,
			AuthorId:        888888,
			Name:            "test_qiniu_01",
			ProgramPb:       nil,
			ThumbnailBase64: "test qiniu 01 create file",
			TemplateType:    2,
			IsNew:           true,
		}
		bs, err := json.Marshal(info)
		if err != nil {
			t.Error(err)
			return
		}

		fileUrl, err1 := QiNiuUploadFile(qiniuAccessKey, qiniuSecretKey, qiniuBucket, "test_qiniu_01", bs)
		t.Log(err1)
		t.Log(fileUrl)
	})

	t.Run("get file", func(t *testing.T) {
		bs, err := QiNiuGetFile(qiniuUrl, "test_qiniu_01")
		t.Log(err)
		t.Log(bs)

		info := &QiNiuTestData{}
		err = json.Unmarshal(bs, info)
		t.Log(err)
		t.Log(fmt.Sprintf("%+v", info))
	})
}
