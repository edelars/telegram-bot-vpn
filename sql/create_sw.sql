use swan;
insert into identities (type, data) VALUES (2,"1235");
SET @last_id_in_identities = LAST_INSERT_ID();
insert into shared_secrets (type,data) VALUES (2,"1235");
SET @last_id_in_shared_secrets = LAST_INSERT_ID();
insert into shared_secret_identity (shared_secret, identity) VALUES (@last_id_in_shared_secrets,@last_id_in_identities);

select * from shared_secret_identity;

select * from identities;

select * from shared_secrets;
select * from shared_secret_identity ;

0x582737343639366436313631373336343733363427
0x582737343639366436313631373336343733363427

update  shared_secrets set  type = 1 where id =1;

truncate shared_secrets;
truncate identities;
truncate shared_secret_identity;