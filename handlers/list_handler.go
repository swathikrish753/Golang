package handlers

import (
	"GolangWorld/models"
	"GolangWorld/services"
	"context"
)

type ListHandler struct {
	ListService services.ListService
}

func NewListHandler(listService services.ListService) *ListHandler {
	return &ListHandler{ListService: listService}
}

func (h *ListHandler) ListAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := h.ListService.ListAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
