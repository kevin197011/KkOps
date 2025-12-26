#!/bin/bash
# Salt Minion 配置脚本 (Ubuntu 24.04)
# 参考: https://docs.saltproject.io/salt/install-guide/en/latest/topics/install-by-operating-system/linux-deb.html

set -e

echo "🚀 开始配置 Salt Minion..."

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

# 安装 Salt Minion
echo "📦 安装 Salt Minion..."
apt-get install -y salt-minion

# 配置 Salt Minion
echo "🔧 配置 Salt Minion..."
mkdir -p /etc/salt/minion.d

# 配置 Master 地址
cat > /etc/salt/minion.d/master.conf << EOF
master: ${MASTER_IP}
EOF

# 配置 Minion ID
cat > /etc/salt/minion.d/id.conf << EOF
id: ${MINION_ID:-salt-minion}
EOF

# 启动服务
echo "🚀 启动 Salt Minion..."
systemctl enable salt-minion
systemctl restart salt-minion

# 等待服务启动并尝试连接
sleep 5

# 验证安装
echo "🔍 验证安装..."
if systemctl is-active --quiet salt-minion; then
  echo "✅ Salt Minion service is running"
  systemctl status salt-minion --no-pager -l | head -5
else
  echo "❌ Salt Minion service failed to start"
  systemctl status salt-minion --no-pager -l
  exit 1
fi

# 显示版本
salt-minion --version || true

echo ""
echo "✅ Salt Minion 配置完成!"
echo "   - Minion ID: ${MINION_ID:-salt-minion}"
echo "   - Master: ${MASTER_IP}"

