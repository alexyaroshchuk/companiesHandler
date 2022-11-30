# Companies handler

## Before start

Setup .env variables(.env.example like example)

## After building the project

1. Need to create user with role - admin (roles have been added for the future)
2. Need to log in into system and copy token
3. Then you will be able to work with companies. You need to pass this token in metadata field like ```{"authorization":"token"}```.
JWT token will be valid for 2 hours.

## System functionality

1. ```Create|Path|Get|Delete``` companies with jwt auth;
2. ```Create|Get``` users;
3. ```Login``` to get JWT token;
4. After store or update company we will produce company to Kafka;

## Some useful commands

1. ```make run``` - to start project, setup all services and migrations
2. ```make lint``` 
3. ```make build-proto``` - generate proto files
4. ```make test``` - in progress
