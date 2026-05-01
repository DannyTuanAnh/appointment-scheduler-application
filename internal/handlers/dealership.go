package handlers

import (
	"net/http"
	"time"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/dto"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/services"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/validation"
	"github.com/gin-gonic/gin"
)

type DealershipHandler struct {
	service services.DealershipService
}

func NewDealershipHandler(service services.DealershipService) *DealershipHandler {
	return &DealershipHandler{service: service}
}

func (d *DealershipHandler) GetAllDealerships(ctx *gin.Context) {
	dealerships, err := d.service.GetAllDealerships(ctx.Request.Context())
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, dealerships)
}

func (d *DealershipHandler) GetDealershipByID(ctx *gin.Context) {
	var req dto.RequestDealershipWithID
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	dealership, err := d.service.GetDealershipByID(ctx.Request.Context(), req.ID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, dealership)
}

func (d *DealershipHandler) CreateDealership(ctx *gin.Context) {
	var req dto.CreateDealershipRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	open_time, err := time.Parse("15:04", req.OpenTime)
	if err != nil {
		utils.ResponseValidator(ctx, gin.H{"error": "Invalid open_time format. Expected HH:mm"})
		return
	}

	close_time, err := time.Parse("15:04", req.CloseTime)
	if err != nil {
		utils.ResponseValidator(ctx, gin.H{"error": "Invalid close_time format. Expected HH:mm"})
		return
	}

	err = d.service.CreateDealership(ctx.Request.Context(), req.Name, open_time, close_time)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusCreated)
}

func (d *DealershipHandler) UpdateDealership(ctx *gin.Context) {
	var reqID dto.RequestDealershipWithID
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var req dto.UpdateDealershipRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if req.Name == nil && req.OpenTime == nil && req.CloseTime == nil {
		utils.ResponseStatusCode(ctx, http.StatusNoContent)
		return
	}

	var openTime, closeTime *time.Time

	if req.OpenTime != nil {
		oTime, err := time.Parse("15:04", *req.OpenTime)
		if err != nil {
			utils.ResponseValidator(ctx, gin.H{"error": "Invalid open_time format. Expected HH:mm"})
			return
		}

		openTime = &oTime
	}

	if req.CloseTime != nil {
		cTime, err := time.Parse("15:04", *req.CloseTime)
		if err != nil {
			utils.ResponseValidator(ctx, gin.H{"error": "Invalid close_time format. Expected HH:mm"})
			return
		}

		closeTime = &cTime
	}

	err := d.service.UpdateDealership(ctx.Request.Context(), reqID.ID, req.Name, openTime, closeTime)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusOK)
}

func (d *DealershipHandler) DeleteDealership(ctx *gin.Context) {
	var req dto.RequestDealershipWithID
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	err := d.service.DeleteDealershipByID(ctx.Request.Context(), req.ID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}

func (d *DealershipHandler) SearchDealerships(ctx *gin.Context) {
	var q dto.SearchDealershipsQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	dealerships, err := d.service.SearchDealershipsByName(ctx.Request.Context(), q.Name)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, dealerships)
}
