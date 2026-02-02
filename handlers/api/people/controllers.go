package api

//
// import (
// 	"context"
// 	"net/http"
// 	"time"
//
// 	"github.com/Infinite-Sum-Games/Am.EventKit/cmd"
// 	db "github.com/Infinite-Sum-Games/Am.EventKit/db/gen"
// 	"github.com/Infinite-Sum-Games/Am.EventKit/models"
// 	"github.com/Infinite-Sum-Games/Am.EventKit/pkg"
// 	"github.com/gin-gonic/gin"
// 	"github.com/jackc/pgx/v5"
// )
//
// func FetchAllPeople(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
//
// 	conn, err := cmd.DBPool.Acquire(ctx)
// 	if pkg.HandleDbAcquireErr(c, err, "PEOPLE") {
// 		return
// 	}
// 	defer conn.Release()
//
// 	q := db.New()
// 	people, err := q.FetchAllPeopleQuery(ctx, conn)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Oops! Something happened. Please try again later",
// 		})
// 		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to fetch people", err)
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "People list fetched successfully",
// 		"people":  people,
// 	})
// 	pkg.Log.SuccessCtx(c)
// }
//
// func AddNewPerson(c *gin.Context) {
// 	req, ok := pkg.ValidateRequest[models.CreateNewPerson](c)
// 	if !ok {
// 		return
// 	}
//
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
//
// 	conn, err := cmd.DBPool.Acquire(ctx)
// 	if pkg.HandleDbAcquireErr(c, err, "PEOPLE") {
// 		return
// 	}
// 	defer conn.Release()
//
// 	q := db.New()
// 	people, err := q.AddNewPersonQuery(ctx, conn, db.AddNewPersonQueryParams{
// 		Name:        req.Name,
// 		PhoneNumber: req.PhoneNumber,
// 		Profession:  pkg.ToPgTextPtr(req.Profession),
// 		Email:       pkg.ToPgTextPtr(req.Email),
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Oops! Something happened. Please try again later",
// 		})
// 		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to add new person", err)
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"message":      "Person added successfully",
// 		"name":         people.Name,
// 		"phone_number": people.PhoneNumber,
// 		"profession":   people.Profession,
// 		"email":        people.Email,
// 	})
// 	pkg.Log.SuccessCtx(c)
// }
//
// func UpdatePersonDetails(c *gin.Context) {
// 	personId, ok := pkg.GrabUuid(c, c.Param("personId"), "PEOPLE", "person")
// 	if !ok {
// 		return
// 	}
//
// 	req, ok := pkg.ValidateRequest[models.UpdatePersonEventRequest](c)
// 	if !ok {
// 		return
// 	}
//
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
//
// 	conn, err := cmd.DBPool.Acquire(ctx)
// 	if pkg.HandleDbAcquireErr(c, err, "PEOPLE") {
// 		return
// 	}
// 	defer conn.Release()
//
// 	q := db.New()
// 	updatedPerson, err := q.UpdatePersonDetailsQuery(ctx, conn,
// 		db.UpdatePersonDetailsQueryParams{
// 			ID:          personId,
// 			Name:        req.Name,
// 			PhoneNumber: req.PhoneNumber,
// 			Profession:  pkg.ToPgTextPtr(req.Profession),
// 			Email:       pkg.ToPgTextPtr(req.Email),
// 		})
// 	if err == pgx.ErrNoRows {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"message": "Person not found",
// 		})
// 		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Person not found for update", err)
// 		return
// 	}
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Oops! Something happened. Please try again later",
// 		})
// 		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to update person details", err)
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"message":      "Person details updated successfully",
// 		"id":           updatedPerson.ID,
// 		"name":         updatedPerson.Name,
// 		"email":        updatedPerson.Email,
// 		"profession":   updatedPerson.Profession,
// 		"phone_number": updatedPerson.PhoneNumber,
// 	})
// 	pkg.Log.SuccessCtx(c)
// }
//
// func DeletePerson(c *gin.Context) {
// 	personId, ok := pkg.GrabUuid(c, c.Param("personId"), "PEOPLE", "person")
// 	if !ok {
// 		return
// 	}
//
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
//
// 	conn, err := cmd.DBPool.Acquire(ctx)
// 	if pkg.HandleDbAcquireErr(c, err, "PEOPLE") {
// 		return
// 	}
// 	defer conn.Release()
//
// 	q := db.New()
//
// 	deletedPersonId, err := q.DeletePersonQuery(ctx, conn, personId)
// 	if err == pgx.ErrNoRows {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"message": "Person not found",
// 		})
// 		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Person not found for deletion", nil)
// 		return
// 	}
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Oops! Something happened. Please try again later",
// 		})
// 		pkg.Log.ErrorCtx(c, "[PEOPLE-ERROR]: Failed to delete person", err)
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"message":           "Person deleted successfully",
// 		"deleted_person_id": deletedPersonId,
// 	})
// 	pkg.Log.SuccessCtx(c)
// }
