package googledbcontroller

import (
	"context"
	"database/sql"
	"fmt"
)

const table = "google_users"

var ctx = context.Background()

func CheckIfUserIdExists(token string, DB *sql.DB) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM %s WHERE g_token = $1)`, table)
	var exists bool
	if err := DB.QueryRowContext(ctx, query, token).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func InsertGoogleUserID(user_id string, g_token string, DB *sql.DB) error {
	query := fmt.Sprintf(`INSERT INTO %s (user_id, g_token) VALUES ($1, $2)`, table)
	if _, err := DB.ExecContext(ctx, query, user_id, g_token); err != nil {
		return err
	}
	return nil
}

func GetGoogleUserID(g_token string, DB *sql.DB) (string, error) {
	query := fmt.Sprintf(`SELECT user_id FROM %s WHERE g_token = $1`, table)
	var user_id string
	if err := DB.QueryRowContext(ctx, query, g_token).Scan(&user_id); err != nil {
		return "", err
	}
	return user_id, nil
}
