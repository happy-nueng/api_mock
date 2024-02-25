package main

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Mock struct {
	Body   string
	Method string
	Path   string
}

func main() {
	content, err := os.ReadFile("mock.txt")
	if err != nil {
		panic(err)
	}

	mock := strings.Split(string(content), "###")
	app := fiber.New()

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
	for _, v := range r {
		println(v.Method, " ", v.Path)
	}
	app.Listen(":8080")
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
