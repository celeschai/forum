A basic forum app written with Go backend, React + Typescript frotend, and a PostgreSQL database. The front and back are connected via a RESTful API using Mux framework for golang. Do note that this is a project designed to maximise learning of the workings of web development, in particular, APIs. Thus, most of the code are implemeted with only standard libraries.

## User cases:
This forum app consists of three basic content forms: threads, posts, and comments. Users can view(GET), create(POST), edit(PATCH), and delete(DELETE) content. Threads are filtered using a single `tag` that forms the feed users see upon login.

## Authentication:
All users have to sign up for an account and remain signed in to use any of its functions. A JWT token is cached on the frontend as a cookie for authentication purposes and it is verified every time the server receives a request for data from the frontend. Users that are logged in remains logged in for a week (this can be customised under `/backend/auth.go` -> `setCookie` -> `MaxAge`), they are automatically redirected to feed. The JWT token is removed upon signing out.

## Database:
Users can only edit and delete content they create, if they attempt to do so on content created by other users, they will be redirected to sign in page. A user's username acts as a foreign key for all content this user creates. Users will be prevented from signing up with an existing username or email. User passwords are encrypted before they are stored. Please refer to the Entity-Relationship Diagram below for more details on the relationships between the different tables in the database and the parts provided for each type fo content.
![entity-relationship-diagram](https://github.com/celeschai/forum/blob/main/entity-relationship-diagram.png)

## Styling:
The SignUp and SignIn pages have been implemented using [MUI](https://mui.com/material-ui/getting-started/templates/). Changes to the MUI template theme can be made in `/frontend/src/SignIn.tsx`. To facilitate my own learning, I styled the rest of the pages with CSS. CSS offers more customisation, you can make your own changes in `/frontend/src/index.css`.

## Docker:
This project is also designed to be run in docker containers. Environment variables that are set in `/.env` are passed into both `./frontend` and `./backend` via the `compose.yml` file under `environement variables`. Backup jobs are automically done and saved on your local machine every day. 
# The file is configured for deployment on [Render](https://render.com/), should you wish to build your docker containers:

 

To start:
```
docker compose up -d
```
This will
1. create docker images, containers, and volumes using `compose.yml`
2. seed the database 
3. start the react frontend
you can view your website on localhost:3000 

## Spinning up the website with public docker images by me 
visit this [repo](https://github.com/celeschai/forum_docker.git)