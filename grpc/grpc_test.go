package grpc

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	ggrpc "google.golang.org/grpc"
)

type noopServer struct{}

func (noopServer) Register(gRPC *ggrpc.Server) error {
	return nil
}

func TestRegisterServerRejectsNil(t *testing.T) {
	manager := NewManager()
	if err := manager.RegisterServer("svc", nil); err == nil {
		t.Fatal("expected nil server registration to fail")
	}
}

func TestRegisterServerRejectsDuplicate(t *testing.T) {
	manager := NewManager()
	if err := manager.RegisterServer("svc", noopServer{}); err != nil {
		t.Fatalf("register server: %v", err)
	}
	if err := manager.RegisterServer("svc", noopServer{}); err == nil {
		t.Fatal("expected duplicate server registration to fail")
	}
}

func TestInitServerRequiresRegistration(t *testing.T) {
	manager := NewManager()
	if _, err := manager.InitServer("missing", "", ""); err == nil {
		t.Fatal("expected init server without registration to fail")
	}
}

func TestRegisterClientAndCloseWithoutTLS(t *testing.T) {
	manager := NewManager()
	address := startTestServer(t, manager, "", "")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := manager.RegisterClient(ctx, "client", address, ""); err != nil {
		t.Fatalf("register client: %v", err)
	}

	conn, err := manager.GetClient("client")
	if err != nil {
		t.Fatalf("get client: %v", err)
	}
	if conn == nil {
		t.Fatal("expected client connection")
	}

	if err := manager.Close(); err != nil {
		t.Fatalf("close manager: %v", err)
	}
	if _, err := manager.GetClient("client"); err == nil {
		t.Fatal("expected closed client to be removed")
	}
}

func TestRegisterClientWithTLS(t *testing.T) {
	certFile, keyFile := writeSelfSignedCert(t)

	manager := NewManager()
	address := startTestServer(t, manager, keyFile, certFile)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := manager.RegisterClient(ctx, "secure-client", address, certFile); err != nil {
		t.Fatalf("register tls client: %v", err)
	}

	if _, err := manager.GetClient("secure-client"); err != nil {
		t.Fatalf("get tls client: %v", err)
	}
}

func TestRegisterClientRejectsDuplicate(t *testing.T) {
	manager := NewManager()
	address := startTestServer(t, manager, "", "")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := manager.RegisterClient(ctx, "client", address, ""); err != nil {
		t.Fatalf("register client: %v", err)
	}
	if err := manager.RegisterClient(ctx, "client", address, ""); err == nil {
		t.Fatal("expected duplicate client registration to fail")
	}
}

func startTestServer(t *testing.T, manager *Manager, keyFile string, certFile string) string {
	t.Helper()

	if err := manager.RegisterServer("svc", noopServer{}); err != nil {
		t.Fatalf("register server: %v", err)
	}

	server, err := manager.InitServer("svc", keyFile, certFile)
	if err != nil {
		t.Fatalf("init server: %v", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	go func() {
		_ = server.Serve(listener)
	}()

	t.Cleanup(func() {
		server.GracefulStop()
		_ = listener.Close()
		_ = manager.Close()
	})

	return listener.Addr().String()
}

func writeSelfSignedCert(t *testing.T) (certFile string, keyFile string) {
	t.Helper()

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate private key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "127.0.0.1",
		},
		NotBefore:             time.Now().Add(-time.Minute),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, privateKey.Public(), privateKey)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		t.Fatalf("marshal private key: %v", err)
	}

	dir := t.TempDir()
	certFile = filepath.Join(dir, "server.crt")
	keyFile = filepath.Join(dir, "server.key")

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privateKeyBytes})

	if err := os.WriteFile(certFile, certPEM, 0o600); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyFile, keyPEM, 0o600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	return certFile, keyFile
}
