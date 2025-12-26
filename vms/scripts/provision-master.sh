#!/bin/bash
# Salt Master 配置脚本 (Ubuntu 24.04)
# 参考: https://docs.saltproject.io/salt/install-guide/en/latest/topics/install-by-operating-system/linux-deb.html

set -e

echo "🚀 开始配置 Salt Master..."

export DEBIAN_FRONTEND=noninteractive

# 更新系统
echo "📦 更新系统包..."
apt-get update
apt-get install -y curl ca-certificates

# 添加 SaltStack 官方仓库
echo "📦 添加 SaltStack 仓库..."
mkdir -p /etc/apt/keyrings

# 下载 Salt Project 公钥
curl -fsSL https://packages.broadcom.com/artifactory/api/security/keypair/SaltProjectKey/public | tee /etc/apt/keyrings/salt-archive-keyring.pgp > /dev/null

# 下载 Salt 源配置
curl -fsSL https://github.com/saltstack/salt-install-guide/releases/latest/download/salt.sources | tee /etc/apt/sources.list.d/salt.sources > /dev/null

# 更新包元数据
apt-get update

# 安装 Salt Master、Minion 和 API
echo "📦 安装 Salt Master、Salt Minion 和 Salt API..."
apt-get install -y salt-master salt-minion salt-api python3-pip
pip3 install cherrypy --break-system-packages

# 配置 Salt Master
echo "🔧 配置 Salt Master..."
mkdir -p /etc/salt/master.d

# Salt API 配置
cat > /etc/salt/master.d/api.conf << 'EOF'
rest_cherrypy:
  port: 8000
  host: 0.0.0.0
  disable_ssl: true

# 启用 Salt API 客户端（重要！）
netapi_enable_clients:
  - local
  - local_async
  - runner
  - runner_async
EOF

# 外部认证配置
cat > /etc/salt/master.d/eauth.conf << 'EOF'
external_auth:
  pam:
    salt:
      - .*
      - "@runner"
      - "@wheel"
EOF

# 自动接受 Minion 密钥
echo 'auto_accept: True' > /etc/salt/master.d/auto_accept.conf

# 配置 Master 接口 IP（监听所有接口）
cat > /etc/salt/master.d/interface.conf << 'EOF'
interface: 0.0.0.0
EOF

# 配置 Salt Minion（指向本地 Master）
echo "🔧 配置 Salt Minion..."
mkdir -p /etc/salt/minion.d

cat > /etc/salt/minion.d/master.conf << EOF
master: 127.0.0.1
id: salt-master
EOF

# 配置防火墙（允许 Salt 端口）
if command -v ufw &> /dev/null; then
  ufw allow 4505/tcp
  ufw allow 4506/tcp
elif command -v firewall-cmd &> /dev/null; then
  firewall-cmd --permanent --add-port=4505/tcp
  firewall-cmd --permanent --add-port=4506/tcp
  firewall-cmd --reload
fi

# 创建 salt 用户（用于 API 认证）
useradd -m -s /bin/bash salt || true
echo 'salt:salt123' | chpasswd

# 启动服务
echo "🚀 启动 Salt 服务..."
systemctl enable salt-master
systemctl enable salt-minion
systemctl enable salt-api
systemctl restart salt-master
systemctl restart salt-api

# 等待 Master 启动后再启动 Minion
sleep 3
systemctl restart salt-minion

# 等待服务启动
sleep 5

# 验证端口监听
echo "🔍 检查端口监听..."
netstat -tlnp 2>/dev/null | grep -E '4505|4506' || ss -tlnp 2>/dev/null | grep -E '4505|4506' || echo "⚠️  无法检查端口状态"

# 验证安装
echo "🔍 验证安装..."
if systemctl is-active --quiet salt-master; then
  echo "✅ Salt Master service is running"
  systemctl status salt-master --no-pager -l | head -5
else
  echo "❌ Salt Master service failed to start"
  systemctl status salt-master --no-pager -l
  exit 1
fi

if systemctl is-active --quiet salt-minion; then
  echo "✅ Salt Minion service is running"
  systemctl status salt-minion --no-pager -l | head -5
else
  echo "❌ Salt Minion service failed to start"
  systemctl status salt-minion --no-pager -l
fi

# 显示版本
salt-master --version || true
salt-minion --version || true

echo ""
echo "✅ Salt Master 配置完成!"
echo "   - Master IP: ${MASTER_IP}"
echo "   - Salt API: http://${MASTER_IP}:8000"
echo "   - API 用户: salt / salt123"

