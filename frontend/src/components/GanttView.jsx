import { useEffect, useRef } from "react";
import Gantt from "frappe-gantt";

const GanttView = ({ tasks }) => {
  const ganttRef = useRef(null);

  useEffect(() => {
    if (ganttRef.current) {
      new Gantt(
        ganttRef.current,
        tasks.map((task) => ({
          id: task.id,
          name: task.title,
          start: task.start_date || new Date(),
          end: task.due_date || new Date(),
          progress: 50,
        })),
        {
          view_mode: "Day",
          language: "es",
        }
      );
    }
  }, [tasks]);

  return <div ref={ganttRef} />;
};

export default GanttView;
