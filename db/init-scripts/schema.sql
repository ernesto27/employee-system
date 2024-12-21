CREATE TABLE employees (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    date_birth DATE,
    image TEXT,
    start_work_date DATE,
    end_work_date DATE,
    active BOOLEAN,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME
);

CREATE TABLE roles (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255)
);

CREATE TABLE employees_roles (
    id INT PRIMARY KEY AUTO_INCREMENT,
    employee_id INT,
    role_id INT,
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE projects (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    description TEXT,
    start_date DATETIME,
    end_date DATETIME,
    active BOOLEAN,
    images TEXT,
    links TEXT,
    contacts TEXT,
    notes TEXT
);

CREATE TABLE employees_projects (
    id INT PRIMARY KEY AUTO_INCREMENT,
    employee_id INT,
    project_id INT,
    start_date DATETIME,
    end_date DATETIME,
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE TABLE technologies (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    description TEXT
);

CREATE TABLE projects_technologies (
    id INT PRIMARY KEY AUTO_INCREMENT,
    project_id INT,
    technology_id INT,
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (technology_id) REFERENCES technologies(id)
);

INSERT INTO employees (name, email, date_birth, image, start_work_date, end_work_date, active, created_at, updated_at) VALUES
('John Doe', 'john.doe@example.com', '1992-01-01', 'image1.jpg', '2022-01-01', '2023-01-01', TRUE, NOW(), NOW()),
('Jane Smith', 'jane.smith@example.com', '1997-01-01', 'image2.jpg', '2021-06-01', '2022-06-01', TRUE, NOW(), NOW());

INSERT INTO employees (name, email, date_birth, image, start_work_date, end_work_date, active, created_at, updated_at) VALUES
('Alice Johnson', 'alice.johnson@example.com', '1994-01-01', 'image3.jpg', '2020-03-01', '2021-03-01', TRUE, NOW(), NOW()),
('Bob Brown', 'bob.brown@example.com', '1987-01-01', 'image4.jpg', '2019-07-01', '2020-07-01', TRUE, NOW(), NOW()),
('Charlie Davis', 'charlie.davis@example.com', '1982-01-01', 'image5.jpg', '2018-05-01 09:00:00', '2019-05-01 17:00:00', TRUE, NOW(), NOW()),
('Diana Evans', 'diana.evans@example.com', '1990-01-01', 'image6.jpg', '2017-09-01 09:00:00', '2018-09-01 17:00:00', TRUE, NOW(), NOW()),
('Ethan Foster', 'ethan.foster@example.com', '1993-01-01', 'image7.jpg', '2016-11-01 09:00:00', '2017-11-01 17:00:00', TRUE, NOW(), NOW()),
('Fiona Green', 'fiona.green@example.com', '1995-01-01', 'image8.jpg', '2015-01-01 09:00:00', '2016-01-01 17:00:00', TRUE, NOW(), NOW()),
('George Harris', 'george.harris@example.com', '1989-01-01', 'image9.jpg', '2014-04-01 09:00:00', '2015-04-01 17:00:00', TRUE, NOW(), NOW()),
('Hannah Ingram', 'hannah.ingram@example.com', '1991-01-01', 'image10.jpg', '2013-06-01 09:00:00', '2014-06-01 17:00:00', TRUE, NOW(), NOW()),
('Ian Jackson', 'ian.jackson@example.com', '1986-01-01', 'image11.jpg', '2012-08-01 09:00:00', '2013-08-01 17:00:00', TRUE, NOW(), NOW()),
('Julia King', 'julia.king@example.com', '1988-01-01', 'image12.jpg', '2011-10-01 09:00:00', '2012-10-01 17:00:00', TRUE, NOW(), NOW()),
('Kevin Lewis', 'kevin.lewis@example.com', '1984-01-01', 'image13.jpg', '2010-12-01 09:00:00', '2011-12-01 17:00:00', TRUE, NOW(), NOW()),
('Laura Martinez', 'laura.martinez@example.com', '1996-01-01', 'image14.jpg', '2009-02-01 09:00:00', '2010-02-01 17:00:00', TRUE, NOW(), NOW()),
('Michael Nelson', 'michael.nelson@example.com', '1985-01-01', 'image15.jpg', '2008-04-01 09:00:00', '2009-04-01 17:00:00', TRUE, NOW(), NOW()),
('Nina Owens', 'nina.owens@example.com', '1983-01-01', 'image16.jpg', '2007-06-01 09:00:00', '2008-06-01 17:00:00', TRUE, NOW(), NOW()),
('Oscar Perez', 'oscar.perez@example.com', '1981-01-01', 'image17.jpg', '2006-08-01 09:00:00', '2007-08-01 17:00:00', TRUE, NOW(), NOW()),
('Paula Quinn', 'paula.quinn@example.com', '1980-01-01', 'image18.jpg', '2005-10-01 09:00:00', '2006-10-01 17:00:00', TRUE, NOW(), NOW()),
('Quincy Roberts', 'quincy.roberts@example.com', '1979-01-01', 'image19.jpg', '2004-12-01 09:00:00', '2005-12-01 17:00:00', TRUE, NOW(), NOW()),
('Rachel Scott', 'rachel.scott@example.com', '1978-01-01', 'image20.jpg', '2003-02-01 09:00:00', '2004-02-01 17:00:00', TRUE, NOW(), NOW()),
('Samuel Thompson', 'samuel.thompson@example.com', '1977-01-01', 'image21.jpg', '2002-04-01 09:00:00', '2003-04-01 17:00:00', TRUE, NOW(), NOW()),
('Tina Underwood', 'tina.underwood@example.com', '1976-01-01', 'image22.jpg', '2001-06-01 09:00:00', '2002-06-01 17:00:00', TRUE, NOW(), NOW()),
('Ulysses Vega', 'ulysses.vega@example.com', '1975-01-01', 'image23.jpg', '2000-08-01 09:00:00', '2001-08-01 17:00:00', TRUE, NOW(), NOW()),
('Victoria White', 'victoria.white@example.com', '1974-01-01', 'image24.jpg', '1999-10-01 09:00:00', '2000-10-01 17:00:00', TRUE, NOW(), NOW()),
('William Xavier', 'william.xavier@example.com', '1973-01-01', 'image25.jpg', '1998-12-01 09:00:00', '1999-12-01 17:00:00', TRUE, NOW(), NOW()),
('Xena Young', 'xena.young@example.com', '1972-01-01', 'image26.jpg', '1997-02-01 09:00:00', '1998-02-01 17:00:00', TRUE, NOW(), NOW()),
('Yvonne Zane', 'yvonne.zane@example.com', '1971-01-01', 'image27.jpg', '1996-04-01 09:00:00', '1997-04-01 17:00:00', TRUE, NOW(), NOW()),
('Zachary Adams', 'zachary.adams@example.com', '1970-01-01', 'image28.jpg', '1995-06-01 09:00:00', '1996-06-01 17:00:00', TRUE, NOW(), NOW()),
('Amber Brown', 'amber.brown@example.com', '1969-01-01', 'image29.jpg', '1994-08-01 09:00:00', '1995-08-01 17:00:00', TRUE, NOW(), NOW()),
('Brian Clark', 'brian.clark@example.com', '1968-01-01', 'image30.jpg', '1993-10-01 09:00:00', '1994-10-01 17:00:00', TRUE, NOW(), NOW()),
('Catherine Davis', 'catherine.davis@example.com', '1967-01-01', 'image31.jpg', '1992-12-01 09:00:00', '1993-12-01 17:00:00', TRUE, NOW(), NOW()),
('David Evans', 'david.evans@example.com', '1966-01-01', 'https://imgs.search.brave.com/inJXWy1ZV0RLz1UvrE9SY1XdowevuzIA6PVSUdPkoQ0/rs:fit:860:0:0:0/g:ce/aHR0cHM6Ly9tLm1l/ZGlhLWFtYXpvbi5j/b20vaW1hZ2VzL00v/TVY1QllqWTRaVFJq/TnpBdE1UUm1NUzAw/WW1VM0xXRTVOV0V0/WmpVMk5UYzJPRFUy/TWpBd1hrRXlYa0Zx/Y0djQC5qcGc', '1991-02-01 09:00:00', '1992-02-01 17:00:00', TRUE, NOW(), NOW());

INSERT INTO roles (name) VALUES
('Developer'),
('Manager'),
('Tester'),
('Designer'),
('DevOps');

INSERT INTO employees_roles (employee_id, role_id) VALUES
(1, 1),
(2, 2);

INSERT INTO projects (name, description, start_date, end_date, active, images, links, contacts, notes) VALUES
('Project Alpha', 'Description of Project Alpha', '2022-01-01 09:00:00', '2022-12-31 17:00:00', TRUE, 'image1.jpg', 'http://example.com', 'contact@example.com', 'Notes about Project Alpha'),
('Project Beta', 'Description of Project Beta', '2021-06-01 09:00:00', '2021-12-31 17:00:00', TRUE, 'image2.jpg', 'http://example.com', 'contact@example.com', 'Notes about Project Beta'),
('Project Gamma', 'Description of Project Gamma', '2023-01-01 09:00:00', '2023-12-31 17:00:00', TRUE, 'image3.jpg', 'http://example.com', 'contact@example.com', 'Notes about Project Gamma'),
('Project Delta', 'Description of Project Delta', '2023-02-01 09:00:00', '2023-12-31 17:00:00', TRUE, 'image4.jpg', 'http://example.com', 'contact@example.com', 'Notes about Project Delta'),
('Project Epsilon', 'Description of Project Epsilon', '2023-03-01 09:00:00', '2023-12-31 17:00:00', TRUE, 'image5.jpg', 'http://example.com', 'contact@example.com', 'Notes about Project Epsilon'),
('Project Zeta', 'Description of Project Zeta', '2023-04-01 09:00:00', '2023-12-31 17:00:00', TRUE, 'image6.jpg', 'http://example.com', 'contact@example.com', 'Notes about Project Zeta'),
('Project Eta', 'Description of Project Eta', '2023-05-01 09:00:00', '2023-12-31 17:00:00', TRUE, 'image7.jpg', 'http://example.com', 'contact@example.com', 'Notes about Project Eta'),
('Project Theta', 'Description of Project Theta', '2023-06-01 09:00:00', '2023-12-31 17:00:00', TRUE, 'image8.jpg', 'http://example.com', 'contact@example.com', 'Notes about Project Theta'),
('Project Iota', 'Description of Project Iota', '2023-07-01 09:00:00', '2023-12-31 17:00:00', TRUE, 'image9.jpg', 'http://example.com', 'contact@example.com', 'Notes about Project Iota'),
('Project Kappa', 'Description of Project Kappa', '2023-08-01 09:00:00', '2023-12-31 17:00:00', TRUE, 'image10.jpg', 'http://example.com', 'contact@example.com', 'Notes about Project Kappa'),
('Project Lambda', 'Description of Project Lambda', '2023-09-01 09:00:00', '2023-12-31 17:00:00', TRUE, 'image11.jpg', 'http://example.com', 'contact@example.com', 'Notes about Project Lambda'),
('Project Mu', 'Description of Project Mu', '2023-10-01 09:00:00', '2023-12-31 17:00:00', TRUE, 'image12.jpg', 'http://example.com', 'contact@example.com', 'Notes about Project Mu');

INSERT INTO employees_projects (employee_id, project_id, start_date, end_date) VALUES
(1, 1, '2022-01-01 09:00:00', '2022-12-31 17:00:00'),
(2, 2, '2021-06-01 09:00:00', '2021-12-31 17:00:00');

INSERT INTO technologies (name, description) VALUES
('Python', 'A high-level programming language'),
('JavaScript', 'A programming language commonly used in web development'),
('Java', 'A high-level, class-based, object-oriented programming language'),
('C++', 'A general-purpose programming language created as an extension of the C programming language'),
('Ruby', 'A dynamic, open source programming language with a focus on simplicity and productivity');

INSERT INTO projects_technologies (project_id, technology_id) VALUES
(1, 1),
(2, 2);

CREATE TABLE employees_technologies (
    id INT PRIMARY KEY AUTO_INCREMENT,
    employee_id INT,
    technology_id INT,
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (technology_id) REFERENCES technologies(id)
);