// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package idp

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/kkops/backend/internal/config"
	"github.com/kkops/backend/internal/model"
	authService "github.com/kkops/backend/internal/service/auth"
)

const sessionUserIDKey = "idp_user_id"

// Handler handles OIDC IdP endpoints
type Handler struct {
	cfg     *config.Config
	db      *gorm.DB
	authSvc *authService.Service
}

// NewHandler creates IdP handler
func NewHandler(cfg *config.Config, db *gorm.DB, authSvc *authService.Service) *Handler {
	return &Handler{cfg: cfg, db: db, authSvc: authSvc}
}

func (h *Handler) issuer() string {
	iss := strings.TrimSuffix(h.cfg.IdP.Issuer, "/")
	if iss == "" {
		iss = "http://localhost:3000/oidc"
	}
	return iss
}

// Discovery returns OIDC discovery document
func (h *Handler) Discovery(c *gin.Context) {
	iss := h.issuer()
	doc := map[string]interface{}{
		"issuer":                                iss,
		"authorization_endpoint":                iss + "/authorize",
		"token_endpoint":                        iss + "/token",
		"userinfo_endpoint":                     iss + "/userinfo",
		"jwks_uri":                              iss + "/jwks",
		"scopes_supported":                      []string{"openid", "profile", "email"},
		"response_types_supported":              []string{"code"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
	}
	c.JSON(http.StatusOK, doc)
}

// JWKS returns public keys
func (h *Handler) JWKS(c *gin.Context) {
	raw, err := getJWKS()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "jwks failed"})
		return
	}
	c.Data(http.StatusOK, "application/json", raw)
}

// Authorize handles GET /oidc/authorize
func (h *Handler) Authorize(c *gin.Context) {
	if !h.cfg.IdP.Enabled {
		c.JSON(http.StatusNotFound, gin.H{"error": "IdP disabled"})
		return
	}
	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")
	responseType := c.Query("response_type")
	scope := c.Query("scope")
	state := c.Query("state")

	if clientID == "" || redirectURI == "" || responseType != "code" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	client, err := h.getClient(clientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_client"})
		return
	}
	if !clientSupportsOIDC(client) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_client_protocol"})
		return
	}
	if !client.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_disabled"})
		return
	}
	uris, _ := parseRedirectURIs(client.RedirectURIs)
	if !contains(uris, redirectURI) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_redirect_uri"})
		return
	}

	// Check session
	session := sessions.Default(c)
	uid := session.Get(sessionUserIDKey)
	if uid == nil {
		// Redirect to IdP login; then come back here
		then := c.Request.URL.String()
		loginURL := "/oidc/login?then=" + url.QueryEscape(then)
		c.Redirect(http.StatusFound, loginURL)
		return
	}
	userID, ok := uid.(int)
	if !ok {
		c.Redirect(http.StatusFound, "/oidc/login?then="+url.QueryEscape(c.Request.URL.String()))
		return
	}

	// Create auth code and redirect back to client
	code, err := randomCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	saveCode(code, &authCodeEntry{
		UserID:      uint(userID),
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		State:       state,
	})

	redir := redirectURI
	if strings.Contains(redir, "?") {
		redir += "&"
	} else {
		redir += "?"
	}
	redir += "code=" + url.QueryEscape(code) + "&state=" + url.QueryEscape(state)
	c.Redirect(http.StatusFound, redir)
}

// Token handles POST /oidc/token (application/x-www-form-urlencoded)
func (h *Handler) Token(c *gin.Context) {
	if !h.cfg.IdP.Enabled {
		c.JSON(http.StatusNotFound, gin.H{"error": "IdP disabled"})
		return
	}
	grantType := c.PostForm("grant_type")
	if grantType != "authorization_code" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported_grant_type"})
		return
	}
	code := c.PostForm("code")
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")
	redirectURI := c.PostForm("redirect_uri")

	if code == "" || clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	client, err := h.getClient(clientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client"})
		return
	}
	if !clientSupportsOIDC(client) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_client_protocol"})
		return
	}
	if !client.Enabled {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "client_disabled"})
		return
	}
	if !checkClientSecret(clientSecret, client.ClientSecret) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client"})
		return
	}
	entry, ok := consumeCode(code)
	if !ok || entry.ClientID != clientID || entry.RedirectURI != redirectURI {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant"})
		return
	}

	user, err := h.getUserByID(entry.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	key, err := getPrivateKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	iss := h.issuer()
	userIDStr := uintToString(entry.UserID)
	idToken, err := signIDToken(iss, clientID, userIDStr, user.Username, user.RealName, user.Email, key, time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	accessToken, err := signAccessToken(iss, clientID, userIDStr, entry.Scope, key, time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   3600,
		"id_token":     idToken,
		"scope":        entry.Scope,
	})
}

// UserInfo handles GET /oidc/userinfo
func (h *Handler) UserInfo(c *gin.Context) {
	if !h.cfg.IdP.Enabled {
		c.JSON(http.StatusNotFound, gin.H{"error": "IdP disabled"})
		return
	}
	auth := c.GetHeader("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
		return
	}
	tokenStr := strings.TrimPrefix(auth, "Bearer ")
	key, err := getPrivateKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	claims, err := verifyAccessToken(tokenStr, &key.PublicKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
		return
	}
	userID := stringToUint(claims.Subject)
	user, err := h.getUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
		return
	}
	body, _ := userInfoResponse(claims.Subject, user.Username, user.RealName, user.Email)
	c.Data(http.StatusOK, "application/json", body)
}

// LoginGet shows IdP login form
func (h *Handler) LoginGet(c *gin.Context) {
	then := c.Query("then")
	if then == "" {
		c.String(http.StatusBadRequest, "missing then")
		return
	}
	html := `<!DOCTYPE html><html><head><meta charset="utf-8"><title>KkOps 登录</title></head><body>
<form method="POST" action="/oidc/login">
<input type="hidden" name="then" value="` + then + `">
<label>用户名 <input name="username" required></label><br>
<label>密码 <input type="password" name="password" required></label><br>
<button type="submit">登录</button>
</form></body></html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// LoginPost validates credentials and sets session, then redirects to then
func (h *Handler) LoginPost(c *gin.Context) {
	then := c.PostForm("then")
	username := c.PostForm("username")
	password := c.PostForm("password")
	if then == "" || username == "" || password == "" {
		c.String(http.StatusBadRequest, "missing then/username/password")
		return
	}
	req := &authService.LoginRequest{Username: username, Password: password}
	_, err := h.authSvc.Login(req)
	if err != nil {
		c.String(http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	// Get user id from token or from auth response
	var user model.User
	if err := h.db.Where("username = ?", username).First(&user).Error; err != nil {
		c.String(http.StatusInternalServerError, "error")
		return
	}
	sess := sessions.Default(c)
	sess.Set(sessionUserIDKey, int(user.ID))
	if err := sess.Save(); err != nil {
		c.String(http.StatusInternalServerError, "error")
		return
	}
	c.Redirect(http.StatusFound, then)
}

func (h *Handler) getClient(clientID string) (*model.OAuth2Client, error) {
	var client model.OAuth2Client
	if err := h.db.Where("client_id = ?", clientID).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func clientSupportsOIDC(c *model.OAuth2Client) bool {
	return c.Protocol == "" || c.Protocol == "oidc"
}

type userInfo struct {
	Username string
	RealName string
	Email    string
}

func (h *Handler) getUserByID(id uint) (*userInfo, error) {
	var u model.User
	if err := h.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &userInfo{Username: u.Username, RealName: u.RealName, Email: u.Email}, nil
}

func checkClientSecret(plain, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}

func randomCode() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func contains(s []string, x string) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}

func uintToString(u uint) string { return fmt.Sprintf("%d", u) }
func stringToUint(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 64)
	return uint(v)
}
