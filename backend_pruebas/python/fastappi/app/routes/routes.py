# app/routes/routes.py
from fastapi import APIRouter
from app.controllers import task_controller

router = APIRouter()

# Registra las rutas del controlador
router.include_router(task_controller.router)
