# app/controllers/task_controller.py
from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from app.models.task import Task as TaskModel
from app.schemas.task import Task as TaskSchema
from app.database.db import get_db
from typing import List, Optional

router = APIRouter()

# Obtener todas las tareas
@router.get("/tasks", response_model=List[TaskSchema])
def get_tasks(db: Session = Depends(get_db)):
    tasks = db.query(TaskModel).all()
    return tasks

# Crear una nueva tarea
@router.post("/tasks", response_model=TaskSchema)
def create_task(progress: float, estimated_time: int, parent_id: Optional[int] = None, db: Session = Depends(get_db)):
    task = TaskModel(progress=progress, estimated_time=estimated_time, parent_id=parent_id)
    db.add(task)
    db.commit()
    db.refresh(task)
    return task

# Calcular el progreso de una tarea
@router.get("/tasks/{task_id}/calculate")
def calculate_task(task_id: int, db: Session = Depends(get_db)):
    task = db.query(TaskModel).filter(TaskModel.id == task_id).first()
    if not task:
        raise HTTPException(status_code=404, detail="Task not found")
    
    subtasks = db.query(TaskModel).filter(TaskModel.parent_id == task_id).all()
    if not subtasks:
        return {"progress": task.progress, "estimated_time": task.estimated_time}

    # Calcular el progreso
    weighted_sum = 0.0
    total_weight = 0.0
    for subtask in subtasks:
        weighted_sum += subtask.progress * subtask.estimated_time
        total_weight += subtask.estimated_time
    
    progress = weighted_sum / total_weight if total_weight > 0 else 0
    return {"progress": progress, "estimated_time": task.estimated_time}
