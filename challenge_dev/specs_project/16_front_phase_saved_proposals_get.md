#Context
Continue to project plan about coding the frontend layer. Now we are integrate the service for obtains saved proposals (GET).

#Role
Expert developer web frontend services

#Requeriments
1. Use html5 and css for pages
2. Use javascript to enable communication between backend service and front service.
3. Use golang for backend service coding if it's necessary.
3. Check if backend service is running, else, run backend service "public proposals" into local machine to be available with frontend layer. (check with /health endpoint).

#Agents - Subagents
1.  Principal agent to run plans, read documents, generate documentation, read md instructions, use deepseek/deepseek-v4-flash model from openrouter.
2.  Subagent to generate coding html css files and coding scripts, coding testing and run sql scripts use deepseek/deepseek-v4-pro model from openrouter.

#Designer description

*   Color palette:
Use the follow color palette to the skeleton front application:
    -primary color: rgb(74, 68, 102) hexadecimal is 4A4466.
    -second color: rgb(110, 173, 188) hexadecimal is 6EADBC.
    -auxiliar colors: rgb(159, 203, 173) hexadecimal is 9FCBAD and rgb(241, 247, 212) hexadecimal is F1F7D4

#Actions

1. Building saved_proposals html page.
        1.1 Create new page for "Saved Proposals".
        1.2 Validate token when load the page. If token expired or is not present from login process, force redirection to login page.
        1.3 When load page validate if the token is present and is not expired, so then, consumes /saved_proposals (GET) endpoint to obtains saved public proposals by current login user id and show public proposals information into a grid equals to dashboard.
        1.4 Build grid for saved proposals equals like dashboard grid page.
        1.5 Omit checkbox for grid.
        1.6 Copy behavior about "Detalles.." button from dashboard page to show all complete information about public proposal saved.
        1.7 Keeps top bar user information.
        1.8 Use suggestions color for the page and use different background vs dashboard page.
        1.9 Create a button navigation into top bar, left side position, for returns  to dashboard page.
        1.10 If the current user does no have saved proposals, shows into grid message "No hay convocatorias guardadas".

2. Add new features into dashboard page:
        2.1  Create floating button into dashboard page, position righ side page above "Guardar..." button. the caption button must be "Ver Convocatorias Guardadas". If the page is scrolling, the buttons keeps position.
        2.2. If button "Ver Convocatorias Guardadas" clicked, so navigate to "Saved Proposals" page.

4. Run up again the frontend service.

#Results expected

1. The application navigates to saved proposals page and list saved proposals by user.

-- 

Do not commit changes yet.

--

Tell me first of all again, just to be sure, what do you plan to do?
