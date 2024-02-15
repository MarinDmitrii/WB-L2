package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	// Создаём папку для сохранения всех файлов и переходим в неё
	if err := os.MkdirAll("downloaded", os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	if err := os.Chdir("downloaded"); err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		os.Exit(1)
	}

	// Получаем URL из командной строки
	flag.Parse()
	args := flag.Args()

	// Проверяем, передан ли URL в качестве аргумента
	if len(args) == 0 {
		fmt.Println("Usage: go run task.go 'URL'")
		os.Exit(1)
	}

	baseURL := args[0]
	baseDir := ""

	// Проверяем, передан ли baseDir в качестве аргумента
	if len(args) >= 2 {
		// baseDir = args[1]
		baseDir = filepath.Join(baseDir, args[1])

		// Создаем папку для сохранения файлов
		if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
			fmt.Printf("Error creating folder: %v\n", err)
		}
	}

	outputDir := filepath.Join(baseDir, getDirectoryName(baseURL))

	var wg sync.WaitGroup
	var mu sync.Mutex

	mapLinks := make(map[string]bool, 1000)
	mapLinks[baseURL] = true

	file, err := Wget(baseURL, baseURL, baseDir, outputDir, &wg, &mu, mapLinks, 0)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Complete downloading URL - %v\nCreated file: %v\n", baseURL, file)
	}
}

func getDirectoryName(URL string) string {
	URL = strings.TrimSuffix(URL, `/`)

	absURL, err := url.Parse(URL)
	if err != nil {
		fmt.Printf("Error rapsing URL %v: %v\n", URL, err)
	}

	directoryName := absURL.Host + strings.ReplaceAll(absURL.Path, `/`, "_")

	return directoryName
}

func Wget(baseURL, URL, baseDir, outputDir string, wg *sync.WaitGroup, mu *sync.Mutex, mapLinks map[string]bool, depth int) (string, error) {
	if depth > 10 {
		return "", nil
	}

	// Скачиваем файл
	file, err := downloadFile(baseURL, URL, baseDir, outputDir, mu)
	if err != nil {
		return "", fmt.Errorf("Error downloanig %v: %v\n", URL, err)
	}

	if file == "" {
		return "", nil
	}

	// Если скачанный файл - html-документ, парсим ссылки и считываем данные, для преобразования ссылок на относительные
	if filepath.Ext(file) == ".html" && depth <= 0 {
		parsedBaseURL, err := url.Parse(baseURL)
		if err != nil {
			fmt.Printf("Error parsing baseURL: %v\n", err)
		}

		parsedURL, err := url.Parse(URL)
		if err != nil {
			fmt.Printf("Error parsing URL: %v\n", err)
		}

		if parsedBaseURL.Host != parsedURL.Host {
			depth += 10
		}

		// Создаем папку для сохранения файлов
		if baseURL == URL {
			if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
				fmt.Printf("Error creating folder: %v\n", err)
			}
		}

		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file %v: %v\n", file, err)
		}

		newData := string(data)

		// Меняем базовую ссылку на относительную
		newURL := "./" + strings.Replace(file, baseDir+`\`, `/`, -1)

		if baseDir != "" {
			newURL = "./" + strings.Replace(file, baseDir+`\`, "", 1)
		}

		urlQuote := `"` + URL + `"`
		newURLQuote := `"` + newURL + `"`

		if strings.Contains(newData, urlQuote) {
			mu.Lock()
			newData = strings.ReplaceAll(newData, urlQuote, newURLQuote)
			mu.Unlock()
		}

		// Получаем HTTP-ответ от сервера
		resp, err := http.Get(URL)
		if err != nil {
			return "", fmt.Errorf("Error getting response from URL: %v\n", err)
		}
		defer resp.Body.Close()

		// Парсим HTML-документ
		links, err := parseHTML(resp.Body)
		if err != nil {
			fmt.Printf("Error parsing URL: %v\n", err)
		}

		// Скачиваем файлы
		for _, link := range links {
			oldLink := link

			// Исключаем якорные ссылки
			if strings.HasPrefix(link, `#`) {
				continue
			}

			if len(link) > 2 {
				baseURL = strings.TrimSuffix(baseURL, `/`)
				link = strings.TrimSuffix(link, `/`)
			}

			parsedBaseURL, err := url.Parse(baseURL)
			if err != nil {
				fmt.Printf("Error parsing baseURL: %v\n", err)
			}

			parsedURL, err := url.Parse(link)
			if err != nil {
				fmt.Printf("Error parsing URL: %v\n", err)
			}

			if parsedURL.Scheme == "" {
				parsedURL = parsedBaseURL.ResolveReference(&url.URL{Path: link})

				tempURL := strings.ReplaceAll(parsedURL.String(), `%3F`, `?`)

				parsedURL, err = url.Parse(tempURL)
				if err != nil {
					fmt.Printf("Error parsing URL: %v\n", err)
				}
			}

			// Исключаем mailto:, tel:, tg: ссылки
			if strings.HasPrefix(parsedURL.String(), "mailto") ||
				strings.HasPrefix(parsedURL.String(), "tel") ||
				strings.HasPrefix(parsedURL.String(), "tg") {
				continue
			}

			// Исключаем базовые ссылки
			if parsedBaseURL.String() == parsedURL.String() && len(link) > 2 {
				continue
			}

			// Исключаем повторяющиеся ссылки
			if (mapLinks[parsedURL.String()] || mapLinks[oldLink]) && len(link) > 2 {
				continue
			}

			if !strings.HasSuffix(baseURL, "/") {
				baseURL = baseURL + "/"
			}

			wg.Add(1)

			mapLinks[oldLink] = true

			go func(link, oldLink string) {
				defer wg.Done()

				// Вызываем Wget для рекурсивного скачивания дополнительных документов
				newFile, err := Wget(baseURL, link, baseDir, outputDir, wg, mu, mapLinks, depth+1)
				if err != nil {
					fmt.Println(err)
				}

				if newFile != "" {
					// Если файл скачан, преобразуем ссылку на относительную
					newLink := "./" + strings.Replace(newFile, baseDir+`\`, `/`, -1)

					if baseDir != "" {
						newLink = "./" + strings.Replace(newFile, baseDir+`\`, "", 1)
					}

					linkQuote := `"` + oldLink + `"`
					newLinkQuote := `"` + newLink + `"`

					if strings.Contains(newData, linkQuote) {
						mu.Lock()
						newData = strings.ReplaceAll(newData, linkQuote, newLinkQuote)
						mu.Unlock()
					}
				}
			}(parsedURL.String(), oldLink)
		}
		wg.Wait()

		// Записываем изменённые данные в файл
		if newData != "" {
			if err := os.WriteFile(file, []byte(newData), os.ModePerm); err != nil {
				fmt.Printf("Error writing data in file %v: %v\n", file, err)
			}
		}
	}

	return file, nil
}

func downloadFile(baseURL, URL, baseDir, outputDir string, mu *sync.Mutex) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Генерируем имя по URL
	fileName, err := getFileName(baseURL, URL, baseDir, outputDir, resp)
	if err != nil {
		return "", err
	}

	// Генерируем уникальное имя в случае, если файл с таким именем уже существует
	fileName, oldFileNames := getUniqueFileName(fileName, mu)

	if fileName == "" {
		return "", nil
	}

	// Сохраняем файл
	err = saveFile(fileName, resp)
	if err != nil {
		return "", err
	}

	// Проверяем, не являются ли старые и новый файлы идентичными
	for _, oldFileName := range oldFileNames {
		if oldFileName == "" {
			continue
		}

		equal, err := compareFiles(fileName, oldFileName)
		if err != nil {
			return "", err
		}

		if equal {
			// Если старый и новый файлы являются идентичными, то удаляем новый файл
			if err = os.Remove(fileName); err != nil {
				return "", err
			}

			fileName = oldFileName
			break
		}
	}

	return fileName, nil
}

func getFileName(baseURL, URL, baseDir, outputDir string, resp *http.Response) (string, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("Error parsing URL %v: %v\n", baseURL, err)
	}

	parsedURL, err := url.Parse(URL)
	if err != nil {
		return "", fmt.Errorf("Error parsing URL %v: %v\n", URL, err)
	}

	contentType := resp.Header.Get("Content-Type")
	contentDisposition := resp.Header.Get("Content-Disposition")

	var fileName string

	// Проверяем тип содержимого и задаём имя файла:
	if strings.Contains(contentDisposition, "filename=") {
		// Используем расширение из Content-Disposition, если оно есть
		results := strings.Split(contentDisposition, "filename=")
		result := strings.Trim(results[len(results)-1], `"`)

		if parsedURL.String() == parsedBaseURL.String() {
			fileName = filepath.Join(baseDir, result)
		} else {
			fileName = filepath.Join(outputDir, result)
		}
	} else {
		if parsedURL.String() == parsedBaseURL.String() {
			fileName = outputDir
		} else if parsedURL.Path == "" || parsedURL.Path == `/` {
			fileName = filepath.Join(outputDir, strings.ReplaceAll(parsedURL.Host, ":", "_"))
		} else {
			fileName = filepath.Join(outputDir, filepath.Base(parsedURL.Path))

			if strings.Contains(fileName, "?") {
				fileName = strings.Split(fileName, "?")[0]
			}
		}

		// Используем расширение из URL, если оно есть
		extension := filepath.Ext(parsedURL.Path)

		if strings.Contains(extension, "?") {
			extension = strings.Split(extension, "?")[0]
		}

		if extension == "" {
			// Если расширение не найдено, извлекаем его из Content-Type
			extensions := strings.Split(contentType, "/")

			if len(extensions) == 2 {
				if strings.Contains(extensions[1], ";") {
					extension = strings.Split(extensions[1], ";")[0]
				} else {
					extension = extensions[1]
				}
			} else {
				return "", fmt.Errorf("Incorrect Content-Type: %v\n", contentType)
			}
		}

		if !strings.Contains(fileName, extension) {
			fileName += "." + extension
		}
	}

	return fileName, nil
}

func getUniqueFileName(fileName string, mu *sync.Mutex) (string, []string) {
	idx := strings.LastIndex(fileName, ".")

	if strings.Count(fileName, ".") == 1 {
		idx = len(fileName) - 1
	}

	fileBody := fileName[:idx]
	fileExt := fileName[idx:]

	if half := len(fileBody) / 2; half > 90 {
		fileBody = fileBody[half:]
		fileName = fileBody + fileExt
	}

	oldFileNames := make([]string, 4)

	// mu.Lock()
	// defer mu.Unlock()

	for i := 1; ; i++ {
		if i > 2 {
			// Too much copies
			return "", nil
		}

		// Проверяем, существует ли файл с таким именем
		_, err := os.Stat(fileName)
		if err == nil {
			// Если файл существует, изменяем имя
			oldFileNames[i-1] = fileName

			num := strconv.Itoa(i)

			fileName = fileBody + "(" + num + ")" + fileExt

			continue
		}

		break
	}

	return fileName, oldFileNames
}

func saveFile(fileName string, resp *http.Response) error {
	if _, err := os.Stat(fileName); err == nil {
		return fmt.Errorf("File %v already exists\n", fileName)
	}

	// Создаём файл
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("Error creating file %v\n", fileName)
	}

	// Читаем содержимое ответа и сохраняем в файл
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("Error copying file %v\n", fileName)
	}
	file.Close()

	return nil
}

func compareFiles(fileName, oldFileName string) (bool, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return false, fmt.Errorf("Error opening file %v\n", fileName)
	}

	data, err := os.ReadFile(file.Name())
	if err != nil {
		return false, fmt.Errorf("Error reading file %v\n", fileName)
	}
	file.Close()

	oldFile, err := os.Open(oldFileName)
	if err != nil {
		return false, fmt.Errorf("Error opening file %v\n", oldFileName)
	}

	oldData, err := os.ReadFile(oldFile.Name())
	if err != nil {
		return false, fmt.Errorf("Error reading file %v\n", oldFileName)
	}
	oldFile.Close()

	ignoreLine := `<!-- page generated`

	stringData := string(data)
	stringOldData := string(oldData)

	if strings.Contains(stringOldData, ignoreLine) && strings.Contains(stringData, ignoreLine) {
		stringOldData = strings.Split(stringOldData, ignoreLine)[0] + `\n`
		stringData = strings.Split(stringData, ignoreLine)[0] + `\n`

		oldData = []byte(stringOldData)
		data = []byte(stringData)
	}

	if bytes.Equal(oldData, data) {
		return true, nil
	}

	return false, nil
}

func parseHTML(body io.Reader) ([]string, error) {
	links := make([]string, 0)
	tokenizer := html.NewTokenizer(body)

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			// Достигнут конец файла
			return links, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()

			for _, attr := range token.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					links = append(links, attr.Val)
				}
			}
		}
	}
}
