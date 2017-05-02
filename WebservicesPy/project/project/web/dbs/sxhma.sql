/* sqlite mode settings, WAL mode is important, otherwise the readers will lock each other and */
/*pragma page_size=4096;
pragma default_cache_size=4000;
pragma journal_mode=WAL;
pragma synchronous=1;*/

CREATE TABLE stations (

    id INTEGER PRIMARY KEY AUTOINCREMENT,
    /* onoma sta8mou */
    name TEXT DEFAULT NULL,
    /* password sta8mou */
    passname TEXT DEFAULT NULL,
    /* gewgrafiko mhkos sta8mou */
    lat FLOAT DEFAULT NULL,
    /* gewgrafiko platos sta8mou */
    long FLOAT DEFAULT NULL
);


CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  /*users mail*/
  mail TEXT DEFAULT NULL,
    /* password */
  password TEXT DEFAULT NULL,
  /*api key*/
  apikey TEXT DEFAULT NULL
);


CREATE TABLE measurements (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
/* tupos rupou */
  rupos TEXT DEFAULT  NULL,
/* times metrhsewn apo csv */
  value REAL DEFAULT NULL,
/*gia na kseroume se poio sta8mo paei*/
  passname TEXT DEFAULT NULL
);
