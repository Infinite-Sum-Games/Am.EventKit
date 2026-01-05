package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/models"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func GetAllHostels(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "GET-HOSTELS") {
		return
	}
	defer conn.Release()

	q := db.New()

	hostels, err := q.GetAllHostelsQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[GET-HOSTELS-ERROR]: Failed to fetch hostels", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hostels fetched successfully",
		"hostels": hostels,
	})
	pkg.Log.SuccessCtx(c)
}

func MapQrStudentId(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.MapQrStudentIdRequest](c)
	if !ok {
		return
	}

	studentId, ok := pkg.GrabUuid(c, req.StudentID, "MAP-QR-STUDENT", "studentID")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "MAP-QR-STUDENT") {
		return
	}
	defer conn.Release()

	q := db.New()

	res, err := q.MapQrStudentIdQuery(ctx, conn, db.MapQrStudentIdQueryParams{
		ID:            studentId,
		HospitalityID: pkg.ToPgText(req.HospitalityId),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[MAP-QR-STUDENT-ERROR]: Failed to map QR code to student ID", err)
		return
	}

	// Check if accommodation exists (can be NULL if student hasn't opted)
	if res.AccommodationID == uuid.Nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Student has not opted for accommodation",
		})
		pkg.Log.WarnCtx(c, "[MAP-QR-STUDENT-WARN]: Student has no accommodation opted")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":                 "QR code mapped successfully",
		"accommodation_id":        res.AccommodationID,
		"has_opted_accommodation": res.HasOptedAccommodation,
	})
	pkg.Log.SuccessCtx(c)
}

func GetAccommodationById(c *gin.Context) {
	accIdStr := c.Param("accId")
	accommodationId, ok := pkg.GrabUuid(c, accIdStr, "GET-ACCOMMODATION", "Accommodation")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "GET-ACCOMMODATION") {
		return
	}
	defer conn.Release()

	q := db.New()

	accommodation, err := q.GetAccommodationByIdQuery(ctx, conn, accommodationId)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Accommodation request not found",
		})
		pkg.Log.WarnCtx(c, "[GET-ACCOMMODATION-WARN]: Accommodation ID does not exist")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[GET-ACCOMMODATION-ERROR]: Failed to fetch accommodation request", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Accommodation request fetched successfully",
		"accommodation": accommodation,
	})
	pkg.Log.SuccessCtx(c)
}

func UpdateAccommodationById(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.UpdateAccommodationByIdRequest](c)
	if !ok {
		return
	}

	checkIn, err := pkg.ParseDateTime(req.CheckInDate, req.CheckInTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid check-in date or time format",
		})
		pkg.Log.ErrorCtx(c, "[UPDATE-ACCOMMODATION-WARN]: Invalid check-in date or time format", err)
		return
	}

	fmt.Println("CHECK IN:", checkIn)

	checkInPgTimestamp := pkg.ToPgTimestamp(checkIn)
	fmt.Println("CHECK IN PG:", checkInPgTimestamp)
	checkOut, err := pkg.ParseDateTime(req.CheckOutDate, req.CheckOutTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid check-out date or time format",
		})
		pkg.Log.ErrorCtx(c, "[UPDATE-ACCOMMODATION-WARN]: Invalid check-out date or time format", err)
		return
	}
	fmt.Println("CHECK OUT:", checkOut)

	checkOutPgTimestamp := pkg.ToPgTimestamp(checkOut)
	fmt.Println("CHECK OUT PG:", checkOutPgTimestamp)

	accommodationIdStr := c.Param("accommodationId")
	accommodationId, ok := pkg.GrabUuid(c, accommodationIdStr, "UPDATE-ACCOMMODATION", "accommodationID")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "UPDATE-ACCOMMODATION") {
		return
	}
	defer conn.Release()

	q := db.New()

	row, err := q.UpdateAccommodationByIdQuery(ctx, conn, db.UpdateAccommodationByIdQueryParams{
		ID:                accommodationId,
		IsMale:            req.IsMale,
		IsHosteller:       req.IsHosteller,
		CollegeRollNumber: req.CollegeRollNumber,
		CollegeName:       req.CollegeName,
		RoomPreference:    req.RoomPreference,
		IsAmritaCampus:    req.IsAmritaCampus,
		CheckIn:           pkg.ToPgTimestamp(checkIn),
		CheckOut:          pkg.ToPgTimestamp(checkOut),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[UPDATE-ACCOMMODATION-ERROR]: Failed to update accommodation request", err)
		return
	}
	if row == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Accommodation request not found",
		})
		pkg.Log.WarnCtx(c, "[UPDATE-ACCOMMODATION-WARN]: Accommodation ID does not exist")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Accommodation request updated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func GateCheckIn(c *gin.Context) {
	_ = c.GetString("hospId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "GATE") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "GATE")

	_ = db.New()

	// ok, err := q.GateCheckInQuery(tx, ctx, hospId)
	// if err != nil {
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "Student checked-in successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func GateCheckOut(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "GATE") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "GATE")

	_ = db.New()

	c.JSON(http.StatusOK, gin.H{
		"message": "Student checked-out successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func GateStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "GATE") {
		return
	}
	defer conn.Release()

	_ = db.New()

	// If has accommodation, check entry time and show exit time
	// If does not have accommodation, show list of logs

	c.JSON(http.StatusOK, gin.H{
		"message": "Student accomodation status fetched successfully",
	})
	pkg.Log.SuccessCtx(c)
}
