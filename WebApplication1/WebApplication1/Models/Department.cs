using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations.Schema;
using System.Linq;
using System.Threading.Tasks;
namespace WebApplication1.Models
{
	[Table("department")]
	public class Department
	{
		[Column("id")]
		public int Id { get; set; }
		[Column("departmentname")]
		public string? departmentName { get; set; }
	}
}
