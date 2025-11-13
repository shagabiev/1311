package main

import "sync"

type CheckResult struct {
	URL       string `json:"url"`
	Available bool   `json:"available"`
}

type Task struct {
	ID      int           `json:"id"`
	URLs    []string      `json:"urls"`
	Results []CheckResult `json:"results"`
	Done    bool          `json:"done"`
	mu      sync.RWMutex
}

type CheckRequest struct {
	Links []string `json:"links"`
}

type ReportRequest struct {
	TaskIDs []int `json:"links_num"`
}
