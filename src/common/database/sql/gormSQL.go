package sql

import (
	"articlebk/src/common/database"
	"strconv"
)

///////////////////////////////////////////////////用户管理////////////////////////////////////////
//用户注册
func RegisterUser(userInfo *database.User) (*database.User, error) {
	db := database.DBSQL
	if err := db.Create(&userInfo).Error; err != nil {
		return userInfo, err
	}
	userInfo.PasswordHash = ""
	userInfo.Articles = nil
	return userInfo, nil
}

//用uid和角色id验证用户是否存在
func IsExistByUidRid(uid, rid string) bool {
	db := database.DBSQL

	user := database.User{}
	uidInt, _ := strconv.Atoi(uid)
	ridInt, _ := strconv.Atoi(rid)
	if err := db.Where("id = ? AND roller_id = ?", uidInt, ridInt).First(&user).Error; err != nil {
		return false
	}
	return true
}

//用户是否存在,username查询(仅用于用户注册和登录时使用)
func UserIsExistByName(username string) bool {
	db := database.DBSQL

	user := database.User{}
	err := db.Raw("select Name from user where name = ?", username).Scan(&user).Error
	if err != nil {
		return false
	}
	return true
}

//用户是否为管理员
func UserIsAdmin(uid string) bool {
	db := database.DBSQL
	user := database.User{}
	id, _ := strconv.Atoi(uid)
	if err := db.Where("id = ? AND roller_id = ?", id, 0).First(&user).Error; err != nil {
		return false
	}
	return true
}

//用户是否存在,用uid查询
func UserIsExistByUid(uid string) bool {
	db := database.DBSQL

	uidint, _ := strconv.Atoi(uid)
	user := database.User{}
	if err := db.Where("id = ?", uidint).First(&user).Error; err != nil {
		return false
	}
	return true
}

//用户登录
func UserLogin(username, password string) (database.User, error) {
	db := database.DBSQL
	user := database.User{}
	err := db.Where("name = ? AND password_hash = ?", username, password).Find(&user).Error
	if err != nil {
		return user, err
	}
	user.PasswordHash = ""
	user.Articles = nil
	return user, nil
}

//用户列表查询,仅管理员可查,可跟实际需求修改
func UserListGet() ([]database.User, error) {
	db := database.DBSQL

	var users []database.User
	if err := db.Select("id, name, roller_id").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

//用户角色变更
func UserRollerUpdate(uid, username, newrolleid string) (database.User, error) {
	db := database.DBSQL
	uidint, _ := strconv.Atoi(uid)
	user := database.User{}
	err := db.Exec("UPDATE user SET roller_id = ? WHERE id = ? AND name = ?", newrolleid, uidint, username).Error
	if err != nil {
		return user, err
	}
	user.Id = uidint
	db.First(&user)
	user.PasswordHash = ""
	user.Articles = nil
	return user, nil
}

//用户密码修改
func UpdateUserPwd(uid, npwd string) error {
	db := database.DBSQL

	uidInt, _ := strconv.Atoi(uid)
	user := database.User{}
	if err := db.Model(&user).Where("id = ?", uidInt).Update("password_hash", npwd).Error; err != nil {
		return err
	}
	return nil
}

//根据uid和密码验证用户是否密码正确
func VerifyUserPwd(uid, opwd string) bool {
	db := database.DBSQL

	uidInt, _ := strconv.Atoi(uid)
	user := database.User{}
	err := db.Where("id = ? AND password_hash = ?", uidInt, opwd).First(&user).Error
	if err != nil {
		return false
	}
	return true
}

//删除用户
func DeleteUser(uid string) error {
	db := database.DBSQL

	uidInt, _ := strconv.Atoi(uid)
	user := database.User{}
	if err := db.Where("id = ?", uidInt).First(&user).Error; err != nil {
		return err
	}
	if err := db.Where("id = ?", uidInt).Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

////////////////////////////////////////////////////////文章管理/////////////////////////////////////

//创建文章
func ArticleAdd(article database.Article, tags []string) (database.Article, error) {
	db := database.DBSQL
	for _, tag := range tags {
		tid, _ := strconv.Atoi(tag)
		var tag database.Tag
		if err := db.Where("id = ?", tid).Find(&tag).Error; err != nil {
			return article, err
		}
		article.Tags = append(article.Tags, &tag)
	}

	if err := db.Create(&article).Error; err != nil {
		return article, err
	}
	return article, nil
}

//查询文章列表

//查询某个标签的文章集合
//SELECT article.* FROM article INNER JOIN article_tag
// ON  article_tag.article_id = article.id WHERE  article_tag.tag_id IN (2);
func GetArticlesByTagId(tagid string) ([]*database.Article, error) {
	db := database.DBSQL

	tid, _ := strconv.Atoi(tagid)
	var articles []*database.Article
	var images []database.Image
	tag := database.Tag{}
	tag.Id = tid
	err := db.Raw(`SELECT article.* FROM article INNER JOIN article_tag ON article_tag.article_id = article.id WHERE  article_tag.tag_id IN (?);`, tid).Scan(&articles).Error
	if err != nil {
		return nil, err
	}
	for _, article := range articles {
		err := db.Where("article_id = ?", article.Id).Find(&images).Error
		if err != nil {
			continue
		}
		article.Images = append(article.Images, images...)
	}
	return articles, nil
}

//查询某专题的文章
func GetArticlesByColumnID(cid int) ([]database.Article, error) {
	db := database.DBSQL
	var articles []database.Article
	//err := db.Where("id = ?", cid).First(&columns).Error
	//if err != nil {
	//	fmt.Println("查不到")
	//}

	err := db.Where("column_id = ?", cid).Find(&articles).Error
	if err != nil {
		return articles, err
	}
	return articles, nil
}

//编辑更新文章
func ArticleUpdate(article database.Article) error {
	db := database.DBSQL

	if err := db.Update(&article).Error; err != nil {
		return err
	}
	return nil
}

//删除文章
func ArticleDel(aid string) (cPath string, err error) {
	db := database.DBSQL

	article := database.Article{}
	aidInt, _ := strconv.Atoi(aid)
	db.Where("id = ?", aidInt).First(&article)

	cPath = article.ContentPath
	if err := db.Delete(&article).Error; err != nil {
		return "", err
	}
	return cPath, nil
}

//文章是否存在
func IsexistArticle(title string) bool {
	db := database.DBSQL

	article := database.Article{}
	err := db.Raw("select title from article where title = ?", title).Scan(&article).Error
	if err != nil {
		return false
	}
	return true
}

//判断用户与文章是否一致
func UserHasArticle(aid, uid int) bool {
	db := database.DBSQL
	article := database.Article{}
	if err := db.Where("id = ? AND user_id = ?", aid, uid).First(&article).Error; err != nil {
		return false
	}
	return true
}

//提交文章
func ArticleSubmit(aid, uid int) error {
	db := database.DBSQL
	article := database.Article{}
	err := db.Model(&article).Where("user_id = ? AND id = ?", uid, aid).Update("status", 1).Error
	if err != nil {
		return err
	}
	return nil
}

//发布文章
func ArticleRelease(aid, uid int) error {
	db := database.DBSQL
	article := database.Article{}
	err := db.Model(&article).Where("user_id = ? AND id = ?", uid, aid).Update("status", 2).Error
	if err != nil {
		return err
	}
	return nil
}

//获取未提交的所有文章
func ArticleWillBeSubmit() (articles []database.Article, err error) {
	db := database.DBSQL
	if err := db.Where("status = ?", 0).Find(articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

//查询用户未提交的文章
func ArticleWillBeSubmitByUid(uid int) (articles []database.Article, err error) {
	db := database.DBSQL
	if err := db.Where("status = ? AND user_id = ?", 0, uid).Find(articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

//查询所有已提交的文章
func ArticleSubmited() (articles []database.Article, err error) {
	db := database.DBSQL
	if err := db.Where("status = ?", 1).Find(articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

//查询用户所提交的文章
func ArticleSubmitedByUid(uid int) (articles []database.Article, err error) {
	db := database.DBSQL
	if err := db.Where("status = ? AND user_id = ?", 1, uid).Find(articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

//查询所有已发布的文章
func ArticleReleased() (articles []database.Article, err error) {
	db := database.DBSQL
	if err := db.Where("status = ?", 2).Find(articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

//查询用户发布的文章
func ArticleReleasedByUid(uid int) (articles []database.Article, err error) {
	db := database.DBSQL
	if err := db.Where("status = ? AND user_id = ?", 2, uid).Find(articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

//获取发布失败的文章
func ArticleReleaseFailed() (articles []database.Article, err error) {
	db := database.DBSQL
	if err := db.Where("status = ?", 3).Find(articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

//查询用户发布的文章
func ArticleReleaseFailedByUid(uid int) (articles []database.Article, err error) {
	db := database.DBSQL
	if err := db.Where("status = ? AND user_id = ?", 3, uid).Find(articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

///////////////////////////////////////////////图片管理///////////////////////////////////////////
func ArticleImageAdd(imgurl, imgPath, aid string) error {
	db := database.DBSQL

	image := database.Image{}
	image.ImageUrl = imgurl
	image.ImagePath = imgPath
	err := db.Exec(`INSERT INTO image (article_id,image_url,image_path) VALUES (?,?,?)`, aid, imgurl, imgPath).Error
	if err != nil {
		return err
	}
	return nil
}

func ArticleImageDelByAid(aid string) []database.Image {
	db := database.DBSQL

	images := []database.Image{}
	aidInt, _ := strconv.Atoi(aid)
	db.Where("article_id = ?", aidInt).Find(&images)
	err := db.Where("article_id = ?", aidInt).Delete(&images).Error
	if err != nil {
		return nil
	}
	return images
}

////////////////////////////////////////////////专题管理////////////////////////////////////////////
//创建专题
func ColumnAdd(name string, Pid int) (database.Columns, error) {
	db := database.DBSQL
	column := database.Columns{}
	column.ColumnName = name
	column.Pid = Pid
	if err := db.Create(&column).Error; err != nil {
		return column, err
	}
	return column, nil
}

//按pid查询专题列表
func ColumnListByPid(pid int) ([]database.Columns, error) {
	db := database.DBSQL

	columns := []database.Columns{}
	err := db.Where("pid = ?", pid).Find(&columns).Error
	if err != nil {
		return nil, err
	}
	return columns, nil
}

//按id查询专题列表
func ColumnListById(id int) ([]database.Columns, error) {
	db := database.DBSQL

	columns := []database.Columns{}
	err := db.Where("pid = ?", id).Find(&columns).Error
	if err != nil {
		return nil, err
	}
	//defer db.Close()
	return columns, nil
}

//更新专题名称
func ColumnCname(sid, newName string) (database.Columns, error) {
	db := database.DBSQL
	sidInt, _ := strconv.Atoi(sid)
	column := database.Columns{}
	err := db.Model(&column).Where("id = ?", sidInt).Update("column_name", newName).Error
	if err != nil {
		return column, err
	}
	column.Articles = nil
	column.Columns = nil
	return column, nil
}

//删除专题
func ColumnDel(sid string) error {
	db := database.DBSQL
	column := database.Columns{}
	sidInt, _ := strconv.Atoi(sid)
	err := db.Where("id = ?", sidInt).Delete(&column).Error
	if err != nil {
		return err
	}
	return nil
}

//是否有子专题
func ColumnHasSub(sid string) bool {
	db := database.DBSQL

	column := database.Columns{}
	sidInt, _ := strconv.Atoi(sid)
	err := db.Where("pid = ?", sidInt).First(&column).Error
	if err != nil {
		return false
	}
	return true
}

//通过专题名判断是否存在
func IsexistColumn(columnname string) bool {
	db := database.DBSQL

	column := database.Columns{}
	err := db.Where("column_name = ?", columnname).First(&column).Error
	if err != nil {
		return false
	}
	return true
}

//通过专题id判断是否存在
func IsExistColumnBySid(sid string) bool {
	db := database.DBSQL

	sidInt, _ := strconv.Atoi(sid)
	column := database.Columns{}
	err := db.Where("id = ?", sidInt).Find(&column).Error
	if err != nil {
		return false
	}
	return true
}

//////////////////////////////////////////标签管理/////////////////////////////////////////
//创建标签
func TagAdd(tname string) (string, error) {
	db := database.DBSQL

	var tag database.Tag
	tag.TagName = tname
	if err := db.Create(&tag).Error; err != nil {
		return "", err
	}
	return strconv.Itoa(tag.Id), nil

}

//标签是否存在
func IsexistTag(tname string) bool {
	db := database.DBSQL

	tag := database.Tag{}
	err := db.Raw("select tag_name from tag where tag_name = ?", tname).Scan(&tag).Error
	if err != nil {
		return false
	}
	return true
}

//使用id查询标签是否存在
func IsexistTagById(tid string) bool {
	db := database.DBSQL

	tag := database.Tag{}
	tidint, _ := strconv.Atoi(tid)
	err := db.Raw("select tag_name from tag where id = ?", tidint).Scan(&tag).Error
	if err != nil {
		return false
	}
	return true
}

//删除标签
func TagDelById(tid string) (string, error) {
	db := database.DBSQL

	tag := database.Tag{}
	tidInt, _ := strconv.Atoi(tid)
	db.Where("id = ?", tidInt).First(&tag)
	tName := tag.TagName
	if err := db.Delete(&tag).Error; err != nil {
		return "", err
	}
	return tName, nil
}

//标签更名
func TagCnameById(tid, tname string) (string, error) {
	db := database.DBSQL
	tag := database.Tag{}
	tidint, _ := strconv.Atoi(tid)
	tag.Id = tidint
	tag.TagName = tname
	if err := db.Save(&tag).Error; err != nil {
		return "", err
	}
	if err := db.Where("id = ?", tidint).First(&tag).Error; err != nil {
		return "", err
	}
	newname := tag.TagName
	return newname, nil
}
