CREATE  TABLE repository ( 
	id                   serial NOT NULL,
	name                 varchar(1000) NOT NULL,
	url                  varchar(1000) NOT NULL,
	created_at           timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	modified_at          timestamptz,
	CONSTRAINT pk_repository PRIMARY KEY ( id )
 );
