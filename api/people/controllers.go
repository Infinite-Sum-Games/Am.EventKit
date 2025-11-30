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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()

	people, err := q.FetchAllPeopleQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to fetch people", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "People list fetched successfully",
		"people":  people,
	})
	pkg.Log.SuccessCtx(c)
}

func FetchPeopleByDepartment(c *gin.Context) {
	dept := c.Param("dept")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()

	people, err := q.FetchPeopleByDepartmentQuery(ctx, conn, dept)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to fetch people by department", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "People list fetched successfully",
		"people":  people,
	})
	pkg.Log.SuccessCtx(c)
}

func FetchPeopleByEvent(c *gin.Context) {
	event := c.Param("event")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()

	people, err := q.FetchPeopleByEventQuery(ctx, conn, event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to fetch people by event", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "People list fetched successfully",
		"people":  people,
	})
	pkg.Log.SuccessCtx(c)

}

func FetchPeopleByDay(c *gin.Context) {
	day := c.Param("day")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()

	people, err := q.FetchPeopleByDayQuery(ctx, conn, day)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to fetch people by day", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "People list fetched successfully",
		"people":  people,
	})
	pkg.Log.SuccessCtx(c)

}

func AddNewPerson(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.AddNewPersonWithEventParams](c)
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
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to begin DB transaction", err)
		return
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to rollback DB transaction", err)
		}
	}()

	q := db.New()

	people, err := q.AddNewPersonQuery(ctx, tx, db.AddNewPersonQueryParams{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Profession:  pkg.ToPgText(req.Profession),
		Email:       pkg.ToPgText(req.Email),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to add new person", err)
		return
	}

	person_to_event_mapping, err := q.MapPersonToEventQuery(ctx, tx, db.MapPersonToEventQueryParams{
		PersonID: req.PersonID,
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
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to commit DB transaction", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":                 "Person added successfully",
		"person_details":          people,
		"person_to_event_mapping": person_to_event_mapping,
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

	tx, err := cmd.DBPool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to begin DB transaction", err)
		return
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to rollback DB transaction", err)
		}
	}()

	q := db.New()

	updatedPerson, err := q.UpdatePersonDetailsQuery(ctx, tx, db.UpdatePersonDetailsQueryParams{
		ID:          personId,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Profession:  pkg.ToPgText(req.Profession),
		Email:       pkg.ToPgText(req.Email),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to update person details", err)
		return
	} else if updatedPerson.ID == uuid.Nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Person not found",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Person not found for update", nil)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to commit DB transaction", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Person details updated successfully",
		"updated_person": updatedPerson,
	})
	pkg.Log.SuccessCtx(c)
}

func DeletePerson(c *gin.Context) {
	id := c.Param("id")
	personId, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Request not processed due to invalid parameters",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Invalid person ID parameter", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to begin DB transaction", err)
		return
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to rollback DB transaction", err)
		}
	}()

	q := db.New()

	deletedPerson, err := q.DeletePersonQuery(ctx, tx, personId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to delete person", err)
		return
	} else if len(deletedPerson) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Person not found",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Person not found for deletion", nil)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to commit DB transaction", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Person deleted successfully",
		"deleted_person": deletedPerson,
	})
	pkg.Log.SuccessCtx(c)
}
