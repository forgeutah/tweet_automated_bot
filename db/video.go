package database

// CREATE TYPE valid_status AS ENUM ('queued', 'sent', 'failed');

// CREATE TABLE IF NOT EXISTS tweets (
// 	id serial,
// 	twitter_username VARCHAR(15),
// 	tweet_text text,
// 	links text,
// 	send_time timestamp,
// 	status valid_status,
// 	created_at TIMESTAMP DEFAULT now()
// );

var create_query = `
		CREATE TABLE IF NOT EXISTS yt_videos (
			id serial,
			twitter_username VARCHAR(15),
			video_title VARCHAR(255),
			video_playlist VARCHAR(255),
			video_url text,
			created_at TIMESTAMP DEFAULT now(),
			last_sent_at TIMESTAMP 
		);

		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'All Types of Golang Types', 'GoWest Conference', 'https://youtu.be/1RYYsLy9bg8');

		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Field Report: Building a game engine for 300 DEFCON hackers to smash', 'GoWest Conference', 'https://youtu.be/ZtP-IggAlB8');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'The Standard Library Bootcamp', 'GoWest Conference', 'https://youtu.be/CyhmhY7aI-s');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Micro Machine Learning in Go', 'GoWest Conference', 'https://youtu.be/Fq5-KmNr_D8');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Go Without Generics, a Retrospective', 'GoWest Conference', 'https://youtu.be/ZtP-IggAlB8');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Functional Programming with Go', 'GoWest Conference', 'https://youtu.be/SSO78QkmMLs');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'So you think you know Go?', 'GoWest Conference', 'https://youtu.be/U_qVSHYgVSE');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Anatomy of a Gopher - Binary Analysis of Go Binaries', 'GoWest Conference', 'https://youtu.be/ou_B5YZzEeU');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Going Serverless', 'GoWest Conference', 'https://youtu.be/M1AvdQVjhx8');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Writing REST Services for the gRPC-curious', 'GoWest Conference', 'https://youtu.be/rbwvY0YpRPI');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Go Templates: Great Library, Bad Rap', 'GoWest Conference', 'https://youtu.be/m8mKqdyD_C8');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Practical Tips to Creating a Great Engineering Culture', 'GoWest Conference', 'https://youtu.be/TxBz7L9ev18');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Callstacks at Giverny: A Go Graphics CLI', 'GoWest Conference', 'https://youtu.be/4LmTe0P7qpg');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'The beauty of Go for building cross-platform graphical applications
', 'GoWest Conference', 'https://youtu.be/VESXrgCW50g');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Building a scalable API platform for live streaming using GoLang', 'GoWest Conference', 'https://youtu.be/jWUhpOnhOQ0');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Go Channels Demystified', 'GoWest Conference', 'https://youtu.be/KdHqx9Bx4jI');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Ent: Making Data Easy in Go', 'GoWest Conference', 'https://youtu.be/NvjvzYacgQg');
		INSERT into videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Build a home monitoring system with Go', 'GoWest Conference', 'https://youtu.be/PolzSYuwaQg');`
