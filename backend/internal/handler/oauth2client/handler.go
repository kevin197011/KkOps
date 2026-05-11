// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package oauth2client

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
)

// Handler handles OAuth2 client (IdP application) CRUD
type Handler struct {
	db *gorm.DB
}

// NewHandler creates the handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// List returns all OAuth2 clients (client_secret never returned)
func (h *Handler) List(c *gin.Context) {
	protocol := strings.TrimSpace(c.Query("protocol"))
	q := h.db.Model(&model.OAuth2Client{})
	if protocol != "" {
		q = q.Where("protocol = ?", protocol)
	}
	var list []model.OAuth2Client
	if err := q.Order("id ASC").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := make([]oauth2ClientResponse, 0, len(list))
	for _, c := range list {
		out = append(out, toResponse(&c))
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// Get returns one OAuth2 client by id (client_secret never returned)
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var client model.OAuth2Client
	if err := h.db.First(&client, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toResponse(&client)})
}

type createRequest struct {
	Name         string   `json:"name" binding:"required"`
	Protocol     string   `json:"protocol"`
	RedirectURIs []string `json:"redirect_uris" binding:"required"`
	Scopes       string   `json:"scopes"`
}

// Create creates an OAuth2 client; returns client with client_secret only in this response
func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clientID, err := randomHex(16)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate client_id"})
		return
	}
	clientSecret, err := randomHex(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate client_secret"})
		return
	}
	redirectURIsJSON, _ := json.Marshal(req.RedirectURIs)
	scopes := strings.TrimSpace(req.Scopes)
	if scopes == "" {
		scopes = "openid profile email"
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash secret"})
		return
	}
	proto := strings.TrimSpace(req.Protocol)
	if proto == "" {
		proto = "oidc"
	}
	client := model.OAuth2Client{
		ClientID:     clientID,
		ClientSecret: string(hash),
		Name:         req.Name,
		Protocol:     proto,
		RedirectURIs: string(redirectURIsJSON),
		Scopes:       scopes,
		Enabled:      true,
	}
	if err := h.db.Create(&client).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := toResponse(&client)
	resp.ClientSecret = clientSecret
	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

type updateRequest struct {
	Name         *string  `json:"name"`
	Protocol     *string  `json:"protocol"`
	RedirectURIs []string `json:"redirect_uris"`
	Scopes       *string  `json:"scopes"`
	Enabled      *bool    `json:"enabled"`
}

// Update updates an OAuth2 client (client_secret never changed here; use RegenerateSecret)
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var client model.OAuth2Client
	if err := h.db.First(&client, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if req.Name != nil {
		client.Name = *req.Name
	}
	if req.Protocol != nil && strings.TrimSpace(*req.Protocol) != "" {
		client.Protocol = strings.TrimSpace(*req.Protocol)
	}
	if req.RedirectURIs != nil {
		redirectURIsJSON, _ := json.Marshal(req.RedirectURIs)
		client.RedirectURIs = string(redirectURIsJSON)
	}
	if req.Scopes != nil {
		client.Scopes = strings.TrimSpace(*req.Scopes)
	}
	if req.Enabled != nil {
		client.Enabled = *req.Enabled
	}
	if err := h.db.Save(&client).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toResponse(&client)})
}

// Delete deletes an OAuth2 client
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.db.Delete(&model.OAuth2Client{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// RegenerateSecret generates a new client_secret and returns it once (plain)
func (h *Handler) RegenerateSecret(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	clientSecret, err := randomHex(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate secret"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash secret"})
		return
	}
	var client model.OAuth2Client
	if err := h.db.First(&client, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	client.ClientSecret = string(hash)
	if err := h.db.Save(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"client_secret": clientSecret}})
}

type oauth2ClientResponse struct {
	ID           uint     `json:"id"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret,omitempty"` // only on create / regenerate
	Name         string   `json:"name"`
	Protocol     string   `json:"protocol"`
	RedirectURIs []string `json:"redirect_uris"`
	Scopes       string   `json:"scopes"`
	Enabled      bool     `json:"enabled"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

func toResponse(c *model.OAuth2Client) oauth2ClientResponse {
	var uris []string
	_ = json.Unmarshal([]byte(c.RedirectURIs), &uris)
	if uris == nil {
		uris = []string{}
	}
	proto := c.Protocol
	if proto == "" {
		proto = "oidc"
	}
	return oauth2ClientResponse{
		ID:           c.ID,
		ClientID:     c.ClientID,
		Name:         c.Name,
		Protocol:     proto,
		RedirectURIs: uris,
		Scopes:       c.Scopes,
		Enabled:      c.Enabled,
		CreatedAt:    c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
