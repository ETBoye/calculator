CREATE TABLE IF NOT EXISTS history (
    calculationId serial PRIMARY KEY,
    sessionId VARCHAR(100) NOT NULL,
    input VARCHAR(1000) NOT NULL,
    outputNum VARCHAR(1000),
    outputDenom VARCHAR(1000),
    outputEstimate VARCHAR(1000),
    error VARCHAR(1000)
);



INSERT INTO history (sessionId, input, outputNum, outputDenom, outputEstimate, error) VALUES 
('init-db-data-session', '0', '0', '1', '0', NULL);
INSERT INTO history (sessionId, input, outputNum, outputDenom, outputEstimate, error) VALUES 
('init-db-data-session', '1', '1', '1', '1', NULL);
INSERT INTO history (sessionId, input, outputNum, outputDenom, outputEstimate, error) VALUES 
('init-db-data-session', '2', '2', '1', '2', NULL);
INSERT INTO history (sessionId, input, outputNum, outputDenom, outputEstimate, error) VALUES 
('init-db-data-session', '3', '3', '1', '3', NULL);
INSERT INTO history (sessionId, input, outputNum, outputDenom, outputEstimate, error) VALUES 
('init-db-data-session', '4', '4', '1', '4', NULL);
INSERT INTO history (sessionId, input, outputNum, outputDenom, outputEstimate, error) VALUES 
('init-db-data-session', '5', '5', '1', '5', NULL);
INSERT INTO history (sessionId, input, outputNum, outputDenom, outputEstimate, error) VALUES 
('init-db-data-session', '6', '6', '1', '6', NULL);
INSERT INTO history (sessionId, input, outputNum, outputDenom, outputEstimate, error) VALUES 
('init-db-data-session', '7', '7', '1', '7', NULL);
INSERT INTO history (sessionId, input, outputNum, outputDenom, outputEstimate, error) VALUES 
('init-db-data-session', '8', '8', '1', '8', NULL);
INSERT INTO history (sessionId, input, outputNum, outputDenom, outputEstimate, error) VALUES 
('init-db-data-session', '9', '9', '1', '9', NULL);
INSERT INTO history (sessionId, input, outputNum, outputDenom, outputEstimate, error) VALUES 
('init-db-data-session', '10', '10', '1', '10', NULL);
INSERT INTO history (sessionId, input, outputNum, outputDenom, outputEstimate, error) VALUES 
('init-db-data-session', '11', '11', '1', '11', NULL);
INSERT INTO history (sessionId, input, outputNum, outputDenom, outputEstimate, error) VALUES 
('init-db-data-session', '12', '12', '1', '12', NULL);
