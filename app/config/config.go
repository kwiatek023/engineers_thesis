package config

type JsonGraphStructure struct {
	Graph struct {
		NofVertices int `json:"nofVertices"`
		Edges       []struct {
			Edge        []int   `json:"edge"`
			Reliability float64 `json:"reliability"`
		} `json:"edges"`
	} `json:"graph"`
}
