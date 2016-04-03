package spyder

import (
	"testing"
	"os"
)

func TestSession(t *testing.T) {
	connString := os.Getenv("MONGODB_URL")
	GetSession(connString)
}

