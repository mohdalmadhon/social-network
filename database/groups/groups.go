package groups

import (
	"database/sql"
	"errors"
	"log"
	"social/internal/helpers"
	"social/internal/models"
	"strconv"
	"strings"
)

func GetGroups(db *sql.DB, userID int) ([]models.Group, error) {
	var groups []models.Group

	rows, err := db.Query(`
		SELECT id,name, users
		FROM user_posts_groups
		WHERE user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var users string
		var group models.Group

		err := rows.Scan(
			&group.ID,
			&group.Name,
			&users,
		)

		if err != nil {
			return nil, err
		}

		log.Println(users)
		if users != "" {
			usersIDsStr := strings.Split(users, ":")
			usersIDs := make([]int, 0, len(usersIDsStr))

			for _, idStr := range usersIDsStr {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return nil, errors.New("failed to get users")
				}

				usersIDs = append(usersIDs, id)
			}

			group.Users, err = helpers.GetPostGroupUsers(db, usersIDs)
			if err != nil {
				return nil, err
			}
		}

		groups = append(groups, group)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func AddGroup(db *sql.DB, group models.NewGroup) error {
	users := helpers.JoinUserIDs(group.Users)
	_, err := db.Exec(`insert into user_posts_groups (name, user_id, users)
					   VALUES (?,?,?)`, group.Name, group.UserID, users)
	return err
}

func DeleteGroup(db *sql.DB, groupID, userID int) error {
	_, err := db.Exec(`DELETE FROM user_posts_groups 
					   WHERE user_id = ? AND id = ?`, userID, groupID)
	return err
}

func UpdateGroup(db *sql.DB, group models.NewGroup) error {
	var current string
	err := db.QueryRow(`select users from user_posts_groups where user_id = ? AND name = ?`, 
		group.UserID, group.Name).Scan(&current)
	if err != nil {
		return err
	}

	users := helpers.JoinUserIDs(group.Users)
	if users == current {
		return nil
	}

	_, err = db.Exec(`
		UPDATE user_posts_groups
		SET users = ?
		WHERE user_id = ? AND name = ?
	`, users, group.UserID, group.Name)
	if err != nil {
		return err
	}
	return nil
} 
