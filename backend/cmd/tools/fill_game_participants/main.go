package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func loadDSN() (string, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}

	host := viper.GetString("service.database.host")
	user := viper.GetString("service.database.user")
	password := viper.GetString("service.database.password")
	port := viper.GetInt("service.database.port")
	dbName := viper.GetString("service.database.name")
	sslMode := viper.GetString("service.database.ssl-mode")
	timeZone := viper.GetString("settings.timezone")

	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s TimeZone=%s",
		user, password, dbName, host, port, sslMode, timeZone,
	), nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: go run ./cmd/tools/fill_game_participants <game_id>")
	}
	gameID := os.Args[1]

	dsn, err := loadDSN()
	if err != nil {
		log.Fatalf("read config: %v", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	var seriesID string
	if err := db.QueryRowContext(ctx, `SELECT series_id::text FROM games WHERE id=$1 AND deleted_at IS NULL`, gameID).Scan(&seriesID); err != nil {
		log.Fatalf("get series_id by game_id: %v", err)
	}

	rows, err := db.QueryContext(ctx, `
SELECT profile_id::text
FROM series_participants
WHERE series_id=$1
ORDER BY created_at ASC
`, seriesID)
	if err != nil {
		log.Fatalf("load series participants: %v", err)
	}
	defer rows.Close()

	participantIDs := make([]string, 0, 32)
	seen := make(map[string]struct{})
	for rows.Next() {
		var pid string
		if err := rows.Scan(&pid); err != nil {
			log.Fatalf("scan participant: %v", err)
		}
		if _, ok := seen[pid]; ok {
			continue
		}
		seen[pid] = struct{}{}
		participantIDs = append(participantIDs, pid)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("iterate participants: %v", err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM game_participants WHERE game_id=$1`, gameID); err != nil {
		log.Fatalf("clear game participants: %v", err)
	}

	for _, pid := range participantIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO game_participants (game_id, profile_id) VALUES ($1,$2)`, gameID, pid); err != nil {
			log.Fatalf("insert game participant %s: %v", pid, err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("commit tx: %v", err)
	}

	log.Printf("done: game_id=%s series_id=%s participants_added=%d\n", gameID, seriesID, len(participantIDs))
}

