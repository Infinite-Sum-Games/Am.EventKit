package tag

import (
	"context"
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func GetAllTags(c *gin.Context) {
	queries := db.New()
	tags, err := queries.ListTags(context.Background(), cmd.DBPool)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Tags list fetched successfully",
		"tags":    tags,
	})
	pkg.Log.SuccessCtx(c)
}
