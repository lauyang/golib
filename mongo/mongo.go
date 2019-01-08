package mongo

/*
* create by caiwy
* 20160405
 */
import (
	//	"fmt"

	"github.com/lauyang/goutils/logs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//const URL = "192.168.3.195:27017" //mongodb连接字符串

var (
	mgoSession *mgo.Session
	dataBase   = "mydb"
	mongoIp    = "ip url"
	mongoPort  = "27017"
)

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error

		logs.Info("connect mongo, ", mongoIp+":"+mongoPort)
		mgoSession, err = mgo.Dial(mongoIp + ":" + mongoPort)
		if err != nil {
			logs.Error("connect mongo fail, ", err)
			panic(err) //直接终止程序运行
		}
	}
	//最大连接池默认为4096
	return mgoSession.Clone()
}

//公共方法
// 根据id获取记录
func GetRecordById(collection string, id interface{}, result interface{}) error {
	session := getSession()
	defer session.Close()
	err := session.DB(dataBase).C(collection).FindId(id).One(result)
	if nil != err && "not found" == err.Error() {
		return nil
	}
	return err
}

// 根据条件获取一条记录
func GetOneRecordByCondition(collection string, selector bson.M, result interface{}) error {
	session := getSession()
	defer session.Close()
	err := session.DB(dataBase).C(collection).Find(selector).One(result)
	if nil != err && "not found" == err.Error() {
		return nil
	}
	return err
}

// 根据条件获取多条记录
func GetAllRecordByCondition(collection string, selector bson.M, result interface{}, sorter ...string) error {
	session := getSession()
	defer session.Close()

	query := session.DB(dataBase).C(collection).Find(selector)
	if len(sorter) > 0 && sorter[0] != "" {
		query = query.Sort(sorter...)
	}

	err := query.All(result)
	if nil != err && "not found" == err.Error() {
		return nil
	}
	return err
}

// 根据条件获取分页记录（返回记录总数）
func GetPageRecordByCondition(collection string, selector bson.M, pageSize int, curPage int, result interface{}) (int, error) {
	session := getSession()
	defer session.Close()
	// 获取记录总数
	nCount := 0
	nCount, err := GetCollectionCountByCondition(collection, selector)
	if err != nil {
		if "not found" == err.Error() {
			return nCount, nil
		}
		return nCount, err
	}

	err = session.DB(dataBase).C(collection).Find(selector).Skip((curPage - 1) * pageSize).All(result)
	if nil != err && "not found" == err.Error() {
		return nCount, nil
	}

	return nCount, err

}

//获取集合中条目
func GetCollectionCount(collection string) (int, error) {
	session := getSession()
	defer session.Close()
	return session.DB(dataBase).C(collection).Count()
}

//获取集合中符合条件的条目
func GetCollectionCountByCondition(collection string, selector bson.M) (int, error) {
	session := getSession()
	defer session.Close()
	return session.DB(dataBase).C(collection).Find(selector).Count()
}

//添加记录到集合中
func AddRecord(collection string, records bson.M) error {
	session := getSession()
	defer session.Close()
	err := session.DB(dataBase).C(collection).Insert(records)
	return err
}

//按id删除记录
func RemoveRecordById(collection string, id interface{}) error {
	session := getSession()
	defer session.Close()
	err := session.DB(dataBase).C(collection).RemoveId(id)
	return err
}

//按条件删除记录集合
func RemoveRecordByCondition(collection string, selector bson.M) (*mgo.ChangeInfo, error) {
	session := getSession()
	defer session.Close()
	changeInfo, err := session.DB(dataBase).C(collection).RemoveAll(selector)
	return changeInfo, err
}

//按id更新记录
func UpdateById(collection string, id interface{}, changer bson.M) error {
	session := getSession()
	defer session.Close()
	err := session.DB(dataBase).C(collection).UpdateId(id, changer)
	return err
}

//按条件更新记录集合
func UpdateByCondition(collection string, selector bson.M, changer bson.M) (*mgo.ChangeInfo, error) {
	session := getSession()
	defer session.Close()
	changeInfo, err := session.DB(dataBase).C(collection).UpdateAll(selector, changer)
	return changeInfo, err
}

/**
 * 执行查询，此方法可拆分做为公共方法
 * [SearchPerson description]
 * @param {[type]} collectionName string [description]
 * @param {[type]} query          bson.M [description]
 * @param {[type]} sort           bson.M [description]
 * @param {[type]} fields         bson.M [description]
 * @param {[type]} skip           int    [description]
 * @param {[type]} limit          int)   (results      []interface{}, err error [description]
 */
func SearchPerson(collectionName string, query bson.M, sort string, fields bson.M, skip int, limit int) (results []interface{}, err error) {
	exop := func(c *mgo.Collection) error {
		t := c.Find(query)
		if sort != "" {
			t = t.Sort(sort)
		}
		if fields != nil {
			t = t.Select(fields)
		}
		if skip > 0 {
			t = t.Skip(skip)
		}
		if limit > 0 {
			t = t.Limit(limit)
		}
		return t.All(&results)
		//		return c.Find(query).Sort(sort).Select(fields).Skip(skip).Limit(limit).All(&results)
	}
	err = WitchCollection(collectionName, exop)
	return
}

//获取collection对象
func WitchCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}
