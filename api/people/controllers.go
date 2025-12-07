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
	"github.com/jackc/pgx/v5"
)

func FetchAllPeople(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "PEOPLE") {
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
	if pkg.HandleDbTxnErr(c, err, "PEOPLE") {
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

	err = tx.Commit(ctx)
	if pkg.HandleDbTxnCommitErr(c, err, "PEOPLE") {
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
	personId, ok := pkg.GrabUuid(c, id, "PEOPLE", "person")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.UpdatePersonEventRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "PEOPLE") {
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
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Person not found for update", err)
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
	personId, ok := pkg.GrabUuid(c, c.Param("id"), "PEOPLE", "person")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "PEOPLE") {
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
