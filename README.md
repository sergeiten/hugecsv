# HUGECSV : Read huge CSV file and insert to MS SQL database

Creates two docker containers (reader and consumer). `Reader` reads csv file(sample.csv) and send to `consumer` by using [protobuf](https://developers.google.com/protocol-buffers/).

## Requirements 

### Docker installation (for MAC users)

Install docker via [homebrew](https://brew.sh) 
```bash
brew cask install docker
```

### 
Create `.env` file (make copy of .env.example):
```.env
DB_HOST=
DB_NAME=
DB_USER=
DB_PASSPORT=
DB_PORT=
```

## Run

```bash
docker-compose up
```



