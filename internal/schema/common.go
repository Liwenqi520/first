package schema

type Response struct {
	Code int64         `json:"code"`
	Msg  string        `json:"msg"`
	Data []interface{} `json:"data"`
}

type PageInfo struct {
	Page  uint `json:"page" form:"page"`
	Limit uint `json:"limit" form:"limit"`
}

type Int64PageInfo struct {
	Page  int64 `json:"page" form:"page"`
	Limit int64 `json:"limit" form:"limit"`
}

type PageInfoSpecial struct {
	Page  uint `json:"currentPage" form:"currentPage"`
	Limit uint `json:"pageSize" pageSize:"pageSize"`
}

type PageResult struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  uint        `json:"currentPage"`
	Limit uint        `json:"pageSize"`
}

type PageResultAndTotalPage struct {
	List      interface{} `json:"list"`
	Total     int64       `json:"total"`
	Page      int64       `json:"currentPage"`
	TotalPage int64       `json:"totalPage"`
	Limit     int64       `json:"pageSize"`
}

type FactoryPageResult struct {
	List        interface{} `json:"list"`
	CompanyName string      `json:"companyName"`
	CompanyLogo string      `json:"companyLogo"`
	Total       int64       `json:"total"`
	Page        uint        `json:"currentPage"`
	Limit       uint        `json:"pageSize"`
}

// HandleDeviceLogSaveParams 设备操作记录请求参数
type HandleDeviceLogSaveParams struct {
	CompanyID   string `json:"company_id"`
	FactoryID   string `json:"factory_id"`   // 厂区ID
	WorkShopID  string `json:"worksho_id"`   // 车间ID
	OperationID string `json:"operation_id"` // 业务ID
	DeviceID    string `json:"device_id"`
	HandleType  int64  `json:"handle_type"` // 1.添加设备 2:修改设备 3:删除设备 4:插入设备 5:拔出设备
	Content     string `json:"content"`
	AdminID     string `json:"admin_id"`
}

const (
	AddDevice     = 1
	ModifyDevice  = 2
	DeleteDevice  = 3
	InsertDevice  = 4
	ExtractDevice = 5
)

var HandleDeviceTypeName = map[int64]string{
	1: "add_device",
	2: "modify_device",
	3: "delete_device",
	4: "insert_device",
	5: "extract_device",
}

func NewTable(page uint, limit uint, defaultPage uint, defaultLimit uint) *PageResult {
	if defaultPage == 0 {
		defaultPage = 1
	}
	if defaultLimit == 0 {
		defaultLimit = 10
	}

	table := &PageResult{
		Page:  defaultPage,
		Limit: defaultLimit,
		Total: 0,
		List:  []interface{}{},
	}
	if page > 0 {
		table.Page = page
	}
	if limit > 0 {
		if limit > 5000 {
			table.Limit = 5000
		} else {
			table.Limit = limit
		}
	}
	return table
}

func NewTotalPageTable(page int64, limit int64, defaultPage int64, defaultLimit int64) *PageResultAndTotalPage {
	if defaultPage == 0 {
		defaultPage = 1
	}
	if defaultLimit == 0 {
		defaultLimit = 10
	}

	table := &PageResultAndTotalPage{
		Page:  defaultPage,
		Limit: defaultLimit,
		Total: 0,
		List:  []interface{}{},
	}
	if page > 0 {
		table.Page = page
	}
	if limit > 0 {
		if limit > 5000 {
			table.Limit = 5000
		} else {
			table.Limit = limit
		}
	}
	return table
}

func NewFactoryTable(page uint, limit uint, defaultPage uint, defaultLimit uint) *FactoryPageResult {
	if defaultPage == 0 {
		defaultPage = 1
	}
	if defaultLimit == 0 {
		defaultLimit = 1000
	}

	table := &FactoryPageResult{
		Page:        defaultPage,
		Limit:       defaultLimit,
		Total:       0,
		List:        []interface{}{},
		CompanyName: "",
	}
	if page > 0 {
		table.Page = page
	}
	if limit > 0 {
		if limit > 5000 {
			table.Limit = 5000
		} else {
			table.Limit = limit
		}
	}
	return table
}
