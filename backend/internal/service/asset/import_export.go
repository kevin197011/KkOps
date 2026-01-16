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

	// Write header
	header := []string{
		"ID", "Hostname", "Project ID",
		"Cloud Platform ID", "Environment ID", "IP", "SSH Port",
		"SSH Key ID", "SSH User", "CPU", "Memory", "Disk",
		"Status", "Description",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Fetch all assets
	var assets []model.Asset
	if err := s.db.Preload("Environment").Preload("CloudPlatform").Find(&assets).Error; err != nil {
		return err
	}

	// Write data
	for _, asset := range assets {
		cloudPlatformName := ""
		if asset.CloudPlatform != nil {
			cloudPlatformName = asset.CloudPlatform.Name
		}
		record := []string{
			strconv.FormatUint(uint64(asset.ID), 10),
			asset.HostName,
			uintPtrToString(asset.ProjectID),
			cloudPlatformName,
			uintPtrToString(asset.EnvironmentID),
			asset.IP,
			strconv.Itoa(asset.SSHPort),
			uintPtrToString(asset.SSHKeyID),
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

	// Optional fields
	if projectIDStr := getValue("Project ID"); projectIDStr != "" {
		if projectID, err := strconv.ParseUint(projectIDStr, 10, 32); err == nil {
			id := uint(projectID)
			req.ProjectID = &id
		}
	}

	// Handle Cloud Platform - try to find by name first, then by ID
	if cloudPlatformStr := getValue("Cloud Platform"); cloudPlatformStr != "" {
		// Try to parse as ID first
		if cloudPlatformID, err := strconv.ParseUint(cloudPlatformStr, 10, 32); err == nil {
			id := uint(cloudPlatformID)
			req.CloudPlatformID = &id
		} else {
			// If not a number, try to find cloud platform by name
			// Note: This requires database lookup, which should be handled in service layer
			// For now, we'll skip and let the service handle it
		}
	}
	// Also support "Cloud Platform ID" column for explicit ID input
	if cloudPlatformIDStr := getValue("Cloud Platform ID"); cloudPlatformIDStr != "" {
		if cloudPlatformID, err := strconv.ParseUint(cloudPlatformIDStr, 10, 32); err == nil {
			id := uint(cloudPlatformID)
			req.CloudPlatformID = &id
		}
	}
	if envIDStr := getValue("Environment ID"); envIDStr != "" {
		if envID, err := strconv.ParseUint(envIDStr, 10, 32); err == nil {
			id := uint(envID)
			req.EnvironmentID = &id
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
