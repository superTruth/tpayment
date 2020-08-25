package fileupload

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
	"tpayment/conf"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

// 申请一个上传文件的URL
func RequestUploadFileUrl(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(UploadFileRequest)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 检查参数是否正确
	if req.FileSize == 0 || req.Tag == "" || req.FileName == "" {
		logger.Warn("parameters miss")
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	filePath := fmt.Sprintf("%v/%v/%v", req.Tag, strings.ReplaceAll(uuid.New().String(), "-", ""), req.FileName)

	// S3
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(conf.GetConfigData().S3Region),
		Credentials: credentials.NewStaticCredentials(conf.GetConfigData().S3Key, conf.GetConfigData().S3Secret, ""),
	}))

	service := s3.New(sess)

	//endpoints.ApNortheast2RegionID
	resp, _ := service.PutObjectRequest(&s3.PutObjectInput{
		ACL:        aws.String("public-read"),
		Bucket:     aws.String(conf.GetConfigData().S3Bucket),
		Key:        aws.String(filePath),
		ContentMD5: &req.Md5,
	})

	//exp := 15 * int64(1+req.FileSize/10/1024/1024)
	//url, err := resp.Presign(time.Minute * time.Duration(exp)) // 15分钟 10M的发送时间
	url, err := resp.Presign(time.Hour)

	if err != nil {
		logger.Error("pre sign fail->", err.Error())
		modules.BaseError(ctx, conf.UnknownError)
		return err
	}

	logger.Info("pre sign url->", url)

	ret := new(UploadFileResponse)
	ret.UploadUrl = base64.StdEncoding.EncodeToString([]byte(url))
	ret.DownloadUrl = "https://" + conf.GetConfigData().S3Bucket + "." + conf.GetConfigData().S3Region + "." + conf.GetConfigData().S3Domain + "/" + filePath
	//ret.Exp = exp

	modules.BaseSuccess(ctx, ret)

	return nil
}

//func test() {
//	fp, _ := os.Open("/Users/truth/project/tpayment/modules/fileupload/requestfileupload.go")
//	defer fp.Close()
//
//	fmt.Println("key->", conf.GetConfigData().S3Key, ", secret->", conf.GetConfigData().S3Secret)
//	// S3
//	sess := session.Must(session.NewSession(&aws.Config{
//		Region:      aws.String(conf.GetConfigData().S3Region),
//		Credentials: credentials.NewStaticCredentials(conf.GetConfigData().S3Key, conf.GetConfigData().S3Secret, ""),
//	}))
//
//	service := s3.New(sess)
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
//	defer cancel()
//
//	_, err := service.PutObjectWithContext(ctx, &s3.PutObjectInput{
//		ACL:    aws.String("public-read"),
//		Bucket: aws.String(conf.GetConfigData().S3Bucket),
//		Key:    aws.String("test"),
//		Body:   fp,
//	})
//
//	if err != nil {
//		fmt.Println("upload fail->", err.Error())
//	} else {
//		fmt.Println("upload success")
//	}
//}
