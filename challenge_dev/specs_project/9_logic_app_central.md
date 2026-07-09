#Context
Continue and improve better the backend phase project, now we go to implement the integration with datos.go.vo service.

#Role
Software engineer developer api services.

#Agents - Subagents
1.  Principal agent to run plans, read documents, generate documentation, read md instructions, use deepseek/deepseek-v4-flash model from openrouter.
2.  Subagent to generate go files and coding scripts and coding testing use deepseek/deepseek-v4-pro model from openrouter.

#tools requeriments
1. Use golang to coding, use the golang version current machine.

#Documentation api external

1. Use technical documentation for integration from working dir into path "integration_doc/datosgovco/documentation.md".
2. For coding integration and improve better code; there is an example file into "integration_doc/datosgovco/test.go".

#Actions Coding

1. Integrate endpoint /public-proposals (GET) to obtains public proposals from datos.gov.co endpoint service.
2. Prepare and configure into .env file a new environment var called "INTEGRATION_URL". The value is the exact base URL endpoint for service datos.gov.co, described into documentation from "integration_doc/datosgovco/documentation.md" file.
3. Use the var "INTEGRATION_URL", to load into service up and calling its external api services.
4. When obtains data from datos.gov.co, only return response data. Don't save into database yet. Just, the target of this endpoint is pull the info from that api external (datos.gov.co) and use for response.
5. The response for the endpoint /public-proposals must be a json object using the same data columns name info for public proposals. Example:

    ```json
    {
      "entidad": "DEPARTAMENTO ADMINISTRATIVO NACIONAL DE ESTADISTICA (DANE)",
      "nit_entidad": "899999027",
      "departamento_entidad": "Distrito Capital de Bogotá",
      "ciudad_entidad": "Bogotá",
      "ordenentidad": "Nacional",
      "id_del_proceso": "CO1.REQ.2577563",
      "referencia_del_proceso": "EDP-545-2022",
      "nombre_del_procedimiento": "Convocatoria DSNFT-0001-FEEC-2024",
      "descripci_n_del_procedimiento": "Aunar esfuerzos entre las partes para ejecutar el proyecto...",
      "fase": "Presentación de oferta",
      "fecha_de_publicacion_del": "2024-02-19T00:00:00.000",
      "fecha_de_ultima_publicaci": "2024-03-12T00:00:00.000",
      "modalidad_de_contratacion": "Licitación pública",
      "precio_base": 500000000.00,
      "estado_del_procedimiento": "Publicado",
      "urlproceso": "https://community.secop.gov.co/..."
    }
    ```

6. keep validation about token JWT. The endpoint only works if Token JWT is available and active.
7. Dont commit and save changes into current branch yet.


---------

Tell me first of all, just to be sure, what do you plan to do?



