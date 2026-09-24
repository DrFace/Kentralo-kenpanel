package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DrFace/Kentralo-kenpanel/core/config"
	"github.com/DrFace/Kentralo-kenpanel/pkg/protocol"
)

func main() {
	socketPath := flag.String("socket", "/run/kenpanel/agent.sock", "Path to Unix domain socket")
	listenAddr := flag.String("listen", "", "Optional mTLS TCP listen address for remote control plane")
	flag.Parse()

	fmt.Printf("============================================================\n")
	fmt.Printf(" KenPanel Privileged Node Daemon v%s\n", config.Version)
	fmt.Printf(" Role: %s\n", protocol.AgentRolePrivilegedDaemon)
	fmt.Printf(" Socket: %s\n", *socketPath)
	if *listenAddr != "" {
		fmt.Printf(" Remote mTLS: %s\n", *listenAddr)
	}
	fmt.Printf("============================================================\n")

	log.Printf("[INFO] KenPanel Node Agent initialized. Awaiting commands on %s...\n", *socketPath)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] Shutting down KenPanel Node Agent...")
}
