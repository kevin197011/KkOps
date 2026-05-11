// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cicd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const cicdHTTPTimeout = 30 * time.Second

// JenkinsClient calls Jenkins HTTP API.
type JenkinsClient struct {
	cfg Config
}

// NewJenkinsClient constructs a Jenkins client from decrypted JSON config.
func NewJenkinsClient(cfgJSON []byte) (*JenkinsClient, error) {
	c, err := parseConfig(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &JenkinsClient{cfg: c}, nil
}

type jenkinsJobsResp struct {
	Jobs []struct {
		Name  string `json:"name"`
		URL   string `json:"url"`
		Color string `json:"color"`
	} `json:"jobs"`
}

func jenkinsBuildPath(jobName string) string {
	segs := strings.Split(strings.Trim(jobName, "/"), "/")
	var b strings.Builder
	for _, seg := range segs {
		if seg == "" {
			continue
		}
		b.WriteString("/job/")
		b.WriteString(url.PathEscape(seg))
	}
	return b.String()
}

// ListPipelines lists top-level jobs.
func (j *JenkinsClient) ListPipelines(ctx context.Context, _ string) ([]Pipeline, error) {
	u := j.cfg.BaseURL + "/api/json?tree=jobs[name,url,color]"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	if j.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+j.cfg.Token)
	}
	cli := &http.Client{Timeout: cicdHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("jenkins list: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("jenkins status %d", resp.StatusCode)
	}
	var wrap jenkinsJobsResp
	if err := json.Unmarshal(body, &wrap); err != nil {
		return nil, fmt.Errorf("decode jenkins jobs: %w", err)
	}
	out := make([]Pipeline, 0, len(wrap.Jobs))
	for _, job := range wrap.Jobs {
		out = append(out, Pipeline{
			ID:     "jenkins:" + job.Name,
			Name:   job.Name,
			Status: job.Color,
			WebURL: strings.TrimSuffix(job.URL, "/"),
		})
	}
	return out, nil
}

// GetPipeline returns metadata for a job name.
func (j *JenkinsClient) GetPipeline(ctx context.Context, jobName string) (*Pipeline, error) {
	list, err := j.ListPipelines(ctx, "")
	if err != nil {
		return nil, err
	}
	want := "jenkins:" + jobName
	for i := range list {
		if list[i].ID == want || list[i].Name == jobName {
			p := list[i]
			return &p, nil
		}
	}
	return nil, fmt.Errorf("jenkins job not found: %s", jobName)
}

// TriggerPipeline starts a build. jobName may contain slashes for folder jobs.
func (j *JenkinsClient) TriggerPipeline(ctx context.Context, jobName string, _ string, vars map[string]string) error {
	basePath := jenkinsBuildPath(jobName)
	suffix := "/build"
	if len(vars) > 0 {
		suffix = "/buildWithParameters"
	}
	buildURL := j.cfg.BaseURL + basePath + suffix
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, buildURL, nil)
	if err != nil {
		return err
	}
	if j.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+j.cfg.Token)
	}
	if len(vars) > 0 {
		q := url.Values{}
		for k, v := range vars {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	cli := &http.Client{Timeout: cicdHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return fmt.Errorf("jenkins trigger: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("jenkins trigger status %d: %s", resp.StatusCode, truncateStr(string(b), 120))
	}
	return nil
}

// GetLogs returns plain-text console output for the latest build.
func (j *JenkinsClient) GetLogs(ctx context.Context, jobName string) (string, error) {
	basePath := jenkinsBuildPath(jobName)
	u := j.cfg.BaseURL + basePath + "/lastBuild/consoleText"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	if j.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+j.cfg.Token)
	}
	cli := &http.Client{Timeout: cicdHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return "", fmt.Errorf("jenkins logs: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("jenkins logs status %d", resp.StatusCode)
	}
	return string(body), nil
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
