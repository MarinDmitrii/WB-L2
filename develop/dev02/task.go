package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
- "a4bc2d5e" => "aaaabccddddde"
- "abcd" => "abcd"
- "45" => "" (некорректная строка)
- "" => ""

Дополнительное задание: поддержка escape - последовательностей
- qwe\4\5 => qwe45 (*)
- qwe\45 => qwe44444 (*)
- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.
Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func unpack(str string) (string, error) {
	// Преобразование строки в руны для корректной обработки символов Unicode
	rs := []rune(str)
	var result strings.Builder

	// Переменная для отслеживания количества цифр в строке
	var count int
	// Флаг для отслеживания символов-экранирования
	var isEscape bool

	for i, v := range rs {
		if isEscape {
			// Если предыдущий символ был экранирован, добавляем текущий символ буквально
			result.WriteRune(v)
			isEscape = false
		} else if unicode.IsDigit(v) {
			// Если текущий символ - цифра, обрабатываем повторение предыдущего символа
			count++
			num, err := strconv.Atoi(string(v))
			if err != nil {
				return "", err
			}

			if i > 0 {
				// Добавляем повторения предыдущего символа
				result.WriteString(strings.Repeat(string(rs[i-1]), num-1))
			} else {
				// Если цифра в начале строки, добавляем текущий символ буквально
				result.WriteRune(v)
			}
		} else if string(v) == `\` {
			// Если текущий символ - обратная косая черта, устанавливаем флаг экранирования
			isEscape = true
		} else {
			// Добавляем символ в результат
			result.WriteRune(v)
		}
	}

	if str == "" {
		// Если строка пуста, возвращаем пустую строку
		return "", nil
	} else if len(rs) == count {
		// Если все символы - цифры, возвращаем ошибку
		return "", fmt.Errorf("Некорректная строка")
	} else {
		// Возвращаем распакованную строку
		return result.String(), nil
	}
}

func main() {
	str, err := unpack(`45`)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(str)
}
