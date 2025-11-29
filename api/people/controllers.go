package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
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
