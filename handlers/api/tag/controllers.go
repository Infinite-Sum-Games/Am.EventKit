package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func FetchEventTags(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "TAG") {
		return
	}
	defer conn.Release()

	q := db.New()

	tags, err := q.ListTagsQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})

		pkg.Log.ErrorCtx(c, "[TAG-ERROR]: Failed to fetch tags", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tags list fetched successfully",
		"tags":    tags,
	})
	pkg.Log.SuccessCtx(c)
}

func CreateEventTag(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, ok := pkg.ValidateRequest[models.CreateTagRequest](c)
	if !ok {
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "TAG") {
		return
	}
	defer conn.Release()

	q := db.New()
	err = q.CreateTagQuery(ctx, conn, db.CreateTagQueryParams{
		Name:         req.Name,
		Abbreviation: req.Abbreviation,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[TAG-ERROR]: Failed to create tag", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tag created successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func EditEventTag(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tagID, ok := pkg.GrabUuid(c, c.Param("tagId"), "TAG", "tag")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.CreateTagRequest](c)
	if !ok {
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "TAG") {
		return
	}
	defer conn.Release()

	q := db.New()
	rows, err := q.UpdateTagByIDQuery(ctx, conn, db.UpdateTagByIDQueryParams{
		ID:           tagID,
		Name:         req.Name,
		Abbreviation: req.Abbreviation,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[TAG-ERROR]: Failed to update tag", err)
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Tag does not exist",
		})
		pkg.Log.ErrorCtx(c, "[TAG-ERROR]: Tag does not exist for update", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tag updated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func DeleteEventTag(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tagID, ok := pkg.GrabUuid(c, c.Param("tagId"), "TAG", "tag")
	if !ok {
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "TAG") {
		return
	}
	defer conn.Release()

	q := db.New()
	rows, err := q.DeleteTagByIDQuery(ctx, conn, tagID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[TAG-ERROR]: Failed to delete tag", err)
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Tag does not exist",
		})
		pkg.Log.ErrorCtx(c, "[TAG-ERROR]: Tag does not exist for delete", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tag deleted successfully",
	})
	pkg.Log.SuccessCtx(c)
}
