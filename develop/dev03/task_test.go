package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
)

func TestGetColumn(t *testing.T) {
	testString := "banana 10 pcs"
	expectedResult := "pcs"
	result := getColumn(testString, 3)
	if result != expectedResult {
		fmt.Errorf("Неожиданный результат: фактический результат = %v\nожидаемый результат = %v ", result, expectedResult)
	}
}

func TestRemoveDuplicates(t *testing.T) {
	testCases := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "c"},
		},
		{
			input:    []string{"1", "2", "2", "3", "4"},
			expected: []string{"1", "2", "3", "4"},
		},
		{
			input:    []string{"apple 1", "apple", "banana 1", "apple 1", "orange", "banana 2"},
			expected: []string{"apple", "banana 1", "orange", "banana 2"},
		},
	}

	for _, v := range testCases {
		result := removeDuplicates(v.input)
		if len(result) == len(v.expected) {
			for i := 0; i < len(result); i++ {
				if result[i] != v.expected[i] {
					t.Errorf("Неожиданный результат: фактический результат = %v\nожидаемый результат = %v", result, v.expected)
				}
			}
		}
	}
}

func TestSort(t *testing.T) {
	goCmdArgs := [][]string{
		{"run", "task.go", "-k=2", "-n"},
		{"run", "task.go", "-r", "-u"},
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
		fileSorted, err := os.Open("_sorted.txt")
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
