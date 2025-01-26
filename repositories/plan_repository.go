package repositories

import (
	"context"

	"github.com/malwarebo/gopay/db"
	"github.com/malwarebo/gopay/models"
)

type PlanRepository struct {
	db *db.DB
}

func NewPlanRepository(db *db.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) Create(ctx context.Context, plan *models.Plan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

func (r *PlanRepository) Update(ctx context.Context, plan *models.Plan) error {
	return r.db.WithContext(ctx).Save(plan).Error
}

func (r *PlanRepository) GetByID(ctx context.Context, id string) (*models.Plan, error) {
	var plan models.Plan
	if err := r.db.WithContext(ctx).First(&plan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepository) List(ctx context.Context) ([]*models.Plan, error) {
	var plans []*models.Plan
	if err := r.db.WithContext(ctx).Where("active = ?", true).Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *PlanRepository) Delete(ctx context.Context, id string) error {
	// Soft delete by setting active = false
	return r.db.WithContext(ctx).Model(&models.Plan{}).Where("id = ?", id).Update("active", false).Error
}
