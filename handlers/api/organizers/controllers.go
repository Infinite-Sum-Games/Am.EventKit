package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Infinite-Sum-Games/Am.EventKit/cmd"
	db "github.com/Infinite-Sum-Games/Am.EventKit/db/gen"
	"github.com/Infinite-Sum-Games/Am.EventKit/models"
	"github.com/Infinite-Sum-Games/Am.EventKit/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetAllOrganizers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ORGANIZER") {
		return
	}
	defer conn.Release()

	q := db.New()
	organizers, err := q.ListOrganizersQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ORGANIZER-ERROR]: Failed to fetch organizers", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Organizers list fetched successfully",
		"organizers": organizers,
	})
	pkg.Log.SuccessCtx(c)
}

func CreateOrganizer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, ok := pkg.ValidateRequest[models.CreateOrganizerRequest](c)
	if !ok {
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ORGANIZER") {
		return
	}
	defer conn.Release()

	q := db.New()

	hashedPassword, err := pkg.Hash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.FatalCtx(c, "[ORGANIZER-FATAL]: Failed to hash password", err)
		return
	}

	err = q.CreateOrganizerQuery(ctx, conn, db.CreateOrganizerQueryParams{
		Name:          req.Name,
		Email:         req.Email,
		Password:      hashedPassword,
		OrgType:       db.OrganizerTypeEnum(req.OrgType),
		StudentHead:   req.StudentHead,
		StudentCoHead: pkg.ToPgText(req.StudentCoHead),
		FacultyHead:   req.FacultyHead,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ORGANIZER-ERROR]: Failed to create organizer", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Organizer created successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func EditOrganizer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orgID, ok := pkg.GrabUuid(c, c.Param("organizerId"), "ORGANIZER", "organizer")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.EditOrganizerRequest](c)
	if !ok {
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ORGANIZER") {
		return
	}
	defer conn.Release()

	q := db.New()
	rows, err := q.UpdateOrganizerByIDQuery(ctx, conn, db.UpdateOrganizerByIDQueryParams{
		ID:            orgID,
		Name:          req.Name,
		Email:         req.Email,
		OrgType:       db.OrganizerTypeEnum(req.OrgType),
		StudentHead:   req.StudentHead,
		StudentCoHead: pkg.ToPgText(req.StudentCoHead),
		FacultyHead:   req.FacultyHead,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ORGANIZER-ERROR]: Failed to update organizer", err)
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Organizer does not exist",
		})
		pkg.Log.ErrorCtx(c, "[ORGANIZER-ERROR]: Organizer does not exist for update", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organizer updated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func ChangeOrganizerPassword(c *gin.Context) {
	orgId, ok := pkg.GrabUuid(c, c.Param("organizerId"), "ORGANIZER", "Organizer")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.ChangeOrganizerPasswordRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ORGANIZER") {
		return
	}
	defer conn.Release()

	passwordHash, err := pkg.Hash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.FatalCtx(c, "[ORGANIZER-FATAL]: Failed to hash password", err)
		return
	}

	q := db.New()
	_, err = q.ChangeOrganizerPasswordQuery(ctx, conn,
		db.ChangeOrganizerPasswordQueryParams{
			ID:       orgId,
			Password: passwordHash,
		})
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Organizer not found",
		})
		pkg.Log.ErrorCtx(c, "[ORGANIZER-ERROR]: Failed to update organizer password", err)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ORGANIZER-ERROR]: Failed to update organizer password", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organizer password updated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func DeleteOrganizer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orgID, ok := pkg.GrabUuid(c, c.Param("organizerId"), "ORGANIZER", "organizer")
	if !ok {
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ORGANIZER") {
		return
	}
	defer conn.Release()

	q := db.New()
	rows, err := q.DeleteOrganizerByIDQuery(ctx, conn, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ORGANIZER-ERROR]: Failed to delete organizer", err)
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Organizer does not exist",
		})
		pkg.Log.ErrorCtx(c, "[ORGANIZER-ERROR]: Organizer does not exist for delete", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organizer deleted successfully",
	})
	pkg.Log.SuccessCtx(c)
}
