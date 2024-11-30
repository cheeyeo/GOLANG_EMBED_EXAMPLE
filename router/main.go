package router

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Channel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	ID        int    `json:"id"`
	ChannelID int    `json:"channel_id"`
	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	Text      string `json:"text"`
}

func CreateUser(c *gin.Context, db *sql.DB) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("CREATING USER: %+v\n", user)

	// Insert user into database
	result, err := db.Exec("INSERT INTO users (username, password) VALUES(?, ?)", user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get ID of newly inserted user
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return ID of newly inserted user
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// Login endpoint
func Login(c *gin.Context, db *sql.DB) {
	// Parse JSON request body into User struct
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Query database for user
	row := db.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", user.Username, user.Password)

	// Get ID of user
	var id int
	err := row.Scan(&id)
	if err != nil {
		// Check if user was not found
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		// Return error if other error occurred
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// Return ID of user
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// Channel creation endpoint
func CreateChannel(c *gin.Context, db *sql.DB) {
	// Parse JSON request body into Channel struct
	var channel Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("CHANNEL: %+v\n", channel)

	// Insert channel into database
	result, err := db.Exec("INSERT INTO channels (name) VALUES (?)", channel.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("DB ERR: %s\n", err.Error())
		return
	}

	// Get ID of newly inserted channel
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return ID of newly inserted channel
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// Channel listing endpoint
func ListChannels(c *gin.Context, db *sql.DB) {
	// Query database for channels
	rows, err := db.Query("SELECT id, name FROM channels")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create slice of channels
	var channels []Channel

	// Required to fix error "database is locked 5 sqlite_busy"
	defer rows.Close()
	// Iterate over rows
	for rows.Next() {
		// Create new channel
		var channel Channel

		// Scan row into channel
		err := rows.Scan(&channel.ID, &channel.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Append channel to slice
		channels = append(channels, channel)
	}

	// Return slice of channels
	c.JSON(http.StatusOK, channels)
}

// Message creation endpoint
func CreateMessage(c *gin.Context, db *sql.DB) {
	// Parse JSON request body into Message struct
	var message Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert message into database
	result, err := db.Exec("INSERT INTO messages (channel_id, user_id, message) VALUES (?, ?, ?)", message.ChannelID, message.UserID, message.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get ID of newly inserted message
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return ID of newly inserted message
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// Message listing endpoint
func ListMessages(c *gin.Context, db *sql.DB) {
	// Parse channel ID from URL
	channelID, err := strconv.Atoi(c.Query("channelID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse optional limit query parameter from URL
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		// Set limit to 100 if not provided
		limit = 100
	}

	// Parse last message ID query parameter from URL. This is used to get messages after a certain message.
	lastMessageID, err := strconv.Atoi(c.Query("lastMessageID"))
	if err != nil {
		// Set last message ID to 0 if not provided
		lastMessageID = 0
	}

	// Query database for messages
	rows, err := db.Query("SELECT m.id, channel_id, user_id, u.username AS user_name, message FROM messages m LEFT JOIN users u ON u.id = m.user_id WHERE channel_id = ? AND m.id > ? ORDER BY m.id ASC LIMIT ?", channelID, lastMessageID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create slice of messages
	var messages []Message

	// Required to fix error "database is locked 5 sqlite_busy"
	defer rows.Close()
	// Iterate over rows
	for rows.Next() {
		// Create new message
		var message Message

		// Scan row into message
		err := rows.Scan(&message.ID, &message.ChannelID, &message.UserID, &message.UserName, &message.Text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Append message to slice
		messages = append(messages, message)
	}

	// Return slice of messages
	c.JSON(http.StatusOK, messages)
}
