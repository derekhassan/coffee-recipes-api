DROP TABLE IF EXISTS recipes;
CREATE TABLE recipes (
  id          INT AUTO_INCREMENT NOT NULL,
  title       VARCHAR(255) NOT NULL,
  steps       TEXT NOT NULL,
  gramsCoffee DECIMAL(10,2) NOT NULL,
  gramsWater  DECIMAL(10,2) NOT NULL,
  waterTempCelsius  DECIMAL(10,2) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO recipes(title,steps,gramsCoffee,gramsWater, waterTempCelsius)
VALUES ('A Great V60 Recipe', 'Pour coffee and win!', 15, 250, 100);