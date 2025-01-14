package model

//AuthUser 后台用户
type AuthUser struct {
	ID             string `gorm:"column:id;type:varchar(30);NOT NULL;comment:主键;" json:"id"`
	RoleID         string `gorm:"column:role_id;type:varchar(50);not null;default:'';comment:用户角色ID;" json:"roleId"`
	CompanyID      string `gorm:"column:company_id;type:varchar(30);not null;default:'';comment:所属公司;" json:"companyId"`
	Sex            int8   `gorm:"column:sex;type:tinyint(2);not null;default:1;comment:性别 1：男 2：女;" json:"sex"`
	UserName       string `gorm:"column:user_name;type:varchar(50);NOT NULL;default:'';comment:用户账号;index" json:"username"`
	NickName       string `gorm:"column:nick_name;type:varchar(50);NOT NULL;default:'';comment:用户昵称;" json:"nickname"`
	Password       string `gorm:"column:password;type:varchar(64);not null;default:'';comment:密码;" json:"password"`
	Image          string `gorm:"column:image;type:varchar(150);not null;default:'';comment:用户头像地址;" json:"avatar"`
	Position       string `gorm:"column:position;type:varchar(50);not null;default:'';comment:职位;" json:"position"`
	MobileAreaCode string `gorm:"column:mobile_area_code;type:varchar(10);not null;default:+86;comment:手机区号信息;index" json:"mobile_area_code"`
	Mobile         string `gorm:"column:mobile;type:varchar(20);not null;default:'';comment:手机号;index" json:"phone"`
	Email          string `gorm:"column:email;type:varchar(128);not null;default:'';comment:用户邮箱;index" json:"email"`
	DepartmentID   string `gorm:"column:department_id;type:varchar(30);not null;default:'';comment:所属部门;" json:"deptId"`
	Status         int8   `gorm:"column:status;type:tinyint(2);not null;default:2;comment:启用状态 -1已删除 1.启用 2:禁用;" json:"status"`
	IsAdmin        int8   `gorm:"column:is_admin;type:tinyint(1);not null;default:0;comment:是否为管理员 1是 0否;" json:"isAdmin"`
	// IsThirdPlatform int8   `gorm:"column:is_third_platform;type:tinyint(1);not null;default:0;comment:是否是第三方开发者 1是 0否;" json:"is_third_platform"`
	AreaPermission  string `gorm:"column:area_permission;type:varchar(3000);not null;default:'';comment:区域权限ID集合;" json:"areaPermission"`
	UserType        int8   `gorm:"column:user_type;type:tinyint(2);not null;default:1;comment:用户类型 1.平台用户 2:班组用户[操作带屏网关] 3:部门负责人;" json:"user_type"`
	TeamID          string `gorm:"column:team_id;type:varchar(128);not null;default:'';comment:班组ID;" json:"team_id"`
	PlatformType    int8   `gorm:"column:platform_type;type:tinyint(2);not null;default:1;comment:平台用户类型 1.酿造管理员 2.酒坛管理员 3.都有;" json:"platform_type"`
	IsProgramManage int64  `gorm:"column:is_program_manage;type:tinyint(1);not null;default:0;comment:是否是小程序管理员 1是 0否;" json:"is_program_manage"`
	ProgramForbit   int64  `gorm:"column:program_forbit;type:tinyint(1);not null;default:0;comment:是否被禁用 1是 0否;" json:"program_forbit"`
	Openid          string `gorm:"column:openid;type:varchar(100);not null;default:'';comment:用户的openid;" json:"openid"`
	HavePlatforms   string `gorm:"column:have_platforms;type:varchar(3000);not null;default:'';comment:拥有的系统集合;" json:"have_platforms"`
	Remark          string `gorm:"column:remark;type:varchar(300);not null;default:'';comment:备注;" json:"remark"`
	CreatedID       string `gorm:"column:created_id;type:varchar(30);not null;default:'';comment:创建人ID;" json:"createdId"`
	DeletedID       string `gorm:"column:deleted_id;type:varchar(30);not null;default:'';comment:删除人;" json:"deletedId"`
	CreatedTime     int64  `gorm:"column:created_time;type:bigint(20);not null;default:0;comment:创建时间;" json:"createdTime"`
	UpdatedTime     int64  `gorm:"column:updated_time;type:bigint(20);not null;default:0;comment:更新时间;" json:"updatedTime"`
	DeletedTime     int64  `gorm:"column:deleted_time;type:bigint(20);not null;default:0;comment:删除时间;" json:"deletedTime"`
}

func init() {
	ModelList = append(ModelList, &AuthUser{})
}

// TableName 表名
func (AuthUser) TableName() string {
	return "auth_user"
}
