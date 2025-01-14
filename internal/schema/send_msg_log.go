package schema

type SendMsgLogSaveParams struct {
	CompanyID string `json:"company_id"`
	SendType  int8   `json:"send_type"`
	Target    string `json:"target"`
	Content   string `json:"content"`
	Status    int8   `json:"status"`
	Res       string `json:"res"`
}
