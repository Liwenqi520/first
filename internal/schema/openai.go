package schema

type OpenAiParams struct {
	Message string `json:"message"`
}

type OpenAiRecharge struct {
	Amount int64  `json:"amount"`  // 充值金额【单位：分】
	UserID string `json:"user_id"` // 用户ID
}
