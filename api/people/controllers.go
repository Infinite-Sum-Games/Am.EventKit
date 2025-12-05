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
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func FetchAllPeople(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.FatalCtx(c, "[PEOPLE-FATAL]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()

	people, err := q.FetchAllPeopleQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to fetch people", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "People list fetched successfully",
		"people":  people,
	})
	pkg.Log.SuccessCtx(c)
}

func AddNewPerson(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.CreateNewPersonWithEventRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.FatalCtx(c, "[PEOPLE-FATAL]: Failed to begin DB transaction", err)
		return
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			pkg.Log.FatalCtx(c, "[PEOPLE-FATAL]: Failed to rollback DB transaction", err)
		}
	}()

	q := db.New()

	people, err := q.AddNewPersonQuery(ctx, tx, db.AddNewPersonQueryParams{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Profession:  pkg.ToPgTextPtr(req.Profession),
		Email:       pkg.ToPgTextPtr(req.Email),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to add new person", err)
		return
	}

	personToEventMapping, err := q.MapPersonToEventQuery(ctx, tx,
		db.MapPersonToEventQueryParams{
			PersonID: people.ID,
			EventID:  req.EventID,
			EventDay: req.EventDay,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to map person to event", err)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.FatalCtx(c, "[PEOPLE-FATAL]: Failed to commit DB transaction", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":                 "Person added successfully",
		"person_details":          people,
		"person_to_event_mapping": personToEventMapping,
	})
	pkg.Log.SuccessCtx(c)
}

func UpdatePersonDetails(c *gin.Context) {
	id := c.Param("id")
	personId, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Request not processed due to invalid parameters",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Invalid person ID parameter", err)
		return
	}

	req, ok := pkg.ValidateRequest[models.UpdatePersonEventRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.FatalCtx(c, "[PEOPLE-FATAL]: Failed to begin DB transaction", err)
		return
	}
	defer conn.Release()

	q := db.New()

	updatedPerson, err := q.UpdatePersonDetailsQuery(ctx, conn,
		db.UpdatePersonDetailsQueryParams{
			ID:          personId,
			Name:        req.Name,
			PhoneNumber: req.PhoneNumber,
			Profession:  pkg.ToPgTextPtr(req.Profession),
			Email:       pkg.ToPgTextPtr(req.Email),
		})
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Person not found",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Person not found for update", nil)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to update person details", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Person details updated successfully",
		"id":           updatedPerson.ID,
		"name":         updatedPerson.Name,
		"email":        updatedPerson.Email,
		"profession":   updatedPerson.Profession,
		"phone_number": updatedPerson.PhoneNumber,
	})
	pkg.Log.SuccessCtx(c)
}

func DeletePerson(c *gin.Context) {
	id := c.Param("id")
	personId, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Request is malformed",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Invalid person ID parameter", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.FatalCtx(c, "[PEOPLE-FATAL]: Failed to acquire connection from DB", err)
		return
	}
	defer conn.Release()

	q := db.New()

	deletedPersonId, err := q.DeletePersonQuery(ctx, conn, personId)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Person not found",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Person not found for deletion", nil)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to delete person", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Person deleted successfully",
		"deleted_person_id": deletedPersonId,
	})
	pkg.Log.SuccessCtx(c)
}
