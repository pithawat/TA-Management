CREATE TABLE professors(
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(100) NOT NULL,
    lastname VARCHAR(150) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE accountants(
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(100) NOT NULL,
    lastname VARCHAR(150) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE courses(
    id SERIAL PRIMARY KEY,
    course_ID VARCHAR(100) NOT NULL,
    course_name VARCHAR(150) NOT NULL,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50) NOT NULL
);

INSERT INTO professors (firstname, lastname) VALUES
('Amnach', 'Khawne'),
('Charoen', 'Vongchumyen'),
('Kietikul', 'Jearanaitanakij'),
('Orachat', 'Chitsobhuk'),
('Rathachai', 'Chawuthai'),
('Sakchai', 'Thipchaksurat'),
('Surin', 'Kittitornkun'),
('Aranya', 'Walairacht'),
('Chompoonuch', 'Tengcharoen'),
('Chutimet', 'Srinilta'),
('Pakorn', 'Watanachaturaporn'),
('Phongsak', 'Keeratiwintakorn'),
('Thanunchai', 'Threepak'),
('Akkradach', 'Watcharapupong'),
('Bundit', 'Pasaya'),
('Kiatnarong', 'Tongprasert'),
('Sorayut', 'Glomglome'),
('Thana', 'Hongsuwan'),
('Jirayu', 'Petchhan'),
('Parinya', 'Ekparinya'),
('Kanut', 'Tangtisanon'),
('Watjanapong', 'Kasemsiri'),
('Jirasak', 'Sittigorn'),
('Pithawat', 'Kitmongkolchai');