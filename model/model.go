package model

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/core/logger"

	_ "github.com/ncruces/go-sqlite3/embed"
	sqlite "github.com/ncruces/go-sqlite3/gormlite"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	gl "gorm.io/gorm/logger"
)

var (
	dao  *Dao
	once sync.Once
	log  logger.Logger
)

type DBProvider interface {
	Db() *gorm.DB
}

type Dao struct {
	db    *gorm.DB
	RawDB *sql.DB
	cache *cache.Cache
}

func (d *Dao) Db() *gorm.DB {
	return d.db
}

func (d *Dao) Begin() *TxDao {
	tx := d.db.Begin()
	if tx.Error != nil {
		panic(tx.Error)
	}
	return &TxDao{Dao: d, txDb: tx}
}

type TxDao struct {
	*Dao
	txDb *gorm.DB
}

func (tx *TxDao) Db() *gorm.DB {
	return tx.txDb
}

func (tx *TxDao) Commit() error {
	return tx.txDb.Commit().Error
}

func (tx *TxDao) Rollback() error {
	return tx.txDb.Rollback().Error
}

var _ DBProvider = (*Dao)(nil)
var _ DBProvider = (*TxDao)(nil)

// NewDao opens SQLite and returns a DAO without mutating package globals.
func NewDao(cfg *config.Config, _ logger.Logger) (*Dao, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Database.Host), &gorm.Config{
		Logger: gl.Default.LogMode(gl.LogLevel(cfg.Log.ORMLogLevel)),
	})
	if err != nil {
		return nil, err
	}

	raw, err := db.DB()
	if err != nil {
		return nil, err
	}

	return &Dao{
		db:    db,
		RawDB: raw,
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}, nil
}

// RegisterDao installs a DAO constructed via NewDao for legacy call sites that
// still use GetDao during the migration.
func RegisterDao(d *Dao, loggerInstance logger.Logger) {
	once.Do(func() {
		dao = d
		log = loggerInstance
	})
}

func InitDao(cfg *config.Config, loggerInstance logger.Logger) {
	once.Do(func() {
		d, err := NewDao(cfg, loggerInstance)
		if err != nil {
			panic(err)
		}
		dao = d
		log = loggerInstance
	})
}

func GetDao() *Dao {
	if dao == nil {
		panic("dao is nil, please call InitDao first")
	}
	return dao
}

type Model struct {
	*gorm.Model
}

func (d *Dao) WriteCache(key string, value any, seconds int64) {
	exp := time.Duration(seconds) * time.Second
	d.cache.Set(key, value, exp)
}

func (d *Dao) GetCache(key string) (any, bool) {
	return d.cache.Get(key)
}

func (d *Dao) AdminGetAllCache() {
	items := d.cache.Items()

	log.Debugw("list cache items",
		"count", len(items),
	)

	for k, v := range items {
		log.Debug(fmt.Sprintf("cache item [ %s --- %v ]", k, v))
	}
}
