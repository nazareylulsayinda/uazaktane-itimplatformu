package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golearn/database"
	"golearn/models"
)

// @Summary Complete a lesson
// @Tags Progress
// @Produce json
// @Param id path int true "Lesson ID"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/lessons/{id}/complete [post]
func CompleteLesson(c *gin.Context) {
	lessonIDStr := c.Param("id")
	lessonID, _ := strconv.Atoi(lessonIDStr)
	userID, _ := c.Get("user_id")

	var lesson models.Lesson
	if err := database.DB.First(&lesson, lessonID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
		return
	}

	var existing models.Progress
	if err := database.DB.Where("user_id = ? AND lesson_id = ?", userID, lesson.ID).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Lesson already completed"})
		return
	}

	progress := models.Progress{
		UserID:   userID.(uint),
		LessonID: lesson.ID,
		CourseID: lesson.CourseID,
	}
	database.DB.Create(&progress)

	c.JSON(http.StatusOK, gin.H{"message": "Lesson completed successfully"})
}

// @Summary Get user progress
// @Tags Progress
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} map[string]interface{}
// @Router /api/my/progress [get]
func GetMyProgress(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var progresses []models.Progress
	database.DB.Where("user_id = ?", userID).Find(&progresses)

	// course id mapping
	courseCompletion := make(map[uint]int)
	for _, p := range progresses {
		courseCompletion[p.CourseID]++
	}

	var result []map[string]interface{}
	for cID, completed := range courseCompletion {
		var course models.Course
		database.DB.First(&course, cID)

		var totalLessons int64
		database.DB.Model(&models.Lesson{}).Where("course_id = ?", cID).Count(&totalLessons)

		percent := 0.0
		if totalLessons > 0 {
			percent = float64(completed) / float64(totalLessons) * 100
		}

		result = append(result, map[string]interface{}{
			"course_id":         cID,
			"course_title":      course.Title,
			"total_lessons":     totalLessons,
			"completed_lessons": completed,
			"percent":           percent,
		})
	}

	if len(result) == 0 {
		c.JSON(http.StatusOK, []map[string]interface{}{})
		return
	}

	c.JSON(http.StatusOK, result)
}
