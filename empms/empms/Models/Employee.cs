using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;
namespace empms.Models
{
	[Table("employee")]
	public class Employee
	{
		[Key]
		[Column("id")]
		public int Id { get; set; }
		[Column("employee_name")]
		public string? employee_name { get; set; }
		[Column("email")]
		public string? email { get; set; }
		[Column("phone")]
		public string? phone { get; set; }
		[Column("address")]
		public string? address { get; set; }
		[Column("departmentid")]
		public int departmentId { get; set; }
	}
}
