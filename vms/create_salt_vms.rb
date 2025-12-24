#!/usr/bin/env ruby
# frozen_string_literal: true

# 使用 OrbStack 创建 Salt Master 和 Minion 虚拟机（基于 Rocky 9）
# 注意: OrbStack 使用默认资源配置，不支持在创建时指定内存/CPU/磁盘

require 'open3'
require 'fileutils'
require 'shellwords'

class SaltVMCreator
  def initialize
    @ssh_key_path = File.expand_path('~/.ssh/id_ed25519.pub')
    @use_ssh_key = File.exist?(@ssh_key_path)
  end

  def run
    puts '🚀 开始创建 Salt 虚拟机...'
    puts 'ℹ️  注意: OrbStack 将使用默认资源配置（通常足够运行 Salt）'
    puts '🔄 自动模式: 如果虚拟机已存在，将自动删除并重新创建'
    puts

    check_orbstack
    check_ssh_key

    create_salt_master
    create_salt_minion

    show_summary
  end

  private

  def check_orbstack
    unless command_exists?('orb')
      puts '❌ 错误: 未找到 orb 命令，请先安装 OrbStack'
      puts '   访问 https://orbstack.dev/ 下载安装'
      exit 1
    end
  end

  def check_ssh_key
    if @use_ssh_key
      puts "✅ 找到 SSH 公钥: #{@ssh_key_path}"
    else
      puts "⚠️  警告: 未找到 SSH 公钥 #{@ssh_key_path}"
      puts '   将使用密码认证（默认: root 或 rocky 用户）'
    end
    puts
  end

  def create_salt_master
    puts '📦 创建 Salt Master 虚拟机...'

    # 检查并删除已存在的虚拟机
    if vm_exists?('salt-master')
      puts '🗑️  删除现有虚拟机 salt-master...'
      run_command('orb delete salt-master -f', allow_failure: true)
      sleep 2
    end

    puts '📦 创建新的 Salt Master 虚拟机...'
    run_command('orb create rocky:9 salt-master')

    puts '⏳ 等待虚拟机启动...'
    sleep 15

    puts '🔧 配置 Salt Master...'
    configure_salt_master

    @master_ip = get_vm_ip('salt-master')
    puts "✅ Salt Master IP: #{@master_ip}"
    puts
  end

  def create_salt_minion
    puts '📦 创建 Salt Minion 虚拟机...'

    # 检查并删除已存在的虚拟机
    if vm_exists?('salt-minion')
      puts '🗑️  删除现有虚拟机 salt-minion...'
      run_command('orb delete salt-minion -f', allow_failure: true)
      sleep 2
    end

    puts '📦 创建新的 Salt Minion 虚拟机...'
    run_command('orb create rocky:9 salt-minion')

    puts '⏳ 等待虚拟机启动...'
    sleep 10

    puts '🔧 配置 Salt Minion...'
    configure_salt_minion

    @minion_ip = get_vm_ip('salt-minion')
    puts "✅ Salt Minion IP: #{@minion_ip}"
    puts
  end

  def configure_salt_master
    script = <<~'SCRIPT'
set -e

# 更新系统
sudo dnf update -y
sudo dnf install -y curl openssh-server ca-certificates

# 添加 SaltStack 官方仓库
# 参考: https://docs.saltproject.io/salt/install-guide/en/latest/topics/install-by-operating-system/linux-rpm.html
curl -fsSL https://github.com/saltstack/salt-install-guide/releases/latest/download/salt.repo | sudo tee /etc/yum.repos.d/salt.repo > /dev/null

# 清理仓库缓存
sudo dnf clean expire-cache

# 安装 Salt Master 和 API
sudo dnf install -y salt-master salt-api python3-pip
sudo pip3 install cherrypy

    SCRIPT

    # 添加 SSH 密钥配置
    script += ssh_key_config

    script += <<~'SCRIPT'

# 配置 Salt Master
sudo mkdir -p /etc/salt/master.d

# Salt API 配置
cat << 'EOFCONF' | sudo tee /etc/salt/master.d/api.conf > /dev/null
rest_cherrypy:
  port: 8000
  host: 0.0.0.0
  disable_ssl: true
EOFCONF

# 外部认证配置
cat << 'EOFCONF' | sudo tee /etc/salt/master.d/eauth.conf > /dev/null
external_auth:
  pam:
    salt:
      - .*
      - "@runner"
      - "@wheel"
EOFCONF

# 自动接受 Minion 密钥
echo 'auto_accept: True' | sudo tee /etc/salt/master.d/auto_accept.conf > /dev/null

# 创建 salt 用户（用于认证）
sudo useradd -m -s /bin/bash salt || true
echo 'salt:salt123' | sudo chpasswd

# 启动服务
sudo systemctl enable salt-master
sudo systemctl enable salt-api
sudo systemctl start salt-master
sudo systemctl start salt-api

# 等待服务启动
sleep 3

echo 'Salt Master installation completed!'

# 验证安装 - 检查服务状态
if sudo systemctl is-active --quiet salt-master; then
  echo '✅ Salt Master service is running'
  sudo systemctl status salt-master --no-pager -l | head -5
else
  echo '❌ Salt Master service failed to start'
  sudo systemctl status salt-master --no-pager -l
  exit 1
fi

# 尝试显示版本（如果命令可用）
if command -v salt-master &> /dev/null; then
  salt-master --version 2>/dev/null || true
elif [ -f /usr/bin/salt-master ]; then
  /usr/bin/salt-master --version 2>/dev/null || true
fi
    SCRIPT

    run_command_with_stdin('salt-master', script)
  end

  def configure_salt_minion
    script = <<~'SCRIPT'
set -e

# 更新系统
sudo dnf update -y
sudo dnf install -y curl openssh-server ca-certificates

# 添加 SaltStack 官方仓库
# 参考: https://docs.saltproject.io/salt/install-guide/en/latest/topics/install-by-operating-system/linux-rpm.html
curl -fsSL https://github.com/saltstack/salt-install-guide/releases/latest/download/salt.repo | sudo tee /etc/yum.repos.d/salt.repo > /dev/null

# 清理仓库缓存
sudo dnf clean expire-cache

# 安装 Salt Minion
sudo dnf install -y salt-minion

    SCRIPT

    # 添加 SSH 密钥配置
    script += ssh_key_config

    # 配置 Salt Minion（需要插入变量）
    script += <<~SCRIPT

# 配置 Salt Minion
sudo mkdir -p /etc/salt/minion.d
echo 'master: #{@master_ip}' | sudo tee /etc/salt/minion.d/master.conf > /dev/null
echo 'id: minion-rocky-01' | sudo tee /etc/salt/minion.d/id.conf > /dev/null

# 启动服务
sudo systemctl enable salt-minion
sudo systemctl start salt-minion

# 等待服务启动
sleep 3

echo 'Salt Minion installation completed!'

# 验证安装 - 检查服务状态
if sudo systemctl is-active --quiet salt-minion; then
  echo '✅ Salt Minion service is running'
  sudo systemctl status salt-minion --no-pager -l | head -5
else
  echo '❌ Salt Minion service failed to start'
  sudo systemctl status salt-minion --no-pager -l
  exit 1
fi

# 尝试显示版本（如果命令可用）
if command -v salt-minion &> /dev/null; then
  salt-minion --version 2>/dev/null || true
elif [ -f /usr/bin/salt-minion ]; then
  /usr/bin/salt-minion --version 2>/dev/null || true
fi
    SCRIPT

    run_command_with_stdin('salt-minion', script)
  end

  def ssh_key_config
    return '' unless @use_ssh_key

    ssh_key_content = File.read(@ssh_key_path).strip
    <<~SCRIPT

# 配置 SSH 密钥
sudo mkdir -p /root/.ssh
sudo chmod 700 /root/.ssh
echo '#{ssh_key_content}' | sudo tee -a /root/.ssh/authorized_keys > /dev/null
sudo chmod 600 /root/.ssh/authorized_keys
    SCRIPT
  end

  def vm_exists?(vm_name)
    stdout, stderr, status = Open3.capture3('orb list 2>&1')
    return false unless status.success?

    stdout.split("\n").any? { |line| line.include?(vm_name) }
  end

  def get_vm_ip(vm_name)
    stdout, stderr, status = Open3.capture3("orb exec #{vm_name} hostname -I")
    unless status.success?
      puts "⚠️  警告: 无法获取 #{vm_name} 的 IP 地址"
      return 'unknown'
    end

    stdout.split.first || 'unknown'
  end

  def run_command(command, allow_failure: false)
    puts "🔧 执行: #{command}" if ENV['DEBUG']
    stdout, stderr, status = Open3.capture3(command)

    unless status.success?
      if allow_failure
        puts "⚠️  命令执行失败（已忽略）: #{command}"
        puts stderr if stderr && !stderr.empty?
        return
      end

      puts "❌ 命令执行失败: #{command}"
      puts "错误输出: #{stderr}" if stderr && !stderr.empty?
      puts "标准输出: #{stdout}" if stdout && !stdout.empty?
      exit 1
    end

    puts stdout if stdout && !stdout.empty? && ENV['DEBUG']
    stdout
  end

  def run_command_with_stdin(vm_name, stdin_data)
    puts "🔧 执行脚本到虚拟机: #{vm_name}" if ENV['DEBUG']

    # 使用 Open3.popen3 正确处理标准输入
    Open3.popen3('orb', 'exec', vm_name, 'bash') do |stdin, stdout, stderr, wait_thr|
      stdin.write(stdin_data)
      stdin.close

      out = stdout.read
      err = stderr.read
      status = wait_thr.value

      unless status.success?
        puts "❌ 命令执行失败: orb exec #{vm_name} bash"
        puts "错误输出: #{err}" if err && !err.empty?
        puts "标准输出: #{out}" if out && !out.empty?
        exit 1
      end

      puts out if out && !out.empty?
      out
    end
  end

  def command_exists?(command)
    system("command -v #{command} > /dev/null 2>&1")
  end

  def show_summary
    puts
    puts '✅ 虚拟机创建完成！'
    puts
    puts '📋 虚拟机信息:'
    puts '   Salt Master:'
    puts "     - 名称: salt-master"
    puts "     - IP: #{@master_ip}"
    puts '     - SSH: orb ssh salt-master'
    puts "     - Salt API: http://#{@master_ip}:8000"
    puts
    puts '   Salt Minion:'
    puts "     - 名称: salt-minion"
    puts "     - IP: #{@minion_ip}"
    puts '     - SSH: orb ssh salt-minion'
    puts '     - Minion ID: minion-rocky-01'
    puts
    puts '💡 提示: 查看虚拟机资源使用情况:'
    puts '   orb info salt-master'
    puts '   orb info salt-minion'
    puts
    puts '🔍 验证连接:'
    puts "   orb ssh salt-master -c 'sudo salt-key -L'"
    puts "   orb ssh salt-master -c \"sudo salt '*' test.ping\""
    puts
  end
end

# 运行主程序
if __FILE__ == $PROGRAM_NAME
  begin
    creator = SaltVMCreator.new
    creator.run
  rescue Interrupt
    puts "\n⚠️  操作被用户中断"
    exit 1
  rescue StandardError => e
    puts "❌ 发生错误: #{e.class}: #{e.message}"
    puts e.backtrace if ENV['DEBUG']
    exit 1
  end
end

