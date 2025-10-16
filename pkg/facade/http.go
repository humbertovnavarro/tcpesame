package facade

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strings"
	"time"
)

var serverHeaders = []string{
	"Apache/2.4.41 (Ubuntu)",
	"nginx/1.18.0 (Ubuntu)",
	"nginx/1.25.2",
	"Microsoft-IIS/10.0",
	"LiteSpeed",
	"Caddy",
	"Apache/2.2.34 (Win32)",
	"cloudflare",
}

var htmlTemplates = []string{
	`<html><head><title>Welcome</title></head><body><h1>It works!</h1></body></html>`,
	`<html><head><title>Index of /</title></head><body><h1>Index of /</h1><hr><pre><a href="/">Parent Directory</a></pre><hr></body></html>`,
	`<html><head><title>Welcome to nginx!</title></head><body><center><h1>Welcome to nginx!</h1></center><hr><center>nginx</center></body></html>`,
	`<html><head><title>403 Forbidden</title></head><body><h1>Forbidden</h1><p>You don't have permission to access this resource.</p><hr></body></html>`,
	`<html><head><title>Under Construction</title></head><body><h1>Site under maintenance</h1></body></html>`,
}

func HTTPFacade(conn net.Conn) {
	maxLineSize := int64(16384)
	limited := io.LimitReader(conn, maxLineSize)
	reader := bufio.NewReader(limited)
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	_, _, _, ok := parseRequestLine(line)
	if !ok {
		return
	}
	for {
		l, err := reader.ReadString('\n')
		if err != nil || l == "\r\n" {
			break
		}
	}
	server := serverHeaders[rand.Intn(len(serverHeaders))]
	body := htmlTemplates[rand.Intn(len(htmlTemplates))]
	date := time.Now().UTC().Format(time.RFC1123)
	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n"+
			"Server: %s\r\n"+
			"Date: %s\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"Content-Length: %d\r\n"+
			"Connection: %s\r\n"+
			"\r\n%s",
		server, date, len(body), "close", body,
	)
	slowWriter := NewSlowWriter(conn, time.Millisecond*500, time.Millisecond*1000, 1, 4)
	slowWriter.Write([]byte(response))
}

func parseRequestLine(line string) (method, path, version string, ok bool) {
	parts := strings.Fields(strings.TrimSpace(line))
	if len(parts) != 3 {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}
