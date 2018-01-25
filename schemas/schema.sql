CREATE DATABASE IF NOT EXISTS `tracker`;

DROP TABLE IF EXISTS `tracker`.`shows`;
CREATE TABLE `tracker`.`shows` (
	id INTEGER NOT NULL AUTO_INCREMENT,
	title VARCHAR(255) NOT NULL,
	wikipedia VARCHAR(255),
	trailer VARCHAR(255),
	finished BOOLEAN DEFAULT false,
	PRIMARY KEY(id)
);

DROP TABLE IF EXISTS `tracker`.`episodes`;
CREATE TABLE `tracker`.`episodes` (
	id INTEGER NOT NULL AUTO_INCREMENT,
	show_id INTEGER NOT NULL,
	season INTEGER NOT NULL,
	episode INTEGER NOT NULL,
	title VARCHAR(255),
	release_date DATE,
	PRIMARY KEY(id),
	UNIQUE KEY(show_id, season, episode)
);
