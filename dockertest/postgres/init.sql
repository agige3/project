CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    user_name varchar(50),
    user_age int,
    group_id int,
    channel_ids int[]
);

insert into users(user_name, user_age, group_id, channel_ids)
values ('Admin', 42, 1, '{1,2,3}'),('User', 42, 2, '{2,3,4}'), ('Custom', 42, 3, '{3,4,5}');
