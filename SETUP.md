#### Troubleshooting

[How to set up a MySQL database in Linux ](https://www.techtarget.com/searchdatacenter/tip/How-to-set-up-a-MySQL-database-in-Linux)

#### How To Solve

```bash
2024/05/13 20:18:51 Error 1698 (28000): Access denied for user 'root'@'localhost'
```

```sql
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'root';
```

#### Open MySQL from terminal

```bash
mysql -u root -p
```

#### Create database

```sql
CREATE DATABASE ecom;
```

#### Show Databases

```sql
SHOW DATABASES;
```

#### Installing MySQL CLI (UBUNTU)

```bash
$ curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
$ echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
$ apt-get update
$ apt-get install -y migrate
```

#### Adding the user table

```bash
make migration add-user-table
```

#### Check creation of user table (MySQL on terminal)

```sql
USE ecom;

SHOW TABLES;
```

#### Adding the product table

```bash
make migration add-product-table
```

#### Adding the order table

```bash
make migration add-order-table
```

#### Adding the order items table (REMEMEBER TO RUN THAT MIGRATION AT LATEST, OTHERWISE \*ERROR)

```bash
make migration add-order-items-table
```

\*ERROR

```bash
2024/05/23 22:14:53 migration failed in line 0: CREATE TABLE IF NOT EXISTS `order_items` (
  `id`        INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `orderId`   INT UNSIGNED NOT NULL,
  `productId` INT UNSIGNED NOT NULL,
  `price`     DECIMAL(10,2) NOT NULL,
  `quantity`  INT NOT NULL,

  PRIMARY KEY (`id`),
  FOREIGN KEY (`orderId`) REFERENCES orders(`id`),
  FOREIGN KEY (`productId`) REFERENCES products(`id`)
) (details: Error 1824 (HY000): Failed to open the referenced table 'orders')
exit status 1
make: *** [Makefile:17: migrate-up] Error 1
```

#### Managing the MySQL DB

I have installed the VSCode extension named `Database Client JDBC`. It is pretty handy!

#### Dockerizing Go and MySQL

The dockerization is not optimal and it can be improved in so many ways. But it works.

1. Create the `Dockerfile` (multi-stage building)

```dockerfile
# Build the application from source
FROM golang:1.18-alpine3.14 AS build-stage
  WORKDIR /app
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . .
  RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
  RUN go test -v ./...

  # This step is kinda a working around. More info below
  RUN go run migration

# Deploy the application binary into a lean image
FROM scratch AS build-release-stage
  WORKDIR /
  COPY --from=build-stage /api /api
  EXPOSE 8080
  ENTRYPOINT ["/api"]
```

1. Create the `docker-compose` file.

```yml
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    restart: always
    ports:
      - 8080:8080
    environment:
      - PUBLIC_HOST=$PUBLIC_HOST
      - PORT=$PORT
      - JWT_SECRET=$JWT_SECRET
      - DB_USER=$DB_USER
      - DB_PASSWORD=$DB_PASSWORD
      - DB_ADDRESS=$DB_ADDRESS
      - DB_HOST=$DB_HOST
      - DB_PORT=$DB_PORT
      - DB_NAME=$DB_NAME
    depends_on:
      - db

  db:
    image: mysql:8.0
    restart: always
    ports:
      - 3366:3306
    environment:
      - MYSQL_ROOT_PASSWORD=$DB_PASSWORD
      - MYSQL_DATABASE=$DB_NAME
      - MYSQL_PASSWORD=$DB_PASSWORD

    volumes:
      - db_data:/var/lib/mysql

volumes:
  db_data:
```

1. Set the `.env` file

```bash
PUBLIC_HOST=http://localhost
PORT=8080
DB_USER=root
DB_PASSWORD=root
DB_ADDRESS=db:3306
DB_HOST=db
DB_PORT=3306
DB_NAME=ecom
JWT_SECRET=ecom-secret
```

- The `db` value is how the `MySQL` service is named on `docker-compose.yml` file. One
  issue is once you set the environment variables, the container will work fine but running
  the application locally will throw the following error:

```bash
2024/07/27 08:45:42 dial tcp: lookup db: Temporary failure in name resolution
exit status 1
```

- If you replace the `db` value for `http://localhost` the application will run fine **locally**
  but the **app** container will be forever restarting and logging the following:

```bash
Unable to connect to mysql server with go and docker - dial tcp 127.0.0.1:3306: connect: connection refused
```

A good to have functionality would be running the application in either container or locally
without changing settings every time.

###### Docker cheatsheet

- Testing building the containers:

```bash
docker compose build --no-cache
```

- If everything went well (no errors on logs) we start the containers (a.k.a. `app` and `db`)

```bash
docker compose up -d
```

- Checking if the containers are running (check the `STATUS` column in the terminal)

```bash
docker ps -a
```

- \*As I pointed out above, I had an issue where the `app` container `STATUS` was always `restarting`
  while the `db` container was up. The command to check detailed logs is:

```bash
docker logs --tail 50 --follow --timestamps <container_id>
```

Once I have fixed those issues and had both containers working I went to test the API calls. On sending `requests/getProducts.rest` it should return an empty array. But it did not, instead I got the following:

```bash
HTTP/1.1 500 Internal Server Error
Content-Type: application/json
Date: Fri, 26 Jul 2024 13:24:34 GMT
Content-Length: 68
Connection: close

{
  "error": "Error 1146 (42S02): Table 'ecom.products' doesn't exist"
}
```

I did not figured out how to properly run the `db migrations` when building the `db container` so I made it
manually, using the approach below.

1. Accessing and executing commands in the `db` container:

```bash
docker exec -it <db_container_name> mysql -uroot -proot
```

2. When the MySQL terminal opened I selected the `ecom` database

```bash
mysql> USE ecom;
```

Once there I copied / paste the SQL instructions from `migrate/<timestamp>_<migration_name>.up.sql`. You may question my methods but never my results.

3. Once done I re tested the `requests/getProducts.rest` and this time I got an empty array
   `[]` as response body as expected (?).

**?\*** How I should migrate my local data to the db container is still a mystery for me. I will
probably learn that later (hopefully).

List on Docker commands I have used most so far, including the ones already mentioned above:

- Executing commands on MySQL container: `docker exec -it <container_name> mysql -uroot -proot`
- Staring containers: `docker compose up -d`
- Stoping containers: `docker compose down`
- Checking containers / images info: `docker ps -a`
- Logging containers: `docker logs --tail 50 --follow --timestamps <container_id>`
- Clean up not used containers / images: `docker system prune --all`
- Stopping a given container: `docker container stop <container_id>`
- Removing a given container: `docker rm <container_id>`
- Removing a given image: `docker rmi <image_id>`
- Testing a fresh docker compose build: `docker compose build --no-cache`
