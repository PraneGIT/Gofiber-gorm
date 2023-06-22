package main

import (
	"net/http"
	"os"

	"github.com/PraneGIT/goFibrePostgres/models"
	"github.com/PraneGIT/goFibrePostgres/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Cannot parse JSON"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Cannot create book",
		})
		return err
	}

	context.Status(http.StatusCreated).JSON(&fiber.Map{
		"message": "Successfully created book",
	})

	return nil

}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	booksList := &[]models.Books{}

	err := r.DB.Find(&booksList).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Cannot get books",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Successfully fetched book",
		"data":    booksList,
	})

	return nil

}

func (r *Repository) GetBook(context *fiber.Ctx) error {
	book := models.Books{}

	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Invalid Id",
		})
		return nil
	}

	err := r.DB.First(&book, id).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Cannot get book",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Successfully fetched book",
		"data":    book,
	})

	return nil

}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {

	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Invalid Id",
		})
		return nil
	}

	err := r.DB.Delete(models.Books{}, id).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Cannot delete book",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Successfully deleted book",
	})

	return nil

}

func (r *Repository) SetupRoutes(app *fiber.App) {

	api := app.Group("/api")
	api.Post("/create_book", r.CreateBook)
	api.Get("/get_books", r.GetBooks)
	api.Get("/get_book/:id", r.GetBook)
	// api.Put("/update_book/:id", r.UpdateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  "disable",
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		panic(err)
	}

	err = models.MigrateBools(db)
	if err != nil {
		panic(err)
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":9090")

}
