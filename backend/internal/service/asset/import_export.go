// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package asset

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/kkops/backend/internal/model"
)

// ExportAssets exports assets to CSV format
func (s *Service) ExportAssets(w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header (readable names for Project, Cloud Platform, Environment, SSH Key)
	header := []string{
		"ID", "Hostname", "Project",
		"Cloud Platform", "Environment", "IP", "SSH Port",
		"SSH Key", "SSH User", "CPU", "Memory", "Disk",
		"Status", "Description",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Fetch all assets with related Project, CloudPlatform, Environment, SSHKey
	var assets []model.Asset
	if err := s.db.Preload("Project").Preload("Environment").Preload("CloudPlatform").Preload("SSHKey").Find(&assets).Error; err != nil {
		return err
	}

	// Write data (use names for Project / Cloud Platform / Environment / SSH Key; empty when nil)
	for _, asset := range assets {
		projectName := ""
		if asset.Project != nil {
			projectName = asset.Project.Name
		}
		cloudPlatformName := ""
		if asset.CloudPlatform != nil {
			cloudPlatformName = asset.CloudPlatform.Name
		}
		environmentName := ""
		if asset.Environment != nil {
			environmentName = asset.Environment.Name
		}
		sshKeyName := ""
		if asset.SSHKey != nil {
			sshKeyName = asset.SSHKey.Name
		}
		record := []string{
			strconv.FormatUint(uint64(asset.ID), 10),
			asset.HostName,
			projectName,
			cloudPlatformName,
			environmentName,
			asset.IP,
			strconv.Itoa(asset.SSHPort),
			sshKeyName,
			asset.SSHUser,
			asset.CPU,
			asset.Memory,
			asset.Disk,
			asset.Status,
			asset.Description,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// ImportAssets imports assets from CSV format
type ImportResult struct {
	Success int
	Failed  int
	Errors  []string
}

func (s *Service) ImportAssets(r io.Reader) (*ImportResult, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, errors.New("failed to read CSV header")
	}

	// Validate header (basic check)
	expectedHeader := []string{"Hostname"}
	headerMap := make(map[string]int)
	for i, col := range header {
		headerMap[col] = i
	}
	for _, reqCol := range expectedHeader {
		if _, exists := headerMap[reqCol]; !exists {
			return nil, errors.New("missing required column: " + reqCol)
		}
	}

	result := &ImportResult{
		Errors: make([]string, 0),
	}

	// Read and import records
	lineNum := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		lineNum++

		// Parse record
		req, parseErr := s.parseImportRecord(record, headerMap)
		if parseErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, "Line "+strconv.Itoa(lineNum)+": "+parseErr.Error())
			continue
		}

		// Create asset
		_, createErr := s.CreateAsset(req)
		if createErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, "Line "+strconv.Itoa(lineNum)+": "+createErr.Error())
			continue
		}

		result.Success++
	}

	return result, nil
}

// parseImportRecord parses a CSV record into CreateAssetRequest
func (s *Service) parseImportRecord(record []string, headerMap map[string]int) (*CreateAssetRequest, error) {
	getValue := func(col string) string {
		if idx, exists := headerMap[col]; exists && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	req := &CreateAssetRequest{}

	// Required fields
	hostname := getValue("Hostname")
	if hostname == "" {
		return nil, errors.New("hostname is required")
	}
	req.HostName = hostname

	// Optional fields: support both ID columns (backward compatible) and name columns (round-trip from export)
	if projectIDStr := getValue("Project ID"); projectIDStr != "" {
		if projectID, err := strconv.ParseUint(projectIDStr, 10, 32); err == nil {
			id := uint(projectID)
			req.ProjectID = &id
		}
	}
	if req.ProjectID == nil {
		if name := getValue("Project"); name != "" {
			var proj model.Project
			if err := s.db.Where("name = ?", name).First(&proj).Error; err == nil {
				req.ProjectID = &proj.ID
			}
		}
	}

	// Cloud Platform: support "Cloud Platform ID" and "Cloud Platform" (name)
	if cloudPlatformIDStr := getValue("Cloud Platform ID"); cloudPlatformIDStr != "" {
		if cloudPlatformID, err := strconv.ParseUint(cloudPlatformIDStr, 10, 32); err == nil {
			id := uint(cloudPlatformID)
			req.CloudPlatformID = &id
		}
	}
	if req.CloudPlatformID == nil {
		if name := getValue("Cloud Platform"); name != "" {
			var cp model.CloudPlatform
			if err := s.db.Where("name = ?", name).First(&cp).Error; err == nil {
				req.CloudPlatformID = &cp.ID
			}
		}
	}

	if envIDStr := getValue("Environment ID"); envIDStr != "" {
		if envID, err := strconv.ParseUint(envIDStr, 10, 32); err == nil {
			id := uint(envID)
			req.EnvironmentID = &id
		}
	}
	if req.EnvironmentID == nil {
		if name := getValue("Environment"); name != "" {
			var env model.Environment
			if err := s.db.Where("name = ?", name).First(&env).Error; err == nil {
				req.EnvironmentID = &env.ID
			}
		}
	}

	req.IP = getValue("IP")

	if sshPortStr := getValue("SSH Port"); sshPortStr != "" {
		if sshPort, err := strconv.Atoi(sshPortStr); err == nil {
			req.SSHPort = sshPort
		}
	}

	if sshKeyIDStr := getValue("SSH Key ID"); sshKeyIDStr != "" {
		if sshKeyID, err := strconv.ParseUint(sshKeyIDStr, 10, 32); err == nil {
			id := uint(sshKeyID)
			req.SSHKeyID = &id
		}
	}
	if req.SSHKeyID == nil {
		if name := getValue("SSH Key"); name != "" {
			var key model.SSHKey
			if err := s.db.Where("name = ?", name).First(&key).Error; err == nil {
				req.SSHKeyID = &key.ID
			}
		}
	}

	req.SSHUser = getValue("SSH User")
	req.CPU = getValue("CPU")
	req.Memory = getValue("Memory")
	req.Disk = getValue("Disk")
	req.Status = getValue("Status")
	req.Description = getValue("Description")

	return req, nil
}

// uintPtrToString converts uint pointer to string
func uintPtrToString(ptr *uint) string {
	if ptr == nil {
		return ""
	}
	return strconv.FormatUint(uint64(*ptr), 10)
}
