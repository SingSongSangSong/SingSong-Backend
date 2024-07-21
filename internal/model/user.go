package model

type User struct {
	ID       int
	Username string
}

func (model *Model) RegisterUser(username string) error {
	_, err := model.db.Exec("INSERT INTO `user` (username) VALUES (?)", username)
	if err != nil {
		return err
	}
	return nil
}

func (model *Model) ListUser() ([]User, error) {
	rows, err := model.db.Query("SELECT id, username FROM `user`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (model *Model) GetUser(username string) (*User, error) {
	var user User
	row := model.db.QueryRow("SELECT * FROM `user` WHERE username = ?", username)
	if err := row.Scan(&user.ID, &user.Username); err != nil {
		return nil, err
	}
	return &user, nil
}
