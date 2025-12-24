# Salt 虚拟机管理

本目录包含创建和管理 Salt Master/Minion 虚拟机的脚本。

## 方式一：Vagrant + VirtualBox（推荐）

使用 Vagrant 和 VirtualBox 创建 Ubuntu 24.04 虚拟机。

### 前置要求

- [VirtualBox](https://www.virtualbox.org/wiki/Downloads)
- [Vagrant](https://www.vagrantup.com/downloads)

### 使用方法

```bash
cd vms

# 创建并启动所有虚拟机
vagrant up

# 仅创建 Salt Master
vagrant up salt-master

# 仅创建 Salt Minion
vagrant up salt-minion

# 查看虚拟机状态
vagrant status

# SSH 登录
vagrant ssh salt-master
vagrant ssh salt-minion

# 停止虚拟机
vagrant halt

# 销毁虚拟机
vagrant destroy -f
```

### 虚拟机信息

| 虚拟机 | IP 地址 | 配置 | 说明 |
|--------|---------|------|------|
| salt-master | 192.168.56.10 | 2 CPU, 2GB RAM | Salt Master + API |
| salt-minion | 192.168.56.11 | 1 CPU, 1GB RAM | Salt Minion |

### Salt API 访问

- URL: http://192.168.56.10:8000
- 用户名: `salt`
- 密码: `salt123`

### 验证 Salt 连接

```bash
# SSH 到 Master
vagrant ssh salt-master

# 查看已接受的 Minion 密钥
sudo salt-key -L

# 测试 Minion 连接
sudo salt '*' test.ping

# 执行命令
sudo salt '*' cmd.run 'hostname'
```

---

## 方式二：OrbStack（macOS）

使用 OrbStack 创建 Rocky 9 虚拟机（适合 ARM Mac）。

### 前置要求

- [OrbStack](https://orbstack.dev/)

### 使用方法

```bash
cd vms

# 运行创建脚本
ruby create_salt_vms.rb

# 调试模式
DEBUG=1 ruby create_salt_vms.rb
```

### OrbStack 常用命令

```bash
# 列出虚拟机
orb list

# SSH 登录
orb ssh salt-master
orb ssh salt-minion

# 查看虚拟机信息
orb info salt-master

# 停止虚拟机
orb stop salt-master salt-minion

# 删除虚拟机
orb delete salt-master salt-minion -f
```

---

## 选择建议

| 场景 | 推荐方式 |
|------|----------|
| Intel Mac / Linux / Windows | Vagrant + VirtualBox |
| Apple Silicon Mac (M1/M2/M3) | OrbStack |
| 需要精确控制资源 | Vagrant + VirtualBox |
| 快速开发测试 | OrbStack |

---

## 目录结构

```
vms/
├── Vagrantfile              # Vagrant 配置文件
├── create_salt_vms.rb       # OrbStack 创建脚本
├── scripts/
│   ├── provision-master.sh  # Salt Master 配置脚本
│   └── provision-minion.sh  # Salt Minion 配置脚本
├── README.md                # 本文件
└── .gitignore               # Git 忽略规则
```

## SSH 密钥

如果存在 `~/.ssh/id_ed25519.pub`，会自动配置到虚拟机的 root 用户，允许直接 SSH 登录：

```bash
# Vagrant 虚拟机
ssh root@192.168.56.10
ssh root@192.168.56.11

# OrbStack 虚拟机
ssh root@<master-ip>
ssh root@<minion-ip>
```

