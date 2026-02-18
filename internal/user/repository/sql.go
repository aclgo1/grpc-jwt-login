package repository

const (
	queryAddUser = `INSERT INTO users (user_id, name, last_name, password, email,
	role, verified, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
	RETURNING user_id, name, last_name, password, email, role, verified, created_at, updated_at`

	queryByID = `select user_id, name, last_name, password, email, role, verified, created_at,
	updated_at from users where user_id=$1`

	queryFindByEmail = `select user_id, name, last_name, password, email, role, verified,
	created_at,updated_at from users where email=$1`

	queryUpdate = `UPDATE "users" SET 
    "name" = COALESCE(NULLIF($1, ''), "name"), 
    "last_name" = COALESCE(NULLIF($2, ''), "last_name"),
    "password" = COALESCE(NULLIF($3, ''), "password"), 
    "email" = COALESCE(NULLIF($4, ''), "email"),
	"role" = COALESCE(NULLIF($5, ''), "role"),
    "verified" = COALESCE(NULLIF($6, ''), "verified"), 
    "updated_at" = COALESCE(NULLIF($7, '')::timestamptz, "updated_at") 
	WHERE user_id = $8
	RETURNING user_id, name, last_name, password, email,
	role, verified, created_at, updated_at;`

	queryDelete = `delete from users where user_id=$1`
)
