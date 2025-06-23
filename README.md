# Template Golang

## Supports

- [x] Clean Architecture
- [x] PostgreSQL
- [x] ~~air~~
  - [x] use `wgo` instead
- [x] Viper config
  - [x] Fix issue about `.env` & `.env.test`
- [x] Docker
- [x] Gorm
- [x] [Swagger-Gin](https://github.com/swaggo/gin-swagger)
- [x] Wire [need to refactor all code to wire]
- [x] Unit Test
- [x] Mock
- [x] [Validator](https://github.com/go-playground/validator)
- [x] [golang-migrate/migrate](https://github.com/golang-migrate/migrate/tree/master?tab=readme-ov-file)
- [ ] [sqlc](https://github.com/sqlc-dev/sqlc)
- [ ] Auth JWT
  - [ ] [goth](https://github.com/markbates/goth)
    - [x] [JWT goth](https://github.com/markbates/goth/issues/310)
    - [x] Permission [admin, staff, user]
      - [x] middleware with role
    - [ ] Refresh token (doing)
      - [x] add exp for JWT
      - [x] add func verify token (1. expired, 2. not exist, 3. valid)
      - [x] implement middleware jwt for call verify func
      - [ ] add refresh token func
    - [ ] Save db
- [ ] Redis
- [ ] Logger system ([zap](https://github.com/uber-go/zap))
- [ ] [Casbin](https://github.com/casbin/casbin)
- [ ] try [WTF](https://github.com/pallat/wtf)
  - [ ] Do corsHandler to support configurable.
- [ ] try [failover](https://github.com/wongnai/lmwn_gomeetup_failover)
- [ ] husky
- [ ] do pkg for common tools

## Design pattern

[clean arch](https://medium.com/@rayato159/how-to-implement-clean-architecture-in-golang-87e9f2c8c5e4)
