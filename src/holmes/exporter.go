package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func Export(c chan int, holmesConfig *HolmesConfig) {
	redisConn := NewRedisConn()
	defer CloseRedisConn(redisConn)
	db := OpenMysql()
	defer CloseMysql(db)
	for {
		_, accesslogLine := redisConn.BlockListRightPop("accesslog_yes", 0)
		if accesslogLine == "" {
			//fmt.Printf("%s now list have no log in accesslog_yes,continue to wait others to add log to list\n", time.Now())
			continue
		}
		accessLog := GetLog(accesslogLine)
		SaveLog(accessLog, db)
	}
	c <- 1
}

func OpenMysql() *sql.DB {
	// dataSourceName := mConf.MySQLUser + ":" + mConf.MySQLPassward + "@tcp(" + mConf.MySQLIP + ":" + mConf.MySQLPort + ")/" + mConf.MySQLDbName + "?charset=utf8"
	dataSourceName := "root:admin@tcp(127.0.0.1:3306)/access_log?charset=utf8"
	//fmt.Println(dataSourceName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("can't connect mysql")
		panic(err.Error())
	}
	return db
}

func CloseMysql(db *sql.DB) {
	db.Close()
}

func SaveLog(accessLog AccessLog, db *sql.DB) {
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO access_log VALUES(?,?,?,?,?, ?,?,?,?,?, ?,?,?,?,?, ?,?,?,?,?, ?,?,?)") // ? = placeholder
	if err != nil {
		panic(err.Error()) //TODO proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave

	_, err = stmtIns.Exec(
		accessLog.Hour,
		accessLog.Year,
		accessLog.Day,
		accessLog.Min,
		accessLog.RequestTime,
		accessLog.UpstreamResponseTime,
		accessLog.RemoteAddr,
		accessLog.UpstreamAddr,
		accessLog.Hostname,
		accessLog.Method,
		accessLog.RequestURI,
		accessLog.HttpCode,
		accessLog.BytesSent,
		accessLog.Referer,
		accessLog.UserAgent,
		accessLog.GzipRatio,
		accessLog.HttpXForwardedFor,
		accessLog.ServerAddr,
		accessLog.GUID,
		accessLog.Sec,
		accessLog.Month,
		accessLog.RequestLen,
		accessLog.ServerPort)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}
