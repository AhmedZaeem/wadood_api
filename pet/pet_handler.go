package pet

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	authmodels "wadood/auth/models"
	authservices "wadood/auth/services"
	"wadood/messages"
	"wadood/pet/models"
	petservice "wadood/pet/services"
)

type Handler struct {
	petService  *petservice.PetService
	userService *authservices.UserService
}

func NewPetHandler(petService *petservice.PetService, userService *authservices.UserService) *Handler {
	return &Handler{
		petService:  petService,
		userService: userService,
	}
}

func (h *Handler) authenticate(c *gin.Context) (*authmodels.User, bool) {
	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, false
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	user, _, err := authservices.FindUserByToken(token)
	if err != nil || user == nil {
		return nil, false
	}
	return user, true
}

func (h *Handler) GetPets(c *gin.Context) {
	lang := c.GetHeader("X-Language")
	if lang == "" {
		lang = c.DefaultQuery("lang", "en")
	}
	user, ok := h.authenticate(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": messages.GetMessage("unauthorized", lang)})
		return
	}
	pets, _, _, err := h.petService.GetPetsByOwner(user.ID, lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": messages.GetMessage("get_pets_failed", lang)})
		return
	}
	if len(pets) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": messages.GetMessage("no_pets_found", lang), "pets": pets})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": messages.GetMessage("get_pets_success", lang), "pets": pets})
}

func (h *Handler) AddPet(c *gin.Context) {
	lang := c.GetHeader("X-Language")
	if lang == "" {
		lang = c.DefaultQuery("lang", "en")
	}
	user, ok := h.authenticate(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": messages.GetMessage("unauthorized", lang)})
		return
	}
	var pet models.Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.GetMessage("invalid_input", lang)})
		return
	}
	// Validate required fields
	if pet.Name == "" || pet.Gender == "" || pet.Type == "" || pet.Age < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.GetMessage("missing_pet_fields", lang)})
		return
	}
	pet.OwnerID = user.ID
	status, msg, newPet, err := h.petService.AddPet(pet, lang)
	if err != nil {
		c.JSON(status, gin.H{"message": msg})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": msg, "pet": newPet})
}

func (h *Handler) EditPet(c *gin.Context) {
	lang := c.GetHeader("X-Language")
	if lang == "" {
		lang = c.DefaultQuery("lang", "en")
	}
	user, ok := h.authenticate(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": messages.GetMessage("unauthorized", lang)})
		return
	}
	var pet models.Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.GetMessage("invalid_input", lang)})
		return
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.GetMessage("invalid_pet_id", lang)})
		return
	}
	pet.ID = id
	pet.OwnerID = user.ID
	status, msg, updatedPet, err := h.petService.EditPet(pet, lang)
	if err != nil {
		c.JSON(status, gin.H{"message": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": msg, "pet": updatedPet})
}

func (h *Handler) DeletePet(c *gin.Context) {
	lang := c.GetHeader("X-Language")
	if lang == "" {
		lang = c.DefaultQuery("lang", "en")
	}
	user, ok := h.authenticate(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": messages.GetMessage("unauthorized", lang)})
		return
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.GetMessage("invalid_pet_id", lang)})
		return
	}
	status, msg, err := h.petService.DeletePet(id, user.ID, lang)
	if err != nil {
		c.JSON(status, gin.H{"message": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}
