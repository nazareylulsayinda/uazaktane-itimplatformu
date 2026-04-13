package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golearn/database"
	"golearn/models"
)

type LessonInput struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	VideoURL string `json:"video_url"`
	Order    int    `json:"order"`
}

// @Summary Get lessons of a course
// @Tags Lesson
// @Produce json
// @Param id path int true "Course ID"
// @Security ApiKeyAuth
// @Success 200 {array} models.Lesson
// @Failure 404 {object} map[string]string
// @Router /api/courses/{id}/lessons [get]
func GetLessons(c *gin.Context) {
	courseID := c.Param("id")
	
	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	var lessons []models.Lesson
	database.DB.Where("course_id = ?", courseID).Order("`order` asc").Find(&lessons)

	c.JSON(http.StatusOK, lessons)
}

// @Summary Add a lesson to a course
// @Tags Lesson
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Param body body handlers.LessonInput true "Lesson Details"
// @Security ApiKeyAuth
// @Success 201 {object} models.Lesson
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/courses/{id}/lessons [post]
func CreateLesson(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, _ := strconv.Atoi(courseIDStr)
	userID, _ := c.Get("user_id")

	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	if course.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this course"})
		return
	}

	var input LessonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lesson := models.Lesson{
		Title:    input.Title,
		Content:  input.Content,
		VideoURL: input.VideoURL,
		Order:    input.Order,
		CourseID: uint(courseID),
	}

	database.DB.Create(&lesson)
	c.JSON(http.StatusCreated, lesson)
}
