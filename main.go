package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"reflect"
	"simple-api/auth"
	"simple-api/middleware"
	"simple-api/models"
)

func rowToStruct(rows *sql.Rows, dest interface{}) error {
	destv := reflect.ValueOf(dest).Elem()
	args := make([]interface{}, destv.Type().Elem().NumField())
	for rows.Next() {
		rowp := reflect.New(destv.Type().Elem())
		rowv := rowp.Elem()

		for i := 0; i < rowv.NumField(); i++ {
			args[i] = rowv.Field(i).Addr().Interface()
		}
		if err := rows.Scan(args...); err != nil {
			return err
		}
		destv.Set(reflect.Append(destv, rowv))
	}
	return nil
}

func setRouter() *gin.Engine {
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Fatal(errEnv)
	}
	conn := os.Getenv("POSTGRESS_URL")
	db, err := gorm.Open("postgres", conn)

	if err != nil {
		log.Fatal(err)
	}
	Migrate(db)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	r.POST("/student", middleware.AuthValid, func(c *gin.Context) {
		postHandler(c, db)
	})

	r.GET("/student", middleware.AuthValid, func(c *gin.Context) {
		getAllHandler(c, db)
	})

	r.GET("/student/:student_id", middleware.AuthValid, func(c *gin.Context) {
		getHandler(c, db)
	})

	r.PUT("/student/:student_id", middleware.AuthValid, func(c *gin.Context) {
		putHandler(c, db)
	})

	r.DELETE("/student/:student_id", middleware.AuthValid, func(c *gin.Context) {
		deleteHandler(c, db)
	})

	r.POST("/login", auth.LoginHalder)

	return r

}

func postHandler(c *gin.Context, db *gorm.DB) {
	var json models.Student

	c.Bind(&json)
	db.Create(&json)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    json,
	})
}
func putHandler(c *gin.Context, db *gorm.DB) {
	var data models.Student
	student_id := c.Param("student_id")
	if db.Find(&data, "stundent_id = ?", student_id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}
	req := data
	c.Bind(&req)
	db.Model(&data).Where("stundent_id = ?", student_id).Update(req)
	c.JSON(http.StatusOK, gin.H{
		"message": "success update",
		"data":    data,
	})

}

func deleteHandler(c *gin.Context, db *gorm.DB) {
	student_id := c.Param("student_id")
	var data models.Student

	if db.Find(&data, "stundent_id = ?", student_id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	db.Delete(&data, "stundent_id = ?", student_id)
	c.JSON(http.StatusOK, gin.H{
		"message": "success delete",
	})
}

func getHandler(c *gin.Context, db *gorm.DB) {
	var data models.Student
	stundentID := c.Param("student_id")

	//db.Where("stundent_id = ?", stundentID).Find(&data).RecordNotFound()
	if db.Find(&data, "stundent_id = ?", stundentID).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    data,
	})
}

func getAllHandler(c *gin.Context, db *gorm.DB) {
	var data []models.Student
	db.Find(&data)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    data,
	})
}
func main() {
	r := setRouter()
	r.Run(":8080")
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Student{})
	data := models.Student{}
	if db.Find(&data).RecordNotFound() {
		fmt.Println("Seeder running")
		seederUser(db)
	} else {
		fmt.Println("Data already exists")
	}
}

func seederUser(db *gorm.DB) {
	data := models.Student{
		Stundent_ID:      1,
		Stundent_name:    "Alan",
		Stundent_age:     20,
		Stundent_address: "jl.babakan sari",
		Stundent_phone:   "08123456789",
	}
	db.Create(&data)
}
