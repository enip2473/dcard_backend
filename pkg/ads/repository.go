package ads

import (
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateAd(ad *Ad) error
	ListActiveAds(offset, limit int, conditions map[string]interface{}) ([]Ad, int64, error)
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

func (repo *AdRepository) ListActiveAds(offset, limit int, conditions map[string]interface{}) ([]Ad, int64, error) {
	var ads []Ad
	var total int64

	query := repo.db.Model(&Ad{}).Where("end_at >= ?", time.Now())

	// Add your conditions logic here
	// For example:
	if age, ok := conditions["age"].(int); ok && age > 0 {
		query = query.Where("conditions @> ?", []map[string]interface{}{
			{
				"ageStart": age - 3,
				"ageEnd":   age + 3,
			},
		})
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Find(&ads).Error
	return ads, total, err
}
