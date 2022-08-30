package resource

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/fraima/key-keeper/internal/config"
)

func (s *resource) checkCSR(cert config.Certificate) {
	csr, err := s.readCertificate(cert.HostPath)
	if csr != nil && time.Until(csr.NotAfter) > cert.RenewBefore {
		return
	}
	if err != nil && !os.IsNotExist(err) {
		zap.L().Error("read csr", zap.String("path", cert.HostPath), zap.Error(err))
	}

	crt, key, err := s.generateCSR(cert)
	if err != nil {
		zap.L().Error(
			"generate csr",
			zap.String("name", cert.Name),
			zap.Error(err),
		)
	}

	if err = s.storeKeyPair(cert.HostPath, crt, key); err != nil {
		zap.L().Error(
			"store csr",
			zap.String("name", cert.Name),
			zap.Error(err),
		)
		return
	}

	for _, command := range cert.Trigger {
		cmd := strings.Split(command, " ")
		err := exec.Command(cmd[0], cmd[1:]...).Run()
		zap.L().Error(
			"csr trigger",
			zap.String("name", cert.Name),
			zap.String("command", command),
			zap.Error(err),
		)
	}
	zap.L().Info("csr generated", zap.String("name", cert.Name))
}

func (s *resource) generateCSR(crt config.Certificate) ([]byte, []byte, error) {
	certData := map[string]interface{}{
		"name":      crt.Name,
		"alt_names": strings.Join(crt.Spec.Hostnames, ","),
		"ip_sans":   strings.Join(crt.Spec.IPAddresses, ","),
	}
	path := path.Join(crt.Vault.Path, "issue", crt.Vault.Role)
	csr, err := s.vault.Write(path, certData)
	if err != nil {
		return nil, nil, fmt.Errorf("generate with path %s : %w", path, err)
	}

	return []byte(csr["certificate"].(string)), []byte(csr["private_key"].(string)), nil
}
