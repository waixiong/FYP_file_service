package model

import (

	//pb "./proto"

	// "github.com/syndtr/goleveldb/leveldb"

	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strings"

	"getitqec.com/server/file/pkg/commons"
	"getitqec.com/server/file/pkg/dao"
	"getitqec.com/server/file/pkg/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/chai2010/webp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// type FileModel struct {
// 	FileDAO dao.IFileDAO
// }

// // InitModel ...
// func InitModel(m commons.MongoDB) FileModelI {
// 	// dao := &dao.UserDAO{}
// 	_fileDao := dao.InitFileDAO(m)
// 	return &FileModel{FileDAO: _fileDao}
// }

// func (m *FileModel) UploadImage(ctx context.Context, name string, data []byte) error {
// 	return m.FileDAO.Create(ctx, name, data)
// }

// func (m *FileModel) DownloadImage(ctx context.Context, name string) ([]byte, error) {
// 	return m.FileDAO.Get(ctx, name)
// }

// need interface

type FileModel struct {
	// signdb        *leveldb.DB
	// authdb        *leveldb.DB
	// sessiondb     *leveldb.DB
	// usersessiondb *leveldb.DB
	// db      *dynamodb.DynamoDB
	FileDAO  dao.IFileDAO
	S3Client *s3.S3
}

//  userprofile
//  storeprofile
//  getitproducts

// InitModel ...
func InitModel(m commons.MongoDB, s *s3.S3) FileModelI {
	// dao := &dao.UserDAO{}
	_fileDao := dao.InitFileDAO(m)
	model := &FileModel{FileDAO: _fileDao, S3Client: s}
	return model
}

func (m *FileModel) UploadFile(ctx context.Context, name string, data []byte) error {
	return m.FileDAO.Create(ctx, name, data)
}

func (m *FileModel) DownloadFile(ctx context.Context, name string) ([]byte, error) {
	return m.FileDAO.Get(ctx, name)
}

// need interface

func (m *FileModel) UploadImage(ctx context.Context, bucketName string, name string, data []byte) (string, error) {
	key := aws.String(name)
	bucket := aws.String("testBuc")
	contentType := aws.String(http.DetectContentType(data))
	// contentType := aws.String("image/webp")
	logger.Log.Debug(fmt.Sprintf("bytes: %d", len(data)))
	logger.Log.Debug(fmt.Sprintf("b: %s", *bucket))
	logger.Log.Debug(fmt.Sprintf("n: %s", *contentType))
	// image/jpeg
	// image/heic -> application/octet-stream
	// image/png
	// image/webp
	var body *bytes.Reader
	var buf bytes.Buffer

	// fmt.Printf("%s Size: %d \n", bucketName, len(data))
	if *contentType == "image/webp" {
		/*
			webp - maintain any config
			jpeg
		*/
		i, err := webp.Decode(bytes.NewReader(data))
		// resize to jpeg
		img := i

		// img, _ = jpeg.Decode(bytes.NewReader(jbuf.Bytes()))
		if err = webp.Encode(&buf, img, &webp.Options{Lossless: true, Quality: 75, Exact: true}); err != nil {
			log.Println(err)
		}
		body = bytes.NewReader(buf.Bytes())
	} else if *contentType == "image/jpeg" {
		/*
			webp - lossy
			jpeg
		*/
		i, err := jpeg.Decode(bytes.NewReader(data))
		// resize
		img := i

		// Encode lossy webp
		if err = webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: 75, Exact: true}); err != nil {
			log.Println(err)
		}
		body = bytes.NewReader(buf.Bytes())
	} else if *contentType == "image/png" {
		/*
			webp - lossless
			jpeg
		*/
		// var buf bytes.Buffer
		i, err := png.Decode(bytes.NewReader(data))
		// resize to jpeg
		img := i

		// img, _ = jpeg.Decode(bytes.NewReader(jbuf.Bytes()))
		// Encode lossless webp
		fmt.Println(len(buf.Bytes()))
		if err = webp.Encode(&buf, img, &webp.Options{Lossless: true, Quality: 75, Exact: true}); err != nil {
			log.Println(err)
		}
		// fmt.Println(len(buf.Bytes()))
		body = bytes.NewReader(buf.Bytes())
	} else {
		return "", status.Errorf(codes.InvalidArgument, "Unsupported image type")
	}
	// logger.Log.Debug(fmt.Sprintf("bytesBody: %v", body))
	*contentType = "image/webp"
	output, err := m.S3Client.PutObject(&s3.PutObjectInput{
		Body:        body, //bytes.NewReader(data),
		Bucket:      bucket,
		Key:         key,
		ContentType: contentType,
	})
	if err != nil {
		fmt.Printf("Failed to upload object %s/%s, %s\n", *bucket, *key, err.Error())
		return "", err
	}
	// jcontentType := aws.String("image/jpeg")
	// key = aws.String(name + ".jpeg")
	// _, err = m.S3Client.PutObject(&s3.PutObjectInput{
	// 	Body:        jbody, //bytes.NewReader(data),
	// 	Bucket:      bucket,
	// 	Key:         key,
	// 	ContentType: jcontentType,
	// })

	// if err != nil {
	// 	fmt.Printf("Failed to upload object %s/%s, %s\n", *bucket, *key, err.Error())
	// 	return "", err
	// }
	logger.Log.Debug(fmt.Sprintf("Successfully uploaded key %s with id %s\n", *key, *output.VersionId))
	return *output.VersionId, nil
	// return "", nil
}

func (m *FileModel) DownloadImage() {}

//// https://stackoverflow.com/questions/25959386/how-to-check-if-a-file-is-a-valid-image
// image formats and magic numbers
var magicTable = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

// mimeFromIncipit returns the mime type of an image file from its first few
// bytes or the empty string if the file does not look like a known file type
func mimeFromIncipit(incipit []byte) string {
	incipitStr := string(incipit)
	for magic, mime := range magicTable {
		if strings.HasPrefix(incipitStr, magic) {
			return mime
		}
	}

	return ""
}

func transfromBackgroundImage(background color.RGBA, src image.Image) *image.RGBA {
	rect := image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy())
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), &image.Uniform{background}, image.ZP, draw.Src)
	draw.Draw(img, img.Bounds(), src, image.ZP, draw.Over)
	return img
}
