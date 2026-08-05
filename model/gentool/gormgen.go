package main

import (
	"os"

	_ "github.com/ncruces/go-sqlite3/embed"
	sqlite "github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gorm.io/gen"
)

func main() {
	dbType := os.Getenv("TIMELOG_GEN_DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}

	var dialector gorm.Dialector
	switch dbType {
	case "postgres":
		dsn := os.Getenv("TIMELOG_GEN_DSN")
		if dsn == "" {
			panic("env TIMELOG_GEN_DSN is not set for postgres")
		}
		dialector = postgres.Open(dsn)
	default:
		dbPath := os.Getenv("TIMELOG_GEN_DB_PATH")
		if dbPath == "" {
			panic("env TIMELOG_GEN_DB_PATH is not set")
		}
		dialector = sqlite.Open(dbPath)
	}

	db, err := gorm.Open(dialector)
	if err != nil {
		panic(err)
	}

	// 创建生成器
	g := gen.NewGenerator(gen.Config{
		OutPath:           "../model/gen",
		ModelPkgPath:      "./model/gen",
		FieldNullable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		WithUnitTest:      false,
		FieldSignable:     false,
		FieldCoverable:    false,
		Mode:              gen.WithoutContext,
	})

	// 使用数据库
	g.UseDB(db)

	// 生成所有模型
	g.ApplyBasic(
		g.GenerateAllTable()...,
	)

	// 执行生成
	g.Execute()
}
