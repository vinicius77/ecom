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
