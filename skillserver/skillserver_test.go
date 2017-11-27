package skillserver

import (
	"testing"
	"crypto/x509"
	"time"
)

func TestIsCertExpired(t *testing.T) {
	cert := &x509.Certificate{
		NotBefore: time.Date(2017, time.January, 10, 0, 0, 0, 0, time.UTC),
		NotAfter: time.Date(2017, time.January, 12, 0, 0, 0, 0, time.UTC),
	}

	now := time.Date(2017, time.January, 14, 0, 0, 0, 0, time.UTC)

	if !isCertExpired(cert, now) {
		t.Error("Cert should have been expired")
	}
}

