package usecase

import (
	"back/lk/internal/entity"
	"back/lk/internal/repo"
	"context"
	"fmt"
)

const (
	ImageTypeRestaurant string = "restaurant"
	ImageTypeFood       string = "food"
)

type RestInterface interface {
	GetInfo(ctx context.Context, id uint64) (entity.Rest, error)
	UploadBaseInfo(ctx context.Context, info entity.BaseInfo, id uint64) (entity.BaseInfo, error)
	UploadLogo(ctx context.Context, file []byte, extention string, mimeType string, restId uint64) error
	UploadDescriptionAndImages(ctx context.Context, info *entity.DescripAndImgs, restId uint64) error
}

type Rest struct {
	repo  repo.RestInterface
	minio repo.RepoMinio
}

func NewRest(r repo.RestInterface, m repo.RepoMinio) RestInterface {
	return &Rest{repo: r, minio: m}
}

func (u *Rest) GetInfo(ctx context.Context, id uint64) (entity.Rest, error) {
	res, err := u.repo.GetBaseInfo(ctx, id)
	if err != nil {
		return entity.Rest{}, err
	}
	schedule, err := u.repo.GetSchedule(ctx, id)
	if err != nil {
		return entity.Rest{}, err
	}
	res.Schedule = schedule
	return res, nil
}

func (u *Rest) UploadBaseInfo(ctx context.Context, info entity.BaseInfo, id uint64) (entity.BaseInfo, error) {
	return u.repo.UploadBaseInfo(ctx, info, id)
}

func (u *Rest) UploadLogo(ctx context.Context, file []byte, extention string, mimeType string, restId uint64) error {
	path := fmt.Sprintf("%s/%d/logo_url%s", ImageTypeRestaurant, restId, extention)
	_, err := u.minio.UploadImage(ctx, file, path, mimeType)
	if err != nil {
		return err
	}
	err = u.repo.PutLogoImage(ctx, path, restId)
	if err != nil {
		return err
	}
	return nil
}

func (u *Rest) UploadDescriptionAndImages(ctx context.Context, info *entity.DescripAndImgs, restId uint64) error {
	if info.Description != nil {
		for i, str := range info.Description {
			err := u.repo.UpdateDescrip(ctx, str, info.DescripIndexes[i], restId)
			if err != nil {
				return err
			}
		}
	}
	for _, img := range info.Img {
		path := ""
		fmt.Println(len(img.Data))
		if len(img.Data) != 0 {
			path = fmt.Sprintf("%s/%d/%d%s", ImageTypeRestaurant, restId, img.Index, img.Ext)
			_, err := u.minio.UploadImage(ctx, img.Data, path, img.MimeType)
			if err != nil {
				return err
			}
		}
		err := u.repo.PutImgUrl(ctx, path, img.Index, restId)
		if err != nil {
			return err
		}
	}
	return nil
}
