package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/RubenPari/clear-songs/internal/application/auth"
	"github.com/RubenPari/clear-songs/internal/infrastructure/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

// Creates auth controller.
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

// Starts.
func (ac *AuthController) Login(c *gin.Context) {
	state := uuid.NewString()
	ac.setOAuthStateCookie(c, state)

	url := ac.loginUC.Execute(state)
	c.Redirect(http.StatusFound, url)
}

// Callback.
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
			logging.LoggerFromGinContext(c).Error("OAuth callback failed",
				zap.String("error_code", re.ErrorCode),
				zap.String("error_description", re.ErrorDescription),
			)
		} else {
			logging.LoggerFromGinContext(c).Error("OAuth callback failed", zap.Error(err))
		}
		ac.JSONInternalError(c, "Error authenticating user")
		return
	}
	c.Redirect(http.StatusFound, redirectURL)
}

// Logs out.
func (ac *AuthController) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	if err := ac.logoutUC.Execute(ctx); err != nil {
		ac.JSONInternalError(c, "Error logging out")
		return
	}

	ac.JSONSuccess(c, gin.H{"message": "User logged out successfully"})
}

// Checks whether auth.
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

// Request is https.
func requestIsHTTPS(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}
	return strings.EqualFold(c.Request.Header.Get("X-Forwarded-Proto"), "https")
}

// Sets oauth state cookie.
func (ac *AuthController) setOAuthStateCookie(c *gin.Context, state string) {
	secure := requestIsHTTPS(c)
	c.SetCookie("oauth_state", state, 10*60, "/", "", secure, true)
}

// Clears oauth state cookie.
func (ac *AuthController) clearOAuthStateCookie(c *gin.Context) {
	secure := requestIsHTTPS(c)
	c.SetCookie("oauth_state", "", -1, "/", "", secure, true)
}
