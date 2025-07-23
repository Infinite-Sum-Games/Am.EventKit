package organizers

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func GetAllOrganizers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to acquire DB connection"})
		pkg.Log.ErrorCtx(c, "[ORGANIZER-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()

	organizers, err := q.ListOrganizersQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch organizers"})
		pkg.Log.ErrorCtx(c, "[ORGANIZER-ERROR]: Failed to fetch organizers", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "Organizers list fetched successfully",
		"organizers": organizers,
	})
	pkg.Log.SuccessCtx(c)

}
