package schema

import "first/internal/model"

type RSA struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type UserRes struct {
	Token    string `json:"token"`
	UserInfo struct {
		ID               string `json:"id"`
		NickName         string `json:"nick_name"`
		CompanyID        string `json:"company_id"`
		Mobile           string `json:"mobile"`
		Email            string `json:"email"`
		Image            string `json:"image"`
		IsAdmin          int8   `json:"is_admin"`
		IsShowCompanyTab int8   `json:"is_show_compay_tab"`
		Permission       string `json:"permission"`
	} `json:"user_info"`
}

type ResetMobileReq struct {
	Mobile      string `json:"mobile"`
	VerifyToken string `json:"verify_token"`
	NewMobile   string `json:"new_mobile"`
	VerifyCode  string `json:"verify_code"`
}

type AccountCodeCheck struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

// UserSaveParams 添加用户参数
type UserSaveParams struct {
	AdminID   string `json:"admin_id"`   // 登录管理员ID
	CompanyID string `json:"company_id"` // 企业ID
	UserID    string `json:"id"`         // 用户ID

	UserName       string `json:"username"`       // 登录账号
	Sex            int8   `json:"sex"`            // 性别
	NickName       string `json:"nickname"`       // 姓名
	Password       string `json:"password"`       // 登录密码
	Position       string `json:"position"`       // 职位
	RoleID         string `json:"roleIds"`        // 角色
	AreaPermission string `json:"areaPermission"` // 区域权限
	Mobile         string `json:"phone"`          // 联系方式
	Email          string `json:"email"`          // 邮箱
	Remark         string `json:"remark"`         // 备注
	Status         int8   `json:"status"`
	DepartmentID   string `json:"deptId"` // 部门ID
}

// UserDeleteParams 用户删除参数
type UserDeleteParams struct {
	ID        string `json:"id"`
	CompanyID string `json:"company_id"`
	AdminID   string `json:"admin_id"`
}

// UserSearchParams 用户搜索参数
type UserSearchParams struct {
	DepartmentID string `json:"deptId"` // 部门ID
	RoleId       string `json:"role_id"`
	Position     string `json:"position"`
	Status       int8   `json:"status"`
	Keyword      string `json:"keyword"`
	Name         string `json:"username"`
	Mobile       string `json:"phone"`
	PageInfo
}

// StoreInfo redis存储的参数
type StoreInfo struct {
	ID             string `json:"id"`
	RoleID         string `json:"role_id"`
	RoleType       int8   `json:"role_type"`
	CompanyID      string `json:"company_id"`
	NickName       string `json:"nick_name"`
	UserName       string `json:"user_name"`
	Image          string `json:"image"`
	Mobile         string `json:"mobile"`
	Email          string `json:"email"`
	Status         int8   `json:"status"`
	IsAdmin        int8   `json:"is_admin"`
	AreaPermission string `json:"area_permission"`
	PlatformType   int8   `json:"platform_type"`  // 1酿造 2酒坛 3两者
	HavePlatforms  string `json:"have_platforms"` // 拥有的平台ID集合，逗号分割
	// FactorgGroupManageButton int64  `json:"factory_group_manage_button"`
	// FactorgGroupSwitch       int64  `json:"factory_group_switch"`
	// WaterQualitySwitch       int64  `json:"water_quality_switch"`
}

// Profile 用户信息
type Profile struct {
	model.AuthUser
	RoleName    string `gorm:"-" json:"role_name"`
	CompanyName string `gorm:"-" json:"company_name"`
	// AreaPermission []RoleAreaList `json:"area_permission"`
}

// LoginParams 登录参数
type LoginParams struct {
	Type       int64  `json:"type"` // 1:账号密码 2:手机验证码 3:邮箱登录
	Account    string `json:"username"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	VerifyCode string `json:"verifyCode"`
}

type SendCodeParams struct {
	Account string `json:"account"`
	Type    int64  `json:"type"` // 2：手机 3：邮箱
}

// LoginRes 登录响应参数
type LoginRes struct {
	Token        string   `json:"token"`
	UserName     string   `json:"username"`
	NickName     string   `json:"nickname"`
	Image        string   `json:"image"`
	Active       bool     `json:"active"`
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	Expires      string   `json:"expires"`
	Avatar       string   `json:"avatar"`
	Roles        []string `json:"role"`
	UserInfo     struct {
		ID             string `json:"id"`
		RoleID         string `json:"role_id"`
		RoleType       int8   `json:"role_type"`
		CompanyID      string `json:"company_id"`
		NickName       string `json:"nick_name"`
		UserName       string `json:"user_name"`
		Image          string `json:"image"`
		Mobile         string `json:"mobile"`
		Email          string `json:"email"`
		Status         int8   `json:"status"`
		IsAdmin        int8   `json:"is_admin"`
		PlatformType   int8   `json:"platform_type"` // 1酿造 2酒坛 3两者
		AreaPermission string `json:"area_permission"`
		HavePlatforms  string `json:"have_platforms"` // 拥有的平台ID集合，逗号分割
		// FactorgGroupManageButton int64  `json:"factory_group_manage_button"`
		// FactorgGroupSwitch       int64  `json:"factory_group_switch"`
		// WaterQualitySwitch       int64  `json:"water_quality_switch"`
	} `json:"user_info"`
	MenuList    []*MenuInfoTree `json:"menu_list,omitempty"`
	Permissions []string        `json:"permissions,omitempty"`
}

type MenuList struct {
	UserName        string          `json:"username"`
	AreaRoutes      []MenuDetail    `json:"areaRoutes"`
	DeviceRoutes    []MenuDetail    `json:"deviceRoutes"`
	DefaultRoutes   []*MenuInfoTree `json:"defaultRoutes"`
	DeviceManRoutes []MenuDetail    `json:"deviceManRoutes"`
}

type MenuDetail struct {
	ID           string       `json:"id"`
	Path         string       `json:"path"`
	Name         string       `json:"name"`
	ShowChildren bool         `json:"showChildren"`
	Component    string       `json:"component,omitempty"`
	Redirect     string       `json:"redirect,omitempty"`
	Meta         Meta         `json:"meta,omitempty"`
	Props        Props        `json:"props,omitempty"`
	Children     []MenuDetail `json:"children,omitempty"`
}

type Meta struct {
	Title       string   `json:"title"`
	Icon        string   `json:"icon,omitempty"`
	Rank        int64    `json:"rank,omitempty"`
	ShowParent  bool     `json:"showParent"`
	Auths       []string `json:"auths,omitempty"`
	IsKeepAlive bool     `json:"keepAlive"`
	ShowLink    bool     `json:"showLink"`
}

type Props struct {
	ID       string `json:"id,omitempty"`
	IsHidden bool   `json:"isHidden,omitempty"`
}

type MenuInfo struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Path      string   `json:"path"`
	PID       string   `json:"pid"`
	Component string   `json:"component"`
	Meta      MetaInfo `json:"meta"`
}

type MenuInfoTree struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Path      string          `json:"path"`
	PID       string          `json:"pid"`
	Component string          `json:"component"`
	Meta      MetaInfo        `json:"meta"`
	Children  []*MenuInfoTree `json:"children,omitempty"`
}

type MetaInfo struct {
	Icon        string   `json:"icon"`
	IsAffix     bool     `json:"fixedTag"`
	IsHide      bool     `json:"hiddenTag"`
	IsIframe    bool     `json:"isIframe"`
	IsKeepAlive bool     `json:"keepAlive"`
	IsLink      bool     `json:"showLink"` // 对应 linke_url
	Title       string   `json:"title"`
	ShowParent  bool     `json:"showParent"`
	Rank        int64    `json:"rank,omitempty"`  // 排序
	Auths       []string `json:"auths,omitempty"` // 权限标识
}

type LoginThirdRes struct {
	Token    string `json:"token"`
	UserInfo struct {
		UserName string `json:"username"`
	} `json:"user_info"`
}

// UserChangeRole 用户切换角色请求参数
type UserChangeRole struct {
	RoleID  string `json:"role_id"`
	UserIDs string `json:"user_ids"`
}

type RoleID struct {
	RoleID   string `json:"id"`
	MenuType int64  `json:"type"`
}

type AesTestParams struct {
	Str string `json:"str"`
}

type RandNumParams struct {
	Max   float64 `json:"max"`
	Min   float64 `json:"min"`
	Digit int64   `json:"digit"`
}

type UserDivisionRoleParams struct {
	UserID string `json:"id"`
	RoleID string `json:"role_id"`
}

type TestParams struct {
	DepartmentID string `json:"department_id"`
	CompanyID    string `json:"company_id"`
}

type BrifeUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UpdateUserImageParams struct {
	ID    string `json:"id"`
	Image string `json:"image"`
	Lang  string `json:"lang"`
}

type CheckUserPhoneParams struct {
	ID    string `json:"id"`       // 用户ID
	Phone string `json:"oldPhone"` // 手机号
	Code  string `json:"code"`     // 验证码
}

type UpdateUserPhoneParams struct {
	ID    string `json:"id"`       // 用户ID
	Phone string `json:"newPhone"` // 手机号
	Code  string `json:"code"`     // 验证码
}

type CheckUserEmailParams struct {
	ID       string `json:"id"`       // 用户ID
	OldEmail string `json:"oldEmail"` // 旧邮箱
	NewEmail string `json:"newEmail"` // 新邮箱
	Code     string `json:"code"`     // 验证码
}

type AuthParams struct {
	UserID string `json:"user_id"` // 加密的
}

type CheckTokenParams struct {
	Token string `json:"token"`
	Lang  string `json:"lang"`
}

type SearchUserListParams struct {
	Keyword string `json:"keyword"`
	PageInfo
}

type UpdateUserInfoParams struct {
	Nickname string `json:"nickname"`
	Lang     string `json:"lang"`
	// Image    string `json:"image"`
}

type UserSlice struct {
	UserCode string `json:"userCode"`
	//用户名
	UserName string `json:"userName"`
	//用户密码
	PassWord string `json:"passWord"`
	//昵称
	NickName string `json:"nickName"`
	//用户状态
	UserStatus int `json:"userStatus"`
}

type CheckQrcodeResParams struct {
	Status   string      `json:"status"`
	Token    string      `json:"token"`
	LoginRes interface{} `json:"login_res"`
}
