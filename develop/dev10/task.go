package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port
go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные из сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

func main() {
	// Определение флагов
	flagTO := flag.String("timeout", "10s", "connection timeout")

	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Println("Usage: go run task.go [--timeout=timeout] host port")
		return
	}

	host := flag.Arg(0)
	port := flag.Arg(1)

	timeOut, err := time.ParseDuration(*flagTO)
	if err != nil {
		fmt.Println("Error parsingg timeout duration:", err)
	}

	conn, err := net.DialTimeout("tcp", host+":"+port, timeOut)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer conn.Close()

	fmt.Println("Connected to", conn.RemoteAddr())

	// Канал для завершения работы
	sigChan := make(chan os.Signal, 1)
	// Передача сигналов ОС в канал для завершения работы
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Канал для чтения STDIN
	inputChan := make(chan string)

	// Чтение STDIN и запись в сокет
	go readFromSTDIN(inputChan)

	// Чтение из сокета и вывод в STDOUT
	go readFromConnection(conn)

	for {
		select {
		case <-sigChan:
			fmt.Println("\nClosing connection...")

			return
		case input, ok := <-inputChan:
			if !ok {
				fmt.Println("\nClosing connection...2")

				return
			}

			if strings.HasSuffix(input, `\n`) {
				input = strings.TrimSuffix(input, `\n`)
			}
			input += "\r\nHost: " + host + `\r\n\r\n`

			_, err := conn.Write([]byte(input))
			if err != nil {
				fmt.Println("Error write to the connection:", err)
			}
		}
	}
}

func readFromSTDIN(inputChan chan string) {
	buf := make([]byte, 4096)

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println(err)
			} else {
				fmt.Println("Error reading from STDIN:", err)
			}

			close(inputChan)

			return
		}

		input := string(buf[:n])

		inputChan <- input
	}
}

func readFromConnection(conn net.Conn) {
	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by remote host")
			} else {
				fmt.Println("Error reading from connection", err)
			}

			break
		}

		fmt.Println(string(buf[:n]))
	}

	os.Exit(0)
}
