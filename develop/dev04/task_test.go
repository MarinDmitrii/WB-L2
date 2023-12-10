package main

import (
	"reflect"
	"testing"
)

func TestSortString(t *testing.T) {
	input := []string{"память", "ток", "лошадь", "орел"}
	expected := []string{"амптья", "кот", "адлошь", "елор"}

	for i := 0; i < len(input); i++ {
		result := sortString(input[i])
		if expected[i] != result {
			t.Errorf("Неожиданный результат: ожидаемый результат - %v, полученный результат - %v\n", expected[i], result)
		}
	}
}

func TestGetAnagram(t *testing.T) {
	input := []string{"ток", "пятак", "тяпка", "кот", "КтО", "столик", "лиСток", "пятка", "слиток", "молоток"}
	expected := map[string][]string{
		"ток":    {"кот", "кто"},
		"пятак":  {"пятка", "тяпка"},
		"столик": {"листок", "слиток"},
	}

	result := getAnagram(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Неожиданный результат\nожидаемый результат - %v\nфактический результат - %v", expected, result)
	}
}
