package tor

import (
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func Start(mu *sync.Mutex) error {
	log.Info("Started tor service.")
	mu.Lock()
	defer mu.Unlock()

	cmd := exec.Command("tor", "--RunAsDaemon", "1")
	err := cmd.Run()
	if err == nil {
		go handleClose(mu)
	}
	return err
}

func Stop(mu *sync.Mutex) error {
	log.Info("Stopped tor service.")
	mu.Lock()
	defer mu.Unlock()

	cmd := exec.Command("killall", "tor")
	err := cmd.Run()
	return err
}

func Restart(mu *sync.Mutex) error {
	log.Info("Restarting tor service...")
	mu.Lock()
	defer mu.Unlock()
	cmd := exec.Command("pkill", "-9", "-f", "tor")
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("tor", "--RunAsDaemon", "1")
	err := cmd.Run()
	return err
}

func handleClose(mu *sync.Mutex) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	for sig := range c {
		switch sig {
		case syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM:
			Stop(mu)
			return
		}
	}
}
