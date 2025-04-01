import { useEffect, useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

import KanbanView from "./KanbanView";
import ListView from "./ListView";
import GanttView from "./GanttView";

const TaskBoard = () => {
  const [tasks, setTasks] = useState([]);
  const [view, setView] = useState("list"); // Estado para la vista activa
  const navigate = useNavigate(); // Para cambiar de página

  useEffect(() => {
    axios
      .get("http://localhost:3000/tasks")
      .then((res) => setTasks(res.data))
      .catch((error) => console.error("❌ Error al obtener tareas:", error));

    // Conectar WebSocket
    const socket = new WebSocket("ws://localhost:3000/ws");

    socket.onmessage = (event) => {
      const newTask = JSON.parse(event.data);
      setTasks((prevTasks) => [...prevTasks, newTask]); // Agregar la nueva tarea en tiempo real
    };

    return () => socket.close();
  }, []);

  return (
    <div className="p-6">
      <button
        onClick={() => navigate("/create-task")}
        className="mb-4 p-2 bg-blue-500 text-white rounded"
      >
        ➕ Crear Nueva Tarea
      </button>

      <div className="grid grid-cols-3 gap-4">
        {["list", "kanban", "gantt"].map((type) => (
          <div key={type} className="p-4 bg-gray-100 rounded">
            <h2 className="text-xl font-bold">{type.toUpperCase()}</h2>
            {tasks
              .filter((t) => t.type_view === type)
              .map((task) => (
                <div key={task.id_task} className="p-2 bg-white shadow my-2">
                  {task.title}
                </div>
              ))}
          </div>
        ))}
      </div>

      {/* Selector de vistas */}
      {/* <div className="flex space-x-4 mb-4">
        {["list", "kanban", "gantt"].map((type) => (
          <button
            key={type}
            className={`p-2 rounded ${
              view === type ? "bg-blue-500 text-white" : "bg-gray-200"
            }`}
            onClick={() => setView(type)}
          >
            {type.toUpperCase()}
          </button>
        ))}
      </div> */}

      {/* Mostrar la vista activa */}
      {/* {view === "list" && <ListView tasks={tasks} />}
      {view === "kanban" && <KanbanView tasks={tasks} />}
      {view === "gantt" && <GanttView tasks={tasks} />} */}
    </div>
  );
};

export default TaskBoard;
