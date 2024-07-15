<h1 align="center">Database Synchronization</h1>

<p align="center">
  <a href="#dart-about">About</a> &#xa0; | &#xa0;
  <a href="#rocket-technologies">Technologies</a> &#xa0; | &#xa0;
  <a href="#how-the-problem-was-solved">Solution</a> &#xa0; | &#xa0;
  <a href="#white_check_mark-requirements">Requirements</a> &#xa0; | &#xa0;
  <a href="#checkered_flag-starting">Starting</a> &#xa0; | &#xa0;
  <a href="#api-documentation">API Doc</a> &#xa0; | &#xa0;
  <a href="https://github.com/BernardoDenkvitts" target="_blank">Author</a>
</p>

<br>

## :dart: About

This project aims to synchronize the databases of three independent Golang Web APIs. Each API manages its own data and creates users independently through its endpoints. However, to ensure data consistency, each 30 seconds a synchronization process happens to replicate the latest users created across all databases.

## :rocket: Technologies

The following tools were used in this project:

- [Golang](https://go.dev/)
- [RabbitMQ](https://www.rabbitmq.com/)
- [Docker](https://www.docker.com/)
- [MySQL](https://www.mysql.com/)
- [PostgresSQL](https://www.postgresql.org/)
- [MongoDB](https://www.mongodb.com/)

## How the problem was solved

The synchronization problem was solved using RabbitMQ and creating users ID (UUID) in application level to keep the same IDs through different applications.

The use of RabbitMQ was extremely important because of the exchanges and queues it provides us.

Setting the exchanges as "fanout" makes possible to send the data to all of the queues that are bound to it. Enabling an application send its data to other applications.

### What happens if the queue receives a user that has already been created in the database ?

This shouldn't happen, but if it does, the received user won't be created again because users ID are primary keys, the downside would be processing an already created user.

## The following image shows how the data is being shared between applications

![Synchronization.png](/Synchronization.png)

As you can see, each application queue receives data from the other applications exchanges, making possible the synchronization.

## :white_check_mark: Requirements

If you would like to test the project, before starting :checkered_flag:, you need to have [Golang](https://go.dev/) and [Docker](https://www.docker.com/) installed.

## :checkered_flag: Starting

```bash
# Clone this project
$ git clone https://github.com/BernardoDenkvitts/Database-Synchronization

# Access
$ cd Database-Synchronization

# To start databases and rabbitmq
$ docker-compose up --build -d

# Starting each application

# Mongo APP
$ cd MongoAPP
$ go run .\cmd\main.go

# MySQL APP
$ cd MySQLAPP
$ go run .\cmd\main.go

# PostgresAPP APP
$ cd MySQLAPP
$ go run .\cmd\main.go

# Mongo APP will initialize in http://localhost:8181
# MySQLAPP will initialize in http://localhost:8080
# PostgresAPP will initialize in http://localhost:8282
```

## API Documentation

<a href="#top">Back to top</a>
