package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func HostelCheckIn(c *gin.Context) {
	personellIdStr, ok := pkg.GrabUserId(c, "HOSTEL")
	if !ok {
		return
	}
	personellId, ok := pkg.GrabUuid(c, personellIdStr, "HOSTEL", "Personell")
	if !ok {
		return
	}

	hospId := c.Param("hospId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "HOSTEL") {
		return
	}
	defer conn.Release()

	q := db.New()
	details, err := q.HostelCheckInQuery(c, conn, db.HostelCheckInQueryParams{
		HospitalityID: pgtype.Text{String: hospId, Valid: true},
		CheckedInBy:   personellId,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Hostel check-in successful",
		"details": details,
	})
	pkg.Log.SuccessCtx(c)
}
