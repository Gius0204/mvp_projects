# app/database/db.py
from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
import os
from dotenv import load_dotenv

# load_dotenv(dotenv_path="../.env")  # Carga las variables de entorno desde el archivo .env

# DATABASE_URL = os.getenv("DATABASE_URL")
# print("HOLA DB")
# print(DATABASE_URL)  # Verificar si se está cargando correctamente

# Configuración de la base de datos
DATABASE_URL ="postgresql://postgres:postgresql123@localhost:5432/db_pruebas"
engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

# Declarative base
Base = declarative_base()

# Obtener sesión de la base de datos
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
