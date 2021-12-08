### Create directory

```bash
mkdir -p $HOME/pg-for-go-devs/data
mkdir -p $HOME/pg-for-go-devs/init
```
### Start Postgres in container
```bash
docker run \
    --rm \
    -d \
    -p 5432:5432 \
    --name postgres \
    -e POSTGRES_PASSWORD=P@ssw0rd \
    -e PGDATA=/var/lib/postgresql/data \
    -v $HOME/pg-for-go-devs/data:/var/lib/postgresql/data \
    -v $HOME/pg-for-go-devs/init:/docker-entrypoint-initdb.d \
    postgres:14.0
```

### Start Postgres if container exist
```bash
docker start postgres
```

### stop Postgres if container exist
```bash
docker start postgres
```

### Connect to DB psql
```bash
psql -h 127.0.0.1 -p 5432 -U snippetbox -d snippetbox
```

### Run command in container
```bash
docker exec \
    -it \
    --user postgres \
    postgres \
    bash
```

### Run sql scripts in container
```bash
cd docker-entrypoint-initdb.d

psql -h 127.0.0.1 -p 5432 -U snippetbox -d snippetbox

\i schema.sql
\i data.sql
```

### Start Redis
docker-compose up -d