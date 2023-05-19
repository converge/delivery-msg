# Delivery MSG

## Overview

This project is a simple delivery message system that uses gRPC and message broker to create and track delivery orders.

## Objective

Objective: have a functional gRPC API that create delivery orders and track the delivery progress.

## Brainstorm

![brainstorm](https://user-images.githubusercontent.com/188671/239596772-94425b5d-c089-4968-880d-b27232a0c7ff.png)

## UI

![UI](https://user-images.githubusercontent.com/188671/239596770-950dd221-495b-44c3-9529-d4988a443119.png)

## Architecture

![architecture](https://user-images.githubusercontent.com/188671/239596763-68085444-4dde-4cea-ae12-3e5417979934.png)

--  

### User Stories

#### gRPC API

- [x] As a user, I want to have a gRPC function that can create a delivery order.
- [x] As a user, I want to have a gRPC function that can update a delivery order.
- [x] As a user, I want to have a persistence layer that can create delivery orders.
- [x] As a user, I want to have a persistence layer that can update delivery orders.

#### message broker microservice

- [x] As a user, I want to have a message broker microservice that can publish messages when a order is created or 
updated.

#### UI/client

- [ ] As a user, I want to have a UI that is subscribed to the delivery topic, and can update the screen based on new
messages.

#### implement logging with bubble tea team

- [x] As a user, I want to have a logging system that can log messages to a file.


### todo:
- [ ] validate inputs
