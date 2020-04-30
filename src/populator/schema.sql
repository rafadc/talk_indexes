CREATE TABLE people_small (
	id INT UNSIGNED auto_increment NOT NULL KEY,
	name varchar(100) NOT NULL,
	surname varchar(200) NOT NULL,
	date_of_birth DATE NOT NULL,
	company varchar(100) NOT NULL
)
ENGINE=InnoDB;

CREATE TABLE people_without_indexes (
	id INT UNSIGNED auto_increment NOT NULL KEY,
	name varchar(100) NOT NULL,
	surname varchar(200) NOT NULL,
	date_of_birth DATE NOT NULL,
	company varchar(100) NOT NULL
)
ENGINE=InnoDB;

