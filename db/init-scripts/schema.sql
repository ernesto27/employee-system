CREATE TABLE employees (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    email VARCHAR(255),
    age INT,
    images TEXT,
    start_work_date DATE,
    end_work_date DATE,
    active BOOLEAN,
    created_at DATETIME,
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

INSERT INTO employees (name, email, age, images, start_work_date, end_work_date, active, created_at, updated_at) VALUES
('John Doe', 'john.doe@example.com', 30, 'image1.jpg', '2022-01-01', '2023-01-01', TRUE, NOW(), NOW()),
('Jane Smith', 'jane.smith@example.com', 25, 'image2.jpg', '2021-06-01', '2022-06-01', TRUE, NOW(), NOW());

INSERT INTO roles (name) VALUES
('Developer'),
('Manager');

INSERT INTO employees_roles (employee_id, role_id) VALUES
(1, 1),
(2, 2);

INSERT INTO projects (name, description, start_date, end_date, active, images, links, contacts, notes) VALUES
('Project Alpha', 'Description of Project Alpha', '2022-01-01 09:00:00', '2022-12-31 17:00:00', TRUE, 'image1.jpg', 'http://example.com', 'contact@example.com', 'Notes about Project Alpha'),
('Project Beta', 'Description of Project Beta', '2021-06-01 09:00:00', '2021-12-31 17:00:00', TRUE, 'image2.jpg', 'http://example.com', 'contact@example.com', 'Notes about Project Beta');

INSERT INTO employees_projects (employee_id, project_id, start_date, end_date) VALUES
(1, 1, '2022-01-01 09:00:00', '2022-12-31 17:00:00'),
(2, 2, '2021-06-01 09:00:00', '2021-12-31 17:00:00');

INSERT INTO technologies (name, description) VALUES
('Python', 'A high-level programming language'),
('JavaScript', 'A programming language commonly used in web development');

INSERT INTO projects_technologies (project_id, technology_id) VALUES
(1, 1),
(2, 2);