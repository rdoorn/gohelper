package signaling

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignals(t *testing.T) {

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	s := New()
	defer s.Close()

	s.Add(func() {
		assert.Nil(t, nil)
		os.Exit(0)
	}, syscall.SIGHUP)

	time.Sleep(1 * time.Second)
	proc.Signal(syscall.SIGHUP)
	time.Sleep(1 * time.Second)

	assert.Fail(t, "not ok")
}
