package main

import (
	"bufio"
	model "convert/pkg"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const csvFileType string = "Csv"
const prnFileType string = "Prn"

var swapIndex = 0 // Смещение индекса по строке, если столбец выровнен по правому краю

func GetUsers(typeFile string, users []model.User) {
	// Создаем файл
	file, err := os.Create("convert" + typeFile + ".html")
	// Если возникла ошибка, выходим из программы
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	// После записи данных, закрываем файл
	defer file.Close()

	// Указываем путь к файлу с шаблоном
	pathTemp := filepath.Join("html", "usersDynamicPage.html")
	// Создаем html-шаблон
	tmpl, err := template.ParseFiles(pathTemp)
	if err != nil {
		fmt.Println("Template error", err.Error())
		return
	}
	// Исполняем именованный шаблон "users", передавая туда массив со списком пользователей
	err = tmpl.ExecuteTemplate(file, "users", users)
	if err != nil {
		fmt.Println("ExecuteTemplate error", err.Error())
		return
	}
}

// Считывает файл и преобразует его содержимое в HTML-код
func GetHTML(fileType string, fileName string) []model.User {
	var html []model.User
	var headerColumnIndexArray []int

	// Открываю файл и сканирую его
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error open file:", err)
		return []model.User{}
	}
	// После чтения данных, закрываем файл
	defer file.Close()

	// Декодер cvs по стандарту ISO8859_1
	decoded := transform.NewReader(file, charmap.ISO8859_1.NewDecoder())
	scanner := bufio.NewScanner(decoded)

	isHeader := true
	for scanner.Scan() {

		row := scanner.Text()
		var users model.User

		// Динамическое определение индексов заголовка
		if isHeader {
			headerColumnIndexArray = getHeaderColumnIndex(row)
			isHeader = false
		} else {
			users, err = getTableLine(fileType, headerColumnIndexArray, row)
		}
		// Проверка ошибок, обнаруженных при обработке строк
		if err != nil {
			fmt.Println("Error row", err)
			return []model.User{}
		}
		html = append(html, users)
	}
	if fileType == prnFileType {
		return html[1:]
	}
	return html
}

// Получаем индексы начала заголовок столбцов для prn
func getHeaderColumnIndex(header string) []int {
	return []int{
		strings.Index(header, "Birthday"),
		strings.Index(header, "Limit"),
		strings.Index(header, "Mobile"),
		strings.Index(header, "Postcode"),
		strings.Index(header, "Address"),
		strings.Index(header, "Name")}
}

func getTableLine(fileType string, headerColumnIndexArray []int, line string) (model.User, error) {
	var value model.User
	var err error

	if strings.EqualFold(fileType, csvFileType) {
		value, err = getCSVToTableRow(regexBody, line)
	}

	if strings.EqualFold(fileType, prnFileType) {
		value = getPRNToTableRow(headerColumnIndexArray, line)
	}

	return value, err
}

// Парсинг строки PRN в html-строку
func getPRNToTableRow(headerColumnIndexArray []int, line string) model.User {
	var s2 []string

	for _, element := range headerColumnIndexArray {
		tmp0 := getChars(line, 0, element-1)
		tmp1 := getChars(line, element, len([]rune(line))-swapIndex)

		line = tmp0 + tmp1
		s2 = append(s2, strings.TrimSpace(tmp1))
		line = line[:len(tmp0)]
	}
	// Возвращение заполненной структуры User
	return model.User{Name: s2[5], Address: s2[4], Postcode: s2[3], Mobile: s2[2], Limit: s2[1], Birthday: s2[0]}
}

// Парсинг строки CSV в html-строку
func getCSVToTableRow(regexp *regexp.Regexp, line string) (model.User, error) {
	var users model.User = model.User{}
	var err error

	matches := regexp.FindStringSubmatch(line)
	if matches == nil {
		err = errors.New("Error while parsing the line")
	} else {
		users = model.User{
			Name:     matches[2],
			Address:  matches[4] + matches[5],
			Postcode: matches[7],
			Mobile:   matches[9],
			Limit:    matches[11],
			Birthday: matches[13]}
	}
	return users, err
}

// Возвращает N количество символов, начинающихся с заданной позиции
func getChars(line string, startPos int, endPos int) string {
	var result string

	//смещение позиции по строке
	if startPos != 0 {
		endPos -= swapIndex
		swapIndex = 0

		for startPos > 0 && startPos < len(line) && string(line[startPos]) != " " {
			startPos--
			swapIndex++
		}
	}

	i := 0
	total := 0
	for _, char := range line {
		if (i >= startPos) && (total <= (endPos - startPos)) {
			result += string(char)
			total++
		}
		i++
	}
	return result
}
