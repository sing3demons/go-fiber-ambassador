package controllers

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/main/database"
	"github.com/sing3demons/main/models"
	"gorm.io/gorm"
)

type Products struct{}

type productsRespones struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (*Products) DB() *gorm.DB {
	return database.GetDB()
}

func (*Products) Chcher() *redis.Client {
	return database.Cache
}

func (tx *Products) FindProducts(c *fiber.Ctx) error {
	var products []models.Product

	if err := tx.DB().Limit(24).Find(&products).Error; err != nil {
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	serializeProduct := []productsRespones{}
	copier.Copy(&serializeProduct, &products)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"products": serializeProduct,
	})
}
func findByID(c *fiber.Ctx) (uint, error) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func (tx *Products) findByProduct(c *fiber.Ctx) (*models.Product, error) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return nil, err
	}

	var product models.Product
	if err := tx.DB().First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (tx *Products) FindProduct(c *fiber.Ctx) error {
	product, err := tx.findByProduct(c)

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	serializeProduct := productsRespones{}
	copier.Copy(&serializeProduct, &product)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"products": serializeProduct,
	})
}

type productsForm struct {
	Title       string  `json:"title" form:"title"`
	Description string  `json:"description" form:"description"`
	Image       string  `json:"image" form:"image"`
	Price       float64 `json:"price" form:"price"`
}

func (tx *Products) CreateProducts(c *fiber.Ctx) error {
	var product models.Product
	var form productsForm
	if err := c.BodyParser(&form); err != nil {
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	copier.Copy(&product, &form)

	if err := tx.DB().Create(&product).Error; err != nil {
		return c.JSON(err)
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (tx *Products) UpdateProduct(c *fiber.Ctx) error {

	product, err := tx.findByProduct(c)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var form productsForm
	if err := c.BodyParser(&form); err != nil {
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	copier.Copy(&product, &form)
	tx.DB().Model(&product).Updates(&product)

	go tx.deleteCache("products::fontend")
	go tx.deleteCache("products::backend")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"products": product,
	})
}

func (tx *Products) deleteCache(key string) {
	time.Sleep(5 * time.Second)
	tx.Chcher().Del(context.Background(), key)

}

func (tx *Products) DeleteProduct(c *fiber.Ctx) error {

	id, err := findByID(c)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := tx.DB().Unscoped().Delete(&models.Product{}, id).Error; err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (tx *Products) ProductsFrontend(c *fiber.Ctx) error {
	var products []models.Product

	ctx := context.Background()
	key := "products::fontend"

	result, err := tx.Chcher().Get(ctx, key).Result()
	if err != nil {
		tx.DB().Limit(10000).Find(&products)

		bytes, err := json.Marshal(products)

		if err != nil {
			panic(err)
		}

		if errKey := tx.Chcher().Set(ctx, key, bytes, 10*time.Second).Err(); errKey != nil {
			panic(errKey)
		}

	} else {
		json.Unmarshal([]byte(result), &products)

	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"products": products,
	})
}

func (tx *Products) ProductsBackend(c *fiber.Ctx) error {
	var products []models.Product

	ctx := context.Background()
	key := "products::backend"

	result, err := tx.Chcher().Get(ctx, key).Result()
	if err != nil {
		tx.DB().Limit(10000).Find(&products)

		bytes, err := json.Marshal(products)

		if err != nil {
			panic(err)
		}

		tx.Chcher().Set(ctx, key, bytes, 10*time.Second)

	} else {
		json.Unmarshal([]byte(result), &products)

	}
	searchedProducts := []models.Product{}

	if s := c.Query("s"); s != "" {
		lower := strings.ToLower(s)
		for _, product := range products {
			if strings.Contains(product.Title, lower) || strings.Contains(strings.ToLower(product.Description), lower) {
				searchedProducts = append(searchedProducts, product)
			}
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"products": searchedProducts,
		})
	} else {
		searchedProducts = products
	}

	if sortParam := c.Query("sort"); sortParam != "" {
		sortLower := strings.ToLower(sortParam)
		if sortLower == "asc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price < searchedProducts[j].Price
			})
		} else if sortLower == "desc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price > searchedProducts[j].Price
			})
		}
	}

	var total = len(searchedProducts)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage := 9

	var data []models.Product

	if total <= page*perPage && total >= (page-1)*perPage {
		data = searchedProducts[(page-1)*perPage : total]
	} else if total >= page*perPage {
		data = searchedProducts[(page-1)*perPage : page*perPage]
	} else {
		data = []models.Product{}
	}

	return c.JSON(fiber.Map{
		"data":      data,
		"total":     total,
		"page":      page,
		"last_page": total/perPage + 1,
	})
}
