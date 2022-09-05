package sdk

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func checkPath(path string) string {
	lenght := len(path)
	if lenght < 1 {
		return ""
	}
	if path[lenght-1] != '/' {
		path = fmt.Sprintf("%s/", path)
	}
	return path
}

func makeAwsService(accessKey, secretKey, regionID string) *s3.S3 {
	awsCnf := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""), // 秘钥
		Region:           aws.String(regionID),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false),
	}
	se := session.Must(session.NewSession(awsCnf))
	return s3.New(se)
}

// upload file  to aws s3
func AwsS3UploadFile(accessKey, secretKey, bucket, regionID, path, file string, body []byte) error {
	if len(file) < 1 || len(body) < 1 {
		return nil
	}

	service := makeAwsService(accessKey, secretKey, regionID)
	ctx, chancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer chancel()

	path = checkPath(path)
	objInput := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("%s%s", path, file)),
		Body:   bytes.NewReader(body),
	}

	_, err := service.PutObjectWithContext(ctx, objInput)
	return err
}

func AwsS3GetFile(accessKey, secretKey, bucket, regionID, path, file string) ([]byte, error) {
	if len(file) < 1 {
		return []byte{}, nil
	}

	service := makeAwsService(accessKey, secretKey, regionID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()

	path = checkPath(path)
	objInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("%s%s", path, file)),
	}

	out, err := service.GetObjectWithContext(ctx, objInput)
	if err != nil {
		return []byte{}, err
	}
	defer out.Body.Close()

	return ioutil.ReadAll(out.Body)
}

func AwsS3DeleteFile(accessKey, secretKey, bucket, regionID, path, file string) error {
	if len(file) < 1 {
		return nil
	}

	service := makeAwsService(accessKey, secretKey, regionID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	path = checkPath(path)
	delInput := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("%s%s", path, file)),
	}
	_, err := service.DeleteObjectWithContext(ctx, delInput)
	return err
}

func AwsS3DeleteFileWaitObjectNotExist(accessKey, secretKey, bucket, regionID, path, file string) error {
	if len(file) < 1 {
		return nil
	}

	service := makeAwsService(accessKey, secretKey, regionID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	path = checkPath(path)
	delInput := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("%s%s", path, file)),
	}
	_, err := service.DeleteObjectWithContext(ctx, delInput)
	if err != nil {
		return err
	}

	existInput := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("%s%s", path, file)),
	}
	err = service.WaitUntilObjectNotExistsWithContext(ctx, existInput)
	return err
}

func AwsS3MultiLandDeleteFile(accessKey, secretKey, bucket, regionID, path string, files []string) error {
	if len(files) < 1 {
		return nil
	}

	service := makeAwsService(accessKey, secretKey, regionID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	path = checkPath(path)
	n, length := 1000, len(files)
	left := length / n
	for i := 0; i <= left; i++ {
		min, max := i*n, (i+1)*n
		if max > length {
			max = length
		}

		delObjects := []*s3.ObjectIdentifier{}
		for _, file := range files[min:max] {
			delObjects = append(delObjects, &s3.ObjectIdentifier{
				Key: aws.String(fmt.Sprintf("%s%s", path, file)),
			})
		}
		delFilesInput := &s3.DeleteObjectsInput{
			Bucket: aws.String(bucket),
			Delete: &s3.Delete{Objects: delObjects, Quiet: aws.Bool(true)},
		}
		_, err := service.DeleteObjectsWithContext(ctx, delFilesInput)
		if err != nil {
			return err
		}
	}

	return nil
}

func AwsS3FileList(accessKey, secretKey, bucket, regionID, path, keyword string) ([]string, error) {
	var objkeys []string
	service := makeAwsService(accessKey, secretKey, regionID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	path = checkPath(path)
	listInput := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(path),
	}

	listCallBack := func(output *s3.ListObjectsOutput, b bool) bool {
		for _, content := range output.Contents {
			pathFile := aws.StringValue(content.Key)
			file := pathFile[len(path):]
			if file == "" || (keyword != "" && !strings.Contains(file, keyword)) {
				continue
			}
			objkeys = append(objkeys, file)
		}
		return true
	}

	err := service.ListObjectsPagesWithContext(ctx, listInput, listCallBack)
	return objkeys, err
}

func AwsS3CopyFile(accessKey, secretKey, bucket, regionID, path, source, target string) error {
	service := makeAwsService(accessKey, secretKey, regionID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	path = checkPath(path)
	copyInput := &s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(url.PathEscape(fmt.Sprintf("%s/%s%s", bucket, path, source))),
		Key:        aws.String(fmt.Sprintf("%s%s", path, target)),
	}
	_, err := service.CopyObjectWithContext(ctx, copyInput)
	return err
}

func AwsS3CopyFileWaitObjectExists(accessKey, secretKey, bucket, regionID, path, source, target string) error {
	service := makeAwsService(accessKey, secretKey, regionID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	path = checkPath(path)
	copyInput := &s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(url.PathEscape(fmt.Sprintf("%s/%s%s", bucket, path, source))),
		Key:        aws.String(fmt.Sprintf("%s%s", path, target)),
	}
	_, err := service.CopyObject(copyInput)
	if err != nil {
		return err
	}

	// 注意一下检查objectExists 会阻塞 (5s*20次 = 1min)
	existInput := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("%s%s", path, target)),
	}
	err = service.WaitUntilObjectExistsWithContext(ctx, existInput)

	return err
}

func AwsS3FileExist(accessKey, secretKey, bucket, regionID, path, file string) bool {
	files, _ := AwsS3FileList(accessKey, secretKey, bucket, regionID, path, file)
	for _, fileName := range files {
		if file == fileName {
			return true
		}
	}
	return false
}
