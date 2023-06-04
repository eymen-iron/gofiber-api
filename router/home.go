package router

import (
	"encoding/json"
	"fmt"
	"github.com/eymen-iron/web-api/db"
	"github.com/gofiber/fiber/v2"
)

type Post struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Url     string `json:"url"`
	Desc    string `json:"desc"`
	Content []struct {
		ID      int         `json:"id"`
		Content interface{} `json:"content"`
		Type    string      `json:"type"`
	} `json:"content"`
	Language string `json:"language"`
	Date     string `json:"date"`
	Redate   string `json:"redate"`
}

func GetAllPosts(c *fiber.Ctx) error {
	db, err := db.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}
	defer db.Close()

	row, err := db.Query("SELECT * FROM posts")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"title":   "Hata",
			"message": "Sql sorgusunda bir hatayla karşılaşıldı",
			"icon":    "warning",
		})
	}
	defer row.Close()
	var posts []Post // Posts kullanılmalı
	for row.Next() {
		var postContent string
		var post Post // Posts kullanılmalı
		err := row.Scan(&post.Id, &post.Title, &post.Url, &post.Desc, &postContent, &post.Language, &post.Date, &post.Redate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"title":   "Hata",
				"message": "Veriler cekilirken bir hata oluştu",
				"icon":    "warning",
			})
			fmt.Println(err)
		}
		err = json.Unmarshal([]byte(postContent), &post.Content)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"title":   "Hata",
				"message": "Json verileri çevirilirken bir hata oluştu",
				"icon":    "warning",
			})
		}
		posts = append(posts, post)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"title":   "Basarili",
		"message": "Veriler başarıyla çekildi",
		"icon":    "success",
		"data":    posts,
	})
}
