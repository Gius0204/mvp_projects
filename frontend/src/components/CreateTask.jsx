import { useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

const CreateTask = () => {
  const [title, setTitle] = useState("");
  const [type, setType] = useState("list");
  const navigate = useNavigate();

  const handleAddTask = (e) => {
    e.preventDefault();

    if (!title.trim()) {
      alert("El título no puede estar vacío");
      return;
    }

    const newTask = { title, type_view: type };

    axios
      .post("http://localhost:3000/tasks", newTask)
      .then(() => {
        navigate("/"); // Volver a la página principal
      })
      .catch((error) => console.error("❌ Error al agregar tarea:", error));
  };

  return (
    <div className="p-6">
      <h2 className="text-lg font-bold mb-4">Crear Nueva Tarea</h2>
      <form onSubmit={handleAddTask} className="p-4 bg-gray-200 rounded">
        <input
          type="text"
          placeholder="Título de la tarea"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className="p-2 border rounded w-full mb-2"
        />
        <select
          value={type}
          onChange={(e) => setType(e.target.value)}
          className="p-2 border rounded w-full mb-2"
        >
          <option value="list">Lista</option>
          <option value="kanban">Kanban</option>
          <option value="gantt">Gantt</option>
        </select>
        <button
          type="submit"
          className="p-2 bg-blue-500 text-white rounded w-full"
        >
          Guardar Tarea
        </button>
      </form>

      <button
        onClick={() => navigate("/")}
        className="mt-4 p-2 bg-gray-500 text-white rounded"
      >
        🔙 Volver
      </button>
    </div>
  );
};

export default CreateTask;
