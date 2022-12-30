package api

type Response[T any] struct {
	Object  string `json:"object"`
	HasMore bool   `json:"has_more"`
	Data    T      `json:"data"`
}
