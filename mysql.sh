docker run --name=mysql --network tracker-network -p 3306:3306 -v mysql-volume:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=rhino -d mysql/mysql-server:latest
