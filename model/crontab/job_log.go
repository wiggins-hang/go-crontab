package crontab

import (
	"go-crontab/model/db_model"

	"gorm.io/gorm"
)

func CreateJob(db *gorm.DB, log *db_model.JobLog) error {
	return db.Create(log).Error
}

func CreateBatchJob(db *gorm.DB, log []*db_model.JobLog) error {
	return db.Create(log).Error
}

type LogPageFilter struct {
	Skip  int
	Limit int
	Name  string
}

func JobLogPageList(db *gorm.DB, pageFilter LogPageFilter) ([]*db_model.JobLog, error) {
	logList := make([]*db_model.JobLog, 0)

	err := db.Where("job_name = ?", pageFilter.Name).Limit(pageFilter.Limit).Offset(pageFilter.Skip).Find(&logList).Error
	return logList, err
}
