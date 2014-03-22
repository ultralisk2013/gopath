///此文件的目的，1 强制使用连接池
///			  2 强制数据库连接对象不能关闭
///使用方法，在主程序加载时使用 调用
///util.DBOpen() 初始化数据库连接
///初始化数据库连接后就可以使用DbMySqlHelper数据库连接对象查询和处理了
///例如：
///		util.MaxOpen=100
///		util.MaxIdle=10
///		util.DbConnString="数据库登陆名:密码@tcp(数据库连接IP:端口)/数据库 库名?charset=utf8"
///		util.DBOpen()
///到这里数据库初始化就完成了
///  在其它地方，只要util.DbMySql 就是数据库连接对象了，不要关闭，关闭了后面的连接就不能用了，这是全局的连接对象。
///如util.DbMySql.Exec("insert  into test(id,idb) values(?,?)", 3, 6)
package util

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type MySql struct {
	*sql.DB
}

var DbMySql MySql         //数据库连接对象
var maxOpenConns int      //应用程序池最大数据库连接数
var maxIdleConns int      //应用程序池空闲连接的最大数目
var dataSourceName string //连接字符串
const (
	DB_PING_TIME = 10 * time.Second //发送ping命令的间
)

///打开数据库连接,一般在主程序加载时使用
func DBOpen(maxOpen, maxIdle int, dataConn string) error {
	//这些是为重连接做准备的
	maxOpenConns = maxOpen
	maxIdleConns = maxIdle
	dataSourceName = dataConn
	//打开Mysql数据库连接
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
		//panic(err)
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	DbMySql.DB = db

	go dbDaemon()
	return nil
}

//数据库守护进程
func dbDaemon() {
	for {
		time.Sleep(DB_PING_TIME)
		ping()
	}
}

///发送ping命令 ,目的是让数据库连接不会在无连接时断开
func ping() {
	err := DbMySql.Ping()
	//如果Ping接出错，则执行下面的操作
	if err != nil {
		panic(err)
	}
}

//干掉Close()方法
func (this *MySql) Close() {

}
