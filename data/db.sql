/*
To populate database:
$ psql -d cattledog -f data/db.sql
*/

DROP TABLE IF EXISTS categories CASCADE;
CREATE TABLE categories (  
  id SERIAL PRIMARY KEY,  
  name VARCHAR(50) NOT NULL
);

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(80) UNIQUE NOT NULL, 
  password VARCHAR(80) NOT NULL
);

DROP TABLE IF EXISTS items CASCADE;
CREATE TABLE items (
  id SERIAL PRIMARY KEY, 
  title VARCHAR(50) NOT NULL,
  description VARCHAR(250) NOT NULL,
  cat_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL
);

INSERT INTO categories (name) VALUES ('Soccer');
INSERT INTO categories (name) VALUES ('Basketball');
INSERT INTO categories (name) VALUES ('Baseball');
INSERT INTO categories (name) VALUES ('Frisbee');
INSERT INTO categories (name) VALUES ('Snowboarding');
INSERT INTO categories (name) VALUES ('Rock Climbing');
INSERT INTO categories (name) VALUES ('Football');
INSERT INTO categories (name) VALUES ('Surfing');
INSERT INTO categories (name) VALUES ('Hockey');
INSERT INTO categories (name) VALUES ('Fishing');

INSERT INTO items (title, description, cat_id, user_id) VALUES ('Boots', 'The latest in comfort, style and cost!', 1, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Indoor ball', 'As used in the NBA', 2, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Outdoor ball', 'Great for local pickup games', 2, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Net', 'The only net for that full swoosh sound', 2, 2);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Wooden bat', 'What a beautiful piece of hickory', 3, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Ball', 'Handmade - childrens fingers make for exquisite stitching', 3, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Frisbee', 'Well...duh!', 4, 2);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Snowboard', 'Boutique board made from carbon fibre, the optimum mix of light-weight strength and profit margin', 5, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Goggles', 'Wrap around UV protection so you can see when you are being ripped off', 5, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Rope', 'Supports a climber up to 150kgs (as long as their fingers do)', 6, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Rocks', 'Do we really sell these???', 6, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Football', 'The good old pigskin', 7, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Helmet', 'Best worn before your time in the concussion protocol', 7, 2);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Shoulder pads', 'Bigger than those in an 80s sitcom', 7, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Surfboard', 'Almost obligatory', 8, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Stick', 'Field or ice?', 9, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Ice skates', 'Well that explains the stick type', 9, 2);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Shin pads', 'Necessary for the times when not faking injury', 1, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Mitt', 'Specifically designed for those with butter fingers', 3, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Wax', 'Protect what is probably your only investment', 8, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Puck', 'Without this hockey could just be called fighting on ice', 9, 2);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Ski poles', 'With these you can create a new hybrid sport', 5, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Metal bat', 'What a beautiful hunk of metal', 3, 1);
INSERT INTO items (title, description, cat_id, user_id) VALUES ('Cleats', 'Engineered by NASA for maximum grip on AstroTurf', 7, 1);