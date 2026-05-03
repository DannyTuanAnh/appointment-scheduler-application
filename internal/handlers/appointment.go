package handlers

import (
	"net/http"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/dto"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/services"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/validation"
	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	service services.AppointmentService
}

func NewAppointmentHandler(service services.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{service: service}
}

func (h *AppointmentHandler) GetAppointment(ctx *gin.Context) {
	var req dto.GetAppointmentRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	slot, err := h.service.GetAppointment(ctx.Request.Context(), req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, slot)
}

// APPOINTMENT MANAGEMENT

func (h *AppointmentHandler) ListAppointments(ctx *gin.Context) {
	rows, err := h.service.ListAppointments(ctx.Request.Context())
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, rows)
}

func (h *AppointmentHandler) ListAppointmentsByDealershipInTimeRange(ctx *gin.Context) {
	var q dto.ListAppointmentsByDealershipQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	rows, err := h.service.ListAppointmentsByDealershipInTimeRange(ctx.Request.Context(), q)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, rows)
}

func (h *AppointmentHandler) ListAppointmentsByTechnicianInTimeRange(ctx *gin.Context) {
	var q dto.ListAppointmentsByTechnicianQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	rows, err := h.service.ListAppointmentsByTechnicianInTimeRange(ctx.Request.Context(), q)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, rows)
}

func (h *AppointmentHandler) ListAppointmentsByServiceBayInTimeRange(ctx *gin.Context) {
	var q dto.ListAppointmentsByServiceBayQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	rows, err := h.service.ListAppointmentsByServiceBayInTimeRange(ctx.Request.Context(), q)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, rows)
}

func (h *AppointmentHandler) ListAppointmentsByServiceInTimeRange(ctx *gin.Context) {
	var q dto.ListAppointmentsByServiceQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	rows, err := h.service.ListAppointmentsByServiceInTimeRange(ctx.Request.Context(), q)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, rows)
}

func (h *AppointmentHandler) SearchAppointments(ctx *gin.Context) {
	var q dto.SearchAppointmentsQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	rows, err := h.service.SearchAppointments(ctx.Request.Context(), q)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, rows)
}

func (h *AppointmentHandler) CreateAppointment(ctx *gin.Context) {
	var req dto.CreateAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	if err := h.service.CreateAppointment(ctx.Request.Context(), req); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusCreated)
}

func (h *AppointmentHandler) CancelledAppointment(ctx *gin.Context) {
	var uri dto.RequestAppointmentWithID
	if err := ctx.ShouldBindUri(&uri); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var req dto.CancelledAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if err := h.service.UpdateAppointmentStatusByID(ctx.Request.Context(), uri.ID, req.Status); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}

func (h *AppointmentHandler) UpdateAppointmentStatus(ctx *gin.Context) {
	var uri dto.RequestAppointmentWithID
	if err := ctx.ShouldBindUri(&uri); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var req dto.UpdateAppointmentStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if err := h.service.UpdateAppointmentStatusByID(ctx.Request.Context(), uri.ID, req.Status); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}

func (h *AppointmentHandler) MarkNoShowAppointments(ctx *gin.Context) {
	var req dto.MarkNoShowAppointmentsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	if err := h.service.MarkNoShowAppointmentsForDealershipInTimeRange(ctx.Request.Context(), req.AppointmentIds); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}
