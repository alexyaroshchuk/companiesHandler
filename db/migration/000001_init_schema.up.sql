CREATE TYPE types AS ENUM ('Corporations', 'NonProfit', 'Cooperative', 'Sole Proprietorship');

CREATE TABLE companies (
    id UUID,
    name varchar(15) unique not null,
    description varchar(3000),
    amount_of_employees int not null,
    registered boolean not null,
    type types not null
);

CREATE TABLE users (
   username varchar(15) unique not null,
   password varchar not null,
   role varchar not null
);
