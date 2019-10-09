package common

const (
	//默认返回
	RECODE_CODE_UNKNOWER = "100"
	////////////////回复状态码及信息定义//////////
	RESP_CODE_OK              = "2000"
	RESP_CODE_GETJSONERR      = "2001"
	RESP_CODE_JSON_DATANULL   = "2002"
	RESP_CODE_DATAISEXISTS    = "2003"
	RESP_CODE_DATAISNOTEXISTS = "2005"
	RESP_CODE_CERATE_ERR      = "2006"
	RESP_CODE_DEL_ERR         = "2007"
	RESP_CODE_UPDATE_ERR      = "2008"
	RESP_CODE_SELECT          = "2009"
	RESP_CODE_NOPERMISSION    = "2010"
	RESP_CODE_SUBMIT_ERR      = "2011"
	///////////////INFO/////////////////////
	RECODE_INFO_UNKNOWERR = "为知错误"
	//返回信息的描述
	RESP_INFO_OK              = "成功"
	RESP_INFO_GETJSONERR      = "无法获取JSON中数据"
	RESP_INFO_JSON_DATANULL   = "JSON中获取的数据为NULL"
	RESP_INFO_DATAISEXISTS    = "数据已存在"
	RESP_INFO_DATAISNOTEXISTS = "数据不存在"
	RESP_INFO_CREATE_ERR      = "新增数据失败"
	RESP_INFO_NOPERMISSION    = "权限不足"
	RESP_INFO_SUBMIT_ERR      = "文章提交错误"
	//////////////日志按业务定义//////////////////

	//文章错误提示
	LOG_ARTICLE_ADD_SUCCESS       = "[ARTICLE_ADD_SUCCESS]: %s"
	LOG_ARTICLE_DEL_SUCCESS       = "[ARTICLE_DEL_SUCCESS]: %s"
	LOG_ARTICLE_SUBMIT_SUCCESS    = "[ARTICLE_SUBMIT_SUCCESS]: %s"
	LOG_ARTICLE_LISTBYTAG_SUCCESS = "[ARTICLE_LISTBYTAG_SUCCESS]: %s"
	LOG_ARTICLE_ADD_ERR           = "[ARTICLE_ADD_ERR]: %s"
	LOG_ARTICLE_DELETE_ERR        = "[ARTICLE_DELETE_ERR]: %s"

	LOG_ARTICLE_EDIT_ERR           = "[ARTICLE_EDIT_ERR]: %s"
	LOG_ARTICLE_LIST_BYTAG_ERR     = "[ARTICLE_LIST_BYTAG_ERR]: %s"
	LOG_ARTICLE_LIST_BYCOLUMN_ERR  = "[ARTICLE_LIST_BYCOLUMN_ERR]: %s"
	LOG_ARTICLE_SUBMIT_ERR         = "[ARTICLE_SUBMIT_ERR]: %s"
	LOG_ARTICLE_GETWILLBESUB_ERR   = "[ARTICLE_GETWILLBESUB_ERR]: %s"
	LOG_ARTICLE_RELEASE_ERR        = "[ARTICLE_RELEASE_ERR]: %s"
	LOG_ARTICLE_GETRELEASELIST_ERR = "[ARTICLE_GETRELEASELIST_ERR]: %s"
	LOG_ARTICLE_RELEASEFAILD_ERR   = "[ARTICLE_RELEASEFAILD_ERR]: %s"
	//文章图片提示
	LOG_IMAGE_ADD_SUCCESS = "[IMAGE_ADD_SUCCESS]: %s"
	LOG_IMAGE_ADD_ERR     = "[IMAGE_ADD_ERR]:"
	LOG_IAMGE_SIZE_ERR    = "图片不可超过10兆"
	LOG_IAMGE_EXT_ERR     = "图片格式仅支持:jpg,png,jpeg"
	//用户提示
	LOG_USER_DEL_SUCCESS           = "[USER_DEL_SUCCESS]: %s"
	LOG_USER_LOGIN_SUCCESS         = "[USER_LOGIN_SUCCESS]: %s"
	LOG_USER_LIST_SUCCESS          = "[USER_LIST_SUCCESS]: %s"
	LOG_USER_UPDATE_ROLLER_SUCCESS = "[USER_UPDATE_ROLLER_SUCCESS]: %s"

	LOG_USER_REGISTER_ERR      = "[USER_REGISTER_ERR]: %s"
	LOG_USER_LOGIN_ERR         = "[USER_LOGIN_ERR]: %s"
	LOG_USER_LIST_ERR          = "[USER_LIST_ERR]: %s"
	LOG_USER_ADD_ERR           = "[USER_ADD_ERR]: %s"
	LOG_USER_DELETE_ERR        = "[USER_DELETE_ERR]: %s"
	LOG_USER_UPDATE_ROLLER_ERR = "[USER_UPDATE_ROLLER_ERR]: %s"
	LOG_USER_UPDATE_PWD_ERR    = "[USER_UPDATE_PWD_ERR]: %s"
	LOG_USER_CNAME_ERR         = "[USER_CNAME_ERR]: %s"
	LOG_USER_ISEXISTS          = "用户已存在"
	LOG_USER_CREATE_OK         = "用户创建成功"

	//专题提示信息
	LOG_COLUMN_ADD_ERR     = "[COLUMN_ADD_ERR]: %s"
	LOG_COLUMN_ADD_SUCCESS = "[COLUMN_ADD_SUCCESS]: %s"
	LOG_COLUMN_DELETE_ERR  = "[COLUMN_DELETE_ERR]: %s"
	LOG_COLUMN_DEL_SUCCESS = "[COLUMN_DEL_SUCCESS]: %s"
	LOG_COLUMN_LIST_ERR    = "[COLUMN_LIST_ERR]: %s"
	LOG_COLUMN_CNAME_ERR   = "[COLUMN_CNAME_ERR]: %s"
	//标签错误提示
	LOG_TAG_ADD_ERR    = "[TAG_ADD_ERR]: %s"
	LOG_TAG_DELETE_ERR = "[TAG_DELETE_ERR]: %s"
	LOG_TAG_CNAME_ERR  = "[TAG_CNAME_ERR]: %s"
	//数据库错误
	DB_CREATE_ERR         = "数据库插入错误"
	DB_UPDATE_ERR         = "数据库更新错误"
	DB_DEL_ERR            = "数据库删除错误"
	DB_SELECT_ERR         = "数据库查询错误"
	DB_ARTICLE_SUBMIT_ERR = "文章提交错误"
)
