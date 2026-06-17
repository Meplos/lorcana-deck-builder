-- Queries.sql

-- name: ListCards :many
SELECT 
    id ,
    "name" ,
    "setId" ,
    "versionId",
    "title" ,
    "cost" ,
    "type" ,
    "number" ,
    "colorMask" ,
    "illustrator" ,
    "keywords" ,
    "lore" ,
    "strength" ,
    "willpower",
    "movement" ,
    "ink" ,
    "characteristics" ,
    "abilities" ,
    "variants" ,
    "rarity" ,
    "language" ,
    "path" ,
    "franchise" ,
    "ordinal" ,
    "formats" 
FROM cards
