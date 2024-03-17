package library

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Initialize default config -> session
var storeSession = session.New(session.Config{
	Expiration:     24 * time.Hour,
	CookiePath:     "/",
	CookieSecure:   true,
	CookieHTTPOnly: true,
})

func SetSession(context *fiber.Ctx, userID string) error {
	sess, err := storeSession.Get(context)
	if err != nil {
		return err
	}

	sess.Set("user_id", userID)
	err = sess.Save()
	if err != nil {
		return err
	}

	return nil
}

func GetAllSession(context *fiber.Ctx) (map[string]interface{}, error) {
	sessionsMap := make(map[string]interface{})
	sess, err := storeSession.Get(context)
	if err != nil {
		return sessionsMap, err
	}

	keys := sess.Keys()

	if len(keys) > 0 {
		for _, key := range keys {
			value := sess.Get(key)
			sessionsMap[key] = value
		}
	}

	return sessionsMap, nil
}

func GetUserID(context *fiber.Ctx) (string, error) {
	maps, err := GetAllSession(context)
	if err != nil {
		return "", fmt.Errorf("get all session : %w", err)
	}
	if maps != nil {
		return maps[`user_id`].(string), nil
	}
	return "", nil
}

func DeleteSession(context *fiber.Ctx) (string, error) {

	userID := ""
	sess, err := storeSession.Get(context)
	if err != nil {
		return userID, err
	}

	keys := sess.Keys()
	if len(keys) > 0 {
		for _, key := range keys {
			value := sess.Get(key)
			userID = value.(string)
		}
	}

	sess.Delete("user_id")
	err = sess.Destroy()
	if err != nil {
		return userID, err
	}

	return userID, nil
}
