package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/RubenPari/clear-songs/internal/application/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

// AuthController is the auth controller using dependency injection
type AuthController struct {
	BaseController
	loginUC    *auth.LoginUseCase
	callbackUC *auth.CallbackUseCase
	logoutUC   *auth.LogoutUseCase
	isAuthUC   *auth.IsAuthUseCase
}

// NewAuthController creates a new AuthController
func NewAuthController(
	loginUC *auth.LoginUseCase,
	callbackUC *auth.CallbackUseCase,
	logoutUC *auth.LogoutUseCase,
	isAuthUC *auth.IsAuthUseCase,
) *AuthController {
	return &AuthController{
		loginUC:    loginUC,
		callbackUC: callbackUC,
		logoutUC:   logoutUC,
		isAuthUC:   isAuthUC,
	}
}

// Login handles GET /auth/login
func (ac *AuthController) Login(c *gin.Context) {
	state := uuid.NewString()
	ac.setOAuthStateCookie(c, state)

	url := ac.loginUC.Execute(state)
	c.Redirect(http.StatusFound, url)
}

// Callback handles GET /auth/callback
func (ac *AuthController) Callback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		ac.JSONValidationError(c, "Authorization code is required")
		return
	}
	state := c.Query("state")
	if cookieState, err := c.Cookie("oauth_state"); err == nil && cookieState != "" {
		if state == "" || state != cookieState {
			ac.JSONValidationError(c, "Invalid OAuth state")
			return
		}
	}
	ac.clearOAuthStateCookie(c)

	ctx := c.Request.Context()
	redirectURL, err := ac.callbackUC.Execute(ctx, code)
	if err != nil {
		var re *oauth2.RetrieveError
		if errors.As(err, &re) {
			log.Printf(
				"ERROR: OAuth callback failed [oauth2_exchange] error_code=%q error_description=%q",
				re.ErrorCode,
				re.ErrorDescription,
			)
		} else {
			log.Printf("ERROR: OAuth callback failed: %v", err)
		}
		ac.JSONInternalError(c, "Error authenticating user")
		return
	}
	c.Redirect(http.StatusFound, redirectURL)
}

// Logout handles GET /auth/logout
func (ac *AuthController) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	if err := ac.logoutUC.Execute(ctx); err != nil {
		ac.JSONInternalError(c, "Error logging out")
		return
	}

	ac.JSONSuccess(c, gin.H{"message": "User logged out successfully"})
}

// IsAuth handles GET /auth/is-auth
func (ac *AuthController) IsAuth(c *gin.Context) {
	ctx := c.Request.Context()
	userInfo, err := ac.isAuthUC.Execute(ctx)
	if err != nil {
		ac.JSONUnauthorized(c)
		return
	}

	ac.JSONSuccess(c, gin.H{
		"user": gin.H{
			"spotify_id":    userInfo.SpotifyID,
			"display_name":  userInfo.DisplayName,
			"email":         userInfo.Email,
			"profile_image": userInfo.ProfileImage,
		},
	})
}

func requestIsHTTPS(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}
	return strings.EqualFold(c.Request.Header.Get("X-Forwarded-Proto"), "https")
}

func (ac *AuthController) setOAuthStateCookie(c *gin.Context, state string) {
	secure := requestIsHTTPS(c)
	c.SetCookie("oauth_state", state, 10*60, "/", "", secure, true)
}

func (ac *AuthController) clearOAuthStateCookie(c *gin.Context) {
	secure := requestIsHTTPS(c)
	c.SetCookie("oauth_state", "", -1, "/", "", secure, true)
}
