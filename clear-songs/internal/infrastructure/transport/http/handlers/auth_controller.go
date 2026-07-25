package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/RubenPari/clear-songs/internal/application/auth"
	"github.com/RubenPari/clear-songs/internal/infrastructure/config"
	"github.com/RubenPari/clear-songs/internal/infrastructure/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// AuthController handles authentication-related HTTP endpoints.
type AuthController struct {
	BaseController
	loginUC    *auth.LoginUseCase
	callbackUC *auth.CallbackUseCase
	logoutUC   *auth.LogoutUseCase
	isAuthUC   *auth.IsAuthUseCase
}

// NewAuthController creates an auth controller with the given use cases.
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

// Login initiates the OAuth flow by generating a state token and redirecting to Spotify.
func (ac *AuthController) Login(c *gin.Context) {
	state := uuid.NewString()
	ac.setOAuthStateCookie(c, state)

	url := ac.loginUC.Execute(state)
	c.Redirect(http.StatusFound, url)
}

// Callback handles the OAuth callback from Spotify. It validates the state parameter,
// exchanges the authorization code for an access token, and redirects to the frontend.
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
	if err := ac.callbackUC.Execute(ctx, code); err != nil {
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

	redirectURL := config.GetFrontendURL() + "/callback"
	c.Redirect(http.StatusFound, redirectURL)
}

// Logout clears the user's OAuth token and ends the session.
func (ac *AuthController) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	if err := ac.logoutUC.Execute(ctx); err != nil {
		ac.JSONInternalError(c, "Error logging out")
		return
	}

	ac.JSONSuccess(c, gin.H{"message": "User logged out successfully"})
}

// IsAuth checks whether the user is currently authenticated and returns user info.
func (ac *AuthController) IsAuth(c *gin.Context) {
	ctx := c.Request.Context()
	userInfo, err := ac.isAuthUC.Execute(ctx)
	if err != nil {
		ac.JSONUnauthorized(c)
		return
	}

	ac.JSONSuccess(c, gin.H{
		"user": gin.H{
			"spotify_id":    userInfo.ID,
			"display_name":  userInfo.DisplayName,
			"email":         userInfo.Email,
			"profile_image": userInfo.ProfileImage,
		},
	})
}

// requestIsHTTPS reports whether the request was made over HTTPS.
func requestIsHTTPS(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}
	return strings.EqualFold(c.Request.Header.Get("X-Forwarded-Proto"), "https")
}

// setOAuthStateCookie stores the OAuth state token in a secure, HTTP-only cookie.
func (ac *AuthController) setOAuthStateCookie(c *gin.Context, state string) {
	secure := requestIsHTTPS(c)
	c.SetCookie("oauth_state", state, 10*60, "/", "", secure, true)
}

// clearOAuthStateCookie removes the OAuth state cookie.
func (ac *AuthController) clearOAuthStateCookie(c *gin.Context) {
	secure := requestIsHTTPS(c)
	c.SetCookie("oauth_state", "", -1, "/", "", secure, true)
}
