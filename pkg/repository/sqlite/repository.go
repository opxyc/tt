package sqlite

import (
	"github.com/opxyc/tt/pkg/tt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sqliteRepository struct {
	client *gorm.DB
}

func (s *sqliteRepository) Find(activityID uint) (*tt.Activity, error) {
	activity := &tt.Activity{}
	err := s.client.Where("id = ?", activityID).First(activity).Error
	return activity, err
}

func (s *sqliteRepository) Store(activity *tt.Activity) error {
	return s.client.Create(activity).Error
}

func (s *sqliteRepository) List(filters *tt.ListFilters) ([]tt.Activity, error) {
	tx := s.client
	if len(filters.Status) != 0 {
		tx = tx.Where("status IN (?)", filters.Status)
	}

	if filters.Date != "" {
		tx = tx.Where("strftime('%Y-%m-%d', created_at) = ?", filters.Date)
	}

	var activities []tt.Activity
	err := tx.Find(&activities).Error
	if err != nil {
		return nil, err
	}

	return activities, nil
}

func (s *sqliteRepository) Update(activityID uint, activity *tt.Activity) error {
	return s.client.Model(&tt.Activity{}).Where("id = ?", activityID).Updates(activity).Error
}

func (s *sqliteRepository) Delete(activityID uint) error {
	return s.client.Delete(&tt.Activity{}, activityID).Error
}

func NewSqliteRespository(sqliteDSN string) (tt.Repository, error) {
	db, err := gorm.Open(sqlite.Open(sqliteDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&tt.Activity{})
	if err != nil {
		return nil, err
	}

	return &sqliteRepository{client: db}, nil
}
