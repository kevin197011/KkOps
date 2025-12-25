package service

import (
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/repository"
	"github.com/kkops/backend/internal/utils"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

type SSHKeyService interface {
	CreateKey(userID uint64, name, username, privateKeyContent string) (*models.SSHKey, error)
	GetKey(id, userID uint64) (*models.SSHKey, error)
	ListKeys(userID uint64, page, pageSize int) ([]models.SSHKey, int64, error)
	UpdateKey(id, userID uint64, name, username string, privateKeyContent string) (*models.SSHKey, error)
	DeleteKey(id, userID uint64) error
	GetDecryptedPrivateKey(id, userID uint64) (string, error) // 获取解密后的私钥（用于SSH连接）
}

type sshKeyService struct {
	keyRepo repository.SSHKeyRepository
}

func NewSSHKeyService(keyRepo repository.SSHKeyRepository) SSHKeyService {
	return &sshKeyService{keyRepo: keyRepo}
}

// validateAndParseSSHKey 验证并解析SSH私钥
func validateAndParseSSHKey(privateKeyContent string) (keyType string, publicKey string, fingerprint string, err error) {
	// 清理私钥内容（移除前后空白）
	privateKeyContent = strings.TrimSpace(privateKeyContent)

	// 解析PEM块
	block, _ := pem.Decode([]byte(privateKeyContent))
	if block == nil {
		return "", "", "", errors.New("invalid SSH key format: not a valid PEM block")
	}

	var keyTypeStr string
	var pubKey ssh.PublicKey

	// 根据密钥类型解析
	switch block.Type {
	case "RSA PRIVATE KEY", "PRIVATE KEY":
		// 尝试解析为RSA密钥
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			// 尝试解析为PKCS8格式
			keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				// 尝试解析为OpenSSH格式
				keyInterface, err := ssh.ParseRawPrivateKey(block.Bytes)
				if err != nil {
					return "", "", "", fmt.Errorf("failed to parse private key: %w", err)
				}
				key, ok := keyInterface.(*rsa.PrivateKey)
				if !ok {
					// 尝试ED25519
					if ed25519Key, ok := keyInterface.(*ed25519.PrivateKey); ok {
						keyTypeStr = "ed25519"
						pubKeyInterface, err := ssh.NewPublicKey((*ed25519Key).Public())
						if err != nil {
							return "", "", "", fmt.Errorf("failed to create public key: %w", err)
						}
						pubKey = pubKeyInterface
					} else {
						return "", "", "", errors.New("unsupported key type")
					}
				} else {
					keyTypeStr = "rsa"
					pubKeyInterface, err := ssh.NewPublicKey(&key.PublicKey)
					if err != nil {
						return "", "", "", fmt.Errorf("failed to create public key: %w", err)
					}
					pubKey = pubKeyInterface
				}
			} else {
				if rsaKey, ok := keyInterface.(*rsa.PrivateKey); ok {
					keyTypeStr = "rsa"
					pubKeyInterface, err := ssh.NewPublicKey(&rsaKey.PublicKey)
					if err != nil {
						return "", "", "", fmt.Errorf("failed to create public key: %w", err)
					}
					pubKey = pubKeyInterface
				} else {
					return "", "", "", errors.New("unsupported key type")
				}
			}
		} else {
			keyTypeStr = "rsa"
			pubKeyInterface, err := ssh.NewPublicKey(&key.PublicKey)
			if err != nil {
				return "", "", "", fmt.Errorf("failed to create public key: %w", err)
			}
			pubKey = pubKeyInterface
		}
	case "OPENSSH PRIVATE KEY":
		// 对于 OpenSSH 格式，使用 ssh.ParsePrivateKey 直接解析整个私钥内容
		signer, err := ssh.ParsePrivateKey([]byte(privateKeyContent))
		if err != nil {
			return "", "", "", fmt.Errorf("failed to parse OpenSSH private key: %w", err)
		}
		// 从 signer 获取公钥
		pubKey = signer.PublicKey()
		// 确定密钥类型
		switch pubKey.Type() {
		case ssh.KeyAlgoRSA, ssh.KeyAlgoRSASHA256, ssh.KeyAlgoRSASHA512:
			keyTypeStr = "rsa"
		case ssh.KeyAlgoED25519:
			keyTypeStr = "ed25519"
		case ssh.KeyAlgoECDSA256, ssh.KeyAlgoECDSA384, ssh.KeyAlgoECDSA521:
			keyTypeStr = "ecdsa"
		default:
			return "", "", "", fmt.Errorf("unsupported OpenSSH key type: %s", pubKey.Type())
		}
	case "EC PRIVATE KEY":
		keyTypeStr = "ecdsa"
		keyInterface, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to parse ECDSA private key: %w", err)
		}
		pubKeyInterface, err := ssh.NewPublicKey(&keyInterface.PublicKey)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to create public key: %w", err)
		}
		pubKey = pubKeyInterface
	default:
		return "", "", "", fmt.Errorf("unsupported key type: %s", block.Type)
	}

	// 生成公钥字符串
	publicKeyBytes := ssh.MarshalAuthorizedKey(pubKey)
	publicKeyStr := strings.TrimSpace(string(publicKeyBytes))

	// 计算指纹（MD5格式，类似OpenSSH）
	md5Fingerprint := md5.Sum(pubKey.Marshal())
	fingerprintStr := ""
	for i, b := range md5Fingerprint {
		if i > 0 {
			fingerprintStr += ":"
		}
		fingerprintStr += fmt.Sprintf("%02x", b)
	}

	return keyTypeStr, publicKeyStr, fingerprintStr, nil
}

func (s *sshKeyService) CreateKey(userID uint64, name, username, privateKeyContent string) (*models.SSHKey, error) {
	// 验证并解析密钥
	keyType, publicKey, fingerprint, err := validateAndParseSSHKey(privateKeyContent)
	if err != nil {
		return nil, fmt.Errorf("invalid SSH key: %w", err)
	}

	// 加密私钥
	encryptedPrivateKey, err := utils.Encrypt(privateKeyContent)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt private key: %w", err)
	}

	// 创建SSH密钥记录
	sshKey := &models.SSHKey{
		UserID:      userID,
		Name:        name,
		Username:    username,
		KeyType:     keyType,
		PrivateKey:  encryptedPrivateKey,
		PublicKey:   publicKey,
		Fingerprint: fingerprint,
	}

	if err := s.keyRepo.Create(sshKey); err != nil {
		return nil, err
	}

	return sshKey, nil
}

func (s *sshKeyService) GetKey(id, userID uint64) (*models.SSHKey, error) {
	key, err := s.keyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 检查用户权限
	if key.UserID != userID {
		return nil, errors.New("unauthorized: key does not belong to user")
	}

	// 不返回私钥
	key.PrivateKey = ""

	return key, nil
}

func (s *sshKeyService) ListKeys(userID uint64, page, pageSize int) ([]models.SSHKey, int64, error) {
	offset := (page - 1) * pageSize
	keys, total, err := s.keyRepo.List(userID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 移除私钥内容
	for i := range keys {
		keys[i].PrivateKey = ""
	}

	return keys, total, err
}

func (s *sshKeyService) UpdateKey(id, userID uint64, name, username string, privateKeyContent string) (*models.SSHKey, error) {
	key, err := s.keyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 检查用户权限
	if key.UserID != userID {
		return nil, errors.New("unauthorized: key does not belong to user")
	}

	key.Name = name
	key.Username = username

	// 如果提供了私钥内容，更新私钥、公钥、指纹和密钥类型
	if privateKeyContent != "" {
		// 验证并解析新的私钥
		keyType, publicKey, fingerprint, err := validateAndParseSSHKey(privateKeyContent)
		if err != nil {
			return nil, fmt.Errorf("invalid SSH key: %w", err)
		}

		// 加密私钥
		encryptedPrivateKey, err := utils.Encrypt(privateKeyContent)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt private key: %w", err)
		}

		// 更新密钥信息
		key.PrivateKey = encryptedPrivateKey
		key.PublicKey = publicKey
		key.Fingerprint = fingerprint
		key.KeyType = keyType
	}

	if err := s.keyRepo.Update(key); err != nil {
		return nil, err
	}

	// 不返回私钥
	key.PrivateKey = ""

	return key, nil
}

func (s *sshKeyService) DeleteKey(id, userID uint64) error {
	key, err := s.keyRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 检查用户权限
	if key.UserID != userID {
		return errors.New("unauthorized: key does not belong to user")
	}

	// 检查密钥是否被使用
	inUse, err := s.keyRepo.CheckKeyInUse(id)
	if err != nil {
		return err
	}
	if inUse {
		return errors.New("cannot delete key: key is in use by one or more hosts")
	}

	return s.keyRepo.Delete(id)
}

func (s *sshKeyService) GetDecryptedPrivateKey(id, userID uint64) (string, error) {
	key, err := s.keyRepo.GetByID(id)
	if err != nil {
		return "", err
	}

	// 检查用户权限
	if key.UserID != userID {
		return "", errors.New("unauthorized: key does not belong to user")
	}

	// 解密私钥
	privateKey, err := utils.Decrypt(key.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt private key: %w", err)
	}

	return privateKey, nil
}

