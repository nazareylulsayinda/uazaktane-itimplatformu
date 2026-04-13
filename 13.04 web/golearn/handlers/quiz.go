package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golearn/database"
	"golearn/models"
)

type QuestionInput struct {
	Text    string `json:"text" binding:"required"`
	OptionA string `json:"option_a" binding:"required"`
	OptionB string `json:"option_b" binding:"required"`
	OptionC string `json:"option_c" binding:"required"`
	OptionD string `json:"option_d" binding:"required"`
	Correct string `json:"correct" binding:"required"` // A, B, C, D
}

type QuizInput struct {
	Questions []QuestionInput `json:"questions" binding:"required"`
}

type QuizSubmit struct {
	Answers map[string]string `json:"answers" binding:"required"` // question_id -> A, B, C, D
}

// @Summary Get quiz of a lesson
// @Tags Quiz
// @Produce json
// @Param id path int true "Lesson ID"
// @Security ApiKeyAuth
// @Success 200 {object} models.Quiz
// @Failure 404 {object} map[string]string
// @Router /api/lessons/{id}/quiz [get]
func GetQuiz(c *gin.Context) {
	lessonID := c.Param("id")
	var quiz models.Quiz
	if err := database.DB.Preload("Questions").Where("lesson_id = ?", lessonID).First(&quiz).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}
	c.JSON(http.StatusOK, quiz)
}

// @Summary Create quiz for a lesson
// @Tags Quiz
// @Accept json
// @Produce json
// @Param id path int true "Lesson ID"
// @Param body body handlers.QuizInput true "Quiz Questions"
// @Security ApiKeyAuth
// @Success 201 {object} models.Quiz
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/lessons/{id}/quiz [post]
func CreateQuiz(c *gin.Context) {
	lessonIDStr := c.Param("id")
	lessonID, _ := strconv.Atoi(lessonIDStr)
	userID, _ := c.Get("user_id")

	var lesson models.Lesson
	if err := database.DB.First(&lesson, lessonID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
		return
	}

	var course models.Course
	database.DB.First(&course, lesson.CourseID)

	if course.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this course"})
		return
	}

	var input QuizInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quiz := models.Quiz{LessonID: uint(lessonID)}
	database.DB.Create(&quiz)

	for _, qi := range input.Questions {
		q := models.Question{
			Text:    qi.Text,
			OptionA: qi.OptionA,
			OptionB: qi.OptionB,
			OptionC: qi.OptionC,
			OptionD: qi.OptionD,
			Correct: qi.Correct,
			QuizID:  quiz.ID,
		}
		database.DB.Create(&q)
	}

	database.DB.Preload("Questions").First(&quiz, quiz.ID)
	c.JSON(http.StatusCreated, quiz)
}

// @Summary Submit quiz answers
// @Tags Quiz
// @Accept json
// @Produce json
// @Param id path int true "Quiz ID"
// @Param body body handlers.QuizSubmit true "Answers"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/quiz/{id}/submit [post]
func SubmitQuiz(c *gin.Context) {
	quizIDStr := c.Param("id")
	userID, _ := c.Get("user_id")

	var input QuizSubmit
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var quiz models.Quiz
	if err := database.DB.Preload("Questions").First(&quiz, quizIDStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	score := 0
	total := len(quiz.Questions)

	for _, question := range quiz.Questions {
		qIDStr := fmt.Sprintf("%d", question.ID)
		if answer, ok := input.Answers[qIDStr]; ok {
			if answer == question.Correct {
				score++
			}
		}
	}

	var percent float64 = 0
	if total > 0 {
		percent = float64(score) / float64(total) * 100
	}

	result := models.QuizResult{
		UserID:  userID.(uint),
		QuizID:  quiz.ID,
		Score:   score,
		Total:   total,
		Percent: percent,
	}
	database.DB.Create(&result)

	c.JSON(http.StatusOK, gin.H{
		"score":   score,
		"total":   total,
		"percent": percent,
		"message": "Quiz submitted successfully",
	})
}
