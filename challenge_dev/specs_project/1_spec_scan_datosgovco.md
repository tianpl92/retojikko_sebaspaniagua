#Context
Discover and document the official public-service documentation on datos.gov.co.

#Role
developer engineer about api services.

#Action
1.  The database.gov.co endpoint uses the Socrata Open Data API (SODA) protocol. Example: https://datos.gov.co/resources/p6dx-8zbt.json (SECOP I - contratos y convocatorias). No API key is required.
2.  Find a skills hermes to help improve the test request using SODA protocol and install it if the system does not have it.
3.  According documentation, find the exact enpoint to consume data about public proposals data (convocatorias públicas).
4.  Test the enpoint with SODA protocol (can you use python for test):
    4.1 Describe the exact endpoint api contract to consume.
    4.2 Describe the input data with the query params or data params required.
    4.2 Describe output response and format response.
    4.4 Explain results test into file called documentaion.md into working dir using the path: "integration_doc/datosgovco/documentation.md".
5.  The documentation result must be generated for using golang language to consume service.
6.  Find a skills hermes to improve better coding using golang.
7.  Create "test.go" file into working dir using path "integration_doc/datosgovco/test.go". Use golang languge to coding, is not neccesary run the go file.

#Agents and subagents.

For this task use deepseek/deepseek-v4-flash and deepseek/deepseek-v4-pro models from openrouter, if the current model dont use any of them, change hermes model by default using deepseek/deepseek-v4-flash.

1.  Configure models into setup models if not exists.
1.  Use deepseek/deepseek-v4-flash model from openrouter for analyse , scan and generate the documentation.
2.  Use deepseek/deepseek-v4-pro model from openrouter for coding and generate the go files test.



#Delivery expected

1. documentation.md file generated about the datos.gov.co endpoint about how to consume services to obtains public proposals data.
2. Create a test.go file for testing manual. But, do not run go file.
3. The output files there are into working dir, path: "integration_doc/datosgovco/".

Tell me first of all, just to be sure, what do you plan to do?
