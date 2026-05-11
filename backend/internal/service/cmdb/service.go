// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package cmdb provides configuration-item (CMDB) persistence and topology graph building.
package cmdb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kkops/backend/internal/integration/provider"
	"github.com/kkops/backend/internal/model"
	"gorm.io/gorm"
)

// Service implements CMDB CRUD and topology graph assembly.
type Service struct {
	db *gorm.DB
}

// NewService constructs a CMDB service.
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateAssetInput is the payload for creating a CMDB asset.
type CreateAssetInput struct {
	Name          string `json:"name" binding:"required"`
	Kind          string `json:"kind" binding:"required"`
	Env           string `json:"env"`
	OwnerUserID   *uint  `json:"owner_user_id"`
	IntegrationID *uint  `json:"integration_id"`
	ExternalRef   string `json:"external_ref"`
	Notes         string `json:"notes"`
}

// UpdateAssetInput updates mutable CMDB fields.
type UpdateAssetInput struct {
	Name          *string `json:"name"`
	Kind          *string `json:"kind"`
	Env           *string `json:"env"`
	OwnerUserID   *uint   `json:"owner_user_id"`
	LabelsJSON    *string `json:"labels"`
	IntegrationID *uint   `json:"integration_id"`
	ExternalRef   *string `json:"external_ref"`
	Notes         *string `json:"notes"`
}

// CreateRelationInput binds relation creation.
type CreateRelationInput struct {
	FromAssetID  uint   `json:"from_asset_id" binding:"required"`
	ToAssetID    uint   `json:"to_asset_id" binding:"required"`
	RelationType string `json:"relation_type" binding:"required"`
	Meta         any    `json:"meta"`
}

// GraphNode is one vertex for topology responses.
type GraphNode struct {
	ID        uint           `json:"id"`
	Kind      string         `json:"kind"` // "cmdb_asset" | "derived"
	Name      string         `json:"name"`
	Subtype   string         `json:"subtype,omitempty"`
	Env       string         `json:"env,omitempty"`
	ExtraMeta map[string]any `json:"extra,omitempty"`
}

// GraphEdge is one dependency edge.
type GraphEdge struct {
	ID           uint   `json:"id,omitempty"`
	FromID       uint   `json:"from"`
	ToID         uint   `json:"to"`
	RelationType string `json:"relation_type"`
}

// TopologyGraph is returned by GET /topology/graph.
type TopologyGraph struct {
	Nodes            []GraphNode    `json:"nodes"`
	Edges            []GraphEdge    `json:"edges"`
	DerivedCounts    map[string]int `json:"derived_counts"`
	PlaceholderHints map[string]any `json:"placeholder_hints,omitempty"`
}

var allowedKinds = map[string]struct{}{
	model.CMDBKindService: {}, model.CMDBKindHost: {}, model.CMDBKindDatabase: {},
	model.CMDBKindCluster: {}, model.CMDBKindOther: {},
}

var allowedRelations = map[string]struct{}{
	model.RelationDependsOn: {}, model.RelationRunsOn: {}, model.RelationCalls: {},
}

func validateKind(k string) error {
	if _, ok := allowedKinds[strings.TrimSpace(k)]; !ok {
		return fmt.Errorf("invalid kind")
	}
	return nil
}

func validateRelationType(t string) error {
	if _, ok := allowedRelations[strings.TrimSpace(t)]; !ok {
		return fmt.Errorf("invalid relation_type")
	}
	return nil
}

// ListAssets returns CMDB assets with optional filters.
func (s *Service) ListAssets(ctx context.Context, kind, env, q string, limit, offset int) ([]model.CMDBAsset, int64, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	q = strings.TrimSpace(q)
	var rows []model.CMDBAsset
	tx := s.db.WithContext(ctx).Model(&model.CMDBAsset{})
	if kind != "" {
		tx = tx.Where("kind = ?", kind)
	}
	if env != "" {
		tx = tx.Where("env = ?", env)
	}
	if q != "" {
		pattern := "%" + q + "%"
		tx = tx.Where("name ILIKE ? OR notes ILIKE ?", pattern, pattern)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := tx.Preload("Owner").Preload("Integration").
		Order("id DESC").Limit(limit).Offset(offset).Find(&rows).Error
	return rows, total, err
}

// GetAsset loads one CMDB asset.
func (s *Service) GetAsset(ctx context.Context, id uint) (*model.CMDBAsset, error) {
	var row model.CMDBAsset
	err := s.db.WithContext(ctx).Preload("Owner").Preload("Integration").First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// CreateAsset persists a new CMDB asset.
func (s *Service) CreateAsset(ctx context.Context, in CreateAssetInput, labelsJSON string) (*model.CMDBAsset, error) {
	if err := validateKind(in.Kind); err != nil {
		return nil, err
	}
	row := model.CMDBAsset{
		Name:          strings.TrimSpace(in.Name),
		Kind:          strings.TrimSpace(in.Kind),
		Env:           strings.TrimSpace(in.Env),
		OwnerUserID:   in.OwnerUserID,
		LabelsJSON:    labelsJSON,
		IntegrationID: in.IntegrationID,
		ExternalRef:   strings.TrimSpace(in.ExternalRef),
		Notes:         in.Notes,
	}
	if row.LabelsJSON == "" {
		row.LabelsJSON = "{}"
	}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// UpdateAsset applies partial updates.
func (s *Service) UpdateAsset(ctx context.Context, id uint, in UpdateAssetInput) (*model.CMDBAsset, error) {
	var row model.CMDBAsset
	if err := s.db.WithContext(ctx).First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if in.Name != nil {
		row.Name = strings.TrimSpace(*in.Name)
	}
	if in.Kind != nil {
		if err := validateKind(*in.Kind); err != nil {
			return nil, err
		}
		row.Kind = strings.TrimSpace(*in.Kind)
	}
	if in.Env != nil {
		row.Env = strings.TrimSpace(*in.Env)
	}
	if in.OwnerUserID != nil {
		row.OwnerUserID = in.OwnerUserID
	}
	if in.LabelsJSON != nil {
		row.LabelsJSON = *in.LabelsJSON
		if row.LabelsJSON == "" {
			row.LabelsJSON = "{}"
		}
	}
	if in.IntegrationID != nil {
		row.IntegrationID = in.IntegrationID
	}
	if in.ExternalRef != nil {
		row.ExternalRef = strings.TrimSpace(*in.ExternalRef)
	}
	if in.Notes != nil {
		row.Notes = *in.Notes
	}
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// DeleteAsset soft-deletes a CMDB asset.
func (s *Service) DeleteAsset(ctx context.Context, id uint) error {
	res := s.db.WithContext(ctx).Delete(&model.CMDBAsset{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ListRelations lists edges with optional filters.
func (s *Service) ListRelations(ctx context.Context, fromID, toID uint, limit int) ([]model.AssetRelation, int64, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	tx := s.db.WithContext(ctx).Model(&model.AssetRelation{})
	if fromID > 0 {
		tx = tx.Where("from_asset_id = ?", fromID)
	}
	if toID > 0 {
		tx = tx.Where("to_asset_id = ?", toID)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.AssetRelation
	err := tx.Preload("FromAsset").Preload("ToAsset").Order("id DESC").Limit(limit).Find(&rows).Error
	return rows, total, err
}

// CreateRelation stores a new edge.
func (s *Service) CreateRelation(ctx context.Context, in CreateRelationInput, metaJSON string) (*model.AssetRelation, error) {
	if in.FromAssetID == in.ToAssetID {
		return nil, fmt.Errorf("self-reference not allowed")
	}
	if err := validateRelationType(in.RelationType); err != nil {
		return nil, err
	}
	var fromCount, toCount int64
	s.db.WithContext(ctx).Model(&model.CMDBAsset{}).Where("id = ?", in.FromAssetID).Count(&fromCount)
	s.db.WithContext(ctx).Model(&model.CMDBAsset{}).Where("id = ?", in.ToAssetID).Count(&toCount)
	if fromCount == 0 || toCount == 0 {
		return nil, fmt.Errorf("from_asset_id or to_asset_id not found")
	}
	row := model.AssetRelation{
		FromAssetID:  in.FromAssetID,
		ToAssetID:    in.ToAssetID,
		RelationType: strings.TrimSpace(in.RelationType),
		MetaJSON:     metaJSON,
	}
	if row.MetaJSON == "" {
		row.MetaJSON = "{}"
	}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// DeleteRelation removes an edge by id.
func (s *Service) DeleteRelation(ctx context.Context, id uint) error {
	res := s.db.WithContext(ctx).Delete(&model.AssetRelation{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// BuildTopologyGraph aggregates CMDB nodes, relation edges, and derived integration hints.
func (s *Service) BuildTopologyGraph(ctx context.Context) (*TopologyGraph, error) {
	var assets []model.CMDBAsset
	if err := s.db.WithContext(ctx).Find(&assets).Error; err != nil {
		return nil, err
	}
	nodes := make([]GraphNode, 0, len(assets)+8)
	for _, a := range assets {
		nodes = append(nodes, GraphNode{
			ID:        a.ID,
			Kind:      "cmdb_asset",
			Name:      a.Name,
			Subtype:   a.Kind,
			Env:       a.Env,
			ExtraMeta: map[string]any{"integration_id": a.IntegrationID, "external_ref": a.ExternalRef},
		})
	}

	var rels []model.AssetRelation
	if err := s.db.WithContext(ctx).Find(&rels).Error; err != nil {
		return nil, err
	}
	edges := make([]GraphEdge, 0, len(rels))
	for _, r := range rels {
		edges = append(edges, GraphEdge{
			ID:           r.ID,
			FromID:       r.FromAssetID,
			ToID:         r.ToAssetID,
			RelationType: r.RelationType,
		})
	}

	// Derived integration placeholders (counts only; no graph DB).
	var k8sCount, harborCount, totalInt int64
	s.db.WithContext(ctx).Model(&model.Integration{}).Count(&totalInt)
	s.db.WithContext(ctx).Model(&model.Integration{}).Where("kind = ? AND enabled = ?", provider.KindKubernetes, true).Count(&k8sCount)
	s.db.WithContext(ctx).Model(&model.Integration{}).Where("kind = ? AND enabled = ?", provider.KindHarbor, true).Count(&harborCount)

	derivedID := uint(1_000_000)
	if k8sCount > 0 {
		nodes = append(nodes, GraphNode{
			ID:      derivedID,
			Kind:    "derived",
			Name:    fmt.Sprintf("Kubernetes integrations (%d)", k8sCount),
			Subtype: "kubernetes_placeholder",
		})
		derivedID++
	}
	if harborCount > 0 {
		nodes = append(nodes, GraphNode{
			ID:      derivedID,
			Kind:    "derived",
			Name:    fmt.Sprintf("Harbor integrations (%d)", harborCount),
			Subtype: "harbor_placeholder",
		})
	}

	return &TopologyGraph{
		Nodes: nodes,
		Edges: edges,
		DerivedCounts: map[string]int{
			"kubernetes_integrations": int(k8sCount),
			"harbor_integrations":     int(harborCount),
			"integrations_total":      int(totalInt),
		},
		PlaceholderHints: map[string]any{
			"harbor_repo_count": "not queried in MVP; connect Harbor integration to browse repositories",
			"note":              "Derived nodes are informational placeholders; CMDB edges use asset_relations only.",
		},
	}, nil
}
