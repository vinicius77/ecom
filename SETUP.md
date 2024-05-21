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
