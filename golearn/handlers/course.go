package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golearn/database"
	"golearn/models"
)

type CourseInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Category    string `json:"category" binding:"required"`
}

// @Summary Get all courses
// @Tags Course
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param category query string false "Category"
// @Param sort query string false "Sort order (desc/asc)" default(desc)
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/courses [get]
func GetCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}
	category := c.Query("category")
	sortOrd := c.DefaultQuery("sort", "desc")

	query := database.DB.Model(&models.Course{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	var total int64
	query.Count(&total)

	orderStr := "id desc"
	if sortOrd == "asc" {
		orderStr = "id asc"
	}

	offset := (page - 1) * limit
	var courses []models.Course
	query.Preload("Teacher").Order(orderStr).Limit(limit).Offset(offset).Find(&courses)

	c.JSON(http.StatusOK, gin.H{
		"data":  courses,
		"page":  page,
		"limit": limit,
		"total": total,
		"total_pages": math.Ceil(float64(total) / float64(limit)),
	})
}

// @Summary Get course by ID
// @Tags Course
// @Produce json
// @Param id path int true "Course ID"
// @Security ApiKeyAuth
// @Success 200 {object} models.Course
// @Failure 404 {object} map[string]string
// @Router /api/courses/{id} [get]
func GetCourse(c *gin.Context) {
	id := c.Param("id")
	var course models.Course
	if err := database.DB.Preload("Teacher").Preload("Lessons").First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	c.JSON(http.StatusOK, course)
}

// @Summary Create a course
// @Tags Course
// @Accept json
// @Produce json
// @Param body body handlers.CourseInput true "Course Details"
// @Security ApiKeyAuth
// @Success 201 {object} models.Course
// @Failure 400 {object} map[string]string
// @Router /api/courses [post]
func CreateCourse(c *gin.Context) {
	var input CourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	
	course := models.Course{
		Title:       input.Title,
		Description: input.Description,
		Category:    input.Category,
		TeacherID:   userID.(uint),
	}

	database.DB.Create(&course)
	c.JSON(http.StatusCreated, course)
}

// @Summary Update a course
// @Tags Course
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Param body body handlers.CourseInput true "Course Details"
// @Security ApiKeyAuth
// @Success 200 {object} models.Course
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/courses/{id} [put]
func UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	var course models.Course
	if err := database.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	if course.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this course"})
		return
	}

	var input CourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course.Title = input.Title
	course.Description = input.Description
	course.Category = input.Category

	database.DB.Save(&course)
	c.JSON(http.StatusOK, course)
}

// @Summary Delete a course
// @Tags Course
// @Produce json
// @Param id path int true "Course ID"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/courses/{id} [delete]
func DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	var course models.Course
	if err := database.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	if course.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this course"})
		return
	}

	database.DB.Delete(&course)
	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}
