CREATE TABLE IF NOT EXISTS customers
(
  id SERIAL NOT NULL
    CONSTRAINT customers_pkey
      PRIMARY KEY,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  birth_date DATE NOT NULL,
  gender VARCHAR(6) NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  address VARCHAR(200)
);

INSERT INTO customers(first_name, last_name, birth_date, gender, email, address)
VALUES ('Jose', 'Castro', '1965-07-30', 'female', 'jose.castro@example.com', '3729 Avenida de Castilla'),
('Samuel', 'Lyshaug', '1971-03-25', 'male', 'samuel.lyshaug@example.com', '5900 Paulus plass'),
('Eeli', 'Annala', '1973-01-11', 'male', 'eeli.annala@example.com', '988 Reijolankatu'),
('Valentim', 'Sales', '1977-03-22', 'male', 'valentim.sales@example.com', '25 Beco dos Namorados'),
('Alicia', 'Gallego', '1977-06-16', 'female', 'alicia.gallego@example.com', '3372 Calle de Tetuán'),
('Ricky', 'Weaver', '1985-09-30', 'male', 'ricky.weaver@example.com', '3628 James St'),
('Sarah', 'Wong', '1986-09-08', 'female', 'sarah.wong@example.com', '6141 Brock Rd'),
('Ava', 'Novak', '1987-04-10', 'female', 'ava.novak@example.com', '6968 Concession Road 23'),
('Iiris', 'Wiita', '1996-06-04', 'female', 'iiris.wiita@example.com', '6490 Pispalan Valtatie'),
('Aatu', 'Laurila', '1998-11-19', 'male', 'aatu.laurila@example.com', '3602 Aleksanterinkatu');
