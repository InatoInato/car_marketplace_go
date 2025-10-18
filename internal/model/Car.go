package model

type Car struct{
	Make string `json:"make"`
	Model string `json:"model"`
	Year int `json:"year"`
	Color string `json:"color"`
	EngineCapacity float64 `json:"engine_capacity"`
}
