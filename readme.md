# Challenger Development Guide

## Architecture
The Challenger App is a microservice based application. 
It consists of:

#### Mobile App

#### Backend services
Consists of an Api-Gateway and multiple other services:
- API-Gateway: handles all incomming requests and directs them to appropriate service. Also handles auth.
- UserService: Handles user actions, friends
- ChallengeService: Handles challenges, private and public.
- ChatService: Handles chats.
- TeamService: Handles users in teams.

#### Databases
Currently the databases in use:
- GraphDB: Neo4j
- SQL: MYSQL

## Getting started
Local development can be achived by either installing the nescessary components or by using Docker.

#### Requirements

##### Langauages:
- Go: https://go.dev/doc/install

##### Databases:
- MYSQL
- Neo4j

#### Important scripts

- `run-all.ps1` Starts all the backend microservices. Make sure to allow go.exe to access private networks in security settings, so you dont have to allow each time. NOTE: ATM only works on windows
