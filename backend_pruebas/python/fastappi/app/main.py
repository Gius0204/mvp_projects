# app/main.py
from fastapi import FastAPI
from app.routes import routes
from app.database.db import engine, Base
from app.models import task

# Crea las tablas de la base de datos (si no existen)
Base.metadata.create_all(bind=engine)

app = FastAPI()

# Incluye las rutas
app.include_router(routes.router)
