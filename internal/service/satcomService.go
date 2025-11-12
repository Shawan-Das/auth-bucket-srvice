package service

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	auth "github.com/rest/api/internal/dbmodel/db_query"
	"github.com/rest/api/internal/model"
)

// CreateSatcomData creates a new satcom data entry
func (s *RESTService) createSatcomData(c *gin.Context) APIResponse {
	var input model.SatcomDataInput
	if !parseInput(c, &input) {
		return BuildResponse400("Invalid input provided")
	}

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	createParams := auth.CreateSatcomDataParams{
		Company:  input.Company,
		Category: input.Category,
		Type:     input.Type,
		Date:     input.Date,
		Time:     input.Time,
		DbPort:   input.DbPort,
		UiPort:   input.UiPort,
		Url:      input.URL,
		Ip:       input.IP,
		Status:   input.Status,
	}

	err := qtx.CreateSatcomData(ctx, createParams)
	if err != nil {
		_asLogger.Errorf("Error creating satcom data: %v", err)
		return BuildResponse500("Failed to create satcom data", err.Error())
	}

	return BuildResponse200("Satcom data created successfully", nil)
}

// GetSatcomDataById retrieves a satcom data entry by ID
func (s *RESTService) getSatcomDataById(c *gin.Context) APIResponse {
	idParam := c.Param("id")
	if idParam == "" {
		return BuildResponse400("ID parameter is required")
	}

	var id int32
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		return BuildResponse400("Invalid ID format")
	}

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	data, err := qtx.GetSatcomDataById(ctx, id)
	if err != nil {
		_asLogger.Errorf("Error getting satcom data: %v", err)
		return BuildResponse404("Satcom data not found", false)
	}

	response := model.SatcomDataResponse{
		ID:       data.ID,
		Company:  data.Company,
		Category: data.Category,
		Type:     data.Type,
		Date:     data.Date,
		Time:     data.Time,
		DbPort:   data.DbPort,
		UiPort:   data.UiPort,
		URL:      data.Url,
		IP:       data.Ip,
		Status:   data.Status,
	}

	return BuildResponse200("Satcom data retrieved successfully", response)
}

// GetAllSatcomData retrieves all satcom data entries
func (s *RESTService) getAllSatcomData(c *gin.Context) APIResponse {
	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	dataList, err := qtx.GetAllSatcomData(ctx)
	if err != nil {
		_asLogger.Errorf("Error getting all satcom data: %v", err)
		return BuildResponse500("Failed to retrieve satcom data", err.Error())
	}

	// Transform to response format
	responseList := make([]model.SatcomDataResponse, 0, len(dataList))
	for _, data := range dataList {
		responseList = append(responseList, model.SatcomDataResponse{
			ID:       data.ID,
			Company:  data.Company,
			Category: data.Category,
			Type:     data.Type,
			Date:     data.Date,
			Time:     data.Time,
			DbPort:   data.DbPort,
			UiPort:   data.UiPort,
			URL:      data.Url,
			IP:       data.Ip,
			Status:   data.Status,
		})
	}

	return BuildResponse200("Satcom data retrieved successfully", responseList)
}

// UpdateSatcomData updates an existing satcom data entry
func (s *RESTService) updateSatcomData(c *gin.Context) APIResponse {
	idParam := c.Param("id")
	if idParam == "" {
		return BuildResponse400("ID parameter is required")
	}

	var id int32
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		return BuildResponse400("Invalid ID format")
	}

	var input model.SatcomDataInput
	if !parseInput(c, &input) {
		return BuildResponse400("Invalid input provided")
	}

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	// Check if record exists
	_, err := qtx.GetSatcomDataById(ctx, id)
	if err != nil {
		_asLogger.Errorf("Error getting satcom data: %v", err)
		return BuildResponse404("Satcom data not found", false)
	}

	updateParams := auth.UpdateSatcomDataParams{
		Company:  input.Company,
		Category: input.Category,
		Type:     input.Type,
		Date:     input.Date,
		Time:     input.Time,
		DbPort:   input.DbPort,
		UiPort:   input.UiPort,
		Url:      input.URL,
		Ip:       input.IP,
		Status:   input.Status,
		ID:       id,
	}

	err = qtx.UpdateSatcomData(ctx, updateParams)
	if err != nil {
		_asLogger.Errorf("Error updating satcom data: %v", err)
		return BuildResponse500("Failed to update satcom data", err.Error())
	}

	return BuildResponse200("Satcom data updated successfully", nil)
}

// DeleteSatcomData deletes a satcom data entry by ID
func (s *RESTService) deleteSatcomData(c *gin.Context) APIResponse {
	idParam := c.Param("id")
	if idParam == "" {
		return BuildResponse400("ID parameter is required")
	}

	var id int32
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		return BuildResponse400("Invalid ID format")
	}

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	// Check if record exists
	_, err := qtx.GetSatcomDataById(ctx, id)
	if err != nil {
		_asLogger.Errorf("Error getting satcom data: %v", err)
		return BuildResponse404("Satcom data not found", false)
	}

	err = qtx.DeleteSatcomData(ctx, id)
	if err != nil {
		_asLogger.Errorf("Error deleting satcom data: %v", err)
		return BuildResponse500("Failed to delete satcom data", err.Error())
	}

	return BuildResponse200("Satcom data deleted successfully", nil)
}
