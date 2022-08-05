create database userma;
-- go_db: old DB

CREATE TABLE users (
    `id` bigInt(20) NOT NULL AUTO_INCREMENT,
    `username` varchar(100) NOT NULL DEFAULT '',
    `nickname` varchar(100) DEFAULT '',
    `picname` varchar(100) DEFAULT '',
    `password` varchar(100) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;