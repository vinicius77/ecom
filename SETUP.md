#### How to set up a MySQL database in Linux

[How to set up a MySQL database in Linux ](https://www.techtarget.com/searchdatacenter/tip/How-to-set-up-a-MySQL-database-in-Linux)

#### How To Solve

```bash
2024/05/13 20:18:51 Error 1698 (28000): Access denied for user 'root'@'localhost'
```

```bash
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'root';
```
