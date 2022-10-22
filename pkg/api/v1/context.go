package v1

import "github.com/jmoiron/sqlx" 

type DownloopContext struct {
	Database *sqlx.DB 
}
