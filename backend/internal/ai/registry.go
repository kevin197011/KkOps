// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package ai

import (
	"sync"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
	integrationsvc "github.com/kkops/backend/internal/service/integration"
)

// Registry resolves Providers by integration id for kind "ai".
type Registry struct {
	svc   *integrationsvc.Service
	db    *gorm.DB
	mu    sync.RWMutex
	cache map[uint]Provider
}

// NewRegistry constructs an AI provider registry.
func NewRegistry(db *gorm.DB, svc *integrationsvc.Service) *Registry {
	return &Registry{
		svc:   svc,
		db:    db,
		cache: map[uint]Provider{},
	}
}

// Invalidate drops cached provider for an integration id (optional hook after updates).
func (r *Registry) Invalidate(id uint) {
	r.mu.Lock()
	delete(r.cache, id)
	r.mu.Unlock()
}

// ProviderForIntegration returns a provider for the given integration row id (kind must be ai).
func (r *Registry) ProviderForIntegration(id uint) (Provider, error) {
	r.mu.RLock()
	if p, ok := r.cache[id]; ok {
		r.mu.RUnlock()
		return p, nil
	}
	r.mu.RUnlock()

	raw, err := r.svc.DecryptConfigForWorker(id)
	if err != nil {
		return nil, err
	}
	p, err := NewProviderFromConfig(raw)
	if err != nil {
		return nil, err
	}
	r.mu.Lock()
	r.cache[id] = p
	r.mu.Unlock()
	return p, nil
}

// DefaultProvider returns the first enabled ai integration's provider and its id.
func (r *Registry) DefaultProvider() (Provider, uint, error) {
	var rows []model.Integration
	if err := r.db.Where("kind = ? AND enabled = ?", "ai", true).Order("id asc").Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, 0, ErrNoAIIntegration
	}
	id := rows[0].ID
	p, err := r.ProviderForIntegration(id)
	if err != nil {
		return nil, 0, err
	}
	return p, id, nil
}

// ListAIIntegrationIDs returns ids of integrations with kind ai.
func (r *Registry) ListAIIntegrationIDs() ([]uint, error) {
	var rows []model.Integration
	if err := r.db.Where("kind = ?", "ai").Order("id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]uint, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.ID)
	}
	return out, nil
}
