package repository

import (
	"hos_schedule/internal/model"

	"gorm.io/gorm"
)

type PatientRepo struct {
	db *gorm.DB
}

func NewPatientRepo(db *gorm.DB) *PatientRepo {
	return &PatientRepo{db: db}
}

func (r *PatientRepo) ListByUser(userID int64) ([]model.Patient, error) {
	var patients []model.Patient
	err := r.db.Where("user_id = ?", userID).Order("is_default DESC").Find(&patients).Error
	return patients, err
}

func (r *PatientRepo) GetByID(id int64) (*model.Patient, error) {
	var patient model.Patient
	err := r.db.First(&patient, id).Error
	return &patient, err
}

func (r *PatientRepo) Create(patient *model.Patient) error {
	return r.db.Create(patient).Error
}

func (r *PatientRepo) Update(patient *model.Patient) error {
	return r.db.Save(patient).Error
}

func (r *PatientRepo) Delete(id int64) error {
	return r.db.Delete(&model.Patient{}, id).Error
}

func (r *PatientRepo) SetDefault(userID, patientID int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Patient{}).Where("user_id = ?", userID).
			Update("is_default", false).Error; err != nil {
			return err
		}
		return tx.Model(&model.Patient{}).Where("id = ? AND user_id = ?", patientID, userID).
			Update("is_default", true).Error
	})
}
