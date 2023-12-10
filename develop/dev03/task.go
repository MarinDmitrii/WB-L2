package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)

Основное.
Поддержать ключ:
-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное.
Поддержать ключи:
-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func getColumn(line string, column int) string {
	fields := strings.Fields(line)
	if column > 0 && column <= len(fields) {
		return fields[column-1]
	}

	return line
}

func removeDuplicates(lines []string) []string {
	unique := make(map[string]bool)
	var result []string

	for _, v := range lines {
		if _, exist := unique[v]; !exist {
			unique[v] = true
			result = append(result, v)
		}
	}

	return result
}

func main() {
	// Определение флагов
	flagK := flag.Int("k", 0, "Указание колонки для сортировки (по умолчанию 0 - вся строка)")
	flagN := flag.Bool("n", false, "Сортировать по числовому значению")
	flagR := flag.Bool("r", false, "Сортировать в обратном порядке")
	flagU := flag.Bool("u", false, "Не выводить повторяющиеся строки")

	flag.Parse()

	// Открытие файла для чтения
	file1, err1 := os.Open("_unsorted.txt")
	if err1 != nil {
		fmt.Println("Ошибка открытия файла:", err1)
		os.Exit(1)
	}
	defer file1.Close()

	// Считывание строк из файла
	var lines []string
	scanner := bufio.NewScanner(file1)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		os.Exit(1)
	}

	// Функция сравнения строк для сортировки
	compare := func(i, j int) bool {
		s1 := getColumn(lines[i], *flagK)
		s2 := getColumn(lines[j], *flagK)

		// Преобразование в числа, если указан флаг -n
		if *flagN {
			num1, err1 := strconv.Atoi(s1)
			num2, err2 := strconv.Atoi(s2)
			if err1 == nil && err2 == nil {
				switch {
				case *flagR && num1 < num2:
					return false
				case *flagR && num1 > num2:
					return true
				case !*flagR && num1 < num2:
					return true
				case !*flagR && num1 > num2:
					return false
				}
			}
		}

		// Сравнение строк
		result := strings.Compare(s1, s2)

		// Применение флага -r
		if *flagR {
			return result > 0
		}

		return result < 0
	}

	// Применение флага -u
	if *flagU {
		lines = removeDuplicates(lines)
	}

	sort.SliceStable(lines, compare)

	// Открытие файла для записи
	file2, err2 := os.Create("_sorted.txt")
	if err2 != nil {
		fmt.Println("Ошибка открытия файла:", err2)
		os.Exit(1)
	}
	defer file2.Close()

	// Записываем строки в файл
	for idx, value := range lines {
		file2.WriteString(value)
		if idx != len(lines)-1 {
			file2.WriteString("\n")
		}
	}
}
