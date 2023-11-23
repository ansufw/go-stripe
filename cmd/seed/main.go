package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/aysf/go-stripe/internal/driver"
)

type config struct {
	db struct {
		dsn string
	}
}

func main() {
	var cfg config

	flag.StringVar(&cfg.db.dsn, "dsn", "ananto:secret@tcp(localhost:3307)/widgets?parseTime=true&tls=false", "(dsn) database string")

	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()

	query := `CREATE TABLE IF NOT EXISTS widgets (
				id INT,
				name VARCHAR(255)
				);`

	_, err = conn.ExecContext(ctx, query)
	if err != nil {
		log.Fatal("error exec query: ", err.Error())
		os.Exit(1)
	}

	// example insert
	// INSERT INTO Customers (CustomerName, ContactName, Address, City, PostalCode, Country)
	// VALUES ('Cardinal', 'Tom B. Erichsen', 'Skagen 21', 'Stavanger', '4006', 'Norway');

	query = `INSERT INTO widgets (id, name)
				VALUES (1, 'fun widget');`

	_, err = conn.ExecContext(ctx, query)
	if err != nil {
		log.Fatal("error exec query: ", err.Error())
		os.Exit(1)
	}

	log.Println("sql init succeed")

}
