package pgxdb

import (
	"context"

	"github.com/AleksandrMac/test_hezzl/pkg/user/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PGXDB struct {
	db *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *PGXDB {
	return &PGXDB{
		db: pool,
	}
}

func (x *PGXDB) Create(ctx context.Context, email string) (id uint64, err error) {
	query := `
INSERT INTO public.user(email)
VALUES ($1) RETURNING id`

	err = x.db.QueryRow(ctx, query, &email).Scan(&id)
	return
}

func (x *PGXDB) Delete(ctx context.Context, id uint64) error {
	query := `
DELETE FROM public.user
WHERE id=$1`

	_, err := x.db.Exec(ctx, query, &id)
	return err
}

func (x *PGXDB) GetList(ctx context.Context, limit, offset uint64) (users []model.UserDB, err error) {
	query := `
SELECT id, email
FROM public.user
LIMIT $1
OFFSET $2`

	rows, err := x.db.Query(ctx, query, &limit, &offset)
	for rows.Next() {
		user := model.UserDB{}
		if err = rows.Scan(&user.ID, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return
}
