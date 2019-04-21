package common

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"time"
	"github.com/emirpasic/gods/utils"
	"os"
)

const dbTable =
	"CREATE TABLE `events` (" +
	"`uid` INTEGER PRIMARY KEY AUTOINCREMENT," +
	"`username` VARCHAR(64) NOT NULL," +
	"`channel` VARCHAR(64) NOT NULL," +
	"`created` DATETIME NOT NULL," +
	"`words` VARCHAR(128) NOT NULL," +
	"`text` TEXT NOT NULL" +
	");"

// DB file name for package scope
var dbFile string

/**
	https://astaxie.gitbooks.io/build-web-application-with-golang/en/05.3.html
 */
func InitDBcheck(initialDBFile string) {
	dbFile = initialDBFile

	db, err := sql.Open("sqlite3", dbFile)
	checkErr(err, true)

	// query
	rows, err := db.Query("SELECT * FROM events")
	if checkErr(err, false) {
		createTables(db)
	}

	/* TODO: Use a structure and avoid redundancies */
	var uid int
	var username string
	var channel string
	var created time.Time
	var words string
	var text string

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&uid, &username, &channel, &created, &words, &text)
		checkErr(err, true)
		//fmt.Println(uid)
		//fmt.Println(username)
		//fmt.Println(channel)
		//fmt.Println(created)
	}

	db.Close()		// DB close
}

/**
	Dynamic table creation.
	@db		sql.DB		DB Connection pointer
 */
func createTables(db *sql.DB) {
	_, err := db.Exec(dbTable)
	checkErr(err, true)
}

func checkErr(err error, interrupt bool) bool {
	result := false
	if err != nil {
		Display(err.Error(), true, true )
		result = true
		if interrupt {
			os.Exit(1)
		}
	}
	return result
}

/**
	Add an alert and useful information to the DB
	@username	string
	@channel	string
	@words		string
	@text		string
 */
func AddAlertToDB(username string, channel string, words string, text string ){
	db, err := sql.Open("sqlite3", dbFile)
	checkErr(err, true)

	Display( "Storing alert from channel " + channel + " detected word: " + words, false, true )
	stmt, err := db.Prepare("INSERT INTO events(username, channel, created, words, text) values(?,?,DateTime('now','localtime'),?,?)")
	checkErr(err, true)

	res, err := stmt.Exec(username, channel, words, text)
	checkErr(err, true)

	id, err := res.LastInsertId()
	checkErr(err, true)

	Display( "New DB alert ID " + utils.ToString(id)   , false, true )
}

/**
	Check if an alert was already triggered on a channel before X minutes.
	@chn	 	string	Channel name
	@timeWait 	int 	Minutes after an alert while no other alerts will be triggered from the same channel
 */
func GetAlertFromDBLastMins(chn string, timeWait int) bool {
	result := false

	db, err := sql.Open("sqlite3", dbFile)
	checkErr(err, true)

	/** TODO: SWITCH TO PREPARED STATEMENTS TO FOLLOW THE BEST PRACTICES */
	queryStr := "SELECT * FROM events WHERE channel = '" + chn + "' AND created > strftime('%Y-%m-%d %H:%M', datetime('now','localtime'), '-" + utils.ToString(timeWait) + " minute')"
	rows, err := db.Query(queryStr)
	if checkErr(err, false) {
		Display( "Error with db select query. "   , false, true )
		Display(queryStr   , false, true )
	}

	/* TODO: Use a structure and avoid redundancies */
	var uid int
	var username string
	var channel string
	var created time.Time
	var words string
	var text string

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&uid, &username, &channel, &created, &words, &text)
		if checkErr(err, false){
			// Alert recently detected on the same channel
			result = false
		}else{
			result = true
		}
	}
	db.Close()		// DB close

	return result
}
