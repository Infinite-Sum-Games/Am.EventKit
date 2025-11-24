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
	defer tx.Rollback(ctx)

	q := db.New()
	password, err := pkg.Hash(req.HashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to hash password", err)
		return
	}

	result, err := q.LoginUserQuery(ctx, tx, db.LoginUserQueryParams{
		Email:    req.Email,
		Password: password,
	})
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No user with given email and password",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Username and password do not exist", err)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to fetch login details", err)
		return
	}

	// If there is no refreshToken then create one, add it to the database and
	// proceed to mint auth token, set the cookies and send back the response
	if result.RefreshToken.String == "" {
		refreshToken, err := pkg.CreateRefreshToken(result.ID.String(), result.Email, true, false)
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

	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to commit transaction", err)
		return
	}

	authToken, err := pkg.CreateAuthToken(result.ID.String(), result.Email, true, false)
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
	// For organizers, email should be like - <department>@cb.amrita.edu
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
	defer tx.Rollback(ctx)

	q := db.New()

	password, err := pkg.Hash(req.HashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to hash password", err)
		return
	}

	result, err := q.LoginOrganizerQuery(ctx, tx, db.LoginOrganizerQueryParams{
		Email:    req.Email,
		Password: password,
	})
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No organizer with given email and password found",
		})

		pkg.Log.WarnCtx(c, "[AUTH-WARN]: Invalid email and password")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})

		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Could not login organizer", err)
		return
	}

	// Check if there is a refreshToken already. If it exists, just mint an
	// auth token, set the cookies and return. Otherwise, mint both and return.
	if result.RefreshToken.String == "" {

		refreshToken, err := pkg.CreateRefreshToken(result.ID.String(), req.Email, false, true)
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
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to commit transaction", err)
		return
	}

	authToken, err := pkg.CreateAuthToken(result.ID.String(), req.Email, false, true)
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
