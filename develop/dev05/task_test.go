package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
)

func TestGrep(t *testing.T) {
	goCmdArgs := [][]string{
		{"run", "task.go", "-A=2", "-n", "b.nana", "_income.txt"},
		{"run", "task.go", "-A=2", "-c", "b.nana", "_income.txt"},
		{"run", "task.go", "-A=1", "-B=1", "1", "_income.txt"},
		{"run", "task.go", "-C=1", "-n", "1", "_income.txt"},
		{"run", "task.go", "-v", "1", "_income.txt"},
		{"run", "task.go", "-i", "honeydew", "_income.txt"},
		{"run", "task.go", "b.nana", "_income.txt"},
		{"run", "task.go", "-F", "b.nana", "_income.txt"},
	}

	for i := range goCmdArgs {
		cmd := exec.Command("go", goCmdArgs[i]...)
		cmd.CombinedOutput()

		// Открытие файла для чтения ожидаемых результатов
		filename := fmt.Sprintf("_expected%v.txt", i+1)
		fileExpected, err := os.Open(filename)
		if err != nil {
			t.Errorf("Ошибка открытия файла: %v", err)
		}
		defer fileExpected.Close()

		// Открытие файла для чтения полученных результатов
		fileSorted, err := os.Open("_result.txt")
		if err != nil {
			t.Errorf("Ошибка открытия файла: %v", err)
		}
		defer fileSorted.Close()

		expected, _ := io.ReadAll(fileExpected)
		result, _ := io.ReadAll(fileSorted)

		if !bytes.Equal(result, expected) {
			t.Errorf("Неожиданный результат: фактический результат = %v\nожидаемый результат = %v\nаргументы - %v", result, expected, goCmdArgs[i])
		}
	}
}
