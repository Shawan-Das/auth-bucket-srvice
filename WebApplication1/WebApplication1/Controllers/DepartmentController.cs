using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using WebApplication1.Models;
using Microsoft.EntityFrameworkCore;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace WebApplication1.Controllers
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

		//Get method
		[HttpGet]
		public async Task<ActionResult<IEnumerable<Department>>> GetDepartmentItems()
		{
			return await _context.departments.ToListAsync();
		}
	}
}
