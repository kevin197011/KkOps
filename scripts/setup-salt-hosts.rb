#!/usr/bin/env ruby
# frozen_string_literal: true

require 'net/http'
require 'json'
require 'uri'

# 配置
API_BASE = 'http://localhost:8080/api/v1'
USERNAME = 'admin'
PASSWORD = 'admin123'

# Salt 主机信息
SALT_MASTER = {
  hostname: 'salt-master',
  ip_address: '192.168.97.3',
  ssh_port: 22,
  salt_minion_id: '',
  os_type: 'Linux',
  os_version: 'Ubuntu 22.04',
  status: 'online',
  description: 'Salt Master 服务器'
}

SALT_MINION = {
  hostname: 'salt-minion',
  ip_address: '192.168.97.7',
  ssh_port: 22,
  salt_minion_id: 'minion-ubuntu-01',
  os_type: 'Linux',
  os_version: 'Ubuntu 22.04',
  status: 'online',
  description: 'Salt Minion 服务器'
}

def make_request(method, path, body: nil, token: nil)
  uri = URI("#{API_BASE}#{path}")
  http = Net::HTTP.new(uri.host, uri.port)
  http.read_timeout = 30
  
  request = case method.upcase
            when 'GET'
              Net::HTTP::Get.new(uri)
            when 'POST'
              Net::HTTP::Post.new(uri)
            when 'PUT'
              Net::HTTP::Put.new(uri)
            when 'DELETE'
              Net::HTTP::Delete.new(uri)
            else
              raise "Unsupported method: #{method}"
            end
  
  request['Content-Type'] = 'application/json'
  request['Authorization'] = "Bearer #{token}" if token
  request.body = body.to_json if body
  
  response = http.request(request)
  
  begin
    parsed = JSON.parse(response.body)
    [response.code.to_i, parsed]
  rescue JSON::ParserError
    [response.code.to_i, { 'error' => response.body }]
  end
end

def login
  puts "正在登录..."
  status, data = make_request('POST', '/auth/login', body: {
    username: USERNAME,
    password: PASSWORD
  })
  
  if status == 200 && data['token']
    puts "✓ 登录成功"
    data['token']
  else
    puts "✗ 登录失败: #{data['error'] || data['message'] || '未知错误'}"
    exit 1
  end
end

def get_or_create_project(token)
  puts "\n正在获取项目列表..."
  status, data = make_request('GET', '/projects?page=1&page_size=100', token: token)
  
  if status == 200 && data['projects']
    projects = data['projects']
    
    # 查找 "Salt 测试" 项目
    salt_project = projects.find { |p| p['name'] == 'Salt 测试' }
    
    if salt_project
      puts "✓ 找到项目: #{salt_project['name']} (ID: #{salt_project['id']})"
      return salt_project['id']
    end
    
    # 如果没有找到，使用第一个项目
    if projects.any?
      project = projects.first
      puts "✓ 使用现有项目: #{project['name']} (ID: #{project['id']})"
      return project['id']
    end
  end
  
  # 创建新项目
  puts "正在创建项目: Salt 测试"
  status, data = make_request('POST', '/projects', body: {
    name: 'Salt 测试',
    description: 'Salt Master 和 Minion 测试项目',
    status: 'active'
  }, token: token)
  
  if status == 201 && data['project']
    puts "✓ 项目创建成功 (ID: #{data['project']['id']})"
    data['project']['id']
  else
    puts "✗ 项目创建失败: #{data['error'] || data['message'] || '未知错误'}"
    exit 1
  end
end

def create_host(token, project_id, host_info)
  puts "\n正在创建主机: #{host_info[:hostname]}"
  
  host_data = {
    project_id: project_id,
    hostname: host_info[:hostname],
    ip_address: host_info[:ip_address],
    ssh_port: host_info[:ssh_port],
    os_type: host_info[:os_type],
    os_version: host_info[:os_version],
    status: host_info[:status],
    description: host_info[:description],
    metadata: '{}'  # 确保 metadata 是有效的 JSON
  }
  
  # 只在有值时才添加 salt_minion_id
  host_data[:salt_minion_id] = host_info[:salt_minion_id] unless host_info[:salt_minion_id].empty?
  
  status, data = make_request('POST', '/hosts', body: host_data, token: token)
  
  if status == 201 && data['host']
    puts "✓ 主机创建成功: #{host_info[:hostname]} (ID: #{data['host']['id']})"
    data['host']['id']
  elsif status == 400 && data['error']&.include?('already exists')
    puts "⚠ 主机已存在: #{host_info[:hostname]}"
    # 尝试查找现有主机
    find_existing_host(token, host_info[:hostname])
  else
    puts "✗ 主机创建失败: #{data['error'] || data['message'] || '未知错误'}"
    nil
  end
end

def find_existing_host(token, hostname)
  status, data = make_request('GET', '/hosts?page=1&page_size=100', token: token)
  
  if status == 200 && data['hosts']
    host = data['hosts'].find { |h| h['hostname'] == hostname }
    if host
      puts "  找到现有主机 ID: #{host['id']}"
      return host['id']
    end
  end
  nil
end

def main
  puts "=========================================="
  puts "Salt 主机自动配置脚本"
  puts "=========================================="
  puts ""
  
  # 登录
  token = login
  
  # 获取或创建项目
  project_id = get_or_create_project(token)
  
  # 创建 Salt Master 主机
  master_id = create_host(token, project_id, SALT_MASTER)
  
  # 创建 Salt Minion 主机
  minion_id = create_host(token, project_id, SALT_MINION)
  
  puts "\n=========================================="
  puts "配置完成！"
  puts "=========================================="
  puts ""
  puts "可以在 http://localhost:3000 查看配置的主机："
  puts "  - Salt Master (ID: #{master_id})" if master_id
  puts "  - Salt Minion (ID: #{minion_id})" if minion_id
  puts ""
end

main if __FILE__ == $PROGRAM_NAME

