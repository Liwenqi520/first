package model

type SysRole struct {
	ID             string `gorm:"column:id;type:varchar(30);not null;comment:主键ID;" json:"id"`
	Name           string `gorm:"column:name;type:varchar(50);not null;default:'';comment:角色名称;" json:"name"`
	CompanyID      string `gorm:"column:company_id;type:varchar(30);not null;default:'0';comment:所属企业ID,默认为0;" json:"company_id"`
	Sort           int64  `gorm:"column:sort;type:int(30);not null;default:0;comment:排序,默认为0;" json:"sort"`
	Status         int64  `gorm:"column:status;type:int(30);not null;default:1;comment:1:启用 2:禁用;" json:"status"`
	MenuPermission string `gorm:"column:menu_permission;type:varchar(2000);not null;default:'';comment:菜单权限ID集合;" json:"menu_permission"`
	AreaPermission string `gorm:"column:area_permission;type:longtext;not null;comment:厂区组织架构权限ID集合;" json:"area_permission"`
	DeptPermission string `gorm:"column:dept_permission;type:varchar(2000);not null;default:'';comment:部门权限ID集合;" json:"dept_permission"`
	AppIDs         string `gorm:"column:app_ids;type:varchar(2000);not null;default:'';comment:应用id集合;" json:"app_ids"`
	DataScope      int64  `gorm:"column:data_scope;type:int(30);not null;default:0;comment:数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）;" json:"data_scope"`
	Identification string `gorm:"column:identification;type:varchar(200);not null;default:'';comment:角色标识;" json:"identification"`
	Desc           string `gorm:"column:desc;type:varchar(200);not null;default:'';comment:角色说明;" json:"desc"`
	CreatedTime    int64  `gorm:"column:created_time;type:bigint(20);not null;default:0;comment:创建时间;" json:"created_time"`
	UpdatedTime    int64  `gorm:"column:updated_time;type:bigint(20);not null;default:0;comment:更新时间;" json:"updated_time"`
	DeletedTime    int64  `gorm:"column:deleted_time;type:bigint(20);not null;default:0;comment:删除时间;" json:"deleted_time"`
}

func init() {
	ModelList = append(ModelList, &SysRole{})
}

// TableName 表名
func (SysRole) TableName() string {
	return "sys_role"
}
