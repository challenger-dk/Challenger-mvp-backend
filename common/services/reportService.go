package services

import (
	"server/common/config"
	"server/common/dto"
	"server/common/models"

	"gorm.io/gorm"
)

func CreateReport(reporterID uint, req dto.ReportCreateDto) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		report := models.Report{
			ReporterID: reporterID,
			TargetID:   req.TargetID,
			TargetType: models.ReportTargetType(req.TargetType),
			Reason:     req.Reason,
			Comment:    req.Comment,
			Status:     "PENDING",
		}

		// Use 'tx' instead of 'config.DB' to ensure this runs inside the transaction
		if err := tx.Create(&report).Error; err != nil {
			return err
		}

		return nil
	})
}
