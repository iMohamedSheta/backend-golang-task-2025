package utils

import (
	"taskgo/pkg/contracts"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func BindToRequestAndExtractFields(c *gin.Context, request contracts.Validatable) error {
	if err := c.ShouldBindBodyWith(request, binding.JSON); err != nil {
		return err
	}

	// Extract raw fields
	var raw map[string]any
	if err := c.ShouldBindBodyWith(&raw, binding.JSON); err != nil {
		return err
	}

	request.SetRequestSentFields(raw)

	return nil
}
