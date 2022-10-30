package v1

import "gorm.io/gorm"

type DownloopContext struct {
	Database *gorm.DB
}
