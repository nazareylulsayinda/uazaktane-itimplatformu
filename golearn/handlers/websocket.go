package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // development için cross origin kabul
	},
}

type Client struct {
	Conn     *websocket.Conn
	Username string
	RoomID   string
}

type Message struct {
	Type     string `json:"type"` // "join", "leave", "message"
	Username string `json:"username"`
	Text     string `json:"text"`
	CourseID string `json:"course_id"`
}

var (
	rooms = make(map[string]map[*Client]bool)
	mutex = sync.Mutex{}
)

// @Summary Connect to websocket room for a course
// @Tags Websocket
// @Param courseId path string true "Course ID"
// @Router /ws/classroom/{courseId} [get]
func WebSocketHandler(c *gin.Context) {
	courseID := c.Param("courseId")
	username := c.Query("username") // Basitlik için query argument'ten username veya tokenden alınabilir. Burada token doğrulandı varsayıyoruz. (Middleware eklenecek)
	if username == "" {
		username = "Anonymous"
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}

	client := &Client{Conn: conn, Username: username, RoomID: courseID}

	mutex.Lock()
	if rooms[courseID] == nil {
		rooms[courseID] = make(map[*Client]bool)
	}
	rooms[courseID][client] = true
	mutex.Unlock()

	// Broadcast Join Message
	broadcast(courseID, Message{Type: "join", Username: client.Username, Text: client.Username + " has joined the classroom", CourseID: courseID})

	defer func() {
		mutex.Lock()
		delete(rooms[courseID], client)
		if len(rooms[courseID]) == 0 {
			delete(rooms, courseID)
		}
		mutex.Unlock()
		client.Conn.Close()
		broadcast(courseID, Message{Type: "leave", Username: client.Username, Text: client.Username + " has left", CourseID: courseID})
	}()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		msg.Type = "message"
		msg.Username = client.Username
		msg.CourseID = courseID
		broadcast(courseID, msg)
	}
}

func broadcast(courseID string, msg Message) {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range rooms[courseID] {
		err := client.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("broadcast error: %v", err)
			client.Conn.Close()
			delete(rooms[courseID], client)
		}
	}
}
