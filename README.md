# Task scheduler

1. В дипломном проекте написан веб-сервер, который реализует функциональность планировщика задач — это аналог TODO-листа. Планировщик хранит задачи; каждая из них содержит дату дедлайна и заголовок с комментарием.
2. Выполнены все задания повышенной сложности, кроме последнего Dockerfile
3. Адрес следует указывать в браузере: http://localhost:7540
4. Запуск локально: `TODO_PORT=7540 TODO_DBFILE=scheduler.db TODO_PASSWORD=12345 go run cmd/main.go`
5. Параметры в tests/settings.go следует использовать:
* Port: 7540
* DBFile = "../scheduler.db"
* FullNextDate = true
* Search = true
* Token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXNzd29yZCI6IjU5OTQ0NzFhYmIwMTExMmFmY2MxODE1OWY2Y2M3NGI0ZjUxMWI5OTgwNmRhNTliM2NhZjVhOWMxNzNjYWNmYzUifQ.YtW2-no0aYnyU5-7zwYp8wTEs8uJBt63W5WgHkZqjnI'
