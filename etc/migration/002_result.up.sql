CREATE TABLE result (
	id                   serial NOT NULL,
	status               varchar(50) DEFAULT 'Queued',
	repo_id              serial NOT NULL,
	findings             jsonb,
	queued_at            timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	scanning_at          timestamptz,
	finished_at          timestamptz,
	CONSTRAINT pk_tbl PRIMARY KEY ( id )
);

ALTER TABLE result ADD CONSTRAINT fk_result_repository FOREIGN KEY ( repo_id ) REFERENCES repository( id );
