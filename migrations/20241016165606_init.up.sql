CREATE TABLE "users" (
                         "id" SERIAL PRIMARY KEY,
                         "first_name" VARCHAR(50) NOT NULL,
                         "last_name" VARCHAR(50) DEFAULT '',
                         "nickname" VARCHAR(50) UNIQUE NOT NULL,
                         "email" VARCHAR(100) UNIQUE NOT NULL,
                         "password" VARCHAR(200) NOT NULL,
                         "auth_method" INTEGER DEFAULT 1,
                         "confirmed" BOOLEAN DEFAULT false,
                         "created_at" BIGINT NOT NULL,
                         "language_id" INTEGER NOT NULL
);

CREATE TABLE "files" (
                         "id" SERIAL PRIMARY KEY,
                         "url" TEXT UNIQUE NOT NULL
);

CREATE TABLE "transcripts" (
                               "id" SERIAL PRIMARY KEY,
                               "text" TEXT NOT NULL,
                               "accurancy" DOUBLE PRECISION,
                               "language_id" INTEGER NOT NULL,
                               "user_id" INTEGER NOT NULL,
                               "file_id" INTEGER NOT NULL,
                               "theme_id" INTEGER NOT NULL
);

CREATE TABLE "themes" (
                              "id" SERIAL PRIMARY KEY,
                              "level" VARCHAR(3) NOT NULL,
                              "topic" TEXT NOT NULL,
                              "question" TEXT NOT NULL,
                              "language_id" INTEGER NOT NULL
);

CREATE TABLE "languages" (
                             "id" SERIAL PRIMARY KEY,
                             "long_name" VARCHAR(50) NOT NULL,
                             "short_name" VARCHAR(10) NOT NULL
);

ALTER TABLE "transcripts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "transcripts" ADD FOREIGN KEY ("file_id") REFERENCES "files" ("id");

ALTER TABLE "transcripts" ADD FOREIGN KEY ("theme_id") REFERENCES "themes" ("id");

ALTER TABLE "transcripts" ADD FOREIGN KEY ("language_id") REFERENCES "languages" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("language_id") REFERENCES "languages" ("id");

ALTER TABLE "themes" ADD FOREIGN KEY ("language_id") REFERENCES "languages" ("id");
