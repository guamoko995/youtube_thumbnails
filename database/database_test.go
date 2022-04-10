package database

import (
	"fmt"
	"testing"
)

// Тест работы базы данных
func TestDatabase(t *testing.T) {
	// Создание базы временной данных (с пустым именем файла)
	b, err := NewDB("")
	if err != nil {
		t.Errorf("NewDB error: %s", err)
	}

	// Добавление записи
	err = b.Add("url1", []byte{0, 1, 2, 3})
	if err != nil {
		t.Errorf("Add error: %s", err)
	}

	// Еще одно добавление записи
	err = b.Add("url2", []byte{4, 3, 2, 1})
	if err != nil {
		t.Errorf("Add error: %s", err)
	}

	// Попытка добавление записи с уже имеющимся ключем
	// !! url должен быть уникальным !! , ожидается ошибка
	err = b.Add("url1", []byte{0, 0})
	// Ожидаемое сообщение ошибки
	exp := "constraint failed: UNIQUE constraint failed: thumbnails.urlVideo (1555)"
	if err == nil {
		t.Errorf("Add error: Expected \"%s\" Added: \"nil\"", exp)
	} else {
		if err.Error() != exp {
			t.Errorf("Add error: Expected \"%s\" Added: \"%s\"", exp, err.Error())
		}
	}

	// Чтение записи
	var img []byte
	img, err = b.Get("url2")
	if err != nil {
		t.Errorf("Get error: %s", err.Error())
	} else {
		exp := "[4 3 2 1]"
		if fmt.Sprint(img) != exp {
			t.Errorf("Ged error: Expected \"%v\" Added: \"%v\"", exp, img)
		}
	}

	// Чтение еще одной записи
	img, err = b.Get("url1")
	if err != nil {
		t.Errorf("Get error: %s", err.Error())
	} else {
		exp := "[0 1 2 3]"
		if fmt.Sprint(img) != exp {
			t.Errorf("Ged error: Expected \"%v\" Added: \"%v\"", exp, img)
		}
	}

	// Попытка чтения не существующей записи
	img, err = b.Get("url3")
	exp = "Key not found"
	if err == nil {
		t.Errorf("Get error: Expected \"%v\" Added: \"%v\"", exp, err.Error())
	} else {
		exp := "[]"
		if fmt.Sprint(img) != exp {
			t.Errorf("Ged error: Expected \"%v\" Added: \"%v\"", exp, img)
		}
	}
}
