package salt

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

	// HTTP 客户端超时设置为 5 分钟，用于长时间运行的请求
	httpTimeout := 5 * time.Minute

	return &Client{
		apiURL:    apiURL,
		username:  cfg.Username,
		password:  cfg.Password,
		eauth:     cfg.EAuth,
		timeout:   time.Duration(cfg.Timeout) * time.Second,
		verifySSL: cfg.VerifySSL,
		httpClient: &http.Client{
			Transport: tr,
			Timeout:   httpTimeout,
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
	// 使用 test.ping 命令来检查 Minion 是否在线
	pingResult, err := c.TestPing(minionID)
	if err != nil {
		return false, err
	}

	// 检查指定 minion 是否在 ping 结果中且返回 true
	if status, exists := pingResult[minionID]; exists {
		return status, nil
	}

	// 如果 minion 不在结果中，说明它不在线或不存在
	return false, nil
}

// ExecuteCommand 执行 Salt 命令
func (c *Client) ExecuteCommand(target string, function string, args []interface{}) (map[string]interface{}, error) {
	cmdBody := map[string]interface{}{
		"client":  "local",
		"tgt":     target,
		"fun":     function,
		"timeout": c.timeout.Seconds(),
	}

	// 只有当 args 不为空时才设置 arg 字段
	if len(args) > 0 {
		// 对于 pkg.install 等函数，需要特殊处理参数格式
		if function == "pkg.install" && len(args) > 0 {
			var finalArgs []string

			// 处理嵌套数组的情况（如 [["nginx","mysql-server"]] -> ["nginx","mysql-server"]）
			if len(args) == 1 {
				if nestedArray, ok := args[0].([]interface{}); ok {
					// 如果第一个参数是数组，展开它
					for _, item := range nestedArray {
						if str, ok := item.(string); ok {
							finalArgs = append(finalArgs, str)
						}
					}
				} else if str, ok := args[0].(string); ok {
					finalArgs = append(finalArgs, str)
				}
			} else {
				// 处理平坦数组的情况
				for _, arg := range args {
					if str, ok := arg.(string); ok {
						finalArgs = append(finalArgs, str)
					}
				}
			}

			if len(finalArgs) > 0 {
				cmdBody["arg"] = strings.Join(finalArgs, ",")
			} else {
				cmdBody["arg"] = args
			}
		} else {
			cmdBody["arg"] = args
		}
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

// GetDiskUsage 获取 Minion 的磁盘使用情况
func (c *Client) GetDiskUsage(target string) (map[string]interface{}, error) {
	return c.ExecuteCommand(target, "disk.usage", nil)
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

// MinionInfo Minion 信息结构
type MinionInfo struct {
	ID          string                 `json:"id"`
	IPAddress   string                 `json:"ip_address"`
	Hostname    string                 `json:"hostname"`
	OS          string                 `json:"os"`
	OSRelease   string                 `json:"os_release"`
	NumCPUs     int                    `json:"num_cpus"`
	MemTotalMB  float64                `json:"mem_total"`
	SaltVersion string                 `json:"salt_version"`
	Grains      map[string]interface{} `json:"grains"`
}

// ListMinions 获取所有已连接的 Minion 列表
// 首先尝试使用 GET /minions 端点，如果失败则尝试使用 runner 客户端
func (c *Client) ListMinions() ([]MinionInfo, error) {
	// 方法1: 尝试使用 GET /minions 端点
	resp, err := c.doRequest("GET", "/minions", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list minions: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// 如果 GET /minions 失败，尝试使用 runner 客户端
		return c.listMinionsViaRunner()
	}

	var minionResp struct {
		Return []map[string]interface{} `json:"return"`
	}
	if err := json.Unmarshal(body, &minionResp); err != nil {
		return nil, fmt.Errorf("failed to decode minions response: %w", err)
	}

	if len(minionResp.Return) == 0 {
		return []MinionInfo{}, nil
	}

	var minions []MinionInfo
	for minionID, data := range minionResp.Return[0] {
		info := MinionInfo{
			ID: minionID,
		}

		if grains, ok := data.(map[string]interface{}); ok {
			info.Grains = grains
			c.extractMinionInfo(&info, grains)
		}

		minions = append(minions, info)
	}

	return minions, nil
}

// listMinionsViaRunner 使用 runner 客户端获取 minion 列表
func (c *Client) listMinionsViaRunner() ([]MinionInfo, error) {
	// 使用 manage.status runner 获取 minion 状态
	cmdBody := map[string]interface{}{
		"client": "runner",
		"fun":    "manage.status",
	}

	resp, err := c.doRequest("POST", "/", cmdBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get minion status via runner: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list minions via runner: status %d, body: %s", resp.StatusCode, string(body))
	}

	var runnerResp struct {
		Return []map[string]interface{} `json:"return"`
	}
	if err := json.Unmarshal(body, &runnerResp); err != nil {
		return nil, fmt.Errorf("failed to decode runner response: %w", err)
	}

	if len(runnerResp.Return) == 0 {
		return []MinionInfo{}, nil
	}

	var minions []MinionInfo
	result := runnerResp.Return[0]

	// manage.status 返回 {"up": ["minion1", "minion2"], "down": ["minion3"]}
	if upMinions, ok := result["up"].([]interface{}); ok {
		for _, m := range upMinions {
			if minionID, ok := m.(string); ok {
				minions = append(minions, MinionInfo{
					ID: minionID,
				})
			}
		}
	}

	return minions, nil
}

// extractMinionInfo 从 grains 中提取 minion 信息
func (c *Client) extractMinionInfo(info *MinionInfo, grains map[string]interface{}) {
	// 提取 IP 地址（优先选择私有 IP）
	if ipv4, ok := grains["ipv4"].([]interface{}); ok {
		var privateIP, publicIP string
		for _, ip := range ipv4 {
			if ipStr, ok := ip.(string); ok && ipStr != "127.0.0.1" {
				if isPrivateIPAddr(ipStr) {
					privateIP = ipStr
				} else {
					publicIP = ipStr
				}
			}
		}
		// 优先使用私有 IP
		if privateIP != "" {
			info.IPAddress = privateIP
		} else if publicIP != "" {
			info.IPAddress = publicIP
		}
	}
	// 提取主机名
	if hostname, ok := grains["host"].(string); ok {
		info.Hostname = hostname
	} else if fqdn, ok := grains["fqdn"].(string); ok {
		info.Hostname = fqdn
	}
	// 提取操作系统信息
	if os, ok := grains["os"].(string); ok {
		info.OS = os
	}
	if osRelease, ok := grains["osrelease"].(string); ok {
		info.OSRelease = osRelease
	}
	// 提取 CPU 核心数
	if numCpus, ok := grains["num_cpus"].(float64); ok {
		info.NumCPUs = int(numCpus)
	}
	// 提取内存信息（MB）
	if memTotal, ok := grains["mem_total"].(float64); ok {
		info.MemTotalMB = memTotal
	}
	// 提取 Salt 版本
	if saltVersion, ok := grains["saltversion"].(string); ok {
		info.SaltVersion = saltVersion
	}
}

// ExecuteState 执行 Salt state.apply
// 使用 local_async 客户端异步执行，然后轮询结果
func (c *Client) ExecuteState(formulaName, target string, pillarData map[string]interface{}) (map[string]interface{}, error) {
	cmdBody := map[string]interface{}{
		"client":  "local_async",
		"tgt":     target,
		"fun":     "state.apply",
		"arg":     []interface{}{formulaName},
		"timeout": 600, // 10 分钟超时
	}

	// 如果 target 包含逗号，说明是多个主机，使用 list 匹配类型
	if strings.Contains(target, ",") {
		cmdBody["tgt_type"] = "list"
	}

	// 添加 Pillar 数据
	if pillarData != nil && len(pillarData) > 0 {
		cmdBody["kwarg"] = map[string]interface{}{
			"pillar": pillarData,
		}
	}

	resp, err := c.doRequest("POST", "/", cmdBody)
	if err != nil {
		return nil, fmt.Errorf("failed to execute state: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("state execution failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	// 解析异步响应获取 JID
	var asyncResp struct {
		Return []struct {
			JID     string   `json:"jid"`
			Minions []string `json:"minions"`
		} `json:"return"`
	}
	if err := json.Unmarshal(body, &asyncResp); err != nil {
		return nil, fmt.Errorf("failed to decode async response: %w", err)
	}

	if len(asyncResp.Return) == 0 || asyncResp.Return[0].JID == "" {
		return nil, fmt.Errorf("no job ID returned from async execution")
	}

	jid := asyncResp.Return[0].JID
	expectedMinions := asyncResp.Return[0].Minions

	// 轮询等待结果，最多等待 10 分钟
	maxWait := 600 * time.Second
	pollInterval := 5 * time.Second
	startTime := time.Now()

	for time.Since(startTime) < maxWait {
		result, err := c.GetJobResult(jid)
		if err != nil {
			time.Sleep(pollInterval)
			continue
		}

		// 检查是否所有 minion 都返回了结果
		if len(result) >= len(expectedMinions) || len(result) > 0 {
			// 返回格式与同步执行一致: {"return": [{minion_id: result, ...}]}
			return map[string]interface{}{"return": []interface{}{result}}, nil
		}

		time.Sleep(pollInterval)
	}

	return nil, fmt.Errorf("timeout waiting for job %s to complete", jid)
}

// GetJobResult 获取 Salt Job 结果（只返回 minion 结果数据）
func (c *Client) GetJobResult(jobID string) (map[string]interface{}, error) {
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
		return nil, fmt.Errorf("failed to get job result: status %d, body: %s", resp.StatusCode, string(body))
	}

	var jobResp struct {
		Return []map[string]interface{} `json:"return"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&jobResp); err != nil {
		return nil, fmt.Errorf("failed to decode job result response: %w", err)
	}

	if len(jobResp.Return) == 0 {
		return nil, fmt.Errorf("empty job result response")
	}

	result := jobResp.Return[0]

	// jobs.lookup_jid 返回格式可能有两种:
	// 1. 新格式（带 highstate outputter）: {"return": [{"outputter": "highstate", "data": {"minion_id": {...}}}]}
	// 2. 旧格式: {"return": [{"minion_id": {...state results...}, ...}]}
	
	// 检查是否是新格式（带 data 字段）
	if data, ok := result["data"].(map[string]interface{}); ok {
		// 新格式：从 data 字段中提取 minion 结果
		minionResults := make(map[string]interface{})
		for key, value := range data {
			minionResults[key] = value
		}
		return minionResults, nil
	}

	// 旧格式：过滤掉元数据字段，只保留 minion 结果
	minionResults := make(map[string]interface{})
	for key, value := range result {
		// 跳过 jobs.lookup_jid 返回的元数据字段
		if key == "outputter" || key == "retcode" {
			continue
		}
		minionResults[key] = value
	}

	return minionResults, nil
}

// ExecuteStateWithArgs 执行 Salt state.apply 带额外参数
func (c *Client) ExecuteStateWithArgs(formulaName, target string, args []interface{}, pillarData map[string]interface{}) (map[string]interface{}, error) {
	cmdBody := map[string]interface{}{
		"client": "local",
		"tgt":    target,
		"fun":    "state.apply",
	}

	// 构建参数列表
	stateArgs := []interface{}{formulaName}
	if len(args) > 0 {
		stateArgs = append(stateArgs, args...)
	}
	cmdBody["arg"] = stateArgs

	// 添加 Pillar 数据
	if pillarData != nil && len(pillarData) > 0 {
		if kwarg, ok := cmdBody["kwarg"].(map[string]interface{}); ok {
			kwarg["pillar"] = pillarData
		} else {
			cmdBody["kwarg"] = map[string]interface{}{
				"pillar": pillarData,
			}
		}
	}

	resp, err := c.doRequest("POST", "/", cmdBody)
	if err != nil {
		return nil, fmt.Errorf("failed to execute state: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("state execution failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// isPrivateIPAddr 判断是否为私有 IP 地址
func isPrivateIPAddr(ip string) bool {
	if len(ip) < 3 {
		return false
	}
	// 简单判断私有 IP 范围
	return ip[:3] == "10." ||
		(len(ip) >= 4 && ip[:4] == "172.") ||
		(len(ip) >= 8 && ip[:8] == "192.168.")
}
