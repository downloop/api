CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users(id UUID NOT NULL DEFAULT uuid_generate_v1(),
				 username VARCHAR(32),
				 created TIMESTAMP DEFAULT current_timestamp);
CREATE TABLE IF NOT EXISTS sessions(id UUID NOT NULL DEFAULT uuid_generate_v1(),
                                    user_id UUID, 
				    start_time TIMESTAMP, 
			            end_time TIMESTAMP,
				    created TIMESTAMP DEFAULT current_timestamp,
				    CONSTRAINT session_pkey PRIMARY KEY(id));
CREATE TABLE IF NOT EXISTS coordinates(session_id UUID, 
				       lat FLOAT, 
				       lon FLOAT, 
				       recorded TIMESTAMP,
				       created TIMESTAMP DEFAULT current_timestamp);
CREATE TABLE IF NOT EXISTS altitude(session_id UUID, 
				    alt FLOAT, 
				    recorded TIMESTAMP,
				    created TIMESTAMP DEFAULT current_timestamp);