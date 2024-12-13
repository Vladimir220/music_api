CREATE TABLE Tracks (
    "Song" VARCHAR(100) NOT NULL,
    "Group_name" VARCHAR(100) NOT NULL,
    "Release_date" DATE NOT NULL,
    "Song_lyrics" TEXT NOT NULL,
    "Link" VARCHAR(255),
    PRIMARY KEY ("Song", "Group_name")
);

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

INSERT INTO Tracks ("Song", "Group_name", "Release_date", "Song_lyrics", "Link") VALUES (
    'Just dance', 'Lady Gaga', 
    '2008-04-08', 
    '[Intro: Lady Gaga & Akon]\nTruth!\nRedOne\nKonvict\nGaga (Oh-oh, eh)\n\n[Verse 1: Lady Gaga & Akon]\nIve had a little bit too much, much (Oh, oh, oh-oh)\nAll of the people start to rush (Start to rush by)\nA dizzy twister dance, cant find my drink or man\nWhere are my keys? I lost my phone, phone (Oh, oh, oh-oh)\n\n[Pre-Chorus: Lady Gaga]\nWhats goin on, on the floor?\nI love this record, baby, but I cant see straight anymore\nKeep it cool, whats the name of this club?\nI cant remember, but its alright, a-alright\n\n[Chorus: Lady Gaga]\nJust dance\nGonna be okay, da-da-doo-doot-n\nJust dance\nSpin that record, babe, da-da-doo-doot-n\nJust dance\nGonna be okay\nDa-da-da-dance, dance, dance\nJust, j-j-just dance\n\n[Verse 2: Lady Gaga & Akon]\nWish I could shut my playboy mouth (Oh, oh, oh-oh)\nHowd I turn my shirt inside out? (Inside out, right)\nControl your poison, babe, roses have thorns, they say\nAnd theyre all gettin hosed tonight (Oh, oh, oh-oh)\n\n[Pre-Chorus: Lady Gaga]\nWhats goin on, on the floor?\nI love this record, baby, but I cant see straight anymore\nKeep it cool, whats the name of this club?\nI cant remember, but its alright, a-alright\n\n[Chorus: Lady Gaga]\nJust dance\nGonna be okay, da-da-doo-doot-n\nJust dance\nSpin that record, babe, da-da-doo-doot-n\nJust dance\nGonna be okay\nDa-da-da-dance, dance, dance\nJust, j-j-just\n\n[Verse 3: Colby ODonis]\nWhen I come through on the dance floor, checking out that catalogue (Hey)\nCant believe my eyes, so many women without a flaw (Hey)\nAnd I aint gon give it up, steady, tryna pick it up like a call (Hey)\nIma hit it, Ima beat it up, latch onto it until tomorrow, yeah\nShorty, I can see that you got so much energy\nThe way you twirlin up them hips round and round\nAnd there is no reason at all why you cant leave here with me\nIn the meantime, stay, let me watch you break it down and\n\n[Chorus: Lady Gaga & Akon]\nDance\nGonna be okay, da-da-doo-doot-n (Oh)\nJust dance (Ooh, yeah)\nSpin that record, babe, da-da-doo-doot-n\nJust dance (Ooh, yeah)\nGonna be okay, da-da-doo-doot-n (Ooh, yeah)\nJust dance (Ooh, yeah)\nSpin that record, babe, da-da-doo-doot-n\nJust dance (Oh)\nGonna be okay, da-da-da-dance (Gonna be okay)\nDance, dance (Yeah)\nJust, j-j-just dance (Oh)\n\n[Interlude: Akon]\nIncredible\nAmazing\nMusic\nWoo!\nLets go!\n\n[Breakdown: Lady Gaga]\nHalf psychotic, sick, hypnotic, got my blueprint, its symphonic\nHalf psychotic, sick, hypnotic, got my blueprint, electronic\nHalf psychotic, sick, hypnotic, got my blueprint, its symphonic\nHalf psychotic, sick, hypnotic, got my blueprint, electronic\n\n[Bridge: Lady Gaga & Akon]\nGo, use your muscle, carve it out, work it, hustle\n(I got it, just stay close enough to get it on)\nDont slow, drive it, clean it, Lysol, bleed it\nSpend the last dough (I got it) in your pock-o (I got it)\n\n[Chorus: Lady Gaga]\nJust dance\nGonna be okay, da-da-doo-doot-n\nJust dance\nSpin that record, babe, da-da-doo-doot-n\nJust dance (Baby)\nGonna be okay, da-da-doo-doot-n\nJust dance\nSpin that record, babe, da-da-doo-doot-n (Oh, baby, yeah)\nJust dance\nGonna be okay (Spin that record, baby, yeah)\nDa-da-da-dance, dance, dance\nJust, j-j-just dance', 
    'https://www.youtube.com/watch?v=2Abk1jAONjw');