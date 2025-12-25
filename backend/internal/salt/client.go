package salt

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kkops/backend/internal/config"
)

// Client Salt API 客户端
type Client struct {
	apiURL     string
	username   string
	password   string
	eauth      string
	timeout    time.Duration
	verifySSL  bool
	token      string
	tokenExp   time.Time
	httpClient *http.Client
}

// NewClient 创建新的 Salt API 客户端
func NewClient(cfg config.SaltConfig) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !cfg.VerifySSL,
		},
	}

	// 规范化 API URL（移除尾部斜杠）
	apiURL := cfg.APIURL
	if len(apiURL) > 0 && apiURL[len(apiURL)-1] == '/' {
		apiURL = apiURL[:len(apiURL)-1]
	}

	return &Client{
		apiURL:    apiURL,
		username:  cfg.Username,
		password:  cfg.Password,
		eauth:     cfg.EAuth,
		timeout:   time.Duration(cfg.Timeout) * time.Second,
		verifySSL: cfg.VerifySSL,
		httpClient: &http.Client{
			Transport: tr,
			Timeout:   time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

// LoginRequest Salt API 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	EAuth    string `json:"eauth"`
}

// LoginResponse Salt API 登录响应
type LoginResponse struct {
	Return []struct {
		Token  string   `json:"token"`
		Expire float64  `json:"expire"` // Salt API 返回的是浮点数（Unix 时间戳，可能包含毫秒）
		Start  float64  `json:"start"`  // Salt API 返回的是浮点数（Unix 时间戳，可能包含毫秒）
		User   string   `json:"user"`
		EAuth  string   `json:"eauth"`
		Perms  []string `json:"perms"`
	} `json:"return"`
}

// MinionStatusResponse Minion 状态响应
type MinionStatusResponse struct {
	Return []map[string]interface{} `json:"return"`
}

// CommandResponse 命令执行响应
type CommandResponse struct {
	Return []map[string]interface{} `json:"return"`
}

// authenticate 认证并获取 token
func (c *Client) authenticate() error {
	// 如果 token 未过期，直接返回
	if c.token != "" && time.Now().Before(c.tokenExp) {
		return nil
	}

	loginReq := LoginRequest{
		Username: c.username,
		Password: c.password,
		EAuth:    c.eauth,
	}

	reqBody, err := json.Marshal(loginReq)
	if err != nil {
		return fmt.Errorf("failed to marshal login request: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/login", c.apiURL), bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return fmt.Errorf("failed to decode login response: %w", err)
	}

	if len(loginResp.Return) == 0 || loginResp.Return[0].Token == "" {
		return fmt.Errorf("invalid login response: no token received")
	}

	c.token = loginResp.Return[0].Token
	// Token 有效期：使用 API 返回的 expire 时间，如果没有则使用 11 小时（提前刷新）
	if loginResp.Return[0].Expire > 0 {
		// expire 是 Unix 时间戳（秒），转换为 time.Time
		expireTime := time.Unix(int64(loginResp.Return[0].Expire), 0)
		// 提前 1 小时刷新 token
		c.tokenExp = expireTime.Add(-1 * time.Hour)
	} else {
		// 如果没有返回过期时间，使用默认的 11 小时
		c.tokenExp = time.Now().Add(11 * time.Hour)
	}

	return nil
}

// doRequest 执行 HTTP 请求
func (c *Client) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	if err := c.authenticate(); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.apiURL, endpoint), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Auth-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// GetMinionStatus 获取 Minion 状态
func (c *Client) GetMinionStatus(minionID string) (bool, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/minions/%s", minionID), nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil // Minion 不存在
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to get minion status: status %d, body: %s", resp.StatusCode, string(body))
	}

	var statusResp MinionStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return false, fmt.Errorf("failed to decode minion status response: %w", err)
	}

	// 如果返回有数据，说明 Minion 存在且在线
	return len(statusResp.Return) > 0, nil
}

// ExecuteCommand 执行 Salt 命令
func (c *Client) ExecuteCommand(target string, function string, args []interface{}) (map[string]interface{}, error) {
	cmdBody := map[string]interface{}{
		"client":  "local",
		"tgt":     target,
		"fun":     function,
		"arg":     args,
		"timeout": c.timeout.Seconds(),
	}

	resp, err := c.doRequest("POST", "/", cmdBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to execute command: status %d, body: %s", resp.StatusCode, string(body))
	}

	var cmdResp CommandResponse
	if err := json.NewDecoder(resp.Body).Decode(&cmdResp); err != nil {
		return nil, fmt.Errorf("failed to decode command response: %w", err)
	}

	if len(cmdResp.Return) == 0 {
		return nil, fmt.Errorf("empty command response")
	}

	return cmdResp.Return[0], nil
}

// TestPing 测试 Minion 连通性
func (c *Client) TestPing(target string) (map[string]bool, error) {
	result, err := c.ExecuteCommand(target, "test.ping", nil)
	if err != nil {
		return nil, err
	}

	pingResult := make(map[string]bool)
	for minion, response := range result {
		if response == true {
			pingResult[minion] = true
		} else {
			pingResult[minion] = false
		}
	}

	return pingResult, nil
}

// GetGrains 获取 Minion 的 Grains 信息
func (c *Client) GetGrains(target string) (map[string]interface{}, error) {
	return c.ExecuteCommand(target, "grains.items", nil)
}

// StateApply 执行 Salt state.apply
// target: 目标 Minion (可以使用 glob 匹配，如 "*" 或 "web*")
// state: State 文件路径或名称
// kwargs: 额外的参数（如 pillar 数据等）
func (c *Client) StateApply(target string, state string, kwargs map[string]interface{}) (map[string]interface{}, error) {
	args := []interface{}{state}
	if kwargs != nil {
		args = append(args, kwargs)
	}
	return c.ExecuteCommand(target, "state.apply", args)
}

// GetJobStatus 获取 Salt Job 状态
// jobID: Salt Job ID
func (c *Client) GetJobStatus(jobID string) (map[string]interface{}, error) {
	cmdBody := map[string]interface{}{
		"client": "runner",
		"fun":    "jobs.lookup_jid",
		"jid":    jobID,
	}

	resp, err := c.doRequest("POST", "/", cmdBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get job status: status %d, body: %s", resp.StatusCode, string(body))
	}

	var jobResp struct {
		Return []map[string]interface{} `json:"return"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&jobResp); err != nil {
		return nil, fmt.Errorf("failed to decode job status response: %w", err)
	}

	if len(jobResp.Return) == 0 {
		return nil, fmt.Errorf("empty job status response")
	}

	return jobResp.Return[0], nil
}

// TestConnection 测试Salt API连接
// 尝试认证并返回连接状态
func (c *Client) TestConnection() error {
	return c.authenticate()
}
