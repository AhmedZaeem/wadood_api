package services

import (
	"database/sql"
	"fmt"
	"net/http"
	"wadood/db"
	"wadood/messages"
	"wadood/pet/models"
)

type PetService struct{}

func (s *PetService) GetPetsByOwner(ownerID int, language string) ([]models.Pet, int, string, error) {
	var pets []models.Pet
	rows, err := db.DB.Query("SELECT id, name, age, type, owner_id FROM pets WHERE owner_id = ?", ownerID)
	if err != nil {
		return nil, http.StatusInternalServerError, messages.GetMessage("database_error", language), err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	for rows.Next() {
		var pet models.Pet
		if err := rows.Scan(&pet.ID, &pet.Name, &pet.Age, &pet.Type, &pet.OwnerID); err != nil {
			return nil, http.StatusInternalServerError, messages.GetMessage("database_error", language), err
		}
		pets = append(pets, pet)
	}
	if err := rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, messages.GetMessage("database_error", language), err
	}
	if len(pets) == 0 {
		return pets, http.StatusOK, messages.GetMessage("no_pets_found", language), nil
	}
	return pets, http.StatusOK, messages.GetMessage("get_pets_success", language), nil
}

func (s *PetService) AddPet(pet models.Pet, language string) (int, string, *models.Pet, error) {
	if pet.Name == "" || pet.Type == "" || pet.Age < 0 {
		return http.StatusBadRequest, messages.GetMessage("invalid_pet_data", language), nil, fmt.Errorf("invalid pet data")
	}
	var exists int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM pets WHERE name = ? AND owner_id = ?", pet.Name, pet.OwnerID).Scan(&exists)
	if err != nil {
		return http.StatusInternalServerError, messages.GetMessage("database_error", language), nil, err
	}
	if exists > 0 {
		return http.StatusBadRequest, messages.GetMessage("pet_name_exists", language), nil, fmt.Errorf("pet name already exists for this owner")
	}
	res, err := db.DB.Exec(
		"INSERT INTO pets (name, age, type, owner_id, gender) VALUES (?, ?, ?, ?, ?)",
		pet.Name, pet.Age, pet.Type, pet.OwnerID, pet.Gender,
	)
	if err != nil {
		return http.StatusInternalServerError, messages.GetMessage("database_error", language), nil, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return http.StatusInternalServerError, messages.GetMessage("database_error", language), nil, err
	}
	pet.ID = int(lastID)
	return http.StatusCreated, messages.GetMessage("pet_added_success", language), &pet, nil
}

func (s *PetService) EditPet(pet models.Pet, language string) (int, string, *models.Pet, error) {
	if pet.ID == 0 || pet.Name == "" || pet.Type == "" || pet.Age < 0 {
		return http.StatusBadRequest, messages.GetMessage("invalid_pet_data", language), nil, fmt.Errorf("invalid pet data")
	}
	var exists int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM pets WHERE name = ? AND owner_id = ? AND id != ?", pet.Name, pet.OwnerID, pet.ID).Scan(&exists)
	if err != nil {
		return http.StatusInternalServerError, messages.GetMessage("database_error", language), nil, err
	}
	if exists > 0 {
		return http.StatusBadRequest, messages.GetMessage("pet_name_exists", language), nil, fmt.Errorf("pet name already exists for this owner")
	}
	res, err := db.DB.Exec("UPDATE pets SET name = ?, age = ?, type = ? WHERE id = ? AND owner_id = ?", pet.Name, pet.Age, pet.Type, pet.ID, pet.OwnerID)
	if err != nil {
		return http.StatusInternalServerError, messages.GetMessage("database_error", language), nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return http.StatusNotFound, messages.GetMessage("pet_not_found", language), nil, fmt.Errorf("pet not found")
	}
	return http.StatusOK, messages.GetMessage("pet_updated_success", language), &pet, nil
}

func (s *PetService) DeletePet(petID, ownerID int, language string) (int, string, error) {
	res, err := db.DB.Exec("DELETE FROM pets WHERE id = ? AND owner_id = ?", petID, ownerID)
	if err != nil {
		return http.StatusInternalServerError, messages.GetMessage("database_error", language), err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return http.StatusNotFound, messages.GetMessage("pet_not_found", language), fmt.Errorf("pet not found")
	}
	return http.StatusOK, messages.GetMessage("pet_deleted_success", language), nil
}
