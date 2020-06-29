package user

import (
	"database/sql"
	"time"

	"github.com/eminetto/clean-architecture-go-v2/domain"

	"github.com/eminetto/clean-architecture-go-v2/domain/entity"
)

//MySQLRepo mysql repo
type MySQLRepo struct {
	db *sql.DB
}

//NewMySQLRepoRepository create new repository
func NewMySQLRepoRepository(db *sql.DB) *MySQLRepo {
	return &MySQLRepo{
		db: db,
	}
}

//Create an user
func (r *MySQLRepo) Create(e *User) (entity.ID, error) {
	stmt, err := r.db.Prepare(`
		insert into user (id, email, password, first_name, last_name, created_at) 
		values(?,?,?,?,?,?)`)
	if err != nil {
		return e.ID, err
	}
	_, err = stmt.Exec(
		e.ID,
		e.Email,
		e.Password,
		e.FirstName,
		e.LastName,
		time.Now().Format("2006-01-02"),
	)
	if err != nil {
		return e.ID, err
	}
	err = stmt.Close()
	if err != nil {
		return e.ID, err
	}
	return e.ID, nil
}

//Get an user
func (r *MySQLRepo) Get(id entity.ID) (*User, error) {
	return getUser(id, r.db)
}

func getUser(id entity.ID, db *sql.DB) (*User, error) {
	stmt, err := db.Prepare(`select id, email, first_name, last_name, created_at from user where id = ?`)
	if err != nil {
		return nil, err
	}
	var u User
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt)
	}
	stmt, err = db.Prepare(`select book_id from book_user where user_id = ?`)
	if err != nil {
		return nil, err
	}
	rows, err = stmt.Query(id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i entity.ID
		err = rows.Scan(&i)
		u.Books = append(u.Books, i)
	}
	return &u, nil
}

//Update an user
func (r *MySQLRepo) Update(e *User) error {
	e.UpdatedAt = time.Now()
	_, err := r.db.Exec("update user set email = ?, password = ?, first_name = ?, last_name = ?, updated_at = ? where id = ?", e.Email, e.Password, e.FirstName, e.LastName, e.UpdatedAt.Format("2006-01-02"), e.ID)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("delete from book_user where user_id = ?", e.ID)
	if err != nil {
		return err
	}
	for _, b := range e.Books {
		_, err := r.db.Exec("insert into book_user values(?,?,?)", e.ID, b, time.Now().Format("2006-01-02"))
		if err != nil {
			return err
		}
	}
	return nil
}

//Search users
func (r *MySQLRepo) Search(query string) ([]*User, error) {
	stmt, err := r.db.Prepare(`select id from user where name like ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var ids []entity.ID
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i entity.ID
		err = rows.Scan(&i)
		if err != nil {
			return nil, err
		}
		ids = append(ids, i)
	}
	if len(ids) == 0 {
		return nil, domain.ErrNotFound
	}
	var users []*User
	for _, id := range ids {
		u, err := getUser(id, r.db)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

//List users
func (r *MySQLRepo) List() ([]*User, error) {
	stmt, err := r.db.Prepare(`select id from user`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var ids []entity.ID
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i entity.ID
		err = rows.Scan(&i)
		if err != nil {
			return nil, err
		}
		ids = append(ids, i)
	}
	if len(ids) == 0 {
		return nil, domain.ErrNotFound
	}
	var users []*User
	for _, id := range ids {
		u, err := getUser(id, r.db)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

//Delete an user
func (r *MySQLRepo) Delete(id entity.ID) error {
	_, err := r.db.Exec("delete from user where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
