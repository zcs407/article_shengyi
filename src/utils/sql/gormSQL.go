package sql

import (
	"articlebk/src/utils/dbtable"
	. "articlebk/src/utils/initinfo"
	"fmt"
	"strconv"
)

var (
	db = DB
)

///////////////////////////////////////////////////用户管理////////////////////////////////////////
//用户注册
func RegisterUser(userInfo *dbtable.User) (uName, uid, rid string, err error) {
	if err := db.Create(&userInfo).Error; err != nil {
		return "", "", "", err
	}
	return userInfo.Name, strconv.Itoa(userInfo.Id), strconv.Itoa(userInfo.RollerId), nil
}

//用uid和角色id验证用户是否存在
func IsExistByUidRid(uid, rid string) bool {
	user := dbtable.User{}
	uidInt, _ := strconv.Atoi(uid)
	ridInt, _ := strconv.Atoi(rid)
	if err := db.Where("id = ? AND roller_id = ?", uidInt, ridInt).First(&user).Error; err != nil {
		return false
	}
	return true
}

//用户是否存在,username查询(仅用于用户注册和登录时使用)
func UserIsExistByName(username string) bool {
	user := dbtable.User{}
	err := db.Raw("select Name from user where name = ?", username).Scan(&user).Error
	if err != nil {
		return false
	}
	return true
}

//用户是否为管理员
func UserIsAdmin(uid string) bool {
	fmt.Println("sql DB ---- ", DB)
	user := dbtable.User{}
	id, _ := strconv.Atoi(uid)
	if err := DB.Where("id = ?", id).First(&user).Error; err != nil {
		return false
	}
	return true
}

//用户是否存在,用uid查询
func UserIsExistByUid(uid string) bool {
	uidint, _ := strconv.Atoi(uid)
	user := dbtable.User{}
	if err := db.Where("id = ?", uidint).First(&user).Error; err != nil {
		return false
	}
	return true
}

//用户登录
func UserLogin(username, password string) (uname, uid, rid string, err error) {
	user := dbtable.User{}
	err = db.Where("name = ? AND password_hash = ?", username, password).Find(&user).Error
	if err != nil {
		return "", "", "", err
	}
	return user.Name, strconv.Itoa(user.Id), strconv.Itoa(user.RollerId), nil
}

//用户列表查询,仅管理员可查,可跟实际需求修改
func UserListGet() ([]dbtable.User, error) {
	var users []dbtable.User
	if err := db.Select("id, name, roller_id").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

//用户角色变更
func UserRollerUpdate(uid, username, newrolleid string) (uname, uuid, urid string, err error) {
	uidint, _ := strconv.Atoi(uid)
	user := dbtable.User{}
	err = db.Exec("UPDATE user SET roller_id = ? WHERE id = ? AND name = ?", newrolleid, uidint, username).Error
	if err != nil {
		return "", "", "", err
	}
	user.Id = uidint
	db.First(&user)
	return user.Name, strconv.Itoa(user.Id), strconv.Itoa(user.RollerId), nil
}

//用户密码修改
func UpdateUserPwd(uid, npwd string) error {
	uidInt, _ := strconv.Atoi(uid)
	user := dbtable.User{}
	if err := db.Model(&user).Where("id = ?", uidInt).Update("password_hash", npwd).Error; err != nil {
		return err
	}
	return nil
}

//根据uid和密码验证用户是否密码正确
func VerifyUserPwd(uid, opwd string) bool {
	uidInt, _ := strconv.Atoi(uid)
	user := dbtable.User{}
	err := db.Where("id = ? AND password_hash = ?", uidInt, opwd).First(&user).Error
	if err != nil {
		return false
	}
	return true
}

//删除用户
func DeleteUser(uid string) error {
	uidInt, _ := strconv.Atoi(uid)
	user := dbtable.User{}
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
func ArticleAdd(article dbtable.Article, tags []string) (string, error) {
	for _, tag := range tags {
		tid, _ := strconv.Atoi(tag)
		var tag dbtable.Tag
		if err := db.Where("id = ?", tid).Find(&tag).Error; err != nil {
			return "", err
		}
		article.Tags = append(article.Tags, &tag)
	}

	if err := db.Create(&article).Error; err != nil {
		return "", err
	}
	return strconv.Itoa(article.Id), nil
}

//查询文章列表

//查询某个标签的文章集合
//SELECT article.* FROM article INNER JOIN article_tag
// ON  article_tag.article_id = article.id WHERE  article_tag.tag_id IN (2);
func GetArticlesByTagId(tagid string) ([]*dbtable.Article, error) {
	tid, _ := strconv.Atoi(tagid)
	var articles []*dbtable.Article
	var images []dbtable.Image
	tag := dbtable.Tag{}
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

func ArticleSelectBySpecial(sid string) {
	var specials []dbtable.Special
	article := dbtable.Article{}
	err := db.Find(&article, sid).Error
	if err != nil {
		fmt.Println("查不到")
	}
	err = db.Model(&article).Related(&specials, "Specials").Error
	if err != nil {
		return
	}
	return
}

//更新文章
func ArticleUpdate(article dbtable.Article) error {
	if err := db.Update(&article).Error; err != nil {
		return err
	}
	return nil
}

//删除文章
func ArticleDel(aid string) (cPath string, err error) {
	article := dbtable.Article{}
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
	article := dbtable.Article{}
	err := db.Raw("select title from article where title = ?", title).Scan(&article).Error
	if err != nil {
		return false
	}
	return true
}

///////////////////////////////////////////////图片管理///////////////////////////////////////////
func ArticleImageAdd(imgurl, imgPath, aid string) error {
	image := dbtable.Image{}
	image.ImageUrl = imgurl
	image.ImagePath = imgPath
	err := db.Exec(`INSERT INTO image (article_id,image_url,image_path) VALUES (?,?,?)`, aid, imgurl, imgPath).Error
	if err != nil {
		return err
	}
	return nil
}

func ArticleImageDelByAid(aid string) []dbtable.Image {
	images := []dbtable.Image{}
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
func SpecialAdd(special *dbtable.Special) (string, error) {
	if err := db.Create(&special).Error; err != nil {
		return "", err
	}
	return strconv.Itoa(special.Id), nil
}

//按pid查询专题列表
func SpecialListByPid(pid int) ([]dbtable.Special, error) {
	specials := []dbtable.Special{}
	err := db.Where("pid = ?", pid).Find(&specials).Error
	if err != nil {
		return nil, err
	}
	return specials, nil
}

//按id查询专题列表
func SpecialListById(id int) ([]dbtable.Special, error) {
	specials := []dbtable.Special{}
	err := db.Where("pid = ?", id).Find(&specials).Error
	if err != nil {
		return nil, err
	}
	//defer db.Close()
	return specials, nil
}

//更新专题名称
func SpecialCname(sid, newName string) error {
	sidInt, _ := strconv.Atoi(sid)
	special := dbtable.Special{}
	err := db.Model(&special).Where("id = ?", sidInt).Update("name = ?", newName).Error
	if err != nil {
		return err
	}
	return nil
}

//删除专题
func SpecialDel(sid string) error {
	special := dbtable.Special{}
	sidInt, _ := strconv.Atoi(sid)
	err := db.Where("id = ?", sidInt).Delete(&special).Error
	if err != nil {
		return err
	}
	return nil
}

//是否有子专题
func SpecialHasSub(sid string) bool {
	special := dbtable.Special{}
	sidInt, _ := strconv.Atoi(sid)
	err := db.Where("pid = ?", sidInt).First(&special).Error
	if err != nil {
		return false
	}
	return true
}

//通过专题名判断是否存在
func IsexistSpecial(specialname string) bool {
	special := dbtable.Special{}
	err := db.Raw("select special_name from column where special_name = ?", specialname).Scan(&special).Error
	if err != nil {
		return false
	}
	return true
}

//通过专题id判断是否存在
func IsExistSpecialBySid(sid string) bool {
	sidInt, _ := strconv.Atoi(sid)
	special := dbtable.Special{}
	err := db.Where("id = ?", sidInt).Find(&special).Error
	if err != nil {
		return false
	}
	return true
}

//////////////////////////////////////////标签管理/////////////////////////////////////////
//创建标签
func TagAdd(tname string) (string, error) {
	var tag dbtable.Tag
	tag.TagName = tname
	if err := db.Create(&tag).Error; err != nil {
		return "", err
	}
	return strconv.Itoa(tag.Id), nil

}

//标签是否存在
func IsexistTag(tname string) bool {
	tag := dbtable.Tag{}
	err := db.Raw("select tag_name from tag where tag_name = ?", tname).Scan(&tag).Error
	if err != nil {
		return false
	}
	return true
}

//使用id查询标签是否存在
func IsexistTagById(tid string) bool {
	tag := dbtable.Tag{}
	tidint, _ := strconv.Atoi(tid)
	err := db.Raw("select tag_name from tag where id = ?", tidint).Scan(&tag).Error
	if err != nil {
		return false
	}
	return true
}

//删除标签
func TagDelById(tid string) (string, error) {
	tag := dbtable.Tag{}
	tidint, _ := strconv.Atoi(tid)
	db.Where("id = ?", tidint).First(&tag)
	tname := tag.TagName
	if err := db.Delete(&tag).Error; err != nil {
		return "", err
	}
	return tname, nil
}

//标签更名
func TagCnameById(tid, tname string) (string, error) {
	tag := dbtable.Tag{}
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

//根据标签查文章
//SELECT article.* FROM article INNER JOIN article_tag
// ON  article_tag.article_id = article.id WHERE  article_tag.tag_id IN (2);
