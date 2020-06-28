package data

import "context"

func GetUserExists(username string) (bool, error) {
	var exists bool
	err := db.QueryRow(context.TODO(), `
SELECT EXISTS(SELECT 1 FROM address WHERE username = $1)
	`, username).
		Scan(&exists)

	return exists, err
}
