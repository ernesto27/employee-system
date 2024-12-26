# Employee Management System

This is an Employee Management System built with Go, providing functionalities to manage employees, projects, roles, and technologies.

## Table of Contents

- [Features](#features)
- [Technologies](#technologies)
- [Setup](#setup)
- [AWS EC2 Deployment](#aws-ec2-deployment)

## Features

- Manage employees, roles, projects, and technologies.
- Create, read, update, and delete operations for employees.
- Assign roles and technologies to employees.
- Manage project details and assign employees to projects.

## Technologies

- Go
- MySQL
- Docker
- HTML/CSS/JavaScript


## Setup

### Prerequisites

- Go 1.23
- Docker


### Environment Variables

copy .env.example to .env and update the values

```bash
DB_HOST=localhost
DB_PORT=3398
DB_USERNAME=root
DB_PASSWORD=root
DB_DATABASE=employees-system

AWS_REGION=
AWS_S3_BUCKET=
AWS_ACCESS_KEY=
AWS_SECRET=

URL="http://localhost:8080"
```

### Local Setup


Create DB using make
```bash
make refresh-db
```

Start server:
```bash
go run main.go
```

The server will be available at http://localhost:8080

### AWS EC2 Deployment

Connect to your EC2 instance
```bash
ssh -i your-key.pem ec2-user@your-ec2-ip
```

Update .env with production values



Start service
```bash
docker-compose -f docker-compose.prod.yml up -d
```

The application will be available at your EC2 instance's public IP on port 9966.

### Database Backups


Setup daily cron job:
```bash
sudo crontab -e
```

Add this line to run backup daily at 3 AM:
```bash
0 3 * * * /path-to-repo/backups/backup.sh
```