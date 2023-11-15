package server

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/kataras/iris/v12/middleware/jwt"
	"os"
	"time"
)

func (s *Server) assignJWTFiles() error {
	if _, err := os.Stat(s.Configuration.JWT.PrivateKey); os.IsNotExist(err) {
		// If the file does not exist && Auto Generation is enabled, we generate it
		if !s.Configuration.JWT.AutoGenerate {
			return fmt.Errorf("error reading JWT files: %v", err)
		}

		log.Info("Auto Generation is enabled so let's generate some JWT files for you!")
		err2 := generateJWTFiles(s.Configuration)
		if err2 != nil {
			return fmt.Errorf("error generating JWT files: %v", err)
		}
	}

	privateKey, publicKey := jwt.MustLoadEdDSA(s.Configuration.JWT.PrivateKey, s.Configuration.JWT.PublicKey)

	s.JWTSigner = jwt.NewSigner(jwt.EdDSA, privateKey, 120*time.Minute) // TODO: Move To Configuration
	s.JWTVerifier = jwt.NewVerifier(jwt.EdDSA, publicKey)

	s.JWTVerifier.WithDefaultBlocklist()

	return nil
}

func generateJWTFiles(cfg *config.Config) error {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("error generating ed25519 key: %v", err)
	}

	b, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("error marshaling private key: %v", err)
	}

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: b,
	}

	err = os.WriteFile(cfg.JWT.PrivateKey, pem.EncodeToMemory(block), 0644)
	if err != nil {
		return fmt.Errorf("error writing key to file: %v", err)
	}

	b, err = x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("error marshaling public key: %v", err)
	}

	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: b,
	}

	err = os.WriteFile(cfg.JWT.PublicKey, pem.EncodeToMemory(block), 0644)
	if err != nil {
		return fmt.Errorf("error writing key to file: %v", err)
	}

	return nil
}
