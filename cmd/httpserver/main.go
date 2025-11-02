package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func getPort() uint16 {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		return 8080
	}
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		log.Printf("Invalid PORT env var, using default port 8080")
		return 8080
	}
	return uint16(port)
}

var port = getPort()

func main() {
	s, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		body := respond200()
		h := response.GetDefaultHeaders(0)
		status := response.StatusOK
		endpoint := req.RequestLine.RequestTarget
		if endpoint == "/yourproblem" {
			status = response.StatusBadRequest
			body = respond400()
		} else if endpoint == "/myproblem" {
			status = response.StatusInternalServerError
			body = respond500()
		} else if strings.HasPrefix(endpoint, "/httpbin/") {
			res, err := http.Get("https://httpbin.org/" + endpoint[len("/httpbin/"):])
			if err != nil {
				body = respond500()
				status = response.StatusInternalServerError
			} else {
				w.WriteStatusLine(response.StatusOK)
				h.Delete("Content-length")
				h.Set("Transfer-Encoding", "chunked")
				h.Replace("Content-Type", "text/plain")
				h.Set("Trailer", "X-Content-SHA256")
				h.Set("Trailer", "X-Content-Length")
				w.WriteHeaders(h)

				fullBody := []byte{}
				for {
					data := make([]byte, 32)
					n, err := res.Body.Read(data)
					if err != nil {
						break
					}
					fullBody = append(fullBody, data[:n]...)
					w.WriteBody(fmt.Appendf(nil, "%x\r\n", n))
					w.WriteBody(data[:n])
					w.WriteBody([]byte("\r\n"))
				}
				w.WriteBody([]byte("0\r\n"))
				trailers := headers.NewHeaders()
				sha := sha256.Sum256(fullBody)
				trailers.Set("X-Content-SHA256", toStr(sha[:]))
				trailers.Set("X-Content-Length", strconv.Itoa(len(fullBody)))
				w.WriteHeaders(trailers)
				return
			}
		} else if endpoint == "/video" {
			f, err := os.ReadFile("assets/vim.mp4")
			if err != nil {
				body = respond500()
				status = response.StatusInternalServerError
				return
			}
			h.Replace("Content-type", "video/mp4")
			h.Replace("Content-length", strconv.Itoa(len(f)))
			w.WriteStatusLine(status)
			w.WriteHeaders(h)
			w.WriteBody(f)
			return
		} else if endpoint == "/json" {
			body = respondJSON()
			h.Replace("Content-type", "application/json")
			h.Replace("Content-length", strconv.Itoa(len(body)))
			w.WriteStatusLine(status)
			w.WriteHeaders(h)
			w.WriteBody(body)
			return
		}
		h.Replace("Content-length", strconv.Itoa(len(body)))
		h.Replace("Content-type", "text/html")

		w.WriteStatusLine(status)
		w.WriteHeaders(h)
		w.WriteBody(body)
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer s.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func respond200() []byte {
	return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
}

func respond400() []byte {
	return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
}

func respond500() []byte {
	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}

func respondJSON() []byte {
	data := response.JsonData{Message: "Success"}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	return jsonBytes
}

func toStr(b []byte) string {
	r := ""
	for _, e := range b {
		r += fmt.Sprintf("%02x", e)
	}
	return r
}
