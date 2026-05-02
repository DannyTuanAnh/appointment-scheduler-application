package handlers

import (
	"net/http"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/dto"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/services"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/validation"
	"github.com/gin-gonic/gin"
)

type ServiceBayHandler struct {
	service services.ServiceBayService
}

func NewServiceBayHandler(service services.ServiceBayService) *ServiceBayHandler {
	return &ServiceBayHandler{service: service}
}

// SERVICE BAYS

func (h *ServiceBayHandler) ListServiceBays(ctx *gin.Context) {
	var q dto.ListServiceBaysQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if q.DealershipID != nil && q.BayTypeID != nil {
		bays, err := h.service.ListServiceBaysByDealershipIDAndTypeID(ctx.Request.Context(), *q.DealershipID, *q.BayTypeID)
		if err != nil {
			utils.ResponseError(ctx, err)
			return
		}

		utils.ResponseSuccess(ctx, http.StatusOK, bays)
		return
	}

	if q.DealershipID != nil {
		bays, err := h.service.ListServiceBaysByDealershipID(ctx.Request.Context(), *q.DealershipID)
		if err != nil {
			utils.ResponseError(ctx, err)
			return
		}

		utils.ResponseSuccess(ctx, http.StatusOK, bays)
		return
	}

	if q.BayTypeID != nil {
		bays, err := h.service.ListServiceBaysByTypeID(ctx.Request.Context(), *q.BayTypeID)
		if err != nil {
			utils.ResponseError(ctx, err)
			return
		}

		utils.ResponseSuccess(ctx, http.StatusOK, bays)
		return
	}

	bays, err := h.service.ListServiceBays(ctx.Request.Context())
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, bays)
}

func (h *ServiceBayHandler) GetServiceBayByID(ctx *gin.Context) {
	var req dto.RequestServiceBayWithID
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	bay, err := h.service.GetServiceBayByID(ctx.Request.Context(), req.ID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, bay)
}

func (h *ServiceBayHandler) SearchServiceBays(ctx *gin.Context) {
	var q dto.SearchServiceBaysQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if q.DealershipID != nil && q.BayTypeID != nil {
		bays, err := h.service.SearchServiceBaysByNameDealershipIDAndTypeID(ctx.Request.Context(), q.Name, *q.DealershipID, *q.BayTypeID)
		if err != nil {
			utils.ResponseError(ctx, err)
			return
		}

		utils.ResponseSuccess(ctx, http.StatusOK, bays)
		return
	}

	if q.DealershipID != nil {
		bays, err := h.service.SearchServiceBaysByNameAndDealershipID(ctx.Request.Context(), q.Name, *q.DealershipID)
		if err != nil {
			utils.ResponseError(ctx, err)
			return
		}

		utils.ResponseSuccess(ctx, http.StatusOK, bays)
		return
	}

	if q.BayTypeID != nil {
		bays, err := h.service.SearchServiceBaysByNameAndTypeID(ctx.Request.Context(), q.Name, *q.BayTypeID)
		if err != nil {
			utils.ResponseError(ctx, err)
			return
		}

		utils.ResponseSuccess(ctx, http.StatusOK, bays)
		return
	}

	bays, err := h.service.SearchServiceBaysByName(ctx.Request.Context(), q.Name)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, bays)
}

func (h *ServiceBayHandler) CreateServiceBay(ctx *gin.Context) {
	var req dto.CreateServiceBayRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	err := h.service.CreateServiceBay(ctx.Request.Context(), req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusCreated)
}

func (h *ServiceBayHandler) UpdateServiceBay(ctx *gin.Context) {
	var reqID dto.RequestServiceBayWithID
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var req dto.UpdateServiceBayRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if req.Name == nil && req.DealershipID == nil && req.BayTypeID == nil && req.IsActive == nil {
		utils.ResponseStatusCode(ctx, http.StatusNoContent)
		return
	}

	err := h.service.UpdateServiceBayByID(ctx.Request.Context(), reqID.ID, req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusOK)
}

func (h *ServiceBayHandler) DeleteServiceBay(ctx *gin.Context) {
	var req dto.RequestServiceBayWithID
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if err := h.service.DeleteServiceBayByID(ctx.Request.Context(), req.ID); err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}

// SERVICE BAY TYPES

func (h *ServiceBayHandler) ListServiceBayTypes(ctx *gin.Context) {
	types, err := h.service.ListServiceBayTypes(ctx.Request.Context())
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, types)
}

func (h *ServiceBayHandler) GetServiceBayTypeByID(ctx *gin.Context) {
	var req dto.RequestServiceBayTypeWithID
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	t, err := h.service.GetServiceBayTypeByID(ctx.Request.Context(), req.ID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, t)
}

func (h *ServiceBayHandler) SearchServiceBayTypes(ctx *gin.Context) {
	var q dto.SearchServiceBayTypesQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	types, err := h.service.SearchServiceBayTypesByName(ctx.Request.Context(), q.Name)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, types)
}

func (h *ServiceBayHandler) CreateServiceBayType(ctx *gin.Context) {
	var req dto.CreateServiceBayTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	_, err := h.service.CreateServiceBayType(ctx.Request.Context(), req.Name)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusCreated)
}

func (h *ServiceBayHandler) UpdateServiceBayType(ctx *gin.Context) {
	var reqID dto.RequestServiceBayTypeWithID
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var req dto.UpdateServiceBayTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	_, err := h.service.UpdateServiceBayTypeByID(ctx.Request.Context(), reqID.ID, req.Name)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusOK)
}

func (h *ServiceBayHandler) DeleteServiceBayType(ctx *gin.Context) {
	var req dto.RequestServiceBayTypeWithID
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if err := h.service.DeleteServiceBayTypeByID(ctx.Request.Context(), req.ID); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}
