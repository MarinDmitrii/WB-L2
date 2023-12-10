package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestCut(t *testing.T) {
	goCmdArgs := [][]string{
		{"run", "task.go", "-f=2,5"},
		{"run", "task.go", "-f=2,5", "-d=;"},
		{"run", "task.go", "-f=2,5", "-d=;", "-s"},
	}

	// Создаем буфер для записи вывода программы
	var stdout bytes.Buffer

	// Задаем входные данные
	income := []string{
		"1 1Penalty for mr. Russel!",
		"2; 2Penalty; for; mr.; Russel!",
		"3;3Penalty;for;mr.;Russel!",
		"4-4Penalty-for-mr.-Russel!",
	}

	// Задаем ожидаемые данные
	expected := []string{
		"1Penalty Russel!\n",
		" 2Penalty; Russel!\n",
		"3Penalty;Russel!\n",
		"\n",
	}

	for i := range goCmdArgs {
		// Запускаем программу и записываем ее вывод в буфер
		cmd := exec.Command("go", goCmdArgs[i]...)
		// cmd.CombinedOutput()
		cmd.Stdin = strings.NewReader(income[i])
		// Очищаем содержимое stdout перед записью
		stdout.Reset()
		cmd.Stdout = &stdout

		err := cmd.Run()
		if err != nil {
			t.Fatalf("Ошибка выполнения команды: %v", err)
		}

		// Сравниваем вывод программы с ожидаемым результатом
		if result := stdout.String(); result != expected[i] {
			t.Errorf("Неожиданный результат №%v:\nфактический = %v\nожидаемый = %v\nаргументы - %v", i, result, expected[i], goCmdArgs[i])
		}
	}
}
