package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// selectFields выбирает указанные поля (колонки) из строки
func selectFields(fields []string, fieldsList string) []string {
	selectedFields := make([]string, 0)
	fieldsIdx := parseFieldsList(fieldsList)

	for _, idx := range fieldsIdx {
		// Нужна ли эта проверка на 0, если она уже произведена в 47 строке?
		if idx > 0 && idx <= len(fields) {
			selectedFields = append(selectedFields, fields[idx-1])
		} else {
			// Нужен ли этот else или улчше отправить пустой слайс?
			selectedFields = append(selectedFields, "")
		}
	}

	return selectedFields
}

// parseFieldsList разбирает строку с номерами полей, разделенными запятыми
func parseFieldsList(fieldsList string) []int {
	fieldsIdx := make([]int, 0)
	fields := strings.Split(fieldsList, ",")

	for _, field := range fields {
		idx := 0

		if n, err := fmt.Sscanf(field, "%d", &idx); err == nil && n > 0 {
			fieldsIdx = append(fieldsIdx, idx)
		}
	}

	return fieldsIdx
}

func main() {
	// Определение флагов
	flagF := flag.String("f", "", "выбрать поля (колонки)")
	flagD := flag.String("d", " ", "использовать другой разделитель")
	flagS := flag.Bool("s", false, "только строки с разделителем")

	flag.Parse()

	// Проверка, что хотя бы один из флагов -f или -s установлен
	if *flagF == "" && !*flagS {
		fmt.Println("Необходимо указать хотя бы один из флагов: -f или -s")
		os.Exit(1)
	}

	// Считываем данные из STDIN
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// Проверка, что строка содержит разделитель (если флаг -s не установлен)
		if *flagS || strings.Contains(line, *flagD) {
			fields := strings.Split(line, *flagD)

			if *flagF != "" {
				// Если указан флаг -f, выбираем указанные поля
				selectedFields := selectFields(fields, *flagF)
				fmt.Println(strings.Join(selectedFields, *flagD))
			} else {
				// Если флаг -f не указан, выводим всю строку
				fmt.Println(line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения данных:", err)
		os.Exit(1)
	}
}
