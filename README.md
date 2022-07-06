# sample-rest-api

This is a demonstration of a go-based backend service which exposes RESTFUL APIs allowing to perform CRUD operations on data representing repositories scanning informations collected by a dedicated client tool. It basically shows how to design and implement a clean RestFul API backend with Go while using its built-in functionnality like interfaces for more flexibility & scalability. The storage layer is designed to support SQL & NoSQL datbases. Implementation has been done for PostgreSQL and MongoDB. Based on the code design, you can easily define or swap the wanted database (PostgreSQL or MongoDB) from the configuration file. A fake mocked database interface called mockdb was also added for live unit testing of the CRUD endpoints. The endpoint listens only on HTTPS and a bash script is available to generate self-signed tls/ssl certificates with openssl. The whole build process is automate when using docker-compose.


## Get-Started

This project could be built & run with docker-compose quickly. You can even use Make Tool for common tasks (check Makefile for available commands).
Otherwise, you can spin up the required database containers with the same docker-compose file and manually run separetely the backend api service.

* First Method

	* Clone the repository

		```shell
		$ git clone https://github.com/jeamon/sample-rest-api.git
		$ cd sample-rest-api
		```

	* Build and Run with docker :

		```
		$ make run-docker
		```

	* Run locally after setting up the database :

		```
		$ make run-local
		```

	* Connect to the server :

		```
		$ https://<server-address>:8080/api/v1/<endpoint>
		```


* Second Method

	* Clone the repository

		```shell
		$ git clone https://github.com/jeamon/sample-rest-api.git
		$ cd sample-rest-api
		```

	* Build and Run with docker :

		```
		$ docker-compose up -d
		```

	* Connect to the server :

		```
		$ https://<server-address>:8080/api/v1/<endpoint>
		```

* Third Method

	* Clone the repository

		```shell
		$ git clone https://github.com/jeamon/sample-rest-api.git
		$ cd sample-rest-api
		```

	* Start all databases or the one set into the configs file :
		
		```
		$ docker-compose up -d postgres mongo pgadmin
		```

	* Create new certificates (update the CN with your server address) and place them under assets/certs folder.

		```shell
		$ openssl req -new -subj "/CN=localhost" -newkey rsa:2048 -nodes -keyout ./server.key -out ./server.csr
		$ openssl x509 -req -days 3650 -in ./server.csr -signkey ./server.key -out ./server.crt
		```

		or (run the script to automatically generate them and move them under assets/certs folder)

		```
		$ chmod +x ./scripts/generate.certs.sh
		$ ./scripts/generate.certs.sh
		```

		or (for testing purpose, rename existing tests certificates)

		```
		$ cp ./assets/certs/test.server.crt ./assets/certs/server.crt
		$ cp ./assets/certs/test.server.key ./assets/certs/server.key
		```

	* Set the configs file with the database address and Run the server.
				
		```
		$ go run main.go
		```

		or (below to fill the build flags : latest git hash and tag ID)

		```
		$ go run -ldflags "-X 'main.GitCommit=$(git rev-list -1 HEAD)' -X 'main.GitTag=$(git describe --tags --abbrev=0)'" main.go
		```
	
* Connect to the pgadmin for managing the postgres database with UI.
	
	```
	http://<server-address>:8081/browser
	```


## Credentials and Settings	
	
The file **<server.config.yml>** contains configurations settings of the server and the databases (postgresql and mongodb). It is loaded at server startup.
The file **<server.config.docker.yml>** contains configurations settings of the server and the databases (postgresql and mongodb) but it is customized to be used when building and running the project with docker-compose.



## Data structure of a scan infos object

This below structure is the core model of a scan infos and its representation into different format (json, bson, sql database). 

```go
type ScanInfos struct {
	ID string `db:"id" json:"id" bson:"id" binding:"required"`

	CompanyID string `db:"company_id" json:"company_id" bson:"company_id" binding:"required"`
	Username  string `db:"username" json:"username" bson:"username" binding:"required"`

	ClientID string `db:"client_id" json:"client_id" bson:"client_id" binding:"required"`

	RepositoryURL string `db:"repository_url" json:"repository_url" bson:"repository_url" binding:"required"`
	CommitID      string `db:"commit_id" json:"commit_id" bson:"commit_id" binding:"required"`
	TagID         string `db:"tag_id" json:"tag_id" bson:"tag_id" binding:"required"`

	Results []string `db:"results" json:"results" bson:"results" binding:"required"`

	StartedAt   int64 `db:"started_at" json:"started_at" bson:"started_at" binding:"required"`
	CompletedAt int64 `db:"completed_at" json:"completed_at" bson:"completed_at" binding:"required"`
	SentAt      int64     `db:"sent_at" json:"sent_at" bson:"sent_at" binding:"required"`
	CreatedAt   time.Time `db:"created_at" json:"-" bson:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"-" bson:"updated_at"`

	Error    string `db:"error" json:"error" bson:"error"`
	Metadata map[string]interface{} `db:"metadata" json:"metadata" bson:"metadata" binding:"required"`
}
```

The fields **<started_at>** **<completed_at>** and **<sent_at>** hold the corresponding value in unix epoch time (in seconds).
The field **<Metadata>** is provided to hold more data if needed. This proactively provides a flexibility into the data structure.


## Endpoints for CRUD Operations

*By default the server runs on localhost (127.0.0.1) and listen on port 8080.*

* Submit a performed scan information

```
[POST] http://<server-address>:<server-port>/api/v1/scaninfos
```


* Fetch specific scan information based on its id

```
[GET] http://<server-address>:<server-port>/api/v1/scaninfos/<scan-id>
```


* Display all available scan information

```
[GET] http://<server-address>:<server-port>/api/v1/scaninfos
```

* Update existing scan information

```
[PUT] http://<server-address>:<server-port>/api/v1/scaninfos
```

* Delete existing scan information from the database

```
[DELETE] http://<server-address>:<server-port>/api/v1/scaninfos/<scan-id>
```


## Others Endpoints

* Quick check of the backend service availability

```
[GET] http://<server-address>:<server-port>/ping
```

* Health check of the global backend system

```
[GET] http://<server-address>:<server-port>/status
```


## POST [Request Body]

Below is the internal structure used to hold request body when submitting a new scan information.

```go
type StoreScanInfosRequest struct {
	CompanyID string `json:"company_id" binding:"required"`
	Username  string `json:"username" binding:"required"`

	ClientID string `json:"client_id" binding:"required"`

	RepositoryUrl string `json:"repository_url" binding:"required"`
	CommitID      string `json:"commit_id" binding:"required"`
	TagID         string `json:"tag_id" binding:"required"`

	Results []string `json:"results" binding:"required"`

	StartedAt   int64 `json:"started_at" binding:"required"`
	CompletedAt int64 `json:"completed_at" binding:"required"`
	SentAt      int64     `json:"sent_at" binding:"required"`

	Error    string `json:"error"`
	Metadata string `json:"metadata" binding:"required"`
}
```

The request body should be sent as JSON data. The fields **<started_at>** **<completed_at>** and **<sent_at>** are the corresponding unix epoch time (in seconds).


* Sample JSON data for POST request body

```json
{
    "company_id": "company-id",
    "username": "jeamon",
    "client_id": "vx.y.z",
    "repository_url": "https://github.com/jeamon/sample-rest-api",
    "commit_id": "d7b8ff1412ebfcde26f9ddfdf9608d1525647958",
    "tag_id": "v1.0.0",
    "results": ["found something x", "found something y", "found something z"],
    "started_at":1655903720,
    "completed_at": 1655903723,
    "sent_at": 1655903725,
    "error": "got an x exception during repository scanning",
    "metadata": {
        "os": "linux",
        "languages": ["go", "bash", "html"],
        "arch": "amd64"
    }
}
```