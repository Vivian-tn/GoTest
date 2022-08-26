package pier

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"

	pico3Thrift "git.in.zhihu.com/thrift-go/pico3_thrift"
)

var (
	maskAppName = MaskAppName(ZaeAppName)
	once        sync.Once
	picoAPI     *PicoAPI
)

const (
	urlSecureProtocol         = "https://"
	strategySchedulerInterval = 10

	channelOfOSS = "https://zhihu-pics-dev.oss-cn-beijing.aliyuncs.com"
	cdnDevDomain = "origin1-zhimg.zhihu.dev"

	maxZHImageLength = 4*4096 - 1       // 知乎图片边长限制 4 * 4096
	maxZHImageSize   = 20*1024*1024 - 1 // 知乎图片大小限制 20M
	maxZHGifSize     = 10*1024*1024 - 1 // 知乎 GIF 大小限制 10M
)

type UploadImageInfo struct {
	Format    string `json:"format"`
	FullUrl   string `json:"full_url"`
	Height    string `json:"height"`
	QRcode    string `json:"qrcode,omitempty"`
	Sequences int64  `json:"sequences,omitempty"`
	Size      string `json:"size"`
	Token     string `json:"token"`
	Width     string `json:"width"`
}

type ImageInfo struct {
	Format    string `json:"format"`
	Height    int64  `json:"height"`
	QRcode    string `json:"qrcode,omitempty"`
	Sequences string `json:"sequences,omitempty"`
	Size      int64  `json:"size"`
	Width     int64  `json:"width"`
}

type PicoAPI struct {
	ossReadClient *OSSReadClient
	pico3RPCImpl  *Pico3RPCImpl

	strategy        *PierStrategy
	domainScheduler DomainScheduler
}

// NewPicoAPI news a PicoAPI for interacting with pico service.
func NewPicoAPI() *PicoAPI {
	once.Do(func() {
		picoAPI = newPicoAPIWithTimeout(3 * time.Second)
		picoAPI.startScheduler(strategySchedulerInterval)
	})
	return picoAPI
}

func newPicoAPIWithTimeout(timeout time.Duration) *PicoAPI {
	ossReadClient := NewOSSReadClient(timeout)
	pico3RPCImpl := NewPico3Client(timeout)

	strategy := NewPierStrategy()
	domainScheduler := getScheduler("weight_random")
	return &PicoAPI{
		ossReadClient:   ossReadClient,
		pico3RPCImpl:    pico3RPCImpl,
		strategy:        strategy,
		domainScheduler: domainScheduler,
	}
}

func (p *PicoAPI) startScheduler(schedulerInterval int) {
	schedulerTicker := time.NewTicker(time.Duration(schedulerInterval) * time.Second)

	go func() {
		for {
			select {
			case <-schedulerTicker.C:
				{
					remoteStrategy, err := p.pico3RPCImpl.fetchScheduleStrategy(p.strategy, "")
					if err != nil {
						continue
					}

					if remoteStrategy.VersionID > 0 && (p.strategy.VersionID != remoteStrategy.VersionID || len(p.strategy.CDNDomains) != len(remoteStrategy.CdnDomains)) {
						strategyLruCache.Add(strategyVersionIDKey, remoteStrategy.VersionID)

						cdnDomains := map[string]CDNDomainParameter{}
						for _, domain := range remoteStrategy.CdnDomains {
							cdnDomains[(*domain).Domain] = CDNDomainParameter{
								Weight: (*domain).Weight,
							}
						}
						p.strategy = &PierStrategy{
							VersionID:      remoteStrategy.VersionID,
							DefaultSpec:    remoteStrategy.DefaultSpec,
							DefaultQuality: remoteStrategy.DefaultQuality,
							SpecMapping:    remoteStrategy.SpecMapping,
							QualityMapping: remoteStrategy.QualityMapping,
							CDNDomains:     cdnDomains,
						}
						strategyLruCache.Add(remoteStrategy.VersionID, p.strategy)

						strVersionId := strconv.FormatInt(remoteStrategy.VersionID, 10)
						PierStatsd.Increment(fmt.Sprintf("pier-go.%s.%s.%s.count", SafePierVersion, SafeZaeAppName, strVersionId))
						PierStatsd.Increment(fmt.Sprintf("pier-go.%s.%s.%s.%s.count", SafePierVersion, SafeZaeAppName, SafeColumn(GetLocalIP()), strVersionId))
					}
				}
			}
		}
	}()
	return
}

func checkSizeValid(format string, size int) (bool, error) {
	isExceedSize := false
	if format != "gif" && size > maxZHImageSize {
		isExceedSize = true
	} else if format == "gif" && size > maxZHGifSize {
		isExceedSize = true
	}
	if isExceedSize {
		PierStatsd.Increment(fmt.Sprintf("pier-go.%s.upload.exceed_image.count", SafePierVersion))
		return false, errors.New("This image is too large, limit: normal size < 20MB, gif size < 10MB")
	}

	return true, nil
}

func checkLengthValid(width int64, height int64) (bool, error) {
	if width > maxZHImageLength || height > maxZHImageLength {
		PierStatsd.Increment(fmt.Sprintf("pier-go.%s.upload.exceed_image.count", SafePierVersion))
		return false, errors.New("This Image is too large, limit: length < 16384")
	}
	return true, nil
}

func (p *PicoAPI) UploadImage(content io.Reader) (*UploadImageInfo, error) {

	PierStatsd.Increment(fmt.Sprintf("pier-go.%s.upload.count", SafePierVersion))
	uploadImageInfo := &UploadImageInfo{}

	imageContent, err := ioutil.ReadAll(content)
	if err != nil {
		return uploadImageInfo, errors.New("Something wrong with the image content")
	}

	imageSize := len(imageContent)
	if imageSize == 0 {
		PierStatsd.Increment(fmt.Sprintf("pier-go.%s.upload.image_blank.count", SafePierVersion))
		return uploadImageInfo, errors.New("The image content can not be blank")
	}

	mime, imageFormat := sniffMIMEAndFormat(imageContent)

	_, sizeErr := checkSizeValid(imageFormat, imageSize)
	if sizeErr != nil {
		return uploadImageInfo, sizeErr
	}

	imageHash := md5SumImage(imageContent)
	imageToken := formatToken(imageHash)

	uploadImageResult, err := p.pico3RPCImpl.UploadImage(imageHash, ZaeAppName, -1)
	if err != nil {
		return uploadImageInfo, err
	}

	if uploadImageResult.UploadFile.State == int16(pico3Thrift.ImageUploadType_UPLOADED_SUCCESS) {
		PierStatsd.Increment(fmt.Sprintf("pier-go.%s.upload.second_upload.count", SafePierVersion))
		// 添加图片 meta 信息
		uploadImageInfo.Format = *uploadImageResult.ImageMeta.Format
		uploadImageInfo.Height = strconv.FormatInt(int64(*uploadImageResult.ImageMeta.Height), 10)
		uploadImageInfo.Width = strconv.FormatInt(int64(*uploadImageResult.ImageMeta.Width), 10)
		uploadImageInfo.Size = strconv.FormatInt(*uploadImageResult.ImageMeta.Size, 10)
	}

	if uploadImageResult.UploadFile.State == int16(pico3Thrift.ImageUploadType_NEW_UPLOAD) {
		PierStatsd.Increment(fmt.Sprintf("pier-go.%s.upload.new_upload.count", SafePierVersion))
		// 通过 OSS 上传图片
		uploadState, err := UploadImageBySTS(
			imageContent,
			imageToken,
			uploadImageResult.UploadToken,
			uploadImageResult.Bucket,
			mime,
		)
		if err != nil {
			return uploadImageInfo, errors.New("Something wrong when upload image to OSS")
		}

		// 更新上传图片的状态
		updateResult, err := p.pico3RPCImpl.UpdateImageUploadStateWithMeta(
			uploadImageResult.UploadFile.ImageID,
			uploadState,
			-1,
			true,
		)
		if err != nil {
			return uploadImageInfo, err
		}

		if uploadState == pico3Thrift.ImageUploadState_SUCCESS {
			if updateResult.ImageMeta.Format != nil {
				intWidth := int64(*updateResult.ImageMeta.Width)
				intHeight := int64(*updateResult.ImageMeta.Height)

				if intWidth <= 0 || intHeight <= 0 {
					PierStatsd.Increment(fmt.Sprintf("pier-go.%s.upload.meta_error.count", SafePierVersion))
					return uploadImageInfo, errors.New("Image meta info invalid")
				}
				_, lengthErr := checkLengthValid(intWidth, intHeight)
				if lengthErr != nil {
					return nil, lengthErr
				}
				// 添加图片 meta 信息
				uploadImageInfo.Format = *updateResult.ImageMeta.Format
				uploadImageInfo.Height = strconv.FormatInt(intHeight, 10)
				uploadImageInfo.Width = strconv.FormatInt(intWidth, 10)
				uploadImageInfo.Size = strconv.FormatInt(*updateResult.ImageMeta.Size, 10)
			}
			PierStatsd.Increment(fmt.Sprintf("pier-go.%s.upload.success.count", SafePierVersion))
		} else {
			PierStatsd.Increment(fmt.Sprintf("pier-go.%s.upload.fail.count", SafePierVersion))
			return uploadImageInfo, errors.New("Upload image to OSS error")
		}
	}

	// 添加图片的 token 及 根据图片的格式生成 URL
	uploadImageInfo.Token = imageToken
	if uploadImageInfo.Format != "" {
		imageFormat = uploadImageInfo.Format
	}
	uploadImageInfo.FullUrl = p.GetFullURL(imageToken, "", imageFormat, "", true, false)
	return uploadImageInfo, nil
}

func (p *PicoAPI) GetImageInfo(token string) (*ImageInfo, error) {
	token, _ = parseToken(token)
	imageInfo, err := p.pico3RPCImpl.GetImageInfo(token)
	PierStatsd.Increment(fmt.Sprintf("pier-go.%s.get_image_info.count", SafePierVersion))
	return imageInfo, err
}

func (p *PicoAPI) DeleteImage(token string) error {
	token, _ = parseToken(token)
	return p.pico3RPCImpl.DeleteImage(token)
}

func (p *PicoAPI) GetScheduler() DomainScheduler {
	return p.domainScheduler
}

func (p *PicoAPI) SetSchedulerMethod(method string) {
	p.domainScheduler = getScheduler(method)
}

func (p *PicoAPI) SetCDNDomains(domains []string) {
	p.domainScheduler.SetDomains(domains)
}

func (p *PicoAPI) GetFullURLDefault(token, size string) string {
	return p.GetFullURL(token, size, "", "", true, false)
}

func (p *PicoAPI) GetFullURL(token, size, fmt, quality string, secure, keepFmt bool) string {
	if token == "" {
		return ""
	}

	if isURL(token) {
		return token
	}

	token, fmtVal := parseToken(token)
	if fmt == "" {
		if fmtVal != "" && keepFmt {
			fmt = fmtVal
		} else {
			fmt = DefaultFmt
		}
	}

	fmt = strings.ToLower(fmt)
	if _, ok := FormatSet[fmt]; !ok {
		fmt = DefaultFmt
	}

	spec := getValidSpec(p.strategy.SpecMapping, size)
	if spec == "" {
		token = token + "." + fmt
	} else {
		token = token + "_" + spec + "." + fmt
	}

	// 域名调度
	domain := cdnDevDomain
	if IsProductionEnv() {
		domain = p.domainScheduler.GetDomain(p.strategy.getCDNDomains(), token)
	}

	quality = getValidQuality(p.strategy.QualityMapping, quality)
	imageUrl := urlSecureProtocol + strings.Join([]string{domain, quality, token}, "/")
	if quality == "" {
		imageUrl = urlSecureProtocol + strings.Join([]string{domain, token}, "/")
	}
	return imageUrl + "?source=" + maskAppName
}

func (p *PicoAPI) GetImageHue(token string) (*ImageHue, error) {
	token, _ = parseToken(token)
	imageHue, err := p.ossReadClient.GetImageHue(token)
	return imageHue, err
}

func (p *PicoAPI) GetImageExif(token string) (*ImageExif, error) {
	token, _ = parseToken(token)
	imageExif, err := p.ossReadClient.GetImageExif(token)
	PierStatsd.Increment(fmt.Sprintf("pier-go.%s.get_image_exif.count", SafePierVersion))
	return imageExif, err
}

func (p *PicoAPI) ListImageSpecs() (*pico3Thrift.ImageSpecsData, error) {
	imageSpecs, err := p.pico3RPCImpl.ListImageSpecs()
	return imageSpecs, err
}

func (p *PicoAPI) URL2Token(url string) (token string, suffix string, err error) {
	return URL2Token(url)
}
