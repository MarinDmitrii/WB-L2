package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func contains(slice []int, element int) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}

	return false
}

func main() {
	// Определение флагов
	flagA := flag.Int("A", 0, "печатать +N строк после совпадения")
	flagB := flag.Int("B", 0, "печатать +N строк до совпадения")
	flagC := flag.Int("C", 0, "печатать ±N строк вокруг совпадения")
	flagc := flag.Bool("c", false, "количество строк")
	flagI := flag.Bool("i", false, "игнорировать регистр")
	flagV := flag.Bool("v", false, "вместо совпадения, исключать")
	flagF := flag.Bool("F", false, "точное совпадение со строкой, не паттерн")
	flagN := flag.Bool("n", false, "печатать номер строки")

	flag.Parse()

	args := flag.Args()

	if len(args) < 2 {
		fmt.Println("Необходимо указать паттерн для поиска.")
		os.Exit(1)
	}

	pattern := args[0]
	fileName := args[1]

	in := strings.Count(pattern, ".")

	// Открытие файла для чтения
	file1, err1 := os.Open(fileName)
	if err1 != nil {
		fmt.Println("Ошибка открытия файла:", err1)
		os.Exit(1)
	}
	defer file1.Close()

	// Считывание строк из файла
	var lines []string
	var result []string
	var lineNum []int
	scanner := bufio.NewScanner(file1)

	for num := 1; scanner.Scan(); num++ {
		line := scanner.Text()
		lines = append(lines, line)

		// Применение флага -F
		if *flagF || in < 1 {
			// Применение флага -i
			if *flagI {
				// Применение флага -v
				if strings.Contains(strings.ToLower(line), strings.ToLower(pattern)) != *flagV {
					lineNum = append(lineNum, num)
				}
			} else {
				if strings.Contains(line, pattern) != *flagV {
					lineNum = append(lineNum, num)
				}
			}
		} else {
			var n []int
			for idx, value := range pattern {
				if value == rune('.') {
					n = append(n, idx)
				}
			}

			newPattern := strings.Replace(pattern, ".", "", -1)
			newLines := strings.Split(line, " ")

			var builder strings.Builder

			for _, newLine := range newLines {
				for i, char := range newLine {
					if !contains(n, i) {
						builder.WriteRune(char)
					}
				}

				// Применение флага -i
				if *flagI {
					// Применение флага -v
					if strings.Contains(strings.ToLower(builder.String()), strings.ToLower(newPattern)) != *flagV {
						lineNum = append(lineNum, num)
						break
					}
				} else {
					if strings.Contains(builder.String(), newPattern) != *flagV {
						lineNum = append(lineNum, num)
						break
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		os.Exit(1)
	}

	// Обработка результатов:
	for _, v := range lineNum {
		// Применение флагов -A, -B
		contextStart := v - *flagB
		contextEnd := v + *flagA

		// Применение флага -C
		if *flagC > 0 {
			contextStart = v - *flagC
			contextEnd = v + *flagC
		}

		if contextStart < 1 {
			contextStart = 1
		}
		if contextEnd > len(lines) {
			contextEnd = len(lines)
		}

		fmt.Println(len(lines))
		fmt.Println(contextEnd)

		for i := contextStart; i <= contextEnd; i++ {
			if *flagN {
				// Применение флага -n
				newline := strconv.Itoa(i) + ":" + lines[i-1]
				result = append(result, newline)
			} else {
				result = append(result, lines[i-1])
			}
		}
	}

	// Открытие файла для записи
	file2, err2 := os.Create("_result.txt")
	if err2 != nil {
		fmt.Println("Ошибка открытия файла:", err2)
		os.Exit(1)
	}
	defer file2.Close()

	// Записываем строки в файл
	if *flagc {
		// Применение флага -c
		file2.WriteString(strconv.Itoa(len(result)))
	} else {
		for idx, value := range result {
			file2.WriteString(value)
			if idx != len(result)-1 {
				file2.WriteString("\n")
			}
		}
	}
}
