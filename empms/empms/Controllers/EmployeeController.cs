using empms.Models;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using StudentAPI.Models;
using Microsoft.AspNetCore.Http;
//using Microsoft.AspNetCore.Mvc;

namespace empms.Controllers
{
	[Route("api/[controller]")]
	[ApiController]
	public class EmployeeController : ControllerBase
	{
		private readonly PostgresContext _context;

		public EmployeeController(PostgresContext context)
		{
			_context = context;
		}


		// GET: api/Employee
		[HttpGet]
		public async Task<ActionResult<IEnumerable<Employee>>> GetEmployee()
		{
			return await _context.Employees.ToListAsync();
		}

		// Get By ID

		[HttpGet("{id}")]    // Search by localhost/api/department/{id}
		public async Task<ActionResult<Employee>> Employee(int id)
		{
			try
			{
				var employee = await _context.Employees.FindAsync(id);

				if (employee == null)
				{
					return NotFound();
				}
				return employee;

			}
			catch (Exception)
			{
				return StatusCode(StatusCodes.Status500InternalServerError, "Error retreiving data from the database");
			}
		}

		// Update Data [PUT operation]

		[HttpPut("{id}")]
		public async Task<IActionResult> PutDepartment(int id, Employee employee)
		{
			if (id != employee.Id)
			{
				return BadRequest();
			}
			if (!EmployeeExists(id))
			{
				return NotFound();
			}

			_context.Entry(employee).State = EntityState.Modified;

			try
			{
				await _context.SaveChangesAsync();
				return Ok("Employee Data Update Successful");
			}
			catch (Exception)
			{
				return StatusCode(StatusCodes.Status500InternalServerError, "Error Updating Employee Data");
			}
		}

		// Create Data [post operation]

		[HttpPost]
		public async Task<ActionResult<Department>> CreateDepartment(Employee employee)
		{
			try
			{
				if (employee == null)
				{
					return BadRequest();
				}
				_context.Employees.Add(employee);
				await _context.SaveChangesAsync();

				return CreatedAtAction(nameof(GetEmployee), new { id = employee.Id }, employee);

			}
			catch (Exception)
			{
				return StatusCode(StatusCodes.Status500InternalServerError, "Error Creating Employee");
			}

		}




		[HttpDelete("{id}")]
		public async Task<IActionResult> DeleteEmployee(int id)
		{
			try
			{
				var deleteEmployee= await _context.Employees.FindAsync(id);
				if (deleteEmployee== null)
				{
					return NotFound($"No Data with id:{id} found");
				}
				else
				{
					_context.Employees.Remove(deleteEmployee);
					await _context.SaveChangesAsync();
					return Ok($"Data delete with if:{id} successful");
				}
			}
			catch (Exception)
			{
				return StatusCode(StatusCodes.Status500InternalServerError, "Error with Deleting Data");
			}
		}

		private bool EmployeeExists(int id)
		{
			return _context.Employees.Any(e => e.Id == id);
		}
	}
}
