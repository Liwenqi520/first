package middleware

import (
	"first/init/mysql"
	"first/internal/app/logic"
	"first/internal/common"
	"first/internal/model"
	"first/internal/schema"

	"strings"

	"github.com/gin-gonic/gin"
)

// PermissionMiddleware 权限中间件
func PermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserSession, _ := common.SessionGet(c)
		Path := c.Request.URL.Path                                  // 获取管理员当前访问的路由 eg: user/list
		if UserSession.RoleID != "0" && UserSession.RoleID != "1" { // 系统管理员直接放行
			right := PermissionCheck(UserSession.RoleID, Path, true)
			if !right {
				c.AbortWithStatusJSON(200,
					schema.Response{
						Code: 340001,
						Msg:  "暂无该功能模块的权限，请联系管理员",
					})
				return
			}
		}
		c.Next()
	}
}

// PermissionCheck 权限校验
func PermissionCheck(RoleID, Path string, IsSystemAdmin bool) bool {
	// RoleID也是全局唯一，我们就不判断所属主体了[RoleID可能是逗号分割的字符串，一个用户可能有多个角色]
	db := mysql.NewDB()
	SystemList := AllRouteList(1)  // 所有权限
	CompanyList := AllRouteList(2) // 企业权限
	CommonList := AllRouteList(3)  // 通用权限
	var RList []schema.RouteList
	if IsSystemAdmin {
		RList = append(RList, SystemList...)
	} else {
		RList = append(RList, CompanyList...)
	}
	var CommonStr []string
	for _, v := range CommonList {
		CommonStr = append(CommonStr, v.Url...)
	}
	RouteList := make(map[string][]string)
	for _, item := range RList {
		RouteList[item.ID] = item.Url
	}
	var Roles []model.SysRole
	RoleIDArr := strings.Split(RoleID, ",")
	db.Model(&model.SysRole{}).Where("id in (?) and deleted_time = 0", RoleIDArr).First(&Roles)
	if len(Roles) == 0 {
		return false
	}
	var PermissionList []string
	for _, role := range Roles {
		permissionStr := role.MenuPermission
		permissionStr = strings.TrimSpace(permissionStr)
		if permissionStr == "" {
			continue
		}
		PermissionList = append(PermissionList, strings.Split(permissionStr, ",")...) // 用户角色权限id集合
	}
	PermissionList = logic.Common{}.RemoveDuplicates(PermissionList) // 去重
	UserPermissonList := CommonStr

	AllPermission := GetMapKey(RouteList)

	for _, v := range PermissionList {
		if InArray(v, AllPermission) {
			UserPermissonList = append(UserPermissonList, RouteList[v]...)
		}
	}
	right := InArray(Path, UserPermissonList)
	return right
}

func AllRouteList(RouteType int64) []schema.RouteList {
	// 全部权限
	RouteAll := []schema.RouteList{
		{ID: "1", Name: "企业管理", Url: []string{}},
		{ID: "2", Name: "部门管理", Url: []string{"/sys-dept/list"}},
		{ID: "20", Name: "新增部门", Url: []string{"/sys-dept/add"}},
		{ID: "21", Name: "编辑部门", Url: []string{"/sys-dept/save"}},
		{ID: "22", Name: "删除部门", Url: []string{"/sys-dept/delete"}},
		{ID: "54", Name: "部门状态切换", Url: []string{"/sys-dept/switch"}},
		{ID: "3", Name: "人员管理", Url: []string{"/user/profile", "/company/profile", "/user/list", "/user/brife-user-list", "/user/get-user-menu-list", "/user/get-user-role-by-id"}},
		{ID: "23", Name: "新增人员", Url: []string{"/user/add"}},
		{ID: "24", Name: "编辑人员", Url: []string{"/user/save", "/user/update-user-image", "/user/modify-nick-name", "/user/check-user-phone", "/user/update-user-phone", "/user/check-user-email", "/user/update-user-email", "/user/send-code", "/company/update-company-logo", "/company/edit-info"}},
		{ID: "25", Name: "删除人员", Url: []string{"/user/delete"}},
		{ID: "43", Name: "人员状态切换", Url: []string{"/user/switch"}},
		{ID: "44", Name: "用户角色分配", Url: []string{"/user/user-division-role"}},
		{ID: "45", Name: "重置密码", Url: []string{"/user/reset-password"}},

		{ID: "4", Name: "角色管理", Url: []string{"/role/role-brief-list", "/role/role-area-permission", "/role/get-user-by-role-id"}},
		{ID: "26", Name: "新增角色", Url: []string{"/sys-role/add"}},
		{ID: "27", Name: "编辑角色", Url: []string{"/sys-role/save"}},
		{ID: "28", Name: "删除角色", Url: []string{"/sys-role/delete"}},
		{ID: "29", Name: "角色权限分配", Url: []string{"/sys-role/division"}},
		{ID: "55", Name: "角色状态切换", Url: []string{"/sys-role/switch"}},

		{ID: "5", Name: "系统设置", Url: []string{}},
		{ID: "8", Name: "操作日志", Url: []string{"/log/list"}},
		{ID: "41", Name: "菜单管理", Url: []string{"/sys-auth-rule/list"}},
		{ID: "30", Name: "新增菜单", Url: []string{"/sys-auth-rule/add"}},
		{ID: "31", Name: "编辑菜单", Url: []string{"/sys-auth-rule/save"}},
		{ID: "32", Name: "删除菜单", Url: []string{"/sys-auth-rule/delete"}},
		{ID: "42", Name: "固件管理", Url: []string{"/app-version-config/config-list"}},
		{ID: "33", Name: "新增固件", Url: []string{"/app-version-config/config-add"}},
		{ID: "34", Name: "编辑固件", Url: []string{"/app-version-config/config-save"}},
		{ID: "35", Name: "删除固件", Url: []string{"/app-version-config/config-delete"}},

		{ID: "6", Name: "厂区管理", Url: []string{}},
		{ID: "7", Name: "厂区结构", Url: []string{"/factory/region", "/factory/get-room-stat-by-factory-id", "/factory/get-factory-asset", "/factory/get-factory-bar-chart", "/factory/get-factory-workshop-data", "/factory/factory-device-list", "factory/factory-device-count", "/workshop/get-operation-type", "/workshop/detail", "/workshop/get-workshop-brief", "/workshop/get-workshop-operation"}},
		{ID: "14", Name: "新增厂区", Url: []string{"/factory/add"}},
		{ID: "15", Name: "编辑厂区", Url: []string{"/factory/save"}},
		{ID: "16", Name: "删除厂区", Url: []string{"/factory/delete"}},
		{ID: "17", Name: "新增车间", Url: []string{"/workshop/add"}},
		{ID: "18", Name: "编辑车间", Url: []string{"/workshop/save"}},
		{ID: "19", Name: "删除车间", Url: []string{"/workshop/delete"}},

		{ID: "9", Name: "班组管理", Url: []string{"/team/list", "/team/manage-operation-list", "/team/team-manage-areas"}},
		{ID: "13", Name: "新增班组", Url: []string{"/team/add"}},
		{ID: "10", Name: "编辑班组", Url: []string{"/team/save"}},
		{ID: "11", Name: "删除班组", Url: []string{"/team/delete"}},
		{ID: "12", Name: "班组区域权限划分", Url: []string{"/team/team-save-manage-areas"}},

		{ID: "36", Name: "设备管理", Url: []string{"/device/list", "/device-type/list", "/brewing/get-device-info", "/overview/get-device-data", "/log/handle-device-log-list"}},
		{ID: "37", Name: "新增设备", Url: []string{"/device/batch-add-device"}},
		{ID: "38", Name: "编辑设备", Url: []string{"/device/save"}},
		{ID: "39", Name: "删除设备", Url: []string{"/device/delete"}},
		{ID: "40", Name: "设备数据导出", Url: []string{"/log/device-log-export"}},

		{ID: "1734511163548215", Name: "部门列表", Url: []string{"/sys-dept/list"}},
		{ID: "1734511219488188", Name: "人员列表", Url: []string{"/user/list"}},
		{ID: "1734511251753566", Name: "角色列表", Url: []string{"/sys-role/list"}},
		{ID: "1734511593951108", Name: "厂区查看", Url: []string{"/factory/list"}},
		{ID: "1734511613418104", Name: "车间查看", Url: []string{"/workshop/list"}},
		{ID: "1734511648674520", Name: "酒坛查看", Url: []string{"/wine-building/list"}},

		// 酒坛相关
		{ID: "56", Name: "组织架构", Url: []string{"/wine-building/append-winejar", "/wine-building/delete-winejar", "/wine-building/edit-winejar", "/wine-overview/overview-stat", "/wine-overview/overview", "/wine-overview/distribute", "/wine-overview/device-list", "/wine-overview/batch-add-device", "/wine-overview/device-update", "/wine-overview/device-delete", "/wine-overview/device-detail", "/wine-overview/room-info", "/wine-overview/summary", "/wine-overview/overview-stat", "/wine-overview/operation-overview", "/wine-overview/winejar-warning-list", "/wine-overview/get-winejar-list-by-room-id", "/wine-overview/winejar-detail", "/wine-overview/get-winejar-history-handle-logs", "/wine-overview/device-overview", "/wine-overview/building-device-distribute", "/wine-overview/operation-device-distribute", "/wine-overview/device-data-list", "/wine-overview/alarm-list"}},
		{ID: "1734403954047949", Name: "添加", Url: []string{"/wine-building/save", "/wine-building/save-room", "/wine-building/batch-save-room", "/wine-building/save-room-door", "/wine-building/save-building"}},
		{ID: "1734405134173801", Name: "编辑", Url: []string{"/wine-building/edit", "/wine-building/edit-room", "/wine-building/batch-edit-room", "/wine-building/edit-room-door", "/wine-building/edit-building"}},
		{ID: "1734405152242365", Name: "删除", Url: []string{"/wine-building/delete", "/wine-building/delete-operation", "/wine-building/delete-row", "/wine-building/delete-room", "/wine-building/edit-building"}},
		{ID: "59", Name: "报警记录", Url: []string{"/wine-overview/alarm-list"}},
		{ID: "60", Name: "设备总览", Url: []string{"/wine-overview/device-overview", "/wine-overview/building-device-distribute", "/wine-overview/operation-device-distribute"}},
		{ID: "61", Name: "设备数据", Url: []string{"/wine-overview/device-data-list"}},
		{ID: "1726812431113807", Name: "操作日志", Url: []string{"/log/list"}},

		{ID: "51", Name: "大数据管理", Url: []string{"/factory/get-factory-bar-chart", "/overview/one", "/overview/two", "/overview/three", "/overview/device-warning", "/weather/today-weather", "/factory/get-factory-asset", "/factory/get-factory-bar-chart", "/factory/get-factory-workshop-data", "/factory/factory-device-count", "/factory/get-factory-bar-chart", "/workshop/get-workshop-brief", "/overview/get-operation-by-workshop-id", "/device/get-device-total", "/overview/get-distribute-by-workshop-id", "/overview/get-distribute-by-operation-id", "/workshop/get-workshop-operation", "/overview/get-current-production-stage", "/brewing/get-area-device-data", "/brewing/get-device-info", "/overview/get-device-data", "overview/get-production-stage-list", "/brewing/get-brewing-stage-position-list", "/overview/get-brewing-history-chart-data", "/overview/get-yeastroom-overturn-log", "/brewing/get-area-device-data", "/get-brewing-yeastroom-position-list", "/device/get-area-device-stat", "/device-type/list", "/factory/factory-device-list"}},

		{ID: "46", Name: "企业列表", Url: []string{"/company/list"}},
		{ID: "47", Name: "添加企业", Url: []string{"/company/add"}},
		{ID: "48", Name: "编辑企业", Url: []string{"/company/save"}},
		{ID: "49", Name: "删除企业", Url: []string{"/company/delete"}},
		{ID: "50", Name: "企业状态切换", Url: []string{"/company/switch"}},

		{ID: "1734511163548215", Name: "部门列表", Url: []string{"/sys-dept/list"}},
		{ID: "1734511219488188", Name: "人员列表", Url: []string{"/user/list"}},
		{ID: "1734511251753566", Name: "角色列表", Url: []string{"/sys-role/list"}},
		{ID: "1734511593951108", Name: "厂区查看", Url: []string{"/factory/list"}},
		{ID: "1734511613418104", Name: "车间查看", Url: []string{"/workshop/list"}},
		{ID: "1734511648674520", Name: "酒坛查看", Url: []string{"/wine-building/list"}},
	}

	// 通用权限
	RouteCommon := []schema.RouteList{
		{ID: "51", Name: "大数据管理", Url: []string{"/factory/get-factory-bar-chart", "/overview/one", "/overview/two", "/overview/three", "/overview/device-warning", "/weather/today-weather", "/factory/get-factory-asset", "/factory/get-factory-bar-chart", "/factory/get-factory-workshop-data", "/factory/factory-device-count", "/factory/get-factory-bar-chart", "/workshop/get-workshop-brief", "/overview/get-operation-by-workshop-id", "/device/get-device-total", "/overview/get-distribute-by-workshop-id", "/overview/get-distribute-by-operation-id", "/workshop/get-workshop-operation", "/overview/get-current-production-stage", "/brewing/get-area-device-data", "/brewing/get-device-info", "/overview/get-device-data", "overview/get-production-stage-list", "/brewing/get-brewing-stage-position-list", "/overview/get-brewing-history-chart-data", "/overview/get-yeastroom-overturn-log", "/brewing/get-area-device-data", "/overview/get-brewing-yeastroom-position-list", "/overview/get-yeastroom-overturn-log-list", "/overview/get-production-stage-list", "/device/get-area-device-stat", "/device-type/list", "/factory/factory-device-list", "/winejar-app/winejar-device-insert-by-system", "/excel/make-device-excel", "/crontab-task/save"}},

		{ID: "51", Name: "用户角色", Url: []string{"/user/profile", "/company/profile", "/sys-role/list-all-role", "/sys-role/role-menu", "/sys-role/role-menu_ids", "/user/get-user-menu-list", "/mp/user-profile", "/user/update-user-image-base64", "/user/save", "/user/reset-password", "/user/update-user-phone", "/user/update-user-email"}},
		{ID: "52", Name: "设备管理", Url: []string{"/device/get-device-list-by-type", "/device/device-total-data", "/device/count-area-device", "/device/get-area-device-sta", "/device/get-factory-device-stat", "/log/handle-device-log-list", "/log/device-log-export"}},
		{ID: "53", Name: "厂区结构", Url: []string{"/device/get-factory-workshop-brief"}},
		{ID: "56", Name: "应用相关", Url: []string{"/application/list", "/application/app-list"}},
		{ID: "57", Name: "部门相关", Url: []string{"/sys-dept/list"}},
		{ID: "58", Name: "app相关", Url: []string{"/mp/device-calibrate-list"}},
		{ID: "59", Name: "小程序相关", Url: []string{"/program/get-device-info-and-data", "/program/get-device-history-data", "/program/device-list", "/program/user-bind-device", "/program/user-unbind-device", "/program/user-weather-station-list", "/program/update-nickname", "/program/update-image"}},
	}

	// 系统权限所独有的权限
	UniRouteSystem := []schema.RouteList{
		{ID: "46", Name: "企业列表", Url: []string{"/company/list"}},
		{ID: "47", Name: "添加企业", Url: []string{"/company/save"}},
		{ID: "48", Name: "编辑企业", Url: []string{"/company/save"}},
		{ID: "49", Name: "删除企业", Url: []string{"/company/delete"}},
		{ID: "50", Name: "禁用企业", Url: []string{"/company/switch"}},

		{ID: "41", Name: "菜单管理", Url: []string{"/sys-auth-rule/list"}},
		{ID: "30", Name: "新增菜单", Url: []string{"/sys-auth-rule/save"}},
		{ID: "31", Name: "编辑菜单", Url: []string{"/sys-auth-rule/save"}},
		{ID: "32", Name: "删除菜单", Url: []string{"/sys-auth-rule/delete"}},
		{ID: "42", Name: "固件管理", Url: []string{"/app-version-config/config-list"}},
		{ID: "33", Name: "新增固件", Url: []string{"/app-version-config/config-save"}},
		{ID: "34", Name: "编辑固件", Url: []string{"/app-version-config/config-save"}},
		{ID: "35", Name: "删除固件", Url: []string{"/app-version-config/config-delete"}},
	}

	// 系统权限
	RouteSystem := []schema.RouteList{
		{ID: "1", Name: "企业管理", Url: []string{}},
		{ID: "2", Name: "部门管理", Url: []string{"/sys-dept/list"}},
		{ID: "20", Name: "新增部门", Url: []string{"/sys-dept/save"}},
		{ID: "21", Name: "编辑部门", Url: []string{"/sys-dept/save"}},
		{ID: "22", Name: "删除部门", Url: []string{"/sys-dept/delete"}},
		{ID: "3", Name: "人员管理", Url: []string{"/user/brife-user-list", "/user/get-user-role-by-id"}},
		{ID: "23", Name: "新增人员", Url: []string{"/user/add"}},
		{ID: "24", Name: "编辑人员", Url: []string{"/user/save", "/user/update-user-image", "/user/modify-nick-name", "/user/check-user-phone", "/user/update-user-phone", "/user/check-user-email", "/user/update-user-email", "/user/send-code", "/company/update-company-logo", "/company/edit-info"}},
		{ID: "25", Name: "删除人员", Url: []string{"/user/delete"}},
		{ID: "43", Name: "人员状态切换", Url: []string{"/user/switch"}},
		{ID: "44", Name: "用户角色分配", Url: []string{"/user/user-division-role"}},
		{ID: "45", Name: "重置密码", Url: []string{"/user/reset-password"}},

		{ID: "4", Name: "角色管理", Url: []string{"/role/role-brief-list", "/role/role-area-permission", "/role/get-user-by-role-id"}},
		{ID: "26", Name: "新增角色", Url: []string{"/sys-role/save"}},
		{ID: "27", Name: "编辑角色", Url: []string{"/sys-role/save"}},
		{ID: "28", Name: "删除角色", Url: []string{"/sys-role/delete"}},
		{ID: "29", Name: "角色权限分配", Url: []string{"/sys-role/division"}},

		{ID: "5", Name: "系统设置", Url: []string{}},
		{ID: "8", Name: "操作日志", Url: []string{"/log/list"}},
		{ID: "41", Name: "菜单管理", Url: []string{"/sys-auth-rule/list"}},
		{ID: "30", Name: "新增菜单", Url: []string{"/sys-auth-rule/save"}},
		{ID: "31", Name: "编辑菜单", Url: []string{"/sys-auth-rule/save"}},
		{ID: "32", Name: "删除菜单", Url: []string{"/sys-auth-rule/delete"}},
		{ID: "42", Name: "固件管理", Url: []string{"/app-version-config/config-list"}},
		{ID: "33", Name: "新增固件", Url: []string{"/app-version-config/config-save"}},
		{ID: "34", Name: "编辑固件", Url: []string{"/app-version-config/config-save"}},
		{ID: "35", Name: "删除固件", Url: []string{"/app-version-config/config-delete"}},

		{ID: "6", Name: "厂区管理", Url: []string{}},
		{ID: "7", Name: "厂区结构", Url: []string{"/factory/region", "/factory/get-room-stat-by-factory-id", "/factory/get-factory-asset", "/factory/get-factory-bar-chart", "/factory/get-factory-workshop-data", "/factory/factory-device-list", "factory/factory-device-count", "/workshop/get-operation-type", "/workshop/detail", "/workshop/get-workshop-brief", "/workshop/get-workshop-operation"}},
		{ID: "14", Name: "新增厂区", Url: []string{"/factory/save"}},
		{ID: "15", Name: "编辑厂区", Url: []string{"/factory/save"}},
		{ID: "16", Name: "删除厂区", Url: []string{"/factory/delete"}},
		{ID: "17", Name: "新增车间", Url: []string{"/workshop/save"}},
		{ID: "18", Name: "编辑车间", Url: []string{"/workshop/save"}},
		{ID: "19", Name: "删除车间", Url: []string{"/workshop/delete"}},

		{ID: "9", Name: "班组管理", Url: []string{"/team/list", "/team/manage-operation-list", "/team/team-manage-areas"}},
		{ID: "13", Name: "新增班组", Url: []string{"/team/save"}},
		{ID: "10", Name: "编辑班组", Url: []string{"/team/save"}},
		{ID: "11", Name: "删除班组", Url: []string{"/team/delete"}},
		{ID: "12", Name: "班组区域权限划分", Url: []string{"/team/team-save-manage-areas"}},

		{ID: "36", Name: "设备管理", Url: []string{"/device/list", "/device-type/list", "/brewing/get-device-info", "/overview/get-device-data", "/log/handle-device-log-list"}},
		{ID: "37", Name: "新增设备", Url: []string{"/device/batch-add-device"}},
		{ID: "38", Name: "编辑设备", Url: []string{"/device/batch-add-device"}},
		{ID: "39", Name: "删除设备", Url: []string{"/device/delete"}},
		{ID: "40", Name: "设备数据导出", Url: []string{"/log/device-log-export"}},

		{ID: "1734511163548215", Name: "部门列表", Url: []string{"/sys-dept/list"}},
		{ID: "1734511219488188", Name: "人员列表", Url: []string{"/user/list"}},
		{ID: "1734511251753566", Name: "角色列表", Url: []string{"/sys-role/list"}},
		{ID: "1734511593951108", Name: "厂区查看", Url: []string{"/factory/list"}},
		{ID: "1734511613418104", Name: "车间查看", Url: []string{"/workshop/list"}},
		{ID: "1734511648674520", Name: "酒坛查看", Url: []string{"/wine-building/list"}},

		{ID: "51", Name: "大数据管理", Url: []string{"/factory/get-factory-bar-chart", "/overview/one", "/overview/two", "/overview/three", "/overview/device-warning", "/weather/today-weather", "/factory/get-factory-asset", "/factory/get-factory-bar-chart", "/factory/get-factory-workshop-data", "/factory/factory-device-count", "/factory/get-factory-bar-chart", "/workshop/get-workshop-brief", "/overview/get-operation-by-workshop-id", "/device/get-device-total", "/overview/get-distribute-by-workshop-id", "/overview/get-distribute-by-operation-id", "/workshop/get-workshop-operation", "/overview/get-current-production-stage", "/brewing/get-area-device-data", "/brewing/get-device-info", "/overview/get-device-data", "overview/get-production-stage-list", "/brewing/get-brewing-stage-position-list", "/overview/get-brewing-history-chart-data", "/overview/get-yeastroom-overturn-log", "/brewing/get-area-device-data", "/get-brewing-yeastroom-position-list", "/device/get-area-device-stat", "/device-type/list", "/factory/factory-device-list"}},

		{ID: "46", Name: "企业列表", Url: []string{"/company/list"}},
		{ID: "47", Name: "添加企业", Url: []string{"/company/save"}},
		{ID: "48", Name: "编辑企业", Url: []string{"/company/save"}},
		{ID: "49", Name: "删除企业", Url: []string{"/company/delete"}},
		{ID: "50", Name: "禁用企业", Url: []string{"/company/switch"}},
	}

	// 企业权限
	RouteCompany := []schema.RouteList{
		{ID: "1", Name: "企业管理", Url: []string{}},
		{ID: "2", Name: "部门管理", Url: []string{"/sys-dept/list"}},
		{ID: "20", Name: "新增部门", Url: []string{"/sys-dept/save"}},
		{ID: "21", Name: "编辑部门", Url: []string{"/sys-dept/save"}},
		{ID: "22", Name: "删除部门", Url: []string{"/sys-dept/delete"}},
		{ID: "3", Name: "人员管理", Url: []string{"/user/profile", "/company/profile", "/user/brife-user-list", "/user/get-user-menu-list", "/user/get-user-role-by-id"}},
		{ID: "23", Name: "新增人员", Url: []string{"/user/add"}},
		{ID: "24", Name: "编辑人员", Url: []string{"/user/save", "/user/update-user-image", "/user/modify-nick-name", "/user/check-user-phone", "/user/update-user-phone", "/user/check-user-email", "/user/update-user-email", "/user/send-code", "/company/update-company-logo", "/company/edit-info"}},
		{ID: "25", Name: "删除人员", Url: []string{"/user/delete"}},
		{ID: "43", Name: "人员状态切换", Url: []string{"/user/switch"}},
		{ID: "44", Name: "用户角色分配", Url: []string{"/user/user-division-role"}},
		{ID: "45", Name: "重置密码", Url: []string{"/user/reset-password"}},

		{ID: "4", Name: "角色管理", Url: []string{"/role/role-brief-list", "/role/role-area-permission", "/role/get-user-by-role-id"}},
		{ID: "26", Name: "新增角色", Url: []string{"/sys-role/save"}},
		{ID: "27", Name: "编辑角色", Url: []string{"/sys-role/save"}},
		{ID: "28", Name: "删除角色", Url: []string{"/sys-role/delete"}},
		{ID: "29", Name: "角色权限分配", Url: []string{"/sys-role/division"}},

		{ID: "5", Name: "系统设置", Url: []string{}},
		{ID: "8", Name: "操作日志", Url: []string{"/log/list"}},

		{ID: "6", Name: "厂区管理", Url: []string{}},
		{ID: "7", Name: "厂区结构", Url: []string{"/factory/region", "/factory/get-room-stat-by-factory-id", "/factory/get-factory-asset", "/factory/get-factory-bar-chart", "/factory/get-factory-workshop-data", "/factory/factory-device-list", "factory/factory-device-count", "/workshop/get-operation-type", "/workshop/detail", "/workshop/get-workshop-brief", "/workshop/get-workshop-operation"}},
		{ID: "14", Name: "新增厂区", Url: []string{"/factory/save"}},
		{ID: "15", Name: "编辑厂区", Url: []string{"/factory/save"}},
		{ID: "16", Name: "删除厂区", Url: []string{"/factory/delete"}},
		{ID: "17", Name: "新增车间", Url: []string{"/workshop/save"}},
		{ID: "18", Name: "编辑车间", Url: []string{"/workshop/save"}},
		{ID: "19", Name: "删除车间", Url: []string{"/workshop/delete"}},

		{ID: "9", Name: "班组管理", Url: []string{"/team/list", "/team/manage-operation-list", "/team/team-manage-areas"}},
		{ID: "13", Name: "新增班组", Url: []string{"/team/save"}},
		{ID: "10", Name: "编辑班组", Url: []string{"/team/save"}},
		{ID: "11", Name: "删除班组", Url: []string{"/team/delete"}},
		{ID: "12", Name: "班组区域权限划分", Url: []string{"/team/team-save-manage-areas"}},

		{ID: "36", Name: "设备管理", Url: []string{"/device/list", "/device-type/list", "/brewing/get-device-info", "/overview/get-device-data", "/log/handle-device-log-list"}},
		{ID: "37", Name: "新增设备", Url: []string{"/device/batch-add-device"}},
		{ID: "38", Name: "编辑设备", Url: []string{"/device/batch-add-device"}},
		{ID: "39", Name: "删除设备", Url: []string{"/device/delete"}},
		{ID: "40", Name: "设备数据导出", Url: []string{"/log/device-log-export"}},

		// 酒坛相关
		{ID: "56", Name: "组织架构", Url: []string{"/wine-building/append-winejar", "/wine-building/delete-winejar", "/wine-building/edit-winejar", "/wine-overview/overview-stat", "/wine-overview/overview", "/wine-overview/distribute", "/wine-overview/device-list", "/wine-overview/batch-add-device", "/wine-overview/device-update", "/wine-overview/device-delete", "/wine-overview/device-detail", "/wine-overview/room-info", "/wine-overview/summary", "/wine-overview/overview-stat", "/wine-overview/operation-overview", "/wine-overview/winejar-warning-list", "/wine-overview/get-winejar-list-by-room-id", "/wine-overview/winejar-detail", "/wine-overview/get-winejar-history-handle-logs", "/wine-overview/device-overview", "/wine-overview/building-device-distribute", "/wine-overview/operation-device-distribute", "/wine-overview/device-data-list", "/wine-overview/alarm-list"}},
		{ID: "1734403954047949", Name: "添加", Url: []string{"/wine-building/save", "/wine-building/save-room", "/wine-building/batch-save-room", "/wine-building/save-room-door", "/wine-building/save-building"}},
		{ID: "1734405134173801", Name: "编辑", Url: []string{"/wine-building/edit", "/wine-building/edit-room", "/wine-building/batch-edit-room", "/wine-building/edit-room-door", "/wine-building/edit-building"}},
		{ID: "1734405152242365", Name: "删除", Url: []string{"/wine-building/delete", "/wine-building/delete-operation", "/wine-building/delete-row", "/wine-building/delete-room", "/wine-building/edit-building"}},
		{ID: "59", Name: "报警记录", Url: []string{"/wine-overview/alarm-list"}},
		{ID: "60", Name: "设备总览", Url: []string{"/wine-overview/device-overview", "/wine-overview/building-device-distribute", "/wine-overview/operation-device-distribute"}},
		{ID: "61", Name: "设备数据", Url: []string{"/wine-overview/device-data-list"}},
		{ID: "1726812431113807", Name: "操作日志", Url: []string{"/log/list"}},

		{ID: "51", Name: "大数据管理", Url: []string{"/factory/get-factory-bar-chart", "/overview/one", "/overview/two", "/overview/three", "/overview/device-warning", "/weather/today-weather", "/factory/get-factory-asset", "/factory/get-factory-bar-chart", "/factory/get-factory-workshop-data", "/factory/factory-device-count", "/factory/get-factory-bar-chart", "/workshop/get-workshop-brief", "/overview/get-operation-by-workshop-id", "/device/get-device-total", "/overview/get-distribute-by-workshop-id", "/overview/get-distribute-by-operation-id", "/workshop/get-workshop-operation", "/overview/get-current-production-stage", "/brewing/get-area-device-data", "/brewing/get-device-info", "/overview/get-device-data", "overview/get-production-stage-list", "/brewing/get-brewing-stage-position-list", "/overview/get-brewing-history-chart-data", "/overview/get-yeastroom-overturn-log", "/brewing/get-area-device-data", "/get-brewing-yeastroom-position-list", "/device/get-area-device-stat", "/device-type/list", "/factory/factory-device-list"}},

		{ID: "1734511163548215", Name: "部门列表", Url: []string{"/sys-dept/list"}},
		{ID: "1734511219488188", Name: "人员列表", Url: []string{"/user/list"}},
		{ID: "1734511251753566", Name: "角色列表", Url: []string{"/sys-role/list"}},
		{ID: "1734511593951108", Name: "厂区查看", Url: []string{"/factory/list"}},
		{ID: "1734511613418104", Name: "车间查看", Url: []string{"/workshop/list"}},
		{ID: "1734511648674520", Name: "酒坛查看", Url: []string{"/wine-building/list"}},
	}
	if RouteType == 1 { // 返回系统权限
		return RouteSystem
	} else if RouteType == 2 { // 返回企业权限
		return RouteCompany
	} else if RouteType == 3 { // 返回通用权限
		return RouteCommon
	} else if RouteType == 4 {
		return UniRouteSystem // 返回系统所独有的权限
	}
	return RouteAll
}

func InArray(need string, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

func GetMapKey(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
