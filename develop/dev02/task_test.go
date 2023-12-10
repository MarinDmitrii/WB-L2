package main

import (
	"fmt"
	"reflect"
	"testing"
)

type TestResult struct {
	result string
	err    error
}

func TestUnpack(t *testing.T) {
	testString := []string{"a4bc2d5e", "abcd", "45", "", `qwe\4\5`, `qwe\45`, `qwe\\5`, "45qwe"}

	expectedResult := []TestResult{
		{result: "aaaabccddddde", err: nil},
		{result: "abcd", err: nil},
		{result: "", err: fmt.Errorf("Некорректная строка")},
		{result: "", err: nil},
		{result: "qwe45", err: nil},
		{result: "qwe44444", err: nil},
		{result: `qwe\\\\\`, err: nil},
		{result: "44444qwe", err: nil},
	}

	for i, v := range testString {
		result, err := unpack(v)

		if result != expectedResult[i].result || !reflect.DeepEqual(err, expectedResult[i].err) {
			t.Errorf("Неожиданный результат: фактический результат = %v и ошибка = %v,\n ожидаемый результат = %v и ошибка = %v",
				result, err.Error(), expectedResult[i].result, expectedResult[i].err)
		}
	}
}
