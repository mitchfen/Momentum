namespace Momentum
{
    public class AppConfig
    {
        public List<string> DailyTasks { get; set; } = new List<string>();
        public string TimeZone { get; set; } = "UTC";

        public DateTime GetToday()
        {
            try
            {
                var tz = TimeZoneInfo.FindSystemTimeZoneById(TimeZone);
                return TimeZoneInfo.ConvertTimeFromUtc(DateTime.UtcNow, tz).Date;
            }
            catch (TimeZoneNotFoundException)
            {
                // Fallback to UTC if timezone is invalid
                return DateTime.UtcNow.Date;
            }
        }
    }
}
