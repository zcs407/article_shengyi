package tag

import (
	. "articlebk/src/common"
	"articlebk/src/common/database/sql"
	"github.com/gin-gonic/gin"
)

func PostTagAdd(ctx *gin.Context) {
	var tag struct {
		TagName string `json:"tag_name"`
	}
	//获取post中的标签名称
	err := ctx.BindJSON(&tag)
	if err != nil {
		Log.Error(LOG_TAG_ADD_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//判空
	tagname := tag.TagName
	if tagname == "" {
		Log.Error(LOG_TAG_ADD_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//判断是否存在
	if sql.IsexistTag(tagname) {
		Log.Error(LOG_TAG_ADD_ERR, RESP_INFO_DATAISEXISTS)
		Resp(ctx, RESP_CODE_DATAISEXISTS, RESP_INFO_DATAISEXISTS, nil)
		return
	}
	//判断数据库返状态
	tid, err := sql.TagAdd(tagname)
	if err != nil {
		Log.Error(LOG_TAG_ADD_ERR, DB_CREATE_ERR, err)
		Resp(ctx, RESP_CODE_CERATE_ERR, DB_CREATE_ERR, nil)
		return
	}
	//返回标签id
	Log.Info(LOG_TAG_ADD_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, tid)
}

func PostTagDel(ctx *gin.Context) {
	var tag struct {
		TagId string `json:"tag_id"`
	}
	//获取post数据
	err := ctx.BindJSON(&tag)
	if err != nil {
		Log.Error(LOG_TAG_DELETE_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//判空
	tagId := tag.TagId
	if tagId == "" {
		Log.Error(LOG_TAG_DELETE_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//判断是否存在
	if !sql.IsexistTagById(tagId) {
		Log.Error(LOG_TAG_DELETE_ERR, RESP_INFO_DATAISEXISTS)
		Resp(ctx, RESP_CODE_DATAISEXISTS, RESP_INFO_DATAISEXISTS, nil)
		return
	}
	//删除标签
	_, err = sql.TagDelById(tagId)
	if err != nil {
		Log.Error(LOG_TAG_DELETE_ERR, DB_CREATE_ERR, err)
		Resp(ctx, RESP_CODE_DEL_ERR, DB_DEL_ERR, nil)
		return
	}
	Log.Info(LOG_TAG_DEL_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, nil)
}

func PostTagCname(ctx *gin.Context) {

	var tag struct {
		TagId   string `json:"tag_id"`
		NewName string `json:"new_name"`
	}
	//获取post数据
	err := ctx.BindJSON(&tag)
	if err != nil {
		Log.Error(LOG_TAG_CNAME_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//判空
	tagId := tag.TagId
	tagNewName := tag.NewName
	if tagId == "" || tagNewName == "" {
		Log.Error(LOG_TAG_CNAME_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//判断是否存在
	if !sql.IsexistTagById(tagId) {
		Log.Error(LOG_TAG_CNAME_ERR, RESP_INFO_DATAISNOTEXISTS)
		Resp(ctx, RESP_CODE_DATAISNOTEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
		return
	}
	tname, err := sql.TagCnameById(tagId, tagNewName)
	if err != nil {
		Log.Error(LOG_TAG_CNAME_ERR, DB_UPDATE_ERR, err)
		Resp(ctx, RESP_CODE_UPDATE_ERR, DB_UPDATE_ERR, nil)
		return
	}
	Log.Info(LOG_TAG_CNAME_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, tname)
}
