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
	"github.com/jackc/pgx/v5/pgtype"
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
		c.JSON(http.StatusOK, gin.H{
			"message":                 "Student has not opted for accommodation",
			"has_opted_accommodation": res.HasOptedAccommodation,
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

	accommodationIdStr := c.Param("accId")
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
	personellIdStr, ok := pkg.GrabUserId(c, "GATE")
	if !ok {
		return
	}
	personellId, ok := pkg.GrabUuid(c, personellIdStr, "GATE", "Personell")
	if !ok {
		return
	}

	hospId := c.Param("hospId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "GATE") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "GATE")

	q := db.New()

	direction, err := q.GateCheckInOutQuery(ctx, tx, db.GateCheckInOutQueryParams{
		HospitalityID: pgtype.Text{
			String: hospId,
			Valid:  true,
		},
		Direction:   db.GateLogDirectionEnumIN,
		PersonellID: personellId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[GATE-ERROR]: Failed to check in", err)
		return
	}

	err = tx.Commit(ctx)
	if pkg.HandleDbTxnCommitErr(c, err, "GATE") {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Student marked successfully",
		"direction": direction,
	})
	pkg.Log.SuccessCtx(c)
}

func GateCheckOut(c *gin.Context) {
	personellIdStr, ok := pkg.GrabUserId(c, "GATE")
	if !ok {
		return
	}
	personellId, ok := pkg.GrabUuid(c, personellIdStr, "GATE", "Personell")
	if !ok {
		return
	}

	hospId := c.Param("hospId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "GATE") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "GATE")

	q := db.New()

	direction, err := q.GateCheckInOutQuery(ctx, tx, db.GateCheckInOutQueryParams{
		HospitalityID: pgtype.Text{
			String: hospId,
			Valid:  true,
		},
		Direction:   db.GateLogDirectionEnumOUT,
		PersonellID: personellId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[GATE-ERROR]: Failed to check out", err)
		return
	}

	err = tx.Commit(ctx)
	if pkg.HandleDbTxnCommitErr(c, err, "GATE") {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Student marked successfully",
		"direction": direction,
	})
	pkg.Log.SuccessCtx(c)
}

func GateStatus(c *gin.Context) {
	hospId := c.Param("hospId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "GATE") {
		return
	}
	defer conn.Release()

	q := db.New()

	resp, err := q.HostelGateStatusQuery(ctx, conn, pgtype.Text{
		String: hospId,
		Valid:  true,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[GATE-ERROR]: Failed to check gate status", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Student accomodation status fetched successfully",
		"student_name":    resp.StudentName,
		"student_email":   resp.StudentEmail,
		"single_check_in": resp.SingleCheckIn,
		"last_check_in":   resp.LastCheckIn,
		"last_check_out":  resp.LastCheckOut,
	})
	pkg.Log.SuccessCtx(c)
}

// If No Accomodation, then day scholar
//   - Everyday One CheckIn.
//
// If Accomodation
//   - LifeTime One CheckIn
func GateCheckInStatus(c *gin.Context) {
	hospId := c.Param("hospId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "GATE-CHECKIN-STATUS") {
		return
	}
	defer conn.Release()

	q := db.New()
	res, err := q.GateCheckStatusQuery(ctx, conn, pgtype.Text{
		String: hospId,
		Valid:  true,
	})
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message":      "Invalid hospitality ID",
			"allow":        false,
			"reason":       "Invalid hospitality ID provided.",
			"name":         res.Name,
			"email":        res.Email,
			"college_name": res.CollegeName,
		})
		pkg.Log.WarnCtx(c, "[GATE-CHECKIN-WARN]: Invalid hospitality ID")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":      "Oops! Something happened. Please try again later",
			"allow":        false,
			"reason":       "Server error.",
			"name":         res.Name,
			"email":        res.Email,
			"college_name": res.CollegeName,
		})
		pkg.Log.ErrorCtx(c, "[GATE-CHECKIN-ERROR]: Failed to check gate status", err)
		return
	}

	// hasAccomodation would be true if payment_status is not NULL or empty
	hasAccomodation := res.AccomodationStatus.Valid && res.AccomodationStatus.String != ""

	if hasAccomodation {
		// Rule: Cannot check-in if already inside
		// A person is inside if they have a valid check-in and either no valid checkout or check-in is after checkout.
		isAlreadyInside := res.LastCheckIn.Valid && (!res.LastCheckOut.Valid || res.LastCheckIn.Time.After(res.LastCheckOut.Time))

		if isAlreadyInside {
			// Rule: Cannot check-in if already inside
			// Already inside is defined as last_check_in > last_check_out
			c.JSON(http.StatusOK, gin.H{
				"message":      "Check-in status",
				"allow":        false,
				"reason":       "Already check-in.",
				"name":         res.Name,
				"email":        res.Email,
				"college_name": res.CollegeName,
			})
			pkg.Log.InfoCtx(c, "[GATE-CHECKIN-STATUS]: Denied check-in for accomodation holder (already inside)")
			return
		}
		// If not inside then fall through to success case
	} else { // does not have accomodation
		// Rule: Day scholars can check-in once per day
		if res.LastCheckIn.Valid {
			now := time.Now()
			lastCheckInTime := res.LastCheckIn.Time

			if lastCheckInTime.Year() == now.Year() && lastCheckInTime.YearDay() == now.YearDay() {
				c.JSON(http.StatusOK, gin.H{
					"message":      "Check-in status",
					"allow":        false,
					"reason":       "Already checked in today",
					"name":         res.Name,
					"email":        res.Email,
					"college_name": res.CollegeName,
				})
				pkg.Log.InfoCtx(c, "[GATE-CHECKIN-STATUS]: Denied check-in for day scholar (already checked in today)")
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Checkin status sent successfully",
		"allow":        true,
		"reason":       "",
		"name":         res.Name,
		"email":        res.Email,
		"college_name": res.CollegeName,
	})
	pkg.Log.SuccessCtx(c)
}

// If No Accomodation, then day scholar
//   - Everyday One CheckOut
//
// If Accomodation
//   - LifeTime One CheckOut
func GateCheckOutStatus(c *gin.Context) {
	hospId := c.Param("hospId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "GATE-CHECKOUT-STATUS") {
		return
	}
	defer conn.Release()

	q := db.New()

	res, err := q.GateCheckStatusQuery(ctx, conn, pgtype.Text{
		String: hospId,
		Valid:  true,
	})
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message":      "Invalid hospitality ID",
			"allow":        false,
			"reason":       "Invalid hospitality ID provided",
			"name":         res.Name,
			"email":        res.Email,
			"college_name": res.CollegeName,
		})
		pkg.Log.WarnCtx(c, "[GATE-CHECKOUT-WARN]: Invalid hospitality ID")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[GATE-CHECKOUT-ERROR]: Failed to check gate status", err)
		return
	}

	// A user can checkout only if they are currently checked in. A user is
	// inside if their last checkin is more recent that their last checkout
	isCurrentlyInside := res.LastCheckIn.Valid && (!res.LastCheckOut.Valid || res.LastCheckIn.Time.After(res.LastCheckOut.Time))

	if !isCurrentlyInside {
		c.JSON(http.StatusOK, gin.H{
			"message":      "Checkout status",
			"allow":        false,
			"reason":       "Not currently checked in.",
			"name":         res.Name,
			"email":        res.Email,
			"college_name": res.CollegeName,
		})
		pkg.Log.InfoCtx(c, "[GATE-CHECKOUT-STATUS]: Denied checkout (not inside)")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Checkout status sent successfully",
		"allow":        true,
		"reason":       "",
		"name":         res.Name,
		"email":        res.Email,
		"college_name": res.CollegeName,
	})
	pkg.Log.SuccessCtx(c)
}
