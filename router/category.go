package router

import (
	"encoding/json"
	"fmt"

	"github.com/eymen-iron/web-api/db"
	"github.com/gofiber/fiber/v2"
)

func GetPostsByCategory(c *fiber.Ctx) error {

	category := c.Params("name")
	db, err := db.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Veri tabanina baglanirken hata ", "success": false})
	}

	defer db.Close()

	const query = `SELECT * FROM posts WHERE post_cat = ?`

	rows, err := db.Query(query, category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Sorgu Yapilirken hata ", "success": false})
	}

	defer rows.Close()

	var posts []Post // Posts kullanılmalı
	for rows.Next() {
		var postContent string
		var post Post // Posts kullanılmalı
		err := rows.Scan(&post.Id, &post.Title, &post.Url, &post.Desc, &postContent, &post.Language, &post.Date, &post.Redate)
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
