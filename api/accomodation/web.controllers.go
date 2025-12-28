package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/models"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func AccomodationExists(c *gin.Context) {
	userIdStr, ok := pkg.GrabUserId(c, "ACCOMODATION")
	if !ok {
		return
	}
	userId, ok := pkg.GrabUuid(c, userIdStr, "ACCOMODATION", "Student")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ACCOMODATION") {
		return
	}
	defer conn.Release()

	q := db.New()
	hasRegistration, err := q.CheckStudentHasTicketQuery(ctx, conn, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ACCOMMODATION-ERROR]: Failed to check if student has a ticket", err)
		return
	}
	if !hasRegistration {
		c.JSON(http.StatusOK, gin.H{
			"has_accommodation": "NOT_REGISTERED",
		})
		pkg.Log.WarnCtx(c, "[ACCOMODATION-WARN]: Attempt to register without a ticket")
		return
	}

	hasAccomodation, err := q.CheckUserAccomodationExistsQuery(ctx, conn, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ACCOMODATION-ERROR]: Failed to check if user has accomodation", err)
		return
	}
	if hasAccomodation {
		c.JSON(http.StatusOK, gin.H{
			"has_accommodation": "FILLED_ACCOMMODATION",
		})
		pkg.Log.WarnCtx(c, "[ACCOMMODATION-WARN]: Already filled accomodation form")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"has_accommodation": "ELIGIBLE",
	})
	pkg.Log.SuccessCtx(c)
}

func AccomodationFormCsrf(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "ACCOMODATION")
	if !ok {
		return
	}

	csrfToken, err := pkg.CreateCsrfToken(email, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[ACCOMODATION-ERROR]: Failed to create CSRF token", err)
		return
	}

	pkg.SetCsrfCookie(c, csrfToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully initiated accomodation form filling",
		"key":     csrfToken,
	})
	pkg.Log.SuccessCtx(c)
}

func AccomodationFormSubmission(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "ACCOMODATION")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.AccomodationForm](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "ACCOMODATION") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "ACCOMODATION")

	q := db.New()
	person, err := q.FetchStudentMetadata(ctx, tx, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ACCOMODATION-ERROR]: Failed to fetch student metadata", err)
		return
	}

	_, err = q.InsertAccomodationFormEntryQuery(ctx, tx,
		db.InsertAccomodationFormEntryQueryParams{
			StudentID:         person.ID,
			Name:              person.Name,
			Email:             email,
			PhoneNumber:       person.PhoneNumber,
			IsMale:            req.IsMale,
			IsHosteller:       req.IsHosteller,
			IsAmritaCampus:    req.IsAmritaCampus,
			CollegeName:       req.CollegeName,
			RoomPreference:    req.RoomPreference,
			CollegeRollNumber: req.CollegeRollNumber,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ACCOMODATION-ERROR]: ", err)
		return
	}

	err = tx.Commit(ctx)
	ok = pkg.HandleDbTxnCommitErr(c, err, "ACCOMODATION")
	if ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Form submitted successfully",
	})
	pkg.Log.SuccessCtx(c)
}
