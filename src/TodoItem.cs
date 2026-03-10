namespace Momentum 
{
    public class TodoItem
    {
        public string Text { get; set; } = string.Empty;
        public string DisplayTitle { get; set; } = string.Empty;
        public List<SubTask> SubTasks { get; set; } = new();
        public bool IsCompleted => SubTasks.All(s => s.IsCompleted);
    }

    public class SubTask
    {
        public string Text { get; set; } = string.Empty;
        public bool IsCompleted { get; set; }
    }
}
