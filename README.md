# Setup information 

# Local data load 

```
cd tests/fixtures
docker exec -it pgsql bash -e /data/insertdata.sh
```

```
sqlc generate
```

```
	"database/sql"
	"fmt"
	"net/url"


	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"

c, err := pgx.ParseConfig(pgURL.String())
	if err != nil {
		return nil, fmt.Errorf("parsing postgres URI: %w", err)
	}

	c.Logger = logrusadapter.NewLogger(logger)
	db := stdlib.OpenDB(*c)

	err = validateSchema(db)
	if err != nil {
		return nil, fmt.Errorf("validating schema: %w", err)
	}

	return &Directory{
		logger:  logger,
		db:      db,
		sb:      squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
		querier: New(db),
	}, nil

```

## Local env setup
1. Install docker
2. Install golang
3. Install sqlc (sudo snap install sqlc)
4. Install make utility
5. Install build essentials in ubuntu
6. Download the code
7. Generate the tool using the following command from the root of the src code directory
```
cd <src>
make tools
./build/exec/hrmstool generate-config configs/local
```
8. Fill up the configuration values
9. Run the local server 
```
./build/exec/hrms -c ./configs/local/config.json
```
10. Run the regression tests under test folder as appropriate

## Test case writing strategy
We will be following regretion test strategy instead of local unit testing strategy. Benefits are 
1. Can change the target system url and test any system with any type of data
2. Will be helpful when make any change in the code. Before deploying we may test the correctness of the change done and then deploy. 
3. We should also be anle check if it breaks anything in the existing system.

## Run Debug Mode
1. use the following code under .vscode/launch.json

{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        

        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/hrms",
            "dlvFlags": ["--check-go-version=false"] ,
            "args": ["-c" ,"./config/connection/dev-config.json"],
            "cwd": "${workspaceFolder}"
        }
    ]
}