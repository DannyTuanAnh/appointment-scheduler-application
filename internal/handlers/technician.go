package handlers

import (
	"net/http"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/dto"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/services"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/validation"
	"github.com/gin-gonic/gin"
)

type TechnicianHandler struct {
	service services.TechnicianService
}

func NewTechnicianHandler(service services.TechnicianService) *TechnicianHandler {
	return &TechnicianHandler{service: service}
}

func (h *TechnicianHandler) CreateTechnician(ctx *gin.Context) {
	var req dto.CreateTechnicianRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	err := h.service.CreateTechnician(ctx.Request.Context(), req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusCreated)
}

func (h *TechnicianHandler) GetTechnicianByID(ctx *gin.Context) {
	var req dto.RequestTechnicianWithID
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	tech, err := h.service.GetTechnicianByID(ctx.Request.Context(), req.ID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, tech)
}

func (h *TechnicianHandler) GetDetailTechnicianByID(ctx *gin.Context) {
	var req dto.RequestTechnicianWithID
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	tech, err := h.service.GetDetailTechnicianByID(ctx.Request.Context(), req.ID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, tech)
}

func (h *TechnicianHandler) UpdateTechnician(ctx *gin.Context) {
	var reqID dto.RequestTechnicianWithID
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var req dto.UpdateTechnicianRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if req.Name == nil && req.Level == nil {
		utils.ResponseStatusCode(ctx, http.StatusNoContent)
		return
	}

	err := h.service.UpdateTechnicianInfoByID(ctx.Request.Context(), reqID.ID, req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusOK)
}

func (h *TechnicianHandler) SetTechnicianOnLeave(ctx *gin.Context) {
	var req dto.RequestTechnicianWithID
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	_, err := h.service.SetTechnicianOnLeave(ctx.Request.Context(), req.ID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusOK)
}

func (h *TechnicianHandler) SetTechnicianBackToWork(ctx *gin.Context) {
	var req dto.RequestTechnicianWithID
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	_, err := h.service.SetTechnicianBackToWork(ctx.Request.Context(), req.ID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusOK)
}

func (h *TechnicianHandler) TransferTechnicianDealership(ctx *gin.Context) {
	var reqID dto.RequestTechnicianWithID
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var req dto.TransferTechnicianDealershipRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	_, err := h.service.TransferTechnicianDealership(ctx.Request.Context(), reqID.ID, req.DealershipID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusOK)
}

func (h *TechnicianHandler) DeleteTechnician(ctx *gin.Context) {
	var req dto.RequestTechnicianWithID
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if err := h.service.DeleteTechnicianByID(ctx.Request.Context(), req.ID); err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}

func (h *TechnicianHandler) ListTechnicians(ctx *gin.Context) {
	var q dto.ListTechniciansQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	rows, err := h.service.ListTechniciansByDealershipID(ctx.Request.Context(), q.DealershipID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, rows)
}

func (h *TechnicianHandler) SearchTechnicians(ctx *gin.Context) {
	var q dto.SearchTechniciansQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if q.DealershipID != nil {
		rows, err := h.service.SearchTechniciansByNameAndDealershipID(ctx.Request.Context(), *q.DealershipID, q.Name)
		if err != nil {
			utils.ResponseError(ctx, err)
			return
		}
		utils.ResponseSuccess(ctx, http.StatusOK, rows)
		return
	}

	rows, err := h.service.SearchTechniciansByName(ctx.Request.Context(), q.Name)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, rows)
}

func (h *TechnicianHandler) FindActiveTechnicians(ctx *gin.Context) {
	var q dto.FindActiveTechniciansQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	ids, err := h.service.FindActiveTechniciansByDealershipWithRequiredSkills(ctx.Request.Context(), q.DealershipID, q.SkillIDs)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, ids)
}

func (h *TechnicianHandler) AddSkillsToTechnician(ctx *gin.Context) {
	var reqID dto.RequestTechnicianWithID
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var req dto.AddSkillsToTechnicianRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if err := h.service.AddSkillsToTechnician(ctx.Request.Context(), reqID.ID, req.SkillIds); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusCreated)
}

func (h *TechnicianHandler) RemoveSkillsFromTechnician(ctx *gin.Context) {
	var reqID dto.RequestTechnicianWithID
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var req dto.RemoveSkillsFromTechnicianRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	_, err := h.service.RemoveSkillsFromTechnician(ctx.Request.Context(), reqID.ID, req.SkillIds)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusOK)
}
