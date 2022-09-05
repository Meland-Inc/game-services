package sdk

import (
	"fmt"
	"testing"
)

var (
	// please set your AwsS3 config
	ossUrl       = ""
	ossAccessKey = ""
	ossSecretKey = ""
	ossBucket    = ""
)

func Test_AliOSS(t *testing.T) {
	t.Run("create file", func(t *testing.T) {
		str := "test create file oss 011"
		body := []byte(str)
		err := ALiOssUpLoadFile(
			ossUrl, ossAccessKey, ossSecretKey, ossBucket,
			"tjtest/test01", body,
		)
		t.Log(err)
	})

	t.Run("get file ", func(t *testing.T) {
		body, err := ALiOssGetFile(
			ossUrl, ossAccessKey, ossSecretKey, ossBucket,
			"tjtest/test01",
		)
		t.Log(err)
		t.Log(string(body))
	})

	t.Run("update file", func(t *testing.T) {
		str := "test create file oss01 update test"
		body := []byte(str)
		err := ALiOssUpLoadFile(
			ossUrl, ossAccessKey, ossSecretKey, ossBucket,
			"tjtest/test01", body,
		)
		t.Log(err)
	})

	t.Run("get update file ", func(t *testing.T) {
		body, err := ALiOssGetFile(
			ossUrl, ossAccessKey, ossSecretKey, ossBucket,
			"tjtest/test01",
		)
		t.Log(err)
		t.Log(string(body))
	})

	t.Run("File exist ", func(t *testing.T) {
		exist, err := ALiOSSFileExist(
			ossUrl, ossAccessKey, ossSecretKey, ossBucket,
			"tjtest/test01",
		)
		t.Log(exist)
		t.Log(err)
	})

	t.Run("delete file", func(t *testing.T) {
		err := ALiOSSDeleteFile(
			ossUrl, ossAccessKey, ossSecretKey, ossBucket,
			"tjtest/test01",
		)
		t.Log(err)
	})
}

func Test_AliOssCopy(t *testing.T) {
	t.Run("create so file", func(t *testing.T) {
		str := "test create copy test  file  oss---001"
		body := []byte(str)
		err := ALiOssUpLoadFile(
			ossUrl, ossAccessKey, ossSecretKey, ossBucket,
			"tjtest/test01", body,
		)
		t.Log(err)
	})

	t.Run("get file ", func(t *testing.T) {
		body, err := ALiOssGetFile(
			ossUrl, ossAccessKey, ossSecretKey, ossBucket,
			"tjtest/test01",
		)
		t.Log(err)
		t.Log(string(body))
	})

	t.Run("COPY file ", func(t *testing.T) {
		c, err := ALiOSSCopyFile(
			ossUrl, ossAccessKey, ossSecretKey, ossBucket,
			"tjtest/test01", "tjtest/test01_copy",
		)
		t.Log(c)
		t.Log(err)
	})

	t.Run("get file ", func(t *testing.T) {
		body, err := ALiOssGetFile(
			ossUrl, ossAccessKey, ossSecretKey, ossBucket,
			"tjtest/test01_copy",
		)
		t.Log(err)
		t.Log(string(body))
	})
}

func Test_ALiOssMultiLand(t *testing.T) {
	t.Run("MultiLand create file", func(t *testing.T) {
		for i := 0; i < 50; i++ {
			str := fmt.Sprintf("test create file %d", i)
			body := []byte(str)
			ALiOssUpLoadFile(
				ossUrl, ossAccessKey, ossSecretKey, ossBucket,
				fmt.Sprintf("tjtest/test%d", i),
				body,
			)
		}
		for i := 0; i < 50; i++ {
			str := fmt.Sprintf("test create file %d", i)
			body := []byte(str)
			ALiOssUpLoadFile(
				ossUrl, ossAccessKey, ossSecretKey, ossBucket,
				fmt.Sprintf("tjtest/c993狼test%d", i),
				body,
			)
		}
	})

	t.Run("get All file list ", func(t *testing.T) {
		files, err := ALiOSSFileList(
			ossUrl, ossAccessKey, ossSecretKey, ossBucket,
			"tjtest/c993",
		)
		t.Log(err)
		t.Log(len(files))
		t.Log(files)
	})

	t.Run("MultiLand Delete File  ", func(t *testing.T) {
		files := []string{}
		for i := 0; i < 50; i++ {
			files = append(files, fmt.Sprintf("tjtest/test%d", i))
		}
		for i := 0; i < 50; i++ {
			files = append(files, fmt.Sprintf("tjtest/c993狼test%d", i))
		}
		err := ALiOSSMultiLandDeleteFile(
			ossUrl, ossAccessKey, ossSecretKey, ossBucket,
			files,
		)
		t.Log(err)
	})
}
