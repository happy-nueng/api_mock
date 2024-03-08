package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type Mock struct {
	Body   string
	Method string
	Path   string
}

type Config struct {
	Port     int
	MockFile string
}

const (
	purple    = lipgloss.Color("99")
	gray      = lipgloss.Color("245")
	lightGray = lipgloss.Color("241")
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
		viper.SetDefault("port", 8080)
		viper.SetDefault("mockfile", "mock.txt")
	}

	conf := Config{}
	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
		return
	}

	content, err := os.ReadFile(conf.MockFile)
	if err != nil {
		log.Fatalf("unable to read file, %v", err)
		return
	}

	mock := strings.Split(string(content), "###")
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	for _, v := range mock {
		m := Convert(v)
		app.Add(m.GetMethod(), m.GetPath(), func(c *fiber.Ctx) error {
			c.Accepts("application/json")
			body, err := m.GetBody()
			if err != nil {
				return err
			}
			return c.JSON(body)
		})
	}
	r := app.GetRoutes()
	rows := [][]string{}
	for _, v := range r {
		rows = append(rows, []string{v.Method, v.Path})
	}
	re := lipgloss.NewRenderer(os.Stdout)
	var (
		HeaderStyle  = re.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
		CellStyle    = re.NewStyle().Padding(0, 1).Width(14)
		OddRowStyle  = CellStyle.Copy().Foreground(gray)
		EvenRowStyle = CellStyle.Copy().Foreground(lightGray)
	)
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return HeaderStyle
			case row%2 == 0:
				return EvenRowStyle
			default:
				return OddRowStyle
			}
		}).
		Headers("Method", "Path").
		Rows(rows...)
	fmt.Println(t)
	app.Listen(fmt.Sprintf(":%v", conf.Port))
}

func Convert(mock string) Mock {
	re := regexp.MustCompile(`([A-Z]+)\s+([^ \n]+)`)
	m := Mock{}
	mock = strings.TrimSpace(mock)
	matches := re.FindStringSubmatch(mock)

	m.Method = strings.TrimSpace(matches[1])
	m.Path = strings.TrimSpace(matches[2])
	body := strings.Replace(mock, m.Method, "", 1)
	body = strings.Replace(body, m.Path, "", 1)
	m.Body = strings.TrimSpace(body)

	return m
}

func (m *Mock) GetMethod() string {
	return m.Method
}

func (m *Mock) GetPath() string {
	return m.Path
}

func (m *Mock) GetBody() (interface{}, error) {
	var body interface{}
	err := json.Unmarshal([]byte(m.Body), &body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
