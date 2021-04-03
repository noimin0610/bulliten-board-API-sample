package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// import する必要はあるが、ソースコード内で参照しないので blank import
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Message struct {
	Name      string `json:"name"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

func mustGetEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

func initSocketConnectionPool() (*sql.DB, error) {
	// [START cloud_sql_postgres_databasesql_create_socket]
	var (
		dbUser                 = mustGetEnv("DB_USER")                  // e.g. 'my-db-user'
		dbPwd                  = mustGetEnv("DB_PASS")                  // e.g. 'my-db-password'
		instanceConnectionName = mustGetEnv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
		dbName                 = mustGetEnv("DB_NAME")                  // e.g. 'my-database'
	)

	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	var dbURI string
	dbURI = fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", dbUser, dbPwd, dbName, socketDir, instanceConnectionName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_postgres_databasesql_create_socket]
}

func configureConnectionPool(dbPool *sql.DB) {
	// [START cloud_sql_postgres_databasesql_limit]

	// Set maximum number of connections in idle connection pool.
	dbPool.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	dbPool.SetMaxOpenConns(7)

	// [END cloud_sql_postgres_databasesql_limit]

	// [START cloud_sql_postgres_databasesql_lifetime]

	// Set Maximum time (in seconds) that a connection can remain open.
	dbPool.SetConnMaxLifetime(1800)

	// [END cloud_sql_postgres_databasesql_lifetime]
}

func AllMessages() ([]Message, error) {
	db, err := initSocketConnectionPool()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT name, text, timestamp FROM messages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.Name, &m.Text, &m.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func AppendMessage(name string, text string) error {
	db, err := initSocketConnectionPool()
	if err != nil {
		return err
	}

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	timestamp := time.Now().In(jst).Format("2006-01-02 15:04:05")
	log.Printf("parameters: name=%s, text=%s, timestamp=%s", name, text, timestamp)

	_, err = db.Exec("INSERT INTO messages(name, text, timestamp) VALUES ($1, $2, $3)", name, text, timestamp)
	if err != nil {
		return err
	}

	return nil
}

func handleGet(w http.ResponseWriter) error {
	m, err := AllMessages()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	j, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.Write(j)
	return nil
}

func Messages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// どのドメインからでも CORS を許可する
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch r.Method {
	case http.MethodGet:
		handleGet(w)

	case http.MethodPost:
		name := r.FormValue("name")
		text := r.FormValue("text")
		if name == "" || text == "" {
			log.Printf("parameters: name=%s, text=%s", name, text)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := AppendMessage(name, text); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodOptions: // POST 前のプリフライトリクエスト用
		w.Header().Set("Allow", "OPTIONS, GET, POST")

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
