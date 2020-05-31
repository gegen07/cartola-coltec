DROP TABLE IF EXISTS formations cascade;
CREATE TABLE formations (
	id          SERIAL PRIMARY KEY,
	goalkeeper  SMALLINT NOT NULL,
	defenders   SMALLINT NOT NULL,
	attackers   SMALLINT NOT null,
	created_at TIMESTAMPTZ(6) NOT null,
	updated_at timestamptz(6) NOT null
);


DROP TABLE IF EXISTS scouts cascade;
CREATE TABLE scouts (
    id SERIAL PRIMARY KEY,
    scout varchar(255) NOT null,
    description text NOT null,
    points decimal not null,
    created_at TIMESTAMPTZ(6) NOT null,
    updated_at timestamptz(6) NOT null
);

DROP TABLE IF EXISTS positions cascade;
CREATE TABLE positions (
    id SERIAL PRIMARY KEY,
    description varchar(255)  NOT NULL,
    created_at TIMESTAMPTZ(6) NOT null,
    updated_at timestamptz(6) NOT null
);

drop table if exists scout_position cascade; 
create table scout_position (
	id serial primary key,
	position_id integer references positions(id),
	scout_id integer references scouts(id),
	created_at TIMESTAMPTZ(6) NOT null,
    updated_at timestamptz(6) NOT null
);

