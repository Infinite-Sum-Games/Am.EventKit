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

func GetAllAccomodationRequests(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ACCOMODATION") {
		return
	}
	defer conn.Release()

	q := db.New()
	requests, err := q.GetAllAccommodationRequestsQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ACCOMMODATION-ERROR]: Failed to fetch accomodation requests", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Accommodation requests returned successfully",
		"requests": requests,
	})
	pkg.Log.SuccessCtx(c)
}

func AccomodationLogin(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.LoginRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "AUTH") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "AUTH")

	q := db.New()
	result, err := q.LoginHospitalityQuery(ctx, tx, req.Email)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Hospitality user not found",
		})
		pkg.Log.WarnCtx(c, "[AUTH-WARN]: Email does not exist")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to fetch login details", err)
		return
	}

	err = pkg.CompareHash(result.Password, req.HashedPassword)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Invalid credentials",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Password verification failed", err)
		return
	}

	if result.RefreshToken.String == "" {
		refreshToken, err := pkg.CreateRefreshToken(result.ID.String(), req.Email, pkg.Roles{
			IsUser:        false,
			IsOrganizer:   false,
			IsAdmin:       false,
			IsHospitality: true,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to create refresh token", err)
			return
		}

		token, err := q.UpdateHospitalityRefreshTokenQuery(ctx, tx, db.UpdateHospitalityRefreshTokenQueryParams{
			RefreshToken: pgtype.Text{
				String: refreshToken,
				Valid:  true,
			},
			Email: req.Email,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Could not add refresh token to DB", err)
			return
		}
		pkg.SetRefreshCookie(c, token.String)
	} else {
		pkg.SetRefreshCookie(c, result.RefreshToken.String)
	}

	err = tx.Commit(ctx)
	if pkg.HandleDbTxnCommitErr(c, err, "AUTH") {
		return
	}

	authToken, err := pkg.CreateAuthToken(result.ID.String(), result.Email, pkg.Roles{
		IsUser:        false,
		IsOrganizer:   false,
		IsAdmin:       false,
		IsHospitality: true,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to create auth token", err)
		return
	}

	pkg.SetAuthCookie(c, authToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Accomodation personnel logged in successful",
		"name":    result.Name,
		"email":   result.Email,
	})
	pkg.Log.SuccessCtx(c)
}

func AccomodationLogout(c *gin.Context) {
	pkg.NullifyCookies(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "Accomodation personnel logged out successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func AccomodationSession(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "SESSION")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "SESSION") {
		return
	}
	defer conn.Release()

	q := db.New()

	result, err := q.FetchHospitalitySessionQuery(ctx, conn, email)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No session found for user",
		})
		pkg.Log.WarnCtx(c, "[SESSION-WARN]: User might deleted but cookies exist")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[SESSION-ERROR]: Failed to fetch user session", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User session obtained successfully",
		"name":    result.Name,
		"email":   result.Email,
	})
	pkg.Log.SuccessCtx(c)
}

func AddHostel(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.AddHostelRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADD-HOSTEL") {
		return
	}
	defer conn.Release()

	q := db.New()

	hostelId, err := q.AddHostelQuery(ctx, conn, db.AddHostelQueryParams{
		RoomCount:   req.RoomCount,
		IsMale:      req.IsMale,
		WardenEmail: pkg.ToPgTextPtr(&req.WardenEmail),
		Latitude:    pkg.ToPgTextPtr(&req.Latitude),
		Longtitude:  pkg.ToPgTextPtr(&req.Longtitude),
		MapUrl:      pkg.ToPgTextPtr(&req.MapUrl),
		HostelName:  req.HostelName,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[ADD-HOSTEL-ERROR]: Failed to add new hostel", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Hostel added successfully",
		"hostel_id": hostelId,
	})
	pkg.Log.SuccessCtx(c)
}

func UpdateHostel(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.UpdateHostelRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "UPDATE-HOSTEL") {
		return
	}
	defer conn.Release()

	q := db.New()

	HostelUuid, ok := pkg.GrabUuid(c, req.HostelID, "UPDATE-HOSTEL", "hostelID")
	if !ok {
		return
	}

	rows, err := q.UpdateHostelQuery(ctx, conn, db.UpdateHostelQueryParams{
		ID:          HostelUuid,
		RoomCount:   req.RoomCount,
		WardenEmail: pkg.ToPgTextPtr(&req.WardenEmail),
		Latitude:    pkg.ToPgTextPtr(&req.Latitude),
		Longtitude:  pkg.ToPgTextPtr(&req.Longtitude),
		MapUrl:      pkg.ToPgTextPtr(&req.MapUrl),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[UPDATE-HOSTEL-ERROR]: Failed to update hostel details", err)
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Hostel not found",
		})
		pkg.Log.WarnCtx(c, "[UPDATE-HOSTEL-WARN]: Hostel ID does not exist")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hostel updated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func DeleteHostel(c *gin.Context) {
	hostelIdStr := c.Param("id")
	hostelId, ok := pkg.GrabUuid(c, hostelIdStr, "DELETE-HOSTEL", "hostelID")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "DELETE-HOSTEL") {
		return
	}
	defer conn.Release()

	q := db.New()

	rows, err := q.DeleteHostelQuery(ctx, conn, hostelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[DELETE-HOSTEL-ERROR]: Failed to delete hostel", err)
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Hostel not found",
		})
		pkg.Log.WarnCtx(c, "[DELETE-HOSTEL-WARN]: Hostel ID does not exist")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hostel deleted successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func AllotHostel(c *gin.Context) {
	accommodationIdStr := c.Param("accommodationId")
	accommodationId, ok := pkg.GrabUuid(c, accommodationIdStr, "ALLOT-HOSTEL", "hostelID")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.AllotHostelRequest](c)
	if !ok {
		return
	}

	hostleIdPgUuid, err := pkg.ToPgUuidPtr(&req.HostelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid hostel ID format",
		})
		pkg.Log.ErrorCtx(c, "[ALLOT-HOSTEL-WARN]: Invalid hostel ID format", err)
		return
	}

	hostelIdUuid, ok := pkg.GrabUuid(c, req.HostelID, "ALLOT-HOSTEL", "hostelID")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "ALLOT-HOSTEL") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "ALLOT-HOSTEL")

	q := db.New()

	hostel, err := q.GetHostelQuery(ctx, tx, hostelIdUuid)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Hostel not found",
		})
		pkg.Log.WarnCtx(c, "[ALLOT-HOSTEL-WARN]: Hostel ID does not exist")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[ALLOT-HOSTEL-ERROR]: Failed to fetch hostel details", err)
		return
	}

	if hostel.RoomCount <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No rooms available in the selected hostel",
		})
		pkg.Log.WarnCtx(c,
			"[ALLOT-HOSTEL-WARN]: No rooms available in the selected hostel")
		return
	}

	row, err := q.DecrementHostelRoomCountQuery(ctx, tx, hostelIdUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c,
			"[ALLOT-HOSTEL-ERROR]: Failed to decrement hostel room count", err)
		return
	}
	if row == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Hostel not found",
		})
		pkg.Log.ErrorCtx(c,
			"[ALLOT-HOSTEL-ERROR]: Hostel ID does not exist", err)
		return
	}

	row, err = q.AllotHostelQuery(ctx, tx, db.AllotHostelQueryParams{
		ID:       accommodationId,
		HostelID: hostleIdPgUuid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[ALLOT-HOSTEL-ERROR]: Failed to allot hostel", err)
		return
	}
	if row == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Accommodation request or Hostel not found",
		})
		pkg.Log.ErrorCtx(c,
			"[ALLOT-HOSTEL-ERROR]: Accommodation ID or Hostel ID does not exist", err)
		return
	}

	err = tx.Commit(ctx)
	if pkg.HandleDbTxnCommitErr(c, err, "ALLOT-HOSTEL") {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hostel allotted successfully",
	})
	pkg.Log.SuccessCtx(c)

}

func AffirmAccommodationPayment(c *gin.Context) {
	accommodationIdStr := c.Param("accommodationId")
	accommodationId, ok := pkg.GrabUuid(c, accommodationIdStr, "AFFIRM-PAYMENT", "Accommodation")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "AFFIRM-PAYMENT") {
		return
	}
	defer conn.Release()

	q := db.New()

	rows, err := q.AffirmAccommodationPaymentQuery(ctx, conn, accommodationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AFFIRM-PAYMENT-ERROR]: Failed to affirm payment", err)
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Accommodation request not found",
		})
		pkg.Log.WarnCtx(c, "[AFFIRM-PAYMENT-WARN]: Accommodation ID does not exist")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Accommodation payment affirmed successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func GetAccommodationById(c *gin.Context) {
	accommodationIdStr := c.Param("accommodationId")
	accommodationId, ok := pkg.GrabUuid(c, accommodationIdStr, "GET-ACCOMMODATION", "accommodationID")
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
