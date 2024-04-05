package ads

import (
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateAd(ad *Ad) error
	FirstOrCreateCountry(code string) (Country, error)
	FirstOrCreatePlatform(name string) (Platform, error)
	ListActiveAds(query Query) ([]Response, error)
}

type AdRepository struct {
	db *gorm.DB
}

func NewAdRepository(db *gorm.DB) *AdRepository {
	return &AdRepository{db: db}
}

func (repo *AdRepository) CreateAd(ad *Ad) error {
	return repo.db.Create(ad).Error
}

func (repo *AdRepository) FirstOrCreateCountry(code string) (Country, error) {
	var country Country
	err := repo.db.FirstOrCreate(&country, Country{Code: code}).Error
	return country, err
}

func (repo *AdRepository) FirstOrCreatePlatform(name string) (Platform, error) {
	var platform Platform
	err := repo.db.FirstOrCreate(&platform, Platform{Name: name}).Error
	return platform, err
}

func (repo *AdRepository) ListActiveAds(query Query) ([]Response, error) {
	var ads []Response

	now := time.Now()
	rawQuery := repo.db.Model(&Ad{}).Select("title, end_at").Where("start_at < ?", now).Where("end_at > ?", now)

	if query.Age != -1 {
		rawQuery = rawQuery.Where("? BETWEEN age_start AND age_end", query.Age)
	}

	if query.Gender != "" {
		rawQuery = rawQuery.Where("gender_target = ? OR gender_target IS NULL", query.Gender)
	}

	if query.Country != "" {
		rawQuery = rawQuery.Joins("JOIN ad_countries ON ad_countries.ad_id = ads.id").
			Where("ad_countries.country_code = ? OR ad_countries.country_code = 'any'", query.Country)
	}

	if query.Platform != "" {
		rawQuery = rawQuery.Joins("JOIN ad_platforms ON ad_platforms.ad_id = ads.id").
			Where("ad_platforms.platform_name = ? OR ad_platforms.platform_name = 'any'", query.Platform)
	}

	err := rawQuery.Offset(query.Offset).Limit(query.Limit).Order("end_at asc").Find(&ads).Error
	return ads, err
}
