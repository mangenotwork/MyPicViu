package db

import (
	"MyPicViu/common/logger"
	"MyPicViu/common/utils"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"strings"
)

const TreeNodeTable = "tree_node"
const TreeNodeKey = "tree_node_key"

var DB *LocalDB
var Tables = []string{TreeNodeTable}
var ISNULL = fmt.Errorf("ISNULL")
var TableNotFound = fmt.Errorf("table notfound")

type LocalDB struct {
	Path   string
	Tables []string
	Conn   *bolt.DB
}

func getDirName(filePath string) string {
	return filePath[0:getLastIndex(filePath, `/`)]
}

func getLastIndex(str, ch string) int {
	return len(str) - len(strings.TrimRight(str, ch))
}

func checkDBFile(dbFilePath string) {

	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {

		dir := getDirName(dbFilePath)

		if dir != "" {
			if err = os.MkdirAll(dir, 0755); err != nil {
				log.Panic(err)
			}
		}

		f, fErr := os.Create(dbFilePath)
		if fErr != nil {
			log.Panic(fErr)
		}

		defer func() {
			_ = f.Close()
		}()

	}

}

func (ldb *LocalDB) Init() {

	checkDBFile(ldb.Path)

	db, err := bolt.Open(ldb.Path, 0600, nil)
	if err != nil {
		logger.Panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	for _, table := range ldb.Tables {

		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(table))

			if b == nil {
				_, err = tx.CreateBucket([]byte(table))
				if err != nil {
					logger.Panic(err)
				}
			}

			return nil
		})

		if err != nil {
			logger.Panic(err)
		}

	}
}

func NewLocalDB(tables []string, path string) *LocalDB {
	return &LocalDB{
		Path:   path,
		Tables: tables,
	}
}

func InitDB(dbPath string) {
	DB = NewLocalDB(Tables, dbPath)
	DB.Init()
	TreeData = buildTreeData()
}

func (ldb *LocalDB) Open() {
	ldb.Conn, _ = bolt.Open(ldb.Path, 0600, nil)
}

func (ldb *LocalDB) Close() {
	_ = ldb.Conn.Close()
}

func (ldb *LocalDB) ClearTable(table string) error {
	ldb.Open()
	defer func() {
		_ = ldb.Conn.Close()
	}()
	return ldb.Conn.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(table))
	})
}

func (ldb *LocalDB) Stats(table string) (bolt.BucketStats, error) {

	var stats bolt.BucketStats

	ldb.Open()

	defer func() {
		_ = ldb.Conn.Close()
	}()

	err := ldb.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b == nil {
			err := ldb.ClearTable(table)
			if err != nil {
				return err
			}
		}

		stats = b.Stats()

		return nil
	})

	return stats, err
}

func (ldb *LocalDB) Get(table, key string, data interface{}) error {

	ldb.Open()
	defer func() {
		_ = ldb.Conn.Close()
	}()

	return ldb.Conn.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(table))
		if b == nil {
			return TableNotFound
		}

		bt := b.Get([]byte(key))
		if len(bt) < 1 {
			return ISNULL
		}

		err := json.Unmarshal(bt, data)
		if err != nil {
			return err
		}

		return nil
	})
}

func (ldb *LocalDB) Set(table, key string, data interface{}) error {

	value, err := utils.AnyToJsonB(data)
	if err != nil {
		return err
	}

	ldb.Open()

	defer func() {
		_ = ldb.Conn.Close()
	}()

	return ldb.Conn.Update(func(tx *bolt.Tx) error {

	R:
		b := tx.Bucket([]byte(table))
		if b == nil {
			_, err = tx.CreateBucket([]byte(table))
			if err != nil {
				return err
			}

			goto R
		}

		err = b.Put([]byte(key), value)
		if err != nil {
			return err
		}

		return nil
	})
}

func (ldb *LocalDB) Delete(table, key string) error {
	ldb.Open()

	defer func() {
		_ = ldb.Conn.Close()

	}()

	return ldb.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b == nil {
			return fmt.Errorf("未获取到表")
		}
		if err := b.Delete([]byte(key)); err != nil {
			return err
		}
		return nil
	})
}
