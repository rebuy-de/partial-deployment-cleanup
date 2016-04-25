package kv

import "time"

type Deployments []Deployment

type Deployment struct {
	Branch    string    `json:"branch"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Distribution map[int]string
