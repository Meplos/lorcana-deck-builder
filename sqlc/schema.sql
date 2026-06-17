--- Schema.sql


CREATE TABLE cards (
    "id" TEXT PRIMARY KEY,
    "name" TEXT NOT NULL,
    "setId" TEXT NOT NULL,
    "versionId" INTEGER NOT NULL,
    "title" TEXT,
    "cost" INTEGER NOT NULL,
    "type" TEXT NOT NULL,
    "number" TEXT NOT NULL,
    "colorMask" INTEGER NOT NULL,
    "illustrator" TEXT,
    "keywords" TEXT,
    "lore" INTEGER,
    "strength" INTEGER,
    "willpower" INTEGER,
    "movement" INTEGER,
    "ink" INTEGER NOT NULL,
    "characteristics" TEXT NOT NULL,
    "abilities" TEXT,
    "variants" TEXT NOT NULL,
    "rarity" TEXT,
    "language" TEXT NOT NULL,
    "path" TEXT,
    "franchise" TEXT,
    "ordinal" INTEGER NOT NULL,
    "formats" TEXT
);

