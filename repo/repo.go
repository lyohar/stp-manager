package repo

import (
	"stpCommon/model"
)

type Repository interface {
	GetExport() (*model.Export, error)
	SetExportStatus(status *model.ExportStatus) error
}

