package schema

type HandleMacheineParams struct {
	ID   string `json:"id"`
	Type int64  `json:"type"`
}

type HandleMacheineRes struct {
	Flag int64 `json:"flag"`
	Data Item  `json:"data"`
}

type Item struct {
	Type string `json:"type"`
	Time int64  `json:"time"`
}
