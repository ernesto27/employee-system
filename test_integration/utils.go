package test_integration

import (
	"employees-system/db"
	"os"

	"github.com/joho/godotenv"
)

type Test struct {
	name            string
	statusExpected  int
	jsonBody        string
	messageExpected string
	countResponse   int
	apiToken        string
	queryParam      string
}

func InitEnvTest() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
}

func GetDBInstanceTest() *db.Mysql {
	host := os.Getenv("DB_HOST_TEST")
	user := os.Getenv("DB_USERNAME_TEST")
	password := os.Getenv("DB_PASSWORD_TEST")
	port := os.Getenv("DB_PORT_TEST")
	database := os.Getenv("DB_DATABASE_TEST")

	dbInstance, err := db.NewMysql(host, user, password, port, database)
	if err != nil {
		panic(err)
	}

	return dbInstance
}
