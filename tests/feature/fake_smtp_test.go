package feature

import (
	"bufio"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func startFakeSMTP(t *testing.T) (int, <-chan string) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = listener.Close() })

	port := listener.Addr().(*net.TCPAddr).Port
	messages := make(chan string, 1)

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(conn)
		writeSMTPLine(writer, "220 localhost Axolotl SMTP test")

		inData := false
		var data strings.Builder

		for {
			line, err := reader.ReadString(byte(10))
			if err != nil {
				return
			}
			line = strings.TrimRightFunc(line, func(r rune) bool {
				return r == 13 || r == 10
			})

			if inData {
				if line == "." {
					messages <- data.String()
					writeSMTPLine(writer, "250 queued")
					inData = false
					continue
				}
				data.WriteString(line)
				data.WriteByte(byte(10))
				continue
			}

			upper := strings.ToUpper(line)
			switch {
			case strings.HasPrefix(upper, "EHLO"):
				separator := string([]byte{13, 10})
				_, _ = writer.WriteString("250-localhost" + separator + "250 OK" + separator)
				_ = writer.Flush()
			case strings.HasPrefix(upper, "HELO"):
				writeSMTPLine(writer, "250 localhost")
			case strings.HasPrefix(upper, "MAIL FROM:"):
				writeSMTPLine(writer, "250 OK")
			case strings.HasPrefix(upper, "RCPT TO:"):
				writeSMTPLine(writer, "250 OK")
			case upper == "DATA":
				writeSMTPLine(writer, "354 End data")
				inData = true
			case upper == "QUIT":
				writeSMTPLine(writer, "221 bye")
				return
			default:
				writeSMTPLine(writer, "250 OK")
			}
		}
	}()

	return port, messages
}

func writeSMTPLine(writer *bufio.Writer, line string) {
	_, _ = writer.WriteString(line + string([]byte{13, 10}))
	_ = writer.Flush()
}
