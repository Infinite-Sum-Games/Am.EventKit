package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func HostelCheckIn(c *gin.Context) {
	personellId, ok := pkg.GrabUserId(c, "HOSTEL")
	if !ok {
		return
	}

	hospId := c.Param("hospId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "") {
		return
	}
	defer conn.Release()

	q := db.New()
	ok, err := q.HostelCheckInQuery(c, conn, db.HostelCheckInQueryParams{
		AccomodationID: ,
		CheckedInBy: personellId,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Hostel check-in successful",
	})
	pkg.Log.SuccessCtx(c)
}
