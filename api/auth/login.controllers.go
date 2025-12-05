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
	"github.com/jackc/pgx/v5/pgtype"
)

func LoginUserCsrf(c *gin.Context) {
	csrfToken, err := pkg.CreateCsrfToken("login.user@amrita.edu", c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to generate CSRF token", err)
		return
	}

	pkg.SetCsrfCookie(c, csrfToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "User login action initiated successfully",
		"key":     csrfToken,
	})
	pkg.Log.SuccessCtx(c)
}

func LoginUser(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.LoginRequest](c)
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
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to initiate DB transaction", err)
		return
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != pgx.ErrTxClosed {
			pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to rollback", rbErr)
		}
	}()

	q := db.New()

	result, err := q.LoginUserQuery(ctx, tx, req.Email)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Invalid credentials",
		})
		pkg.Log.WarnCtx(c, "[AUTH-WARN]: Email does not exist")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to fetch login details", err)
		return
	}

	// Password verification
	err = pkg.CompareHash(result.Password, req.HashedPassword)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Invalid credentials",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Password verification failed", err)
		return
	}

	// If there is no refreshToken then create one, add it to the database and
	// proceed to mint auth token, set the cookies and send back the response
	if result.RefreshToken.String == "" {
		refreshToken, err := pkg.CreateRefreshToken(result.ID.String(), result.Email, pkg.Roles{
			IsUser:      true,
			IsOrganizer: false,
			IsAdmin:     false,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later.",
			})
			pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to create refresh token", err)
			return
		}

		token, err := q.UpdateRefreshTokenQuery(ctx, tx, db.UpdateRefreshTokenQueryParams{
			Email: req.Email,
			RefreshToken: pgtype.Text{
				String: refreshToken,
				Valid:  true,
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later.",
			})
			pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Could not add refresh token to DB", err)
			return
		}
		pkg.SetRefreshCookie(c, token.String)
	} else {
		pkg.SetRefreshCookie(c, result.RefreshToken.String)
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to commit transaction", err)
		return
	}

	authToken, err := pkg.CreateAuthToken(result.ID.String(), result.Email, pkg.Roles{
		IsUser:      true,
		IsOrganizer: false,
		IsAdmin:     false,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to create auth token", err)
		return
	}

	pkg.SetAuthCookie(c, authToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
		"name":    result.Name,
		"email":   result.Email,
	})
	pkg.Log.SuccessCtx(c)
}

func LoginOrganizerCsrf(c *gin.Context) {
	csrfToken, err := pkg.CreateCsrfToken("login.organizer@amrita.edu", c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		return
	}

	pkg.SetCsrfCookie(c, csrfToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Organizer login action initiated successfully",
		"key":     csrfToken,
	})
	pkg.Log.SuccessCtx(c)
}

func LoginOrganizer(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.LoginRequest](c)
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

		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to initiate DB transaction", err)
		return
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != pgx.ErrTxClosed {
			pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to rollback", rbErr)
		}
	}()

	q := db.New()

	result, err := q.LoginOrganizerQuery(ctx, tx, req.Email)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No organizer with given credentials exist",
		})

		pkg.Log.WarnCtx(c, "[AUTH-WARN]: No organizer with given email")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Could not login organizer", err)
		return
	}

	err = pkg.CompareHash(result.Password, req.HashedPassword)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No organizer with given credentials exist",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Password hash does not match", err)
	}

	// Check if there is a refreshToken already. If it exists, just mint an
	// auth token, set the cookies and return. Otherwise, mint both and return.
	if result.RefreshToken.String == "" {

		refreshToken, err := pkg.CreateRefreshToken(result.ID.String(), req.Email, pkg.Roles{
			IsUser:      false,
			IsOrganizer: true,
			IsAdmin:     false,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later.",
			})

			pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: ", err)
			return
		}

		token, err := q.UpdateOrganizerRefreshTokenQuery(ctx, tx,
			db.UpdateOrganizerRefreshTokenQueryParams{
				Email: req.Email,
				RefreshToken: pgtype.Text{
					String: refreshToken,
					Valid:  true,
				},
			})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later.",
			})
			pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Could not add refresh token to DB", err)
			return
		}
		pkg.SetRefreshCookie(c, token.String)
	} else {
		pkg.SetRefreshCookie(c, result.RefreshToken.String)
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to commit transaction", err)
		return
	}

	authToken, err := pkg.CreateAuthToken(result.ID.String(), req.Email, pkg.Roles{
		IsUser:      false,
		IsOrganizer: true,
		IsAdmin:     false,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to create auth token", err)
		return
	}

	pkg.SetAuthCookie(c, authToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Organizer logged in successfully",
	})
	pkg.Log.SuccessCtx(c)
}
