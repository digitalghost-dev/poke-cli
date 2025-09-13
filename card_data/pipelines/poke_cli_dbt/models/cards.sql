{{ config(materialized='table') }}

SELECT  id, image, name, "localId", category
FROM {{ source('staging', 'cards') }}
