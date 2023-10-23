# Simple Person Data Enrichment Application

## Overview

This is a simple Go application that enriches and stores person data in a MySQL database.


## How to Set Up and Run the Application

### 1. Create the Database

Before using the application, you need to set up a MySQL database and create a table named "People" with the following columns:

You can create the table by running the following SQL command:

```sql
CREATE TABLE People (
    Name VARCHAR(255),
    Surname VARCHAR(255),
    Patronymic VARCHAR(255),
    Age INT,
    Gender VARCHAR(255),
    Nationality VARCHAR(255)
);
```

### 2. Configure Database Connection

Create a `.env` file in the same directory as the application with the following content:

```env
DB_USERNAME=root
DB_PASSWORD=your_password
DB_HOST=127.0.0.1:3306
DB_NAME=testdb
```

Replace `your_password` with your actual MySQL database password.

### 3. Run the Application

Compile and run the Go application using your preferred method. The application will fetch example person data, enrich it with age, gender, and nationality information from external APIs, and insert the results into the "People" table in the MySQL database.
