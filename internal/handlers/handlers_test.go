package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"serviceCatalog/internal/models"
)

type HandlerTestSuite struct {
	suite.Suite
	db      *gorm.DB
	handler *Handler
	router  *gin.Engine
	// Store test data IDs
	testServiceID uint
}

func (s *HandlerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	dsn := "host=localhost user=postgres password=postgres dbname=servicecatalog_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db

	err = db.AutoMigrate(&models.Service{}, &models.Version{})
	if err != nil {
		s.T().Fatal(err)
	}

	s.handler = NewHandler(db)
	s.router = gin.Default()
	s.router.GET("/services", s.handler.ListServices)
	s.router.GET("/services/:id", s.handler.GetService)
	s.router.GET("/services/:id/versions", s.handler.GetServiceVersions)
}

func (s *HandlerTestSuite) SetupTest() {
	// Clean up existing data
	s.db.Exec("TRUNCATE TABLE versions CASCADE")
	s.db.Exec("TRUNCATE TABLE services CASCADE")
	s.db.Exec("ALTER SEQUENCE services_id_seq RESTART WITH 1")
	s.db.Exec("ALTER SEQUENCE versions_id_seq RESTART WITH 1")

	// Create test service
	service := models.Service{
		Name:        "Test Service",
		Description: "Test Description",
	}
	result := s.db.Create(&service)
	if result.Error != nil {
		s.T().Fatal(result.Error)
	}
	s.testServiceID = service.ID

	// Create test version
	version := models.Version{
		ServiceID: service.ID,
		Number:    "1.0.0",
	}
	result = s.db.Create(&version)
	if result.Error != nil {
		s.T().Fatal(result.Error)
	}

	// Verify data was created
	var count int64
	s.db.Model(&models.Service{}).Count(&count)
	if count != 1 {
		s.T().Fatal("Test service was not created")
	}
}

func (s *HandlerTestSuite) TearDownSuite() {
	db, err := s.db.DB()
	if err != nil {
		s.T().Fatal(err)
	}
	db.Close()
}

func (s *HandlerTestSuite) TestListServices() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services", nil)
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), 200, w.Code)

	var response ListServicesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), 1, len(response.Services))
	assert.Equal(s.T(), "Test Service", response.Services[0].Name)
	assert.Equal(s.T(), 1, response.Services[0].Versions)
}

func (s *HandlerTestSuite) TestGetService() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services/1", nil)
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), 200, w.Code)

	var service models.ServiceResponse
	err := json.Unmarshal(w.Body.Bytes(), &service)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), "Test Service", service.Name)
	assert.Equal(s.T(), 1, service.Versions)
}

func (s *HandlerTestSuite) TestGetServiceVersions() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services/1/versions", nil)
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), 200, w.Code)

	var versions []models.Version
	err := json.Unmarshal(w.Body.Bytes(), &versions)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), 1, len(versions))
	assert.Equal(s.T(), "1.0.0", versions[0].Number)
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
