package pier

import (
	"context"
	"time"

	client "git.in.zhihu.com/go/base/tzone"
	pico3Thrift "git.in.zhihu.com/thrift-go/pico3_thrift"
)

const (
	serviceName = "Pico3Service"
	targetName  = "pico3-image-service"
)

var appLocalIP = GetLocalIP()

type Pico3RPCImpl struct {
	pico3Client *pico3Thrift.Pico3ServiceClient
}

func NewPico3Client(timeout time.Duration) *Pico3RPCImpl {
	newClient := client.NewClient(
		serviceName,
		client.TargetName(targetName),
		client.Timeout(timeout),
	)
	pico3Client := pico3Thrift.NewPico3ServiceClient(newClient)
	return &Pico3RPCImpl{
		pico3Client: pico3Client,
	}
}

func (p *Pico3RPCImpl) UploadImage(imageHash, source string, memberID int64) (*pico3Thrift.ImageUploadResult_, error) {
	ctx := context.Background()
	uploadImageParam := pico3Thrift.ImageDataParam{
		ImageHash: imageHash,
		Source:    source,
	}
	if memberID > 0 {
		uploadImageParam.MemberID = &memberID
	}

	uploadResult, err := p.pico3Client.GetUploadToken(ctx, &uploadImageParam)
	return uploadResult, err
}

func (p *Pico3RPCImpl) UpdateImageUploadStateWithMeta(imageId string, imageUploadState pico3Thrift.ImageUploadState, memberID int64, isNeedMeta bool) (*pico3Thrift.UpdateImageStateResult_, error) {
	ctx := context.Background()
	imageStateParam := pico3Thrift.UpdateImageStateParam{
		ImageID:          imageId,
		ImageUploadState: imageUploadState,
		IsNeedMeta:       &isNeedMeta,
	}
	if memberID > 0 {
		imageStateParam.MemberID = &memberID
	}

	updateResult, err := p.pico3Client.UpdateImageUploadStateWithMeta(ctx, &imageStateParam)
	return updateResult, err
}

func (p *Pico3RPCImpl) GetImageInfo(token string) (*ImageInfo, error) {
	ctx := context.Background()
	param := pico3Thrift.GetImageInfoParam{
		Token: token,
	}

	imageInfo := &ImageInfo{}
	result, err := p.pico3Client.GetImageInfo(ctx, &param)
	if err != nil {
		return imageInfo, err
	}

	if result.Format != nil {
		imageInfo.Format = *result.Format
	}
	if result.Size != nil {
		imageInfo.Size = *result.Size
	}
	if result.Height != nil {
		imageInfo.Height = int64(*result.Height)
	}
	if result.Width != nil {
		imageInfo.Width = int64(*result.Width)
	}
	return imageInfo, err
}

func (p *Pico3RPCImpl) DeleteImage(token string) error {
	ctx := context.Background()
	err := p.pico3Client.ImageDelete(ctx, token)
	return err
}

func (p *Pico3RPCImpl) ListImageSpecs() (*pico3Thrift.ImageSpecsData, error) {
	ctx := context.Background()
	imageSpecsData, err := p.pico3Client.ListImageSpecs(ctx)
	return imageSpecsData, err
}

func (p *Pico3RPCImpl) fetchScheduleStrategy(currentStrategy *PierStrategy, userIP string) (*pico3Thrift.StrategyObject, error) {
	ctx := context.Background()

	domains := []string{}
	for domain, _ := range currentStrategy.CDNDomains {
		domains = append(domains, domain)
	}

	param := pico3Thrift.ScheduleStrategyParam{
		VersionID:   currentStrategy.getVersionID(),
		AppName:     ZaeAppName,
		PierVersion: SafePierVersion,
		AppLocalIP:  &appLocalIP,
		Domains:     domains,
		UserIP:      &userIP,
	}
	strategy, err := p.pico3Client.FetchScheduleStrategy(ctx, &param)
	return strategy, err
}
