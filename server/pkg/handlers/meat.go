package handlers

import (
	"context"

	"butcher/pkg/models"

	"github.com/alpha-omega-corp/_core/app"
)

type MeatHandler struct {
	repository app.Repository[models.Meat]
}

func (h *MeatHandler) New(deps app.Deps) {
	h.repository = app.NewRepository[models.Meat](deps.DB)
}

func (h *MeatHandler) GetAll(ctx context.Context) (*[]models.Meat, error) {
	return h.repository.GetAll(ctx, nil)
}

func (h *MeatHandler) GetOne(ctx context.Context, id int64) (*models.Meat, error) {
	return h.repository.GetOne(ctx, id)
}

func (h *MeatHandler) Create(ctx context.Context, data *models.Meat) (*models.Meat, error) {
	return h.repository.Create(ctx, data)
}

func (h *MeatHandler) Update(ctx context.Context, data *models.Meat) (*models.Meat, error) {
	return h.repository.Update(ctx, data)
}

func (h *MeatHandler) Delete(ctx context.Context, id int64) error {
	return h.repository.Delete(ctx, id)
}

func (h *MeatHandler) SetImage(ctx context.Context, id int64, imagePath string) (*models.Meat, error) {
	meat, err := h.repository.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	meat.ImagePath = imagePath
	return h.repository.Update(ctx, meat)
}
