package handlers

import (
	"taskgo/pkg/contracts"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Handler represents the base handler for all handlers
type Handler struct{}

// BindAndExtract handles JSON binding + raw field extraction
func (h *Handler) BindBodyAndExtractToRequest(gin *gin.Context, req contracts.Validatable) error {
	if err := gin.ShouldBindBodyWith(req, binding.JSON); err != nil {
		return err
	}

	var raw map[string]any
	if err := gin.ShouldBindBodyWith(&raw, binding.JSON); err != nil {
		return err
	}

	req.SetRequestSentFields(raw)
	return nil
}
