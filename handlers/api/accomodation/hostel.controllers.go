package api

//
// import (
// 	"context"
// 	"net/http"
// 	"time"
//
// 	"github.com/Infinite-Sum-Games/Am.EventKit/cmd"
// 	db "github.com/Infinite-Sum-Games/Am.EventKit/db/gen"
// 	"github.com/Infinite-Sum-Games/Am.EventKit/pkg"
// 	"github.com/gin-gonic/gin"
// 	"github.com/jackc/pgx/v5/pgtype"
// )
//
// func HostelCheckIn(c *gin.Context) {
// 	personellIdStr, ok := pkg.GrabUserId(c, "HOSTEL")
// 	if !ok {
// 		return
// 	}
// 	personellId, ok := pkg.GrabUuid(c, personellIdStr, "HOSTEL", "Personell")
// 	if !ok {
// 		return
// 	}
//
// 	hospId := c.Param("hospId")
//
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
//
// 	conn, err := cmd.DBPool.Acquire(ctx)
// 	if pkg.HandleDbAcquireErr(c, err, "HOSTEL") {
// 		return
// 	}
// 	defer conn.Release()
//
// 	q := db.New()
// 	details, err := q.HostelCheckInQuery(c, conn, db.HostelCheckInQueryParams{
// 		HospitalityID: pgtype.Text{String: hospId, Valid: true},
// 		CheckedInBy:   personellId,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Oops! Something happened. Please try again later",
// 		})
// 		pkg.Log.ErrorCtx(c, "[HOSTEL-ERROR]: Failed to hostel check-in", err)
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Hostel check-in successful",
// 		"details": details,
// 	})
// 	pkg.Log.SuccessCtx(c)
// }
