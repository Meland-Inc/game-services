package sdk

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func makeOssService(ossUrl, accessKey, secretKey, bucket string) (*oss.Bucket, error) {
	client, err := oss.New(ossUrl, accessKey, secretKey)
	if err != nil {
		return nil, err
	}

	return client.Bucket(bucket)
}

func ALiOssUpLoadFile(ossUrl, accessKey, secretKey, bucket string, file string, body []byte) error {
	service, err := makeOssService(ossUrl, accessKey, secretKey, bucket)
	if err != nil || service == nil {
		return err
	}

	// 上传文件流。
	return service.PutObject(file, bytes.NewReader(body))
}

func ALiOssGetFile(ossUrl, accessKey, secretKey, bucket string, file string) ([]byte, error) {
	service, err := makeOssService(ossUrl, accessKey, secretKey, bucket)
	if err != nil || service == nil {
		return nil, err
	}

	// 上传文件流。
	fd, err := service.GetObject(file)
	if err != nil {
		strErr := err.Error()
		if -1 != strings.Index(strErr, "StatusCode=404") {
			return nil, err
		}
		return nil, err
	}

	bs, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func ALiOSSCopyFile(ossUrl, accessKey, secretKey, bucket string, source, target string) (*oss.CopyObjectResult, error) {
	service, err := makeOssService(ossUrl, accessKey, secretKey, bucket)
	if err != nil || service == nil {
		return nil, err
	}

	cor, err := service.CopyObject(source, target)
	if err != nil {
		return nil, err
	}
	return &cor, err
}

func ALiOSSFileExist(ossUrl, accessKey, secretKey, bucket string, file string) (bool, error) {
	service, err := makeOssService(ossUrl, accessKey, secretKey, bucket)
	if err != nil || service == nil {
		return false, err
	}

	return service.IsObjectExist(file)
}

func ALiOSSDeleteFile(ossUrl, accessKey, secretKey, bucket string, file string) error {
	service, err := makeOssService(ossUrl, accessKey, secretKey, bucket)
	if err != nil || service == nil {
		return err
	}
	if file == "" {
		return nil
	}

	return service.DeleteObject(file)
}

func ALiOSSMultiLandDeleteFile(ossUrl, accessKey, secretKey, bucket string, files []string) error {
	service, err := makeOssService(ossUrl, accessKey, secretKey, bucket)
	if err != nil || service == nil {
		return err
	}

	_, err = service.DeleteObjects(files, oss.DeleteObjectsQuiet(true))
	return err
}

func ALiOSSFileList(ossUrl, accessKey, secretKey, bucket string, keyWorld string) ([]string, error) {
	service, err := makeOssService(ossUrl, accessKey, secretKey, bucket)
	if err != nil || service == nil {
		return nil, err
	}

	marker := oss.Marker("")
	prefix := oss.Prefix(keyWorld)

	var files []string
	for {
		lsRes, err := service.ListObjects(marker, prefix, oss.Delimiter("/"))
		if err != nil {
			return nil, err
		}

		for _, object := range lsRes.Objects {
			files = append(files, object.Key)
		}

		prefix = oss.Prefix(lsRes.Prefix)
		marker = oss.Marker(lsRes.NextMarker)
		if !lsRes.IsTruncated {
			break
		}
	}

	return files, err
}
