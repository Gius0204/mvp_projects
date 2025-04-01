import { BrowserRouter, Routes, Route } from "react-router-dom";
import TaskBoard from "./TaskBoard";
import CreateTask from "./CreateTask";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<TaskBoard />} />
        <Route path="/create-task" element={<CreateTask />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
