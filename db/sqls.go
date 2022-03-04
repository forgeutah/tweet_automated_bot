package database

var videoInsert = `
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'All Types of Golang Types', 'GoWest Conference', 'https://youtu.be/1RYYsLy9bg8', '2020', '@carson_ops');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'Field Report: Building a game engine for 300 DEFCON hackers to smash', 'GoWest Conference', 'https://youtu.be/ZtP-IggAlB8', '2020', '@astockwell and @WarOnShrugs');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'The Standard Library Bootcamp', 'GoWest Conference', 'https://youtu.be/CyhmhY7aI-s', '2020', '@JeremyCMorgan');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'Micro Machine Learning in Go', 'GoWest Conference', 'https://youtu.be/Fq5-KmNr_D8', '2020', 'Joshua Bowles');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'Go Without Generics, a Retrospective', 'GoWest Conference', 'https://youtu.be/UqEtKx_9Wc8', '2020', '@lostluck');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'Functional Programming with Go', 'GoWest Conference', 'https://youtu.be/SSO78QkmMLs', '2020', 'Dylan Meeus');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'So you think you know Go?', 'GoWest Conference', 'https://youtu.be/U_qVSHYgVSE', '2020', '@corylanou');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'Anatomy of a Gopher - Binary Analysis of Go Binaries', 'GoWest Conference', 'https://youtu.be/ou_B5YZzEeU', '2020', 'Alex Useche');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'Going Serverless', 'GoWest Conference', 'https://youtu.be/M1AvdQVjhx8', '2020', '@bogaczio');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'Writing REST Services for the gRPC-curious', 'GoWest Conference', 'https://youtu.be/rbwvY0YpRPI', '2020', '@JohanBrandhorst');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'Go Templates: Great Library, Bad Rap', 'GoWest Conference', 'https://youtu.be/m8mKqdyD_C8', '2021', '@carson_ops');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'Practical Tips to Creating a Great Engineering Culture', 'GoWest Conference', 'https://youtu.be/TxBz7L9ev18', '2021', '@clintberry');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'Callstacks at Giverny: A Go Graphics CLI', 'GoWest Conference', 'https://youtu.be/4LmTe0P7qpg', '2021', '@juliecoding');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'The beauty of Go for building cross-platform graphical applications', 'GoWest Conference', 'https://youtu.be/VESXrgCW50g', '2021', '@andydotxyz');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'Building a scalable API platform for live streaming using GoLang', 'GoWest Conference', 'https://youtu.be/jWUhpOnhOQ0', '2021', 'Amit Mishra');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'Go Channels Demystified', 'GoWest Conference', 'https://youtu.be/KdHqx9Bx4jI', '2021', '@moficodes');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'Ent: Making Data Easy in Go', 'GoWest Conference', 'https://youtu.be/NvjvzYacgQg', '2021', '@DmitryVinnik');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url, conference_year, presenter_twitter_username) VALUES ('gowestconf', 'Build a home monitoring system with Go', 'GoWest Conference', 'https://youtu.be/PolzSYuwaQg', '2021', '@dlsniper');`
