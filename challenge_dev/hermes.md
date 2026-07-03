#Hermes md

This file contains the guide for the Hermes agent to work with the code in this repository.

#Role
Development Engineer using Go, HTML5-CSS with JavaScript and PostgreSQL.

#Context

This repository will contain the "Public Calls Portal" program, a web application that connects to the "datos.gov.co" services and allows users to browse, filter, and save public calls for proposals.

#Use Cases

#General Technical Requirements

-The program is a web application.

-The base structure should be divided into two main subprojects (applications):

    1. Backend
    2. Frontend

-The backend project will be used to run a service that will receive requests from different endpoints using RESTful JSON (using GET, POST, and PUT protocols according to development engineering guidelines).
-The frontend project will be used to run the portal's visual layer using HTML5 and JavaScript to communicate with the backend services.

#Workflow guide
A. The project will be develop in two phases. First, the backend phase, and then the frontend phase.

    A-1. Backend Phase:

    -The backend project will be divided into the following phases:

        1. Creation of the database structure using PostgreSQL.

        2. Development of the application with Go to receive and process requests, connecting to the respective database.

        3. Creation of unit tests for each functionality.

    A-2. Phase Frontend:

    -The frontend project will be divided into the following phases:

        1. Generation of the login page

        2. Dashboard: main page with submenus for each functionality.

        3. Integration with the backend to utilize the endpoints.

#Project Structure

-Backend

    -Database (SQL scripts for creating the schema, tables, and initial data)
    -Repository (All the logic for connecting to the database and functions for searching, creating, and/or modifying records)
    -Service (Main logic for running the application)
    -Handler (Input functions for endpoint requests, which will call the functions in the service)
    -Router (To structure the different service paths (endpoints), which will be linked (called) by the handler functions that process the request).

-Frontend

    -Pages (HTML5 pages)
    -Scripts (JavaScript files)
    -Sheets (CSS stylesheets)







