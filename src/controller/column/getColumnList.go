package column

import (
	"articlebk/src/common/dbtable"
	"articlebk/src/common/sql"
	"github.com/gin-gonic/gin"

	"log"
)

func GetColumnList(ctx *gin.Context) {
	loger := log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
	var resp struct {
		Code        string            `json:"code"`
		Info        string            `json:"info"`
		SpecialList []dbtable.Special `json:"special_list"`
	}
	specials, err := sql.SpecialListByPid(0)
	loger.Println("specials,", specials)
	if err != nil {
		resp.Code = "501"
		resp.Info = "查询pid为0的数据错误"
		ctx.JSON(200, resp)
		return
	}
	for _, special := range specials {
		getAllColumn(&special)
		resp.SpecialList = append(resp.SpecialList, special)
	}
	resp.Code = "200"
	resp.Info = "以获取到专题列表"
	loger.Println("获取到的专题列表是", resp)
	ctx.JSON(200, resp)
}

func getAllColumn(special *dbtable.Special) {
	loger := log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
	sps, err := sql.SpecialListById(special.Id)
	if err != nil {
		loger.Println("[get_special_list_info]", ":已找结点专题")
		return
	}
	for _, v := range sps {
		getAllColumn(&v)
		special.Specials = append(special.Specials, v)
	}
	return
}
