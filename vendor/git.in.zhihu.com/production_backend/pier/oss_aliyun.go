package pier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	lru "github.com/hashicorp/golang-lru"

	pico3Thrift "git.in.zhihu.com/thrift-go/pico3_thrift"
)

const (
	maxCacheSize   = (1 << 20) * 8
	OSSReadChannel = "http://zhihu-pics.oss-cn-beijing.aliyuncs.com"
)

var (
	lruCache, _ = lru.New(maxCacheSize)

	ErrToken = errors.New("The image token does no exist.")
)

func UploadImageBySTS(content []byte, imageName string, uploadToken *pico3Thrift.UploadTokenObject, bucket *pico3Thrift.OSSBucketObject, mime string) (pico3Thrift.ImageUploadState, error) {
	uploadState := pico3Thrift.ImageUploadState_UPLOAD_FAIL

	client, err := oss.New(bucket.Endpoint, uploadToken.AccessKey, uploadToken.AccessSecret, oss.SecurityToken(uploadToken.AccessToken))
	if err != nil {
		return uploadState, err
	}

	bucketClient, err := client.Bucket(bucket.BucketName)
	if err != nil {
		return uploadState, err
	}

	if mime == "" {
		err = bucketClient.PutObject(imageName, bytes.NewReader(content))
	} else {
		err = bucketClient.PutObject(imageName, bytes.NewReader(content), oss.ContentType(mime))
	}
	if err == nil {
		uploadState = pico3Thrift.ImageUploadState_SUCCESS
	} else {
		PierStatsd.Increment(fmt.Sprintf("pier-go.%s.upload_oss.fail.count", SafePierVersion))
	}
	return uploadState, nil
}

type OSSReadClient struct {
	readClient *http.Client
}

func NewOSSReadClient(timeout time.Duration) *OSSReadClient {
	return &OSSReadClient{
		readClient: &http.Client{
			Timeout: timeout,
		},
	}
}

type ImageHue struct {
	RGB string `json:"RGB"`
}

func (o *OSSReadClient) GetImageHue(token string) (*ImageHue, error) {
	var (
		metaKey = "image/average-hue"
		hueKey  = token + metaKey
	)

	imageHueCache, ok := lruCache.Get(hueKey)
	if ok {
		return imageHueCache.(*ImageHue), nil
	}

	imageHue := &ImageHue{}
	infoUrl := fmt.Sprintf("%s/%s?x-oss-process=%s", OSSReadChannel, token, metaKey)
	resp, err := o.readClient.Get(infoUrl)

	if err != nil {
		return imageHue, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return imageHue, ErrToken
	}

	err = json.NewDecoder(resp.Body).Decode(imageHue)

	lruCache.Add(hueKey, imageHue)
	return imageHue, err
}

type ExifFieldValue struct {
	Value string `json:"Value"`
}

type ImageExif struct {
	Compression                 ExifFieldValue
	DateTime                    ExifFieldValue
	ExifTag                     ExifFieldValue
	FileSize                    ExifFieldValue
	Format                      ExifFieldValue
	GPSLatitude                 ExifFieldValue
	GPSLatitudeRef              ExifFieldValue
	GPSLongitude                ExifFieldValue
	GPSLongitudeRef             ExifFieldValue
	GPSMapDatum                 ExifFieldValue
	GPSTag                      ExifFieldValue
	GPSVersionID                ExifFieldValue
	ImageHeight                 ExifFieldValue
	ImageWidth                  ExifFieldValue
	JPEGInterchangeFormat       ExifFieldValue
	JPEGInterchangeFormatLength ExifFieldValue
	Orientation                 ExifFieldValue
	ResolutionUnit              ExifFieldValue
	Software                    ExifFieldValue
	XResolution                 ExifFieldValue
	YResolution                 ExifFieldValue
}

func (o *OSSReadClient) GetImageExif(token string) (*ImageExif, error) {
	var (
		metaKey = "image/info"
		exifKey = token + metaKey
	)

	imageExifCache, ok := lruCache.Get(exifKey)
	if ok {
		return imageExifCache.(*ImageExif), nil
	}

	imageExif := &ImageExif{}
	infoUrl := fmt.Sprintf("%s/%s?x-oss-process=%s", OSSReadChannel, token, metaKey)
	resp, err := o.readClient.Get(infoUrl)

	if err != nil {
		return imageExif, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return imageExif, ErrToken
	}

	err = json.NewDecoder(resp.Body).Decode(imageExif)

	lruCache.Add(exifKey, imageExif)
	return imageExif, err
}
