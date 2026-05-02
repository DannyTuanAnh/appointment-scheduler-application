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
