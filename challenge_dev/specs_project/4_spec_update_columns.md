#Context
Update document generation schema instructions located in database/2_database_spec.md with the next features.


#Actions steps

1. Read "integration_doc/datos.gov.co/documentation.md" file, analyse the section "Available Columns (59 total)" key columns. That information means all posible fields expected about public proposals.
2. Update "database/2_database_spec.md" instructions, update table creation instructions about "public_calls_proposals" and update columns using key columns from previos step.
3  keep into table the `id` column uuid like primary key.
4. Remove call_name, call_code, call_date, issuer_name, issuer_id, issuer_phone, issuer_email specifications columns.
5. Add new columns specifications about key columns from datos.gov.co integracion "Available columns".
6. Describe the type object necessary for each columns.
7. Describe wich columns are required and not.
8. Update other instructions that are linked with the changes.
9. Dont modify section "Result final", keep instruction.

#Delivery results

1. Just update 2_database_spec.md instructions about public_calls_proposals using columns from datos.gov.co documentation.



Tell me first of all, just to be sure, what do you plan to do?
