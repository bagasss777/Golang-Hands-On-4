package main

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt"
	"log"
	"time"
	//"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func main() {
	//var produk1 = Product{"mie", "food", 1000, "mini market", 10000, int(time.April)}

	app := fiber.New()
	addJWTRoute(app)
	app.Listen(":3000")

}

var signingKey = []byte("secret")

type UserRequest struct {
	Username string `json:"user"`
	Password string `json:"password"`
}

type Product struct {
	Nama_Product      string    `json:"Nama Produk"`
	Jenis_Produk      string    `json:"Jenis Produk"`
	Harga_Produk      int       `json:"Harga Produk"`
	Tempat_Pembelian  string    `json:"Tempat Pembelian"`
	Nomor_Barcode     int       `json:"Nomor Barcode"`
	Tanggal_Pembelian time.Time `json:"Tanggal Pembelian"`
}

func addJWTRoute(app *fiber.App) {
	produk := []Product{
		Product{"mie", "food", 1000, "mini market", 10000, time.Now()},
		Product{"susu", "minuman", 50000, "mini market", 575857857, time.Now()},
		Product{"jus", "food", 1000, "mini market", 10000, time.Now()},
	}

	apiGroup := app.Group("/api")
	apiGroup.Post("/login", func(c *fiber.Ctx) (err error) {
		var req UserRequest
		err = c.BodyParser(&req)
		if err != nil {
			log.Printf("Error in parsing the JSON request: %v.", err)
			return
		}

		if req.Username != "admin" || req.Password != "4dm1n" {
			err = c.SendStatus(fiber.StatusUnauthorized)
			return
		}

		signJwt := jwt.New(jwt.SigningMethodHS256)

		claims := signJwt.Claims.(jwt.MapClaims)
		claims["name"] = "Admin"
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		token, err := signJwt.SignedString(signingKey)
		if err != nil {
			err = c.SendStatus(fiber.StatusInternalServerError)
			return
		}

		err = c.JSON(fiber.Map{"token": token})
		return
	})

	apiGroup.Use("/produk", jwtware.New(jwtware.Config{
		SigningKey: signingKey,
	}))
	apiGroup.Get("/produk", func(c *fiber.Ctx) error {
		return c.JSON(produk)
	})
}
