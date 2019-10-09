package controller

import (
	. "articlebk/src/common"
	"articlebk/src/common/database"
	"articlebk/src/common/database/sql"
	"github.com/gin-gonic/gin"

	"strconv"
)

func PostColumnAdd(ctx *gin.Context) {

	var column struct {
		columnName string `json:"column_name"`
		ParentId   string `json:"parent_id"`
	}

	if err := ctx.BindJSON(&column); err != nil {
		Log.Error(LOG_COLUMN_ADD_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	columnName := column.columnName
	columnParentId := column.ParentId
	if columnName == "" {
		Log.Error(LOG_COLUMN_ADD_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//如果没有父级id,则属于一级菜单,默认父级id是0
	if columnParentId == "" {
		columnParentId = "0"
	}
	columnPid, _ := strconv.Atoi(columnParentId)
	if sql.IsexistColumn(columnName) {
		Log.Error(LOG_COLUMN_ADD_ERR, RESP_INFO_DATAISEXISTS)
		Resp(ctx, RESP_CODE_DATAISEXISTS, RESP_INFO_DATAISEXISTS, nil)
		return
	}
	//如果pid非0且不存在,则无法添加
	//if columnParentId != "0" && !sql.IsExistColumnBySid(columnParentId) {
	//	Log.Error(LOG_COLUMN_ADD_ERR, RESP_INFO_DATAISEXISTS)
	//	Resp(ctx, RESP_CODE_DATAISEXISTS, RESP_INFO_DATAISEXISTS, nil)
	//	return
	//}

	col, err := sql.ColumnAdd(columnName, columnPid)
	if err != nil {
		Log.Error(LOG_COLUMN_ADD_ERR, DB_CREATE_ERR, err)
		Resp(ctx, RESP_CODE_CERATE_ERR, DB_CREATE_ERR, nil)
		return
	}
	Log.Info(LOG_COLUMN_ADD_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, col)
}

func DeleteColumn(ctx *gin.Context) {
	var columnInfo struct {
		ColumnId string `json:"column_id"`
	}

	err := ctx.BindJSON(&columnInfo)
	if err != nil {
		Log.Error(LOG_COLUMN_DELETE_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//获取专题信息判空
	sid := columnInfo.ColumnId
	if sid == "" {
		Log.Error(LOG_COLUMN_DELETE_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//判断专题是否存在
	if !sql.IsExistColumnBySid(sid) {
		Log.Error(LOG_COLUMN_DELETE_ERR, RESP_INFO_DATAISNOTEXISTS)
		Resp(ctx, RESP_CODE_DATAISEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
		return
	}
	//判断该专题是否有子专题,如果有则禁止删除
	if sql.ColumnHasSub(sid) {
		Log.Error(LOG_COLUMN_DELETE_ERR, "请先删除子专栏")
		Resp(ctx, RESP_CODE_DEL_ERR, "请先删除子专栏", nil)
		return
	}
	//删除专题
	if err := sql.ColumnDel(sid); err != nil {
		Log.Error(LOG_COLUMN_DELETE_ERR, DB_DEL_ERR)
		Resp(ctx, RESP_CODE_DEL_ERR, DB_DEL_ERR, nil)
		return
	}
	Log.Info(LOG_COLUMN_DEL_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, nil)
}

func GetColumnList(ctx *gin.Context) {
	columnList := []database.Column{}
	columns, err := sql.ColumnListByPid(0)
	if err != nil {
		Log.Error(LOG_COLUMN_LIST_ERR, DB_SELECT_ERR, err)
		Resp(ctx, RESP_CODE_SELECT, DB_SELECT_ERR, nil)
		return
	}
	for _, column := range columns {
		getAllColumn(&column)
		columnList = append(columnList, column)
	}
	Log.Info(LOG_COLUMN_LIST_ERR, DB_SELECT_ERR, err)
	Resp(ctx, RESP_CODE_SELECT, DB_SELECT_ERR, columnList)
}

func getAllColumn(column *database.Column) {
	sps, err := sql.ColumnListById(column.Id)
	if err != nil {
		Log.Error(LOG_COLUMN_LIST_ERR, DB_SELECT_ERR, err)
		return
	}
	for _, v := range sps {
		getAllColumn(&v)
		column.Columns = append(column.Columns, v)
	}
}

func PutColumnCname(ctx *gin.Context) {
	var columnInfo struct {
		columnId      string `json:"column_id"`
		columnNewName string `json:"column_new_name"`
	}

	err := ctx.BindJSON(&columnInfo)
	if err != nil {
		Log.Error(LOG_COLUMN_CNAME_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//获取专题信息判空
	sid := columnInfo.columnId
	newName := columnInfo.columnNewName
	if sid == "" || newName == "" {
		Log.Error(LOG_COLUMN_CNAME_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//判断专题是否存在
	if !sql.IsExistColumnBySid(sid) {
		Log.Error(LOG_COLUMN_CNAME_ERR, RESP_INFO_DATAISNOTEXISTS)
		Resp(ctx, RESP_CODE_DATAISEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
		return
	}
	//修改专题名
	col, err := sql.ColumnCname(sid, newName)
	if err != nil {
		Log.Error(LOG_COLUMN_CNAME_ERR, DB_UPDATE_ERR, err)
		Resp(ctx, RESP_CODE_UPDATE_ERR, DB_UPDATE_ERR, nil)
		return
	}
	Log.Info(LOG_COLUMN_ADD_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, col)
}
