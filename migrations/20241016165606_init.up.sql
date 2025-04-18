CREATE TABLE "users" (
                         "id" SERIAL PRIMARY KEY,
                         "nickname" VARCHAR(50) UNIQUE NOT NULL,
                         "created_at" BIGINT NOT NULL,
                         "language_id" INTEGER NOT NULL
);

CREATE TABLE "files" (
                         "id" SERIAL PRIMARY KEY,
                         "key" TEXT UNIQUE NOT NULL
);

CREATE TABLE "transcripts" (
                               "id" SERIAL PRIMARY KEY,
                               "text" TEXT default '',
                               "accuracy" DOUBLE PRECISION,
                               "language_id" INTEGER NOT NULL,
                               "user_id" INTEGER NOT NULL,
                               "file_id" INTEGER,
                               "theme_id" INTEGER NOT NULL,
                               "suggestion" JSONB default '{}'
);

CREATE TABLE "themes" (
                              "id" SERIAL PRIMARY KEY,
                              "level" VARCHAR(3) NOT NULL,
                              "topic_id" INTEGER NOT NULL,
                              "question" TEXT UNIQUE NOT NULL,
                              "language_id" INTEGER NOT NULL
);

CREATE TABLE "languages" (
                             "id" SERIAL PRIMARY KEY,
                             "long_name" VARCHAR(50) NOT NULL,
                             "short_name" VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS topics(
                                     "id" SERIAL PRIMARY KEY,
                                     "language_id" INTEGER NOT NULL,
                                     "name" text NOT NULL
);


ALTER TABLE "transcripts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "transcripts" ADD FOREIGN KEY ("theme_id") REFERENCES "themes" ("id");

ALTER TABLE "transcripts" ADD FOREIGN KEY ("language_id") REFERENCES "languages" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("language_id") REFERENCES "languages" ("id");

ALTER TABLE "themes" ADD FOREIGN KEY ("language_id") REFERENCES "languages" ("id");

ALTER TABLE "themes" ADD FOREIGN KEY ("topic_id") REFERENCES "topics" ("id");


