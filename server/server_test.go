package main

import (
	"testing"
)

// Тест получения thumbnail'ов с Youtube
func TestGetImgFromYoutube(t *testing.T) {
	// Ожидаемый размер изображения к видеоролику с адресом ниже 22793 байта
	// !!!НЕ НАДЕЖНЫЙ ТЕСТ!!! но, возможно, лучше чем ничего,
	// ПРОВАЛИТСЯ если ролик на ютубе УДАЛЯТ или ИЗМЕНЯТ изображение и в т.п. случаях
	url := "https://youtu.be/Qy1iSu9lO34"
	exp := 22793

	img, err := getImgFromYoutube(url)
	if err != nil {
		t.Errorf("Unable to get image from youtube:   %s/n", err)
	}
	if l := len(img); l != exp {
		t.Errorf("Unable to get image from youtube: Expected size image = %v Added size image = %v\n", exp, l)
	}

	// Ожидаемый размер изображения к видеоролику с не существующим URL 1097 байта
	// Чуть более надежный тест, но тоже !!!НЕ НАДЕЖНЫЙ ТЕСТ!!! но, возможно, лучше чем ничего.
	// Провалится если на Youtube поменяютизображение для не существующих роликов
	url = "https://youtu.be/net_takoy_url"
	exp = 1097

	img, err = getImgFromYoutube(url)
	if err != nil {
		t.Errorf("Unable to get image from youtube:   %s/n", err)
	}
	if l := len(img); l != exp {
		t.Errorf("Unable to get image from youtube: Expected size image = %v Added size image = %v\n", exp, l)
	}
}
