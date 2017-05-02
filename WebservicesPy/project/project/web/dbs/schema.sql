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
  apikey TEXT DEFAULT NULL,
  /*request's counters*/
  stations_request INT DEFAULT NULL,
  rupos_request INT DEFAULT NULL,
  range_request INT DEFAULT NULL
);


CREATE TABLE measurements (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
/* tupos rupou */
  rupos TEXT DEFAULT  NULL,
/* times metrhsewn apo csv */
  value REAL DEFAULT NULL,
/*gia na kseroume se poio sta8mo paei*/
  passname TEXT DEFAULT NULL,
/*hmeromhnia metrhsewn*/
  dt TEXT ISO8601 strings DEFAULT NULL,
/*24 times gia ka8e wra ths meras*/
  one INT DEFAULT NULL,
  two INT DEFAULT NULL,
  three INT DEFAULT NULL,
  four INT DEFAULT NULL,
  five INT DEFAULT NULL,
  six INT DEFAULT NULL,
  seven INT DEFAULT NULL,
  eight INT DEFAULT NULL,
  nine INT DEFAULT NULL,
  ten INT DEFAULT NULL,
  eleven INT DEFAULT NULL,
  twelve INT DEFAULT NULL,
  thirteen INT DEFAULT NULL,
  fourteen INT DEFAULT NULL,
  fifteen INT DEFAULT NULL,
  sixteen INT DEFAULT NULL,
  seventeen INT DEFAULT NULL,
  eighteen INT DEFAULT NULL,
  nineteen INT DEFAULT NULL,
  twenty INT DEFAULT NULL,
  twentyone INT DEFAULT NULL,
  twentytwo INT DEFAULT NULL,
  twentythree INT DEFAULT NULL,
  twentyfour INT DEFAULT NULL
);
