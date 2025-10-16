package facade

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

var sshVersions = []string{
	"7.6p1", "7.9p1", "8.0p1", "8.2p1", "8.4p1", "8.6p1", "8.9p1", "9.0p1", "9.3p1",
}

var osTags = []string{
	"Ubuntu-4ubuntu0.3",
	"Debian-1",
	"Debian-5",
	"Raspbian-2",
	"Alpine-3",
	"Fedora-1",
	"FreeBSD",
	"OpenBSD",
	"Arch",
}

func generateRandomBanner() string {
	version := sshVersions[rand.Intn(len(sshVersions))]
	osTag := osTags[rand.Intn(len(osTags))]
	return fmt.Sprintf("SSH-2.0-OpenSSH_%s %s\r\n", version, osTag)
}

func SSHFacade(conn net.Conn) {
	var banner = generateRandomBanner()
	var slowWriter = NewSlowWriter(conn, time.Millisecond*100, time.Millisecond*1000, 1, 128)
	slowWriter.Write([]byte(banner))
}
