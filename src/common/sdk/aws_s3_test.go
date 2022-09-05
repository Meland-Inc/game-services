package sdk

import (
	"fmt"
	"testing"
)

var (
	// please set your AwsS3 config
	awsAccessKey = ""
	awsSecretKey = ""
	awsBucket    = ""
	awsRegionID  = ""
	awsFilePath  = ""
)

func Test_AwsS3(t *testing.T) {
	t.Run("create file", func(t *testing.T) {
		str := "test create file 011"
		body := []byte(str)
		err := AwsS3UploadFile(
			awsAccessKey, awsSecretKey, awsBucket, awsRegionID, awsFilePath,
			"test01", body)
		t.Log(err)
	})

	t.Run("get file ", func(t *testing.T) {
		body, err := AwsS3GetFile(awsAccessKey, awsSecretKey, awsBucket, awsRegionID, awsFilePath, "test01")
		t.Log(err)
		t.Log(string(body))
	})

	t.Run("update file", func(t *testing.T) {
		str := "test create file 01  update test"
		body := []byte(str)
		err := AwsS3UploadFile(awsAccessKey, awsSecretKey, awsBucket, awsRegionID, awsFilePath, "test01", body)
		t.Log(err)
	})

	t.Run("get updated file ", func(t *testing.T) {
		body, err := AwsS3GetFile(awsAccessKey, awsSecretKey, awsBucket, awsRegionID, awsFilePath, "test01")
		t.Log(err)
		t.Log(string(body))
	})

	t.Run("File exist ", func(t *testing.T) {
		exist := AwsS3FileExist(awsAccessKey, awsSecretKey, awsBucket, awsRegionID, awsFilePath, "test01")
		t.Log(exist)
	})

	t.Run("delete file", func(t *testing.T) {
		err := AwsS3DeleteFile(awsAccessKey, awsSecretKey, awsBucket, awsRegionID, awsFilePath, "test01")
		t.Log(err)
	})
}

func Test_AwsS3Copy(t *testing.T) {
	t.Run("create file", func(t *testing.T) {
		str := "test create copy test  file 001"
		body := []byte(str)
		err := AwsS3UploadFile(awsAccessKey, awsSecretKey, awsBucket, awsRegionID, awsFilePath, "ctest01", body)
		t.Log(err)
	})

	t.Run("COPY file ", func(t *testing.T) {
		err := AwsS3CopyFile(awsAccessKey, awsSecretKey, awsBucket, awsRegionID, awsFilePath, "ctest01", "ctest01_copy")
		t.Log(err)
	})
}

func Test_AwsS3MultiLand(t *testing.T) {
	t.Run("MultiLand create file", func(t *testing.T) {
		for i := 0; i < 5000; i++ {
			str := fmt.Sprintf("test create file %d", i)
			body := []byte(str)
			AwsS3UploadFile(awsAccessKey, awsSecretKey, awsBucket, awsRegionID, awsFilePath, fmt.Sprintf("test_%d", i), body)
		}
	})

	t.Run("get file list ", func(t *testing.T) {
		files, err := AwsS3FileList(awsAccessKey, awsSecretKey, awsBucket, awsRegionID, awsFilePath, "")
		t.Log(err)
		t.Log(len(files))
	})

	t.Run("MultiLand Delete File  ", func(t *testing.T) {
		files := []string{}
		for i := 0; i < 5000; i++ {
			files = append(files, fmt.Sprintf("test_%d", i))
		}
		err := AwsS3MultiLandDeleteFile(awsAccessKey, awsSecretKey, awsBucket, awsRegionID, awsFilePath, files)
		t.Log(err)
	})
}
