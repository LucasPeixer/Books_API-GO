package db

import (
	"database/sql"
	"log"
)

func connectDB() (*sql.DB, error) {
	connecta := "host=vividly-climbing-mallard.data-1.use1.tembo.io port=5432 user=postgres password=8fjHOEvlIsib4fhH dbname=postgres sslmode=require"

	db, err := sql.Open("postgres", connecta)

	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao pingar o banco de dados: ", err)
	}

	log.Println("Conexão estabelecida com sucesso!")

	return db, nil
}