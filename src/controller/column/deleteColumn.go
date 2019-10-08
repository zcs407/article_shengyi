package column

import (
	"articlebk/src/utils/sql"
	"github.com/gin-gonic/gin"
	"log"
)

func DeleteColumn(ctx *gin.Context) {
	loger := log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
	var specialInfo struct {
		SpecialId string `json:"special_id"`
	}
	var resp struct {
		Code string `json:"code"`
		Info string `json:"info"`
	}
	err := ctx.BindJSON(&specialInfo)
	if err != nil {
		resp.Code = "404"
		resp.Info = "无法获取post中的专题信息"
		loger.Println("[special_delete_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//获取专题信息判空
	sid := specialInfo.SpecialId
	if sid == "" {
		resp.Code = "403"
		resp.Info = "post中的special_id不能为空"
		loger.Println("[special_delete_err]", resp)
		ctx.JSON(200, resp)
		return
	}
	//判断专题是否存在
	if !sql.IsExistSpecialBySid(sid) {
		resp.Code = "501"
		resp.Info = "该专题不存在"
		loger.Println("[special_delete_err]", resp)
		ctx.JSON(200, resp)
		return
	}
	//判断该专题是否有子专题,如果有则禁止删除
	if sql.SpecialHasSub(sid) {
		resp.Code = "405"
		resp.Info = "该专题有子专题,请先删除子专题"
		loger.Println("[special_delete_err]", resp)
		ctx.JSON(200, resp)
		return
	}
	//删除专题
	if err := sql.SpecialDel(sid); err != nil {
		resp.Code = "500"
		resp.Info = "数据库无法删除该专题"
		loger.Println("[special_delete_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	resp.Code = "200"
	resp.Info = "专题删除成功"
	loger.Println("[special_delete_err]", resp)
	ctx.JSON(200, resp)
}
