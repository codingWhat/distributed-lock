
# 第一种实现方式， unique key + last_lock_time
CREATE  TABLE  tbl_lock_info_v1 (
    id bigint not null auto_increment,
    holder varchar(200) not null default  '',
    lock_key  varchar(200) not null  default '',
    ttl int not null default  0,
    last_lock_time int null default  0 ,
    primary  key (id),
    unique  unk (lock_key)
) engine=innodb default charset=utf8mb4;

#第二种实现方式
CREATE  TABLE tbl_lock_info_v2 (
    id bigint not null  not null auto_increment,
    lock_key varchar(200) not null  default  '',
    primary key (id),
    unique unk(lock_key)
) engine=innodb default charset=utf8mb4;