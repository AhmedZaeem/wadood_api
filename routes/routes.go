package routes

import (
	"github.com/gin-gonic/gin"
	"wadood/auth"
	authservice "wadood/auth/services"
	"wadood/pet"
	petservice "wadood/pet/services"
)

func Setup() *gin.Engine {
	r := gin.Default()
	petService := &petservice.PetService{}
	userService := &authservice.UserService{}
	petHandler := pet.NewPetHandler(petService, userService)

	r.POST("/register", auth.Register)
	r.POST("/login", auth.Login)
	r.POST("/logout", auth.Logout)
	r.GET("/get_pets", petHandler.GetPets)
	r.POST("/add_pet", petHandler.AddPet)
	r.PUT("/edit_pet/:id", petHandler.EditPet)
	r.DELETE("/delete_pet/:id", petHandler.DeletePet)
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Wadood online")
	})
	return r
}
