package db

import (
	"chesss/config"
	"chesss/pkg/ent"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
)

func InitializeDB() *ent.Client {
	dbConfig := config.Conf.Database
	dsn := dbConfig.UserName + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port) +
		")/" + dbConfig.Database + "?charset=" + dbConfig.Charset + "&parseTime=True&loc=Local"
	client, err := ent.Open(dbConfig.Driver, dsn)
	if err != nil {
		log.Printf("failed at creating ent client with db %s, err: %v", dsn, err)
		return nil
	}

	return client
}
