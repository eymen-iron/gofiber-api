package router

import (
	"encoding/json"
	"github.com/eymen-iron/web-api/db"
	"github.com/gofiber/fiber/v2"
)

func GetPostByName(c *fiber.Ctx) error {
	db, err := db.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}
	defer db.Close()

	name := c.Params("name")

	const query = `SELECT * FROM posts WHERE post_uri = ?`

	rows, err := db.Query(query, name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Sorgu Yapilirken hata ", "success": false})
	}
	defer rows.Close()
	var post Post
	for rows.Next() {
		var postContent string
		var pst Post // Posts kullanılmalı
		err := rows.Scan(&pst.Id, &pst.Title, &pst.Url, &pst.Desc, &postContent, &pst.Language, &pst.Date, &pst.Redate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"title":   "Hata",
				"message": "Veriler cekilirken bir hata oluştu",
				"icon":    "warning",
			})
		}
		err = json.Unmarshal([]byte(postContent), &pst.Content)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"title":   "Hata",
				"message": "Json verileri çevirilirken bir hata oluştu",
				"icon":    "warning",
			})
		}
		post = pst

	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"title":   "Basarili",
		"message": "Veriler başarıyla çekildi",
		"icon":    "success",
		"data":    post,
	})

}
