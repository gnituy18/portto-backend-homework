create table blocks (
  block_num bigint not null unique,
  block_hash char(66) not null unique,
  block_time bigint not null,
  parent_hash char(66) not null unique,
  primary key(block_num)
);

create index block_num_index on blocks(block_num);
create index block_hash_index on blocks(block_hash);

create table transactions (
  tx_hash char(66) not null unique,
  from_acc char(66),
  to_acc char(66),
  nonce bigint,
  data text,
  value bigint,
  logs text,
  primary key(tx_hash)
);
