// Command migrate-sqlite-to-pg copies all data from a SQLite database file
// into a Postgres database whose schema was created by model/migrations/postgres.
//
// Usage:
//
//	go run ./cmd/migrate-sqlite-to-pg <sqlite-file> [postgres-dsn]
//
// If postgres-dsn is omitted, the MIGRATE_DSN environment variable is used.
// The target tables are truncated first, so the destination is fully replaced
// by the source data. Identity sequences are reset to MAX(id) afterwards.
package main

import (
	"fmt"
	"os"

	"github.com/blacksheepaul/timelog/model/gen"

	_ "github.com/ncruces/go-sqlite3/embed"
	sqlite "github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

// tables lists every table in FK-safe insert order (parents before children).
var tables = []struct {
	name string
	dst  any
}{
	{"categories", &[]gen.Category{}},
	{"tasks", &[]gen.Task{}},
	{"timelogs", &[]gen.Timelog{}},
	{"constraints", &[]gen.Constraint{}},
	{"metrics", &[]gen.Metric{}},
	{"metric_records", &[]gen.MetricRecord{}},
	{"webauthn_credentials", &[]gen.WebauthnCredential{}},
	{"temp_passwords", &[]gen.TempPassword{}},
	{"audit_logs", &[]gen.AuditLog{}},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: migrate-sqlite-to-pg <sqlite-file> [postgres-dsn]")
		os.Exit(1)
	}
	sqlitePath := os.Args[1]

	dsn := os.Getenv("MIGRATE_DSN")
	if len(os.Args) >= 3 {
		dsn = os.Args[2]
	}
	if dsn == "" {
		fmt.Fprintln(os.Stderr, "postgres DSN not provided and MIGRATE_DSN is not set")
		os.Exit(1)
	}

	src, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{Logger: glogger.Default.LogMode(glogger.Silent)})
	if err != nil {
		panic(fmt.Errorf("open sqlite: %w", err))
	}

	dst, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: glogger.Default.LogMode(glogger.Silent)})
	if err != nil {
		panic(fmt.Errorf("open postgres: %w", err))
	}

	// Replace target contents entirely, including the baseline seed categories.
	names := make([]string, len(tables))
	for i, t := range tables {
		names[i] = t.name
	}
	truncate := fmt.Sprintf("TRUNCATE %s CASCADE", joinNames(names))
	if err := dst.Exec(truncate).Error; err != nil {
		panic(fmt.Errorf("truncate target: %w", err))
	}

	for _, t := range tables {
		// Unscoped: also copy soft-deleted rows.
		if err := src.Unscoped().Table(t.name).Find(t.dst).Error; err != nil {
			panic(fmt.Errorf("read %s: %w", t.name, err))
		}
		n := countOf(t.dst)
		if n == 0 {
			fmt.Printf("%-24s 0 rows\n", t.name)
			continue
		}
		if err := dst.Create(t.dst).Error; err != nil {
			panic(fmt.Errorf("write %s: %w", t.name, err))
		}
		fmt.Printf("%-24s %d rows\n", t.name, n)

		// Explicit IDs were inserted; advance the identity sequence.
		reset := fmt.Sprintf(
			"SELECT setval(pg_get_serial_sequence('%[1]s', 'id'), COALESCE((SELECT MAX(id) FROM %[1]s), 1), (SELECT MAX(id) FROM %[1]s) IS NOT NULL)",
			t.name,
		)
		if err := dst.Exec(reset).Error; err != nil {
			panic(fmt.Errorf("reset sequence %s: %w", t.name, err))
		}
	}

	fmt.Println("import done")
}

func joinNames(names []string) string {
	out := names[0]
	for _, n := range names[1:] {
		out += ", " + n
	}
	return out
}

func countOf(v any) int {
	switch s := v.(type) {
	case *[]gen.Category:
		return len(*s)
	case *[]gen.Task:
		return len(*s)
	case *[]gen.Timelog:
		return len(*s)
	case *[]gen.Constraint:
		return len(*s)
	case *[]gen.Metric:
		return len(*s)
	case *[]gen.MetricRecord:
		return len(*s)
	case *[]gen.WebauthnCredential:
		return len(*s)
	case *[]gen.TempPassword:
		return len(*s)
	case *[]gen.AuditLog:
		return len(*s)
	}
	return 0
}
