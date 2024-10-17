-- Create keyspace
CREATE KEYSPACE IF NOT EXISTS swipe_mission WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

-- Use the keyspace
USE swipe_mission;

-- Create session table
CREATE TABLE IF NOT EXISTS session (
    user_id UUID,
    object_type TEXT,
    progress MAP<TEXT, TEXT>,
    PRIMARY KEY (user_id)
);

-- Create feedback table
CREATE TABLE IF NOT EXISTS feedback (
    user_id UUID,
    object_type TEXT,
    image_name TEXT,
    attitude TEXT,
    value BOOLEAN,
    PRIMARY KEY ((user_id), image_name, attitude)
);
