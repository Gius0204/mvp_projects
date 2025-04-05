# backends

## backend

CREATE EXTENSION IF NOT EXISTS "pgcrypto"; -- Necesario para generar UUIDs aleatorios

-- PERSON TABLE
CREATE TABLE Person (
personId UUID PRIMARY KEY DEFAULT gen_random_uuid(),
firstName VARCHAR(255) NOT NULL,
lastName VARCHAR(255) NOT NULL,

    phone VARCHAR(50) NULL,
    address TEXT NULL,
    birthDate DATE NULL,
    nationality VARCHAR(100) NULL,
    gender VARCHAR(50) NULL,
    capacityPerWeek FLOAT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

-- USER ACCOUNT TABLE
CREATE TABLE UserAccount (
userId UUID PRIMARY KEY DEFAULT gen_random_uuid(),
personId UUID UNIQUE NOT NULL REFERENCES Person(personId) ON DELETE SET NULL,
email VARCHAR(255) UNIQUE NOT NULL,
passwordHash TEXT NOT NULL,
nickname VARCHAR(255) NULL,
userType VARCHAR(50) CHECK (userType IN ('Admin', 'Member', 'Guest')) NOT NULL,
timezone VARCHAR(100) NOT NULL,
isActive BOOLEAN DEFAULT TRUE,
createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- AREA TABLE
CREATE TABLE Area (
areaId UUID PRIMARY KEY DEFAULT gen_random_uuid(),
title VARCHAR(255) NOT NULL,
description VARCHAR(255) NULL
);

CREATE TABLE AreaUser (
areaId UUID NOT NULL REFERENCES Area(areaId) ON DELETE CASCADE,
userId UUID NOT NULL REFERENCES UserAccount(userId) ON DELETE CASCADE,
PRIMARY KEY (areaId, userId)
);

-- PROJECT TABLE
CREATE TABLE Project (
projectId UUID PRIMARY KEY DEFAULT gen_random_uuid(),
title VARCHAR(255) NOT NULL,
description TEXT NULL,
creatorUserId UUID NOT NULL REFERENCES UserAccount(userId) ON DELETE SET NULL,
isPrivate BOOLEAN DEFAULT FALSE,
createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ProjectUser (
projectId UUID NOT NULL REFERENCES Project(projectId) ON DELETE CASCADE,
userId UUID NOT NULL REFERENCES UserAccount(userId) ON DELETE CASCADE,
projectRole VARCHAR(50) CHECK (projectRole IN ('Admin', 'Leader', 'Editor', 'Commenter', 'Viewer')) NOT NULL,
PRIMARY KEY (projectId, userId)
);

CREATE TABLE ProjectArea (
projectId UUID NOT NULL REFERENCES Project(projectId) ON DELETE CASCADE,
areaId UUID NOT NULL REFERENCES Area(areaId) ON DELETE CASCADE,
projectRole VARCHAR(50) CHECK (projectRole IN ('Admin', 'Editor', 'Commenter', 'Viewer')) NOT NULL,
PRIMARY KEY (projectId, areaId)
);

-- STATUS TABLE
CREATE TABLE Status (
statusId UUID PRIMARY KEY DEFAULT gen_random_uuid(),
title VARCHAR(255) NOT NULL UNIQUE,
color VARCHAR(50) NOT NULL DEFAULT 'gray',
isCompletedStatus BOOLEAN DEFAULT FALSE
);

-- TASK TABLE
CREATE TABLE Task (
taskId UUID PRIMARY KEY DEFAULT gen_random_uuid(),
projectId UUID NOT NULL REFERENCES Project(projectId) ON DELETE CASCADE,
title VARCHAR(255) NOT NULL,
description TEXT NULL,
creatorUserId UUID NOT NULL REFERENCES UserAccount(userId) ON DELETE SET NULL,
statusId UUID NOT NULL REFERENCES Status(statusId) ON DELETE SET NULL,
priority INT CHECK (priority IN (1,2,3,4)) NULL,
startDate DATE NULL,
dueDate DATE NULL,
estimatedHours FLOAT NULL,
isMilestone BOOLEAN DEFAULT FALSE,
parentTaskId UUID NULL REFERENCES Task(taskId) ON DELETE SET NULL,
createdDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
lastModifiedDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
completedDate TIMESTAMP NULL
);

-- TASK ASSIGNEE TABLE
CREATE TABLE TaskUserAssign (
taskId UUID NOT NULL REFERENCES Task(taskId) ON DELETE CASCADE,
userId UUID NOT NULL REFERENCES UserAccount(userId) ON DELETE CASCADE,
PRIMARY KEY (taskId, userId)
);

-- SUBTASK RELATION TABLE
CREATE TABLE Subtask (
parentTaskId UUID NOT NULL REFERENCES Task(taskId) ON DELETE CASCADE,
childTaskId UUID NOT NULL REFERENCES Task(taskId) ON DELETE CASCADE,
PRIMARY KEY (parentTaskId, childTaskId)
);

-- DEPENDENCY TABLE
CREATE TABLE Dependency (
dependencyId UUID PRIMARY KEY DEFAULT gen_random_uuid(),
predecessorTaskId UUID NOT NULL REFERENCES Task(taskId) ON DELETE CASCADE,
successorTaskId UUID NOT NULL REFERENCES Task(taskId) ON DELETE CASCADE,
dependencyType VARCHAR(2) CHECK (dependencyType IN ('FS', 'FF', 'SS', 'SF')) NOT NULL,
lagDurationMinutes INT NOT NULL DEFAULT 0
);

SELECT _ FROM Project;
SELECT _ FROM person;
SELECT _ FROM useraccount;
SELECT _ FROM status;
SELECT \* FROM task;

INSERT INTO Status (title, color, isCompletedStatus) VALUES
('Sin comenzar', 'gray', FALSE),
('En proceso', 'blue', FALSE),
('Listo', 'green', TRUE);

//debo cambiar atributo status de task sea null, sino imposible cambiarlo con UUID oknose

SELECT \* FROM Project WHERE projectId = 'aa0a887f-82a5-4b2a-96f9-f96e86ef429d';

SELECT conname, confupdtype, confdeltype
FROM pg_constraint
WHERE conrelid = 'Task'::regclass;

SELECT _ FROM UserAccount WHERE userId = 'f089576d-9550-443e-9223-c567cfdac3e4';
SELECT _ FROM Status WHERE statusId = '8135903d-0cfb-45ab-ba1b-3206e0bbf733';

SELECT projectId::TEXT FROM Project;

SELECT conname, contype, conrelid::regclass, confrelid::regclass
FROM pg_constraint
WHERE conrelid = 'Task'::regclass;

INSERT INTO Task (
projectId, title, description, creatorUserId, statusId, priority,
startDate, dueDate, estimatedHours, isMilestone, parentTaskId
) VALUES (
'aa0a887f-82a5-4b2a-96f9-f96e86ef429d',
'Implementar autenticación',
'Desarrollar el sistema de login con OAuth2.',
'f089576d-9550-443e-9223-c567cfdac3e4',
'8135903d-0cfb-45ab-ba1b-3206e0bbf733',
2,
'2025-04-01',
'2025-04-10',
12.5,
FALSE,
NULL
);

SELECT \* FROM Task;

/_ prueba con hash tocado
_/

-- Insertar en la tabla Person
INSERT INTO Person (personId, firstName, lastName, phone, address, birthDate, nationality, gender, capacityPerWeek)
VALUES
('550e8400-e29b-41d4-a716-446655440000', 'Juan', 'Pérez', '123456789', 'Calle Falsa 123', '1990-05-15', 'Peruano', 'Masculino', 40.0);

-- Insertar en la tabla UserAccount
INSERT INTO UserAccount (userId, personId, email, passwordHash, nickname, userType, timezone, isActive)
VALUES
('661e8400-e29b-41d4-a716-446655440111', '550e8400-e29b-41d4-a716-446655440000', 'jperez@example.com', '$2a$10$Qr.pAtYQd/yaMejpU34LcO5jD9g74BRQpZxDiI1u89356b9eM20ea', 'Juanito', 'Admin', 'America/Lima', TRUE);

## backend_preubas: go_fiber + ent

CREATE TABLE tasks (
id SERIAL PRIMARY KEY,
parent_id INTEGER REFERENCES tasks(id),
progress DECIMAL(5, 2) DEFAULT 0.00, -- Porcentaje de progreso (0.00 - 100.00)
estimated_time INTEGER, -- Tiempo estimado en segundos
created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tasks_parent_id ON tasks (parent_id);

-- Insertar datos de prueba en la tabla Task
INSERT INTO tasks (progress, estimated_time)
VALUES (25.50, 3600); -- 1 hora
INSERT INTO tasks (parent_id, progress, estimated_time)
VALUES (1, 75.00, 1800); -- 30 minutos

SELECT \* FROM tasks;
