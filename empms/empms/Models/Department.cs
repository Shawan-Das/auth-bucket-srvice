using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace empms.Models
{
    [Table("department")]
    public class Department
    {
            [Key]
            [Column("id")]
            public int Id { get; set; }
            [Column("department_name")]
            public string? DepartmentName { get; set; }
    }
}
