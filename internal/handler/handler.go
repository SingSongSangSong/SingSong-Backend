package handler

import "http-practice/internal/model"

type Handler struct {
	model *model.Model
}

func NewHandler(model *model.Model) (*Handler, error) {
	handler := &Handler{model: model}
	return handler, nil
}
