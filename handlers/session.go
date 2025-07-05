package handlers

import (
	"context"
	"time"

	"session-service/config"
	"session-service/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var ctx = context.Background()
var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: config.GetRedisAddr(),
	})
}

// POST /api/v1/session
func CreateSession(c *gin.Context) {
	var input struct {
		UserID string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Datos inválidos"})
		return
	}

	sessionID := uuid.New().String()
	now := time.Now()
	expiration := now.Add(24 * time.Hour)

	session := models.Session{
		SessionID: sessionID,
		UserID:    input.UserID,
		CreatedAt: now,
		ExpiresAt: expiration,
	}

	key := "session:" + sessionID
	err := redisClient.HSet(ctx, key, map[string]interface{}{
		"user_id":    session.UserID,
		"created_at": session.CreatedAt.Format(time.RFC3339),
		"expires_at": session.ExpiresAt.Format(time.RFC3339),
	}).Err()

	if err != nil {
		c.JSON(500, gin.H{"error": "No se pudo crear la sesión"})
		return
	}

	redisClient.ExpireAt(ctx, key, expiration)

	c.JSON(201, session)
}

// GET /api/v1/session/user/:user_id
func GetSessionsByUser(c *gin.Context) {
	userID := c.Param("user_id")
	keys, _ := redisClient.Keys(ctx, "session:*").Result()

	var sessions []models.Session

	for _, key := range keys {
		data, _ := redisClient.HGetAll(ctx, key).Result()
		if data["user_id"] == userID {
			createdAt, _ := time.Parse(time.RFC3339, data["created_at"])
			expiresAt, _ := time.Parse(time.RFC3339, data["expires_at"])
			sessions = append(sessions, models.Session{
				SessionID: key[len("session:"):],
				UserID:    data["user_id"],
				CreatedAt: createdAt,
				ExpiresAt: expiresAt,
			})
		}
	}

	c.JSON(200, sessions)
}

// DELETE /api/v1/session/:session_id
func DeleteSession(c *gin.Context) {
	sessionID := c.Param("session_id")
	key := "session:" + sessionID
	redisClient.Del(ctx, key)
	c.JSON(200, gin.H{"message": "Sesión eliminada"})
}

// DELETE /api/v1/session/user/:user_id
func DeleteSessionsByUser(c *gin.Context) {
	userID := c.Param("user_id")
	keys, _ := redisClient.Keys(ctx, "session:*").Result()

	for _, key := range keys {
		data, _ := redisClient.HGetAll(ctx, key).Result()
		if data["user_id"] == userID {
			redisClient.Del(ctx, key)
		}
	}

	c.JSON(200, gin.H{"message": "Sesiones eliminadas"})
}
