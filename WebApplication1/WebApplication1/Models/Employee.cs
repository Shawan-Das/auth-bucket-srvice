using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations.Schema;
using System.Linq;
using System.Threading.Tasks;
namespace WebApplication1.Models
{
	[Table("employee")]
	public class Employee
	{
		[Column("id")]
		public int Id { get; set; }
		[Column("Employee Name")]
		public string? employeeName { get; set; }
		[Column("Department")]
		public string? deparment { get; set; }
		[Column("Date of Joining")]
		public string? dateOfJoining { get; set; }
	}
}
