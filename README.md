# snippetbox
Site for notes. 

Frontend + backend on Go. 

Clone from my repository https://github.com/White-AK111/Go/tree/main/src/snippetbox for further development.


Already DONE:
- Standard library net/http web server;
- Data model - pkg models;
- MySQL strorage - pkg mysql;
- HTML pages geterated by templates;


TODO:
- Change storage to PostgreSQL and write new data model;
- Add authentication;
- Unit tests;
- Hosting in WEB;
- CI\CD by GtHub Actions;
- ...

HOW TO MIGRATION (init database):
1. Migration must run on empty database, for create database run script "init.sql" from: "/pkg/models/postgres/scripts";
2. Download migrate tool: "go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest";
3. Go to catalog: "/pkg/models/postgres/migrations";
4. Run: "migrate -database "postgresql://snippetbox:P@ssw0rd@localhost:5432/snippetbox?sslmode=disable" -path migrations up".