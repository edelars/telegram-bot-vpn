use swan;
insert into identities (type, data) VALUES (2,"1235");
SET @last_id_in_identities = LAST_INSERT_ID();
insert into shared_secrets (type,data) VALUES (2,"1235");
SET @last_id_in_shared_secrets = LAST_INSERT_ID();
insert into shared_secret_identity (shared_secret, identity) VALUES (@last_id_in_shared_secrets,@last_id_in_identities);

select * from shared_secret_identity;

select id from identities where data = "aasd";

select * from shared_secrets;
select shared_secret from shared_secret_identity where identity = 1;