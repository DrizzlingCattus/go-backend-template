package main

import (
	"fmt"
	gorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	echo "github.com/labstack/echo"
	"io"
	"net/http"
	"os"
)

var db *gorm.DB

type Cat struct {
	// struct embedded를 살펴볼 것
	gorm.Model
	Name string `gorm:"default:'default-cat'" json:"name"`
	Type string `gorm:"default:'default-type'" json:"type"`
}

func GetCats(c echo.Context) error {
	var cats []Cat
	db.Find(&cats)

	catMaps := make(map[int]map[string]string)
	for i, cat := range cats {
		catMaps[i] = map[string]string{
			"name": cat.Name,
			"type": cat.Type,
		}
	}
	return c.JSON(http.StatusOK, catMaps)
}

func GetCat(c echo.Context) error {
	// query param
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	// path param
	dataType := c.Param("data")
	if dataType == "string" {
		// String, JSON, HTML ... etc 모두 Response를 내부적으로 사용
		return c.String(http.StatusOK, fmt.Sprintf("[cat]: %s\n[type]: %s", catName, catType))
	} else if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	} else {
		return c.JSON(http.StatusOK, map[string]string{
			"error": "format dismatch",
		})
	}
}

func AddCat(c echo.Context) error {
	citty := new(Cat)
	// Bind body data to model
	if err := c.Bind(citty); err != nil {
		return c.String(http.StatusBadRequest, "fail")
	}
	db.Create(&Cat{Name: citty.Name, Type: citty.Type})
	return c.String(http.StatusOK, "success")
}

func TestGORM() {
	// CRUD - create
	db.Create(&Cat{Name: "cat1", Type: "water"})
	db.Create(&Cat{Name: "cat2", Type: "water"})

	// CRUD - read
	var cat1 Cat
	db.First(&cat1, "type = ?", "water") // find cat , id 1
	err := fmt.Sprintf("cat1 %s %s\n", cat1.Name, cat1.Type)
	io.WriteString(os.Stdout, err)

	// CRUD - update
	db.Model(&cat1).Update("type", "fire")
	// check update
	var cat2 Cat
	db.First(&cat2, "type = ?", "water") // select from Cat where type = 'fire'
	err = fmt.Sprintf("cat2 %s %s\n", cat2.Name, cat2.Type)
	io.WriteString(os.Stdout, err)

	// CRUD - delete
	db.Delete(&cat1)
	db.Delete(&cat2)
}

func main() {
	fmt.Println("server start!")

	tdb, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("fail to connect db")
	}
	db = tdb
	defer db.Close()

	// migrate the schema
	db.AutoMigrate(&Cat{})

	TestGORM()
	// make new echo server
	e := echo.New()
	// TODO request에 로그 붙히는 방법 -> 아마 미들웨어에 붙힐듯
	e.GET("/cats", GetCats)
	// e.GET("/cats/:data", GetCats)
	e.POST("/cat", AddCat)
	e.Logger.Fatal(e.Start(":8000"))
}
