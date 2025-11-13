package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	srv    *http.Server
	client *http.Client
}

func NewServer() *Server {
	mux := http.NewServeMux()
	s := &Server{
		srv: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
		client: &http.Client{
			Timeout: 15 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}

	mux.HandleFunc("/check", s.checkHandler)
	mux.HandleFunc("/report", s.reportHandler)
	return s
}

func (s *Server) checkHandler(w http.ResponseWriter, r *http.Request) {
	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if len(req.Links) == 0 {
		http.Error(w, "No links provided", http.StatusBadRequest)
		return
	}

	results := make(map[string]string)
	client := s.client

	for _, raw := range req.Links {
		url := normalizeURL(raw)

		resp, err := client.Head(url)
		if err != nil || resp.StatusCode == http.StatusMethodNotAllowed {
			resp, err = client.Get(url)
		}
		if err != nil {
			results[raw] = "not available"
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			results[raw] = "available"
		} else {
			results[raw] = "not available"
		}
	}

	// Сохраняем задачу
	taskID := NextTaskID()
	task := &Task{
		ID:   taskID,
		URLs: req.Links,
		Done: true,
	}
	for u, st := range results {
		task.Results = append(task.Results, CheckResult{
			URL:       normalizeURL(u),
			Available: st == "available",
		})
	}
	SaveTask(task)

	resp := map[string]any{
		"links":     results,
		"links_num": taskID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) reportHandler(w http.ResponseWriter, r *http.Request) {
	var req ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	tasks := GetTasks(req.TaskIDs)
	if len(tasks) == 0 {
		http.Error(w, "No tasks found", http.StatusNotFound)
		return
	}

	pdfData := GeneratePDF(tasks)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="report.pdf"`)
	w.Write(pdfData)
}

func normalizeURL(s string) string {
	s = strings.TrimSpace(s)
	if !strings.Contains(s, "://") {
		s = "https://" + s
	}
	return s
}

func (s *Server) Start() error {
	log.Println("Server starting on :8080")
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Graceful shutdown...")
	return s.srv.Shutdown(ctx)
}
