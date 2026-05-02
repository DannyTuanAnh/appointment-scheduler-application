package handlers

import (
	"net/http"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/dto"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/services"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/validation"
	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	service services.ServiceService
}

func NewServiceHandler(service services.ServiceService) *ServiceHandler {
	return &ServiceHandler{service: service}
}

// SERVICES

func (h *ServiceHandler) ListServices(ctx *gin.Context) {
	rows, err := h.service.ListServices(ctx.Request.Context())
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, rows)
}

func (h *ServiceHandler) GetServiceDetailByID(ctx *gin.Context) {
	var req dto.RequestServiceWithID
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	row, err := h.service.GetServiceDetailByID(ctx.Request.Context(), req.ID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, row)
}

func (h *ServiceHandler) SearchServices(ctx *gin.Context) {
	var q dto.SearchServicesQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	rows, err := h.service.SearchServicesByName(ctx.Request.Context(), q.Name)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, rows)
}

func (h *ServiceHandler) CreateService(ctx *gin.Context) {
	var req dto.CreateServiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if err := h.service.CreateService(ctx.Request.Context(), req); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusCreated)
}

func (h *ServiceHandler) UpdateService(ctx *gin.Context) {
	var reqID dto.RequestServiceWithID
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var req dto.UpdateServiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	// if empty payload -> no content
	if req.RequiredBayTypeID == nil && req.Name == nil && req.AnticipatedMinutes == nil {
		utils.ResponseStatusCode(ctx, http.StatusNoContent)
		return
	}

	if err := h.service.UpdateServiceByID(ctx.Request.Context(), reqID.ID, req); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusOK)
}

func (h *ServiceHandler) DeleteService(ctx *gin.Context) {
	var req dto.RequestServiceWithID
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if err := h.service.DeleteServiceByID(ctx.Request.Context(), req.ID); err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}

// SERVICE REQUIREMENTS

func (h *ServiceHandler) AddSkillRequirementsToService(ctx *gin.Context) {
	var reqID dto.RequestServiceWithID
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var req dto.AddServiceRequirementsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if err := h.service.AddSkillRequirementsToService(ctx.Request.Context(), reqID.ID, req.SkillIds); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusCreated)
}

func (h *ServiceHandler) RemoveSkillRequirementsFromService(ctx *gin.Context) {
	var reqID dto.RequestServiceWithID
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var req dto.RemoveServiceRequirementsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	_, err := h.service.RemoveSkillRequirementsFromService(ctx.Request.Context(), reqID.ID, req.SkillIds)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusOK)
}
