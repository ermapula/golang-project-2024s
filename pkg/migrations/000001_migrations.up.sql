CREATE TABLE IF NOT EXISTS publishers (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    headquarters text NOT NULL,
    website text NOT NULL
);

CREATE TABLE IF NOT EXISTS games (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    genres text[] NOT NULL,
    price double precision NOT NULL DEFAULT 0,
    release_date timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    publisher_id bigserial NOT NULL REFERENCES publishers (id),
    version integer NOT NULL DEFAULT 1
);

INSERT INTO publishers (name, headquarters, website) 
	VALUES
		('Electronic Arts', 'Redwood City, California, USA', 'https://www.ea.com'),
		('Ubisoft', 'Montreuil, France', 'https://www.ubisoft.com'),
		('Nintendo', 'Kyoto, Japan', 'https://www.nintendo.com'),
		('Activision Blizzard', 'Santa Monica, California, USA', 'https://www.activisionblizzard.com'),
		('FromSoftware', 'Tokyo, Japan', 'https://www.fromsoftware.jp');
		
INSERT INTO games (title, genres, release_date, price, publisher_id) 
	VALUES 
		('Battlefield V', '{"fps", "action"}', '2018-11-20', 59.99, 1),
		('Assassin''s Creed Unity', '{"action-adventure", "stealth"}', '2014-11-11', 39.99, 2),
		('The Legend of Zelda: Breath of the Wild', '{"action-adventure", "open world", "puzzle"}', '2017-03-03', 59.99, 3),
		('Call of Duty: Warzone', '{"battle royale", "fps"}', '2020-03-10', 0, 4),
		('Elden Ring', '{"action-rpg", "souls-like"}', '2022-02-25', 59.99, 5),
		('Apex Legends', '{"fps", "battle royale"}', '2019-02-04', 0, 1),
		('Far Cry 6', '{"action-adventure", "fps", "open world"}', '2021-10-07', 59.99, 2),
		('Super Mario Odyssey', '{"platform", "action-adventure"}', '2017-10-27', 49.99, 3),
		('Call of Duty: Ghosts', '{"fps"}', '2014-03-25', 59.99, 4),
		('Dark Souls III', '{"action-rpg", "souls-like"}', '2019-03-22', 59.99, 5),
		('Super Mario Sunshine', '{"platform", "action-adventure"}', '2002-07-19', 19.99, 3),
		('Far Cry 5', '{"action-adventure", "fps", "open world"}', '2018-03-27', 59.99, 2),
		('The Legend of Zelda: Ocarina of Time', '{"action-adventure", "open world", "puzzle"}', '1998-11-23', 39.99, 3),
		('Call of Duty: Black Ops 4', '{"fps"}', '2018-10-12', 59.99, 4),
		('Bloodborne', '{"action-rpg", "souls-like"}', '2015-03-24', 19.99, 5),
		('Super Mario 64', '{"platform", "action-adventure"}', '1996-06-23', 9.99, 3),
		('Far Cry 4', '{"action-adventure", "fps", "open world"}', '2014-11-18', 59.99, 2),
		('The Legend of Zelda: A Link to the Past', '{"action-adventure", "open world", "puzzle"}', '1992-04-13', 19.99, 3),
		('Call of Duty: Modern Warfare', '{"fps"}', '2019-10-25', 59.99, 4),
		('Demon''s Souls', '{"action-rpg", "souls-like"}', '2020-11-12', 69.99, 5),
		('Super Mario World', '{"platform", "action-adventure"}', '1990-11-21', 9.99, 3),
		('Far Cry 3', '{"action-adventure", "fps", "open world"}', '2012-12-04', 19.99, 2),
		('The Legend of Zelda: Link''s Awakening', '{"action-adventure", "open world", "puzzle"}', '1993-06-06', 19.99, 3),
		('Call of Duty: Black Ops Cold War', '{"fps"}', '2020-11-13', 59.99, 4),
		('Sekiro: Shadows Die Twice', '{"action-adventure", "souls-like"}', '2019-03-22', 59.99, 5),
		('Super Mario Bros.', '{"platform", "action-adventure"}', '1985-09-13', 9.99,),
		('Far Cry 2', '{"action-adventure", "fps", "open world"}', '2008-10-21', 19.99, 2),
		('The Legend of Zelda: Majora''s Mask', '{"action-adventure", "open world", "puzzle"}', '2000-04-27', 39.99, 3),
		('Call of Duty: Infinite Warfare', '{"fps"}', '2016-11-04', 59.99, 4),
		('Dark Souls II', '{"action-rpg", "souls-like"}', '2014-03-11', 19.99, 5),
		('Super Mario Bros. 3', '{"platform", "action-adventure"}', '1988-10-23', 9.99, 3),
		('Far Cry', '{"action-adventure", "fps", "open world"}', '2004-03-23', 19.99, 2),
		('The Legend of Zelda: Twilight Princess', '{"action-adventure", "open world", "puzzle"}', '2006-11-19', 39.99, 3),
		('Call of Duty: Advanced Warfare', '{"fps"}', '2014-11-04', 59.99, 4),
		('Dark Souls', '{"action-rpg", "souls-like"}', '2011-09-22', 19.99, 5),
		('Super Mario Bros. 2', '{"platform", "action-adventure"}', '1988-09-13', 9.99, 3),
		('Far Cry 3: Blood Dragon', '{"action-adventure", "fps", "open world"}', '2013-04-30', 14.99, 2),
		('The Legend of Zelda: Skyward Sword', '{"action-adventure", "open world", "puzzle"}', '2011-11-18', 39.99, 3),
		('Call of Duty: World at War', '{"fps"}', '2008-11-11', 59.99, 4),
		('Bloodborne: The Old Hunters', '{"action-rpg", "souls-like"}', '2015-11-24', 19.99, 5),
		('Super Mario Galaxy', '{"platform", "action-adventure"}', '2007-11-01', 19.99, 3),
		('Far Cry Primal', '{"action-adventure", "fps", "open world"}', '2016-02-23', 19.99, 2),
		('The Legend of Zelda: The Wind Waker', '{"action-adventure", "open world", "puzzle"}', '2002-12-13', 39.99, 3),
		('Call of Duty: Black Ops III', '{"fps"}', '2015-11-06', 59.99, 4),
		('Sekiro: Shadows Die Twice - Game of the Year Edition', '{"action-adventure", "souls-like"}', '2020-10-28', 69.99, 5),
		('Super Mario Galaxy 2', '{"platform", "action-adventure"}', '2010-05-23', 19.99, 3),
		('Far Cry New Dawn', '{"action-adventure", "fps", "open world"}', '2019-02-15', 39.99, 2),
		('The Legend of Zelda: Phantom Hourglass', '{"action-adventure", "open world", "puzzle"}', '2007-06-23', 29.99, 3),
		('Call of Duty: Black Ops II', '{"fps"}', '2012-11-13', 59.99, 4),
		('Sekiro: Shadows Die Twice - Game of the Year Edition', '{"action-adventure", "souls-like"}', '2020-10-28', 69.99, 5),
		('Super Mario 3D World', '{"platform", "action-adventure"}', '2013-11-22', 59.99, 3),
		('Far Cry 6', '{"action-adventure", "fps", "open world"}', '2021-10-07', 59.99, 2),
		('The Legend of Zelda: Spirit Tracks', '{"action-adventure", "open world", "puzzle"}', '2009-12-07', 29.99, 3);
