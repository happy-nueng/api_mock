package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Mock struct {
	Method   string
	Url      string
	Response interface{}
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
		fmt.Println(m.GetMethod(), m.GetUrl())
		app.Add(m.GetMethod(), m.GetUrl(), func(c *fiber.Ctx) error {
			return c.JSON(m.GetResponse())
		})
	}

	app.Listen(":8080")
}

func Convert(mock string) Mock {
	m := Mock{}
	mock = strings.TrimSpace(mock)
	mocks := strings.Split(mock, "\n")

	mm := strings.Split(strings.TrimSpace(mocks[0]), " ")

	m.Method = strings.TrimSpace(mm[0])
	m.Url = strings.TrimSpace(mm[1])
	jsonText := strings.TrimSpace(strings.Join(mocks[1:], "\n"))
	json.Unmarshal([]byte(jsonText), &m.Response)
	return m
}

func (m *Mock) GetMethod() string {
	return m.Method
}

func (m *Mock) GetUrl() string {
	return m.Url
}

func (m *Mock) GetResponse() interface{} {
	return m.Response
}
