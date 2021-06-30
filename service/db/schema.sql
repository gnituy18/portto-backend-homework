create database hw;

create table Block (
  block_num bigint not null unique,
  block_hash char(64) not null unique,
  block_time bigint not null,
  parent_hash char(64) not null,
  PRIMARY KEY (block_num)
);

create index block_num_index on Block(block_num);
