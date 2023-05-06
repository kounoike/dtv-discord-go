CREATE DATABASE IF NOT EXISTS `airflow`;
CREATE USER IF NOT EXISTS 'airflow'@'%' IDENTIFIED BY 'airflow';
GRANT ALL ON `airflow`.* TO 'airflow'@'%';
