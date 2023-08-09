using empms.Models;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using StudentAPI.Models;

namespace empms.Controllers
{
    [Route("api/[controller]")]  
    [ApiController]
    public class DepartmentController : ControllerBase
    {
        private readonly PostgresContext _context;

		public DepartmentController(PostgresContext context)
        {
            _context = context;
        }

        // GET: api/DepartmentDetails
        [HttpGet]             // Search by localhost/api/department
		public async Task<ActionResult<IEnumerable<Department>>> GetDepartment()
        {
            return await _context.Departments.ToListAsync();
        }

        // Get By ID

        [HttpGet("{id}")]    // Search by localhost/api/department/{id}
		public async Task<ActionResult<Department>> Department(int id)
        {
            try
            {
				var department = await _context.Departments.FindAsync(id);

				if (department == null)
				{
					return NotFound();
				}
				return department;

			}
            catch (Exception)
            {
                return StatusCode(StatusCodes.Status500InternalServerError, "Error retreiving data from the database");
            }
        }

        // Update Data [PUT operation]

        [HttpPut("{id}")]
		public async Task<IActionResult> PutDepartment(int id, Department department)
		{
			if (id != department.Id)
			{
				return BadRequest();
			}
			if (!DepartmentExists(id))
			{
				return NotFound();
			}

			_context.Entry(department).State = EntityState.Modified;

			try
			{
				await _context.SaveChangesAsync();
                return Ok("Data Update Successful");
			}
			catch (Exception)
			{
				return StatusCode(StatusCodes.Status500InternalServerError, "Error Updating Data");
			}
		}

        // Create Data [post operation]

		[HttpPost]
        public async Task<ActionResult<Department>> CreateDepartment(Department department)
        {
            try
            {
                if(department == null)
                {
                    return BadRequest();
                }
				_context.Departments.Add(department);
				await _context.SaveChangesAsync();

				return CreatedAtAction(nameof(GetDepartment), new { id = department.Id }, department);

			}
            catch (Exception)
            {
                return StatusCode(StatusCodes.Status500InternalServerError, "Error Creating Data");
            }
            
        }



        //Delete Data [DELETE operation]
        [HttpDelete("{id}")]
        public async Task<IActionResult> DeleteDepartment(int id)
        {
            try
            {
				var deleteDepartment = await _context.Departments.FindAsync(id);
				if (deleteDepartment == null)
				{
					return NotFound($"No Data with id:{id} found");
				}
                else
                {
					_context.Departments.Remove(deleteDepartment);
					await _context.SaveChangesAsync();

/*					try
					{
						var deleteUser = await _context.Employees.FirstOrDefaultAsync(e => e.departmentId == id);
						_context.Employees.Remove(deleteUser);
						await _context.SaveChangesAsync();

						if (deleteUser != null)
						{
							_context.Employees.Remove(deleteUser);
							await _context.SaveChangesAsync();
							return Ok($"User: {deleteUser.employee_name} Delete Successful");
						}
					}
					catch
					{
					}*/
					//return Ok("User Delete Successful");

					return Ok($"Data delete :{deleteDepartment.Id} successful");
				}
			}
            catch(Exception)
            {
				return StatusCode(StatusCodes.Status500InternalServerError, "Error with Deleting Data");
			}
        }

        // Check for Data
        private bool DepartmentExists(int id)
        {
            return _context.Departments.Any(e => e.Id == id);
        }
    }
}
