CREATE TABLE Groups (
    "Id" SERIAL PRIMARY KEY,
    "Group_name" VARCHAR(100) NOT NULL UNIQUE
);

CREATE INDEX idx_group_name ON Groups("Group_name");

CREATE TABLE Tracks (
    "Group_id" INT NOT NULL,
    "Song" VARCHAR(100) NOT NULL,
    "Release_date" DATE NOT NULL,
    "Song_lyrics" TEXT NOT NULL,
    "Link" VARCHAR(255) NOT NULL,
    PRIMARY KEY ("Song", "Group_id"),
    FOREIGN KEY ("Group_id") REFERENCES Groups("Id") ON DELETE CASCADE
);

CREATE INDEX idx_group_id ON Tracks("Group_id");


CREATE OR REPLACE FUNCTION check_release_date()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW."Release_date" > CURRENT_DATE THEN
        RAISE EXCEPTION 'Дата релиза не может быть больше сегодняшнего дня.';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_release_date_trigger
BEFORE INSERT ON Tracks
FOR EACH ROW
EXECUTE FUNCTION check_release_date();