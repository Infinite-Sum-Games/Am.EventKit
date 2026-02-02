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
// 	"github.com/jackc/pgx/v5"
// )
//
// func SecurityCheck(c *gin.Context) {
// 	hospId := c.Param("hospId")
//
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
//
// 	conn, err := cmd.DBPool.Acquire(ctx)
// 	if pkg.HandleDbAcquireErr(c, err, "GET-STUDENT-DETAILS") {
// 		return
// 	}
// 	defer conn.Release()
//
// 	q := db.New()
//
// 	studentDetails, err := q.GetStudentDetailsForSecurityQuery(ctx, conn, pkg.ToPgText(hospId))
// 	if err == pgx.ErrNoRows {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"message": "Accommodation request not found",
// 		})
// 		pkg.Log.WarnCtx(c, "[GET-STUDENT-DETAILS-WARN]: Accommodation ID does not exist")
// 		return
// 	}
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Oops! Something happened. Please try again later.",
// 		})
// 		pkg.Log.ErrorCtx(c, "[GET-STUDENT-DETAILS-ERROR]: Failed to fetch student details", err)
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"message":        "Student details fetched successfully",
// 		"name":           studentDetails.StudentName,
// 		"email":          studentDetails.StudentEmail,
// 		"college_name":   studentDetails.CollegeName,
// 		"college_roll":   studentDetails.CollegeRollNumber,
// 		"hostel_name":    studentDetails.HostelName,
// 		"check_in-date":  studentDetails.CheckInDate,
// 		"check_out_date": studentDetails.CheckOutDate,
// 	})
// }
