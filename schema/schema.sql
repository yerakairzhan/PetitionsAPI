create table Users(
    id serial primary key,
    username varchar not null,
    password_hash varchar not null,
    created_at timestamp default current_timestamp
);

create table Petitions (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    number_votes INTEGER DEFAULT 0,
    user_id INTEGER NOT NULL REFERENCES Users(id) ON DELETE CASCADE
);


CREATE TABLE votes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    petition_id INT NOT NULL REFERENCES petitions(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);