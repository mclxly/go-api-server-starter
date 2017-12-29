package models

type User struct {
    Id string `db:"id"`
    Nickname  string `db:"nickname"`
    Email     string
}