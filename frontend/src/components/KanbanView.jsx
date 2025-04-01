import { DndContext, closestCorners } from "@dnd-kit/core";
import {
  useSortable,
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";

const KanbanView = ({ tasks }) => {
  const groupedTasks = tasks.reduce((acc, task) => {
    (acc[task.status] = acc[task.status] || []).push(task);
    return acc;
  }, {});

  return (
    <DndContext collisionDetection={closestCorners}>
      <div className="grid grid-cols-3 gap-4">
        {Object.entries(groupedTasks).map(([status, tasks]) => (
          <div key={status} className="p-4 bg-gray-100 rounded">
            <h2 className="text-xl font-bold">{status.toUpperCase()}</h2>
            <SortableContext
              items={tasks.map((task) => task.id)}
              strategy={verticalListSortingStrategy}
            >
              {tasks.map((task) => (
                <KanbanCard key={task.id} task={task} />
              ))}
            </SortableContext>
          </div>
        ))}
      </div>
    </DndContext>
  );
};

const KanbanCard = ({ task }) => {
  const { attributes, listeners, setNodeRef, transform, transition } =
    useSortable({ id: task.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      className="p-2 bg-white shadow my-2 cursor-grab"
    >
      {task.title}
    </div>
  );
};

export default KanbanView;
