package profiles

import (
	"database/sql"
	"errors"
	database "social/database/users"
	"social/internal/models"
)

func GetUserAbout(db *sql.DB, userID int) (models.UserAbout, error) {
	var userProfile models.UserAbout
	err := db.QueryRow(`
		SELECT work,hobbies,education,intrests,travel,website,linkedin,instgram,twitter
		FROM user_about
		WHERE user_id = ?;
	`, userID).Scan(
		&userProfile.Work,
		&userProfile.Hobbies,
		&userProfile.Education,
		&userProfile.Intrests,
		&userProfile.Travel,
		&userProfile.Website,
		&userProfile.Linkedin,
		&userProfile.Instgram,
		&userProfile.Twitter,
	)

	if err != nil {
		return models.UserAbout{}, err
	}

	err = db.QueryRow(`select about from profile where user_id = ?`, userID).Scan(&userProfile.Bio)
	return userProfile, err
}

func UpdateUserAbout(db *sql.DB, userID int, UserAbout models.UserAbout) error {
	result, err := db.Exec(`
		UPDATE user_about
		SET work = ?,hobbies = ?,education = ?, intrests = ?, travel = ?, website = ?, linkedin = ?, instgram = ?, twitter = ?
		WHERE user_id = ? 
	`, UserAbout.Work, UserAbout.Hobbies, UserAbout.Education, UserAbout.Intrests, UserAbout.Travel,
		UserAbout.Website, UserAbout.Linkedin, UserAbout.Instgram, UserAbout.Twitter, userID)
	if err != nil {
        return err
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rows == 0 {
        return sql.ErrNoRows
    }

    return nil
}

func IsPrivate(db *sql.DB, userID int) (bool, error) {
	var IsPrivate int
	err := db.QueryRow(`
		SELECT is_private from profile where user_id = ?
	`,userID).Scan(&IsPrivate)

	if err != nil {
		return false, err
	}

	if IsPrivate == 1 {
		return true, nil
	} else if IsPrivate == 0 {
		return false, nil
	} else {
		return false, errors.New("Invalid data")
	}
}

// this function check the followers satus
//@return 0 for [REQUESTED] and 1 for [following]
func CheckFollower(db *sql.DB, followerID, targetID int) (int, error) {
	var status int
	err := db.QueryRow(`
		select status from user_followers where follower_id = ? AND target_id = ?
	`, followerID, targetID).Scan(&status)
	return status, err
}

func GetUserData(db *sql.DB, userID int) (models.UserData, error) {
	userData, err := database.GetUserData(db, userID)
	if err != nil {
		return models.UserData{}, err
	}	
	
	userAbout, err := GetUserAbout(db, userID)
	if err != nil {
		return models.UserData{}, err
	}
	
	userData.About = userAbout
	return userData, nil
}

func GetPrivateProfileData(db *sql.DB, userID int) (models.UserData, error) {
	var userData models.UserData
	err := db.QueryRow(`select first_name, last_name from user where id = ?`,userID ).Scan(
		&userData.UserInfo.FirstName,
		&userData.UserInfo.LastName,
	)
	if err != nil {
		return models.UserData{}, err
	}

	err = db.QueryRow(`select num_of_following, num_of_followers, num_of_posts, about, avatar_path 
						from profile where user_id = ?`, userID).Scan(
							&userData.NumOfFollowing,
							&userData.NumOfFollowers,
							&userData.NumOfPosts,
							&userData.About.Bio,
							&userData.UserInfo.Avatar,
						)
	return userData, err
}
