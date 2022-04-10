package database

import (
	"database/sql"
	"sync"

	_ "modernc.org/sqlite"
)

// add - Заголовок запроса на добавление
// записи в базу данных
var add = `
INSERT INTO thumbnails (
    urlVideo, img
) VALUES (
    ?, ?
)
`

// get - Заголовок запроса извлечения изображения (Thumbail'а)
// по url видеоролика Youtube из базы данных
var get = ` 
SELECT img FROM thumbnails
WHERE urlVideo == ?
`

// schemaSQL - заголовок схемы базы данных
var schemaSQL = `
CREATE TABLE IF NOT EXISTS thumbnails (
    urlVideo TEXT PRIMARY KEY,
    img BLOB
);

CREATE INDEX IF NOT EXISTS thumbnails_urlVideo ON thumbnails(urlVideo);
`

// Thumbail - запись в базе данных
type Thumbail struct {
	urlVideo string // url видиоролика с Youtube
	img      []byte // соответствующий ролику еhumbail
}

// DB - это база данных thumbail'ов.
type DB struct {
	mu  sync.Mutex
	sql *sql.DB

	// Предварительно скомпеллированные запросы
	add *sql.Stmt // Добавление записи
	get *sql.Stmt // Извлечение изображения
}

// NewDB создает/открывает базу данных.
func NewDB(dbFile string) (*DB, error) {
	db := DB{}
	var err error

	// создает/открывает файл базы данных
	db.sql, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	// создает/проверяет наличие таблицы Thumbail,
	// соответствующей schemaSQL
	if _, err = db.sql.Exec(schemaSQL); err != nil {
		return nil, err
	}

	// предварительная компеляция заголовка add
	db.add, err = db.sql.Prepare(add)
	if err != nil {
		return nil, err
	}

	// предварительная компеляция заголовка get
	db.get, err = db.sql.Prepare(get)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// Add - добавляет запись в базу данных
func (db *DB) Add(urlVideo string, img []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(db.add).Exec(urlVideo, img)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Get - возвращает thumbail по url видеоролика Youtube
func (db *DB) Get(urlVideo string) ([]byte, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	tx, err := db.sql.Begin()
	if err != nil {
		return nil, err
	}

	resp, err := tx.Stmt(db.get).Query(urlVideo)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	img := make([]byte, 0)

	for resp.Next() {
		err = resp.Scan(&img)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		return img, tx.Commit()
	}
	return nil, tx.Commit()
}
