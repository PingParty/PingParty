package db

import (
	"database/sql"

	"github.com/PingParty/PingParty/models"
)

func (d *DB) GetExistingThirdPartyUser(loginType, id string) (*models.User, error) {
	const q = `SELECT * FROM users WHERE LoginType = ? AND LoginID = ?`
	u := &models.User{}
	err := d.d.QueryRow(q, loginType, id).Scan(u)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}
