# BookStore Cloud

Simple online bookstore platform, broken into independent microservices. Each microservice handles one bounded context and communicates via REST or message queues (event-driven communication).


# High-Level Architecture

## Microservices:

User Service – manages accounts, authentication, and profiles.

Book Catalog Service – stores book data, search, and categories.

Order Service – manages shopping cart and order processing.

Payment Service – simulates payment gateway interaction.

Notification Service – sends order confirmation emails or messages.

## Shared Components:

API Gateway – single entry point for clients.

Message Broker – asynchronous communication (RabbitMQ).

Database per Service – no shared DB, each service owns its data.
