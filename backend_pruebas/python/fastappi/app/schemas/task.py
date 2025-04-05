# app/schemas/task.py
from pydantic import BaseModel
from datetime import datetime
from typing import Optional

class TaskBase(BaseModel):
    progress: float
    estimated_time: int
    created_at: datetime
    updated_at: datetime
    parent_id: Optional[int] = None

class Task(TaskBase):
    id: int

    class Config:
        orm_mode = True
