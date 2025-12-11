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

func GetAllDependencies(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "DEPENDENCY") {
		return
	}
	defer conn.Release()

	q := db.New()

	deps, err := q.ListAllEventDependencies(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DEPENDENCY-ERROR]: Fetch failed", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Event dependencies fetched",
		"dependencies": deps,
	})
	pkg.Log.SuccessCtx(c)
}

func AddDependency(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, ok := pkg.ValidateRequest[models.EventDependencyRequest](c)
	if !ok {
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "DEPENDENCY") {
		return
	}
	defer conn.Release()

	q := db.New()

	err = q.AddEventDependency(ctx, conn, db.AddEventDependencyParams{
		StartEvent: req.StartEvent,
		EndEvent:   req.EndEvent,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DEPENDENCY-ERROR]: Add failed", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event dependency added",
	})
	pkg.Log.SuccessCtx(c)
}

func RemoveDependency(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, ok := pkg.ValidateRequest[models.EventDependencyRequest](c)
	if !ok {
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "DEPENDENCY") {
		return
	}
	defer conn.Release()

	q := db.New()

	rows, err := q.RemoveEventDependency(ctx, conn, db.RemoveEventDependencyParams{
		StartEvent: req.StartEvent,
		EndEvent:   req.EndEvent,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DEPENDENCY-ERROR]: Delete failed", err)
		return
	}

	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Dependency does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event dependency removed",
	})
	pkg.Log.SuccessCtx(c)
}
