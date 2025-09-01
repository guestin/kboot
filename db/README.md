# Single Datasource Config

```toml
[db]
# only support sqlite or postgres
type = "postgres"
dsn = """host=localhost
port=5432
user=postgres
password=xxxxxx
dbname=xxx
sslmode=disable
TimeZone=Asia/Shanghai"""
timezone = "Asia/Shanghai"
debug = false
```

## Usage

```
db:=db.ORM() //same as db.ORM("default")
// then use it 
```

# Multi Datasource Config

```toml
[db]
# only support sqlite or postgres
type = "postgres"
dsn = """host=localhost
port=5432
user=postgres
password=xxxxxx
dbname=xxx
sslmode=disable
TimeZone=Asia/Shanghai"""
timezone = "Asia/Shanghai"
debug = false

# define ds1 ,ds1 is a custom name ,can be everything you like
[db.ds1]
# only support sqlite or postgres
type = "postgres"
dsn = """host=localhost
port=5432
user=postgres
password=xxxxxx
dbname=xxx
sslmode=disable
TimeZone=Asia/Shanghai"""
timezone = "Asia/Shanghai"
debug = false

# define ds2 ,ds2 is a custom name ,can be everything you like
[db.ds2]
# only support sqlite or postgres
type = "postgres"
dsn = """host=localhost
port=5432
user=postgres
password=xxxxxx
dbname=xxx
sslmode=disable
TimeZone=Asia/Shanghai"""
timezone = "Asia/Shanghai"
debug = false
```

```
// get the default db
db:=db.ORM() //same as db.ORM("default")
// then use it 

//get ds1
db1:=db.ORM("ds1")
// then use it 

//get ds2
db2:=db.ORM("ds2")
// then use it 
```