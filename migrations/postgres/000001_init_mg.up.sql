CREATE TABLE Persons (
    person_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255),
    age INT,
    gender VARCHAR(255),
    isDeleted BOOL DEFAULT FALSE
);

CREATE TABLE Nationalities (
    nationality_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    nationality_name VARCHAR(2) UNIQUE
);

CREATE TABLE PersonNationalities (
    person_id INT,
    nationality_id INT,
    probability FLOAT,
    PRIMARY KEY (person_id, nationality_id),
    FOREIGN KEY (person_id) REFERENCES Persons(person_id),
    FOREIGN KEY (nationality_id) REFERENCES Nationalities(nationality_id)
);