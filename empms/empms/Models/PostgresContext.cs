using System;
using System.Collections.Generic;
using empms.Models;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata;

namespace StudentAPI.Models
{
    public partial class PostgresContext : DbContext
    {
        public PostgresContext()
        {
        }
        public PostgresContext(DbContextOptions<PostgresContext> options)
           : base(options)
        {
        }

        public virtual DbSet<Department> Departments { get; set; } = null!;

        protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
        {
            if (!optionsBuilder.IsConfigured)
            {
                optionsBuilder.UseNpgsql("Host=localhost;Database=empms;Username=postgres;Password=postgres");
            }
        }

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            modelBuilder.HasPostgresExtension("adminpack")
                .HasAnnotation("Relational:Collation", "English_United States.1252");
        }
        partial void OnModelCreatingPartial(ModelBuilder modelBuilder);
    }

}
