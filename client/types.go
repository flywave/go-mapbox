package client

import "time"

type ListStyle struct {
	Version  int64     `json:"version,omitempty"`
	Name     string    `json:"name,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	ID       string    `json:"id,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
	Owner    string    `json:"owner,omitempty"`
}
