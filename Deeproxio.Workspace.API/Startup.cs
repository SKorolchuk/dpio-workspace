using System;
using System.IO;
using AutoMapper;
using Deeproxio.UserManagement.API.Extensions;
using FluentValidation.AspNetCore;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Diagnostics.HealthChecks;
using Prometheus;

namespace Deeproxio.Workspace.API
{
    public class Startup
    {
        public IConfiguration Configuration { get; }

        public Startup(IWebHostEnvironment env)
        {
            var builder = new ConfigurationBuilder()
                .SetBasePath(Directory.GetCurrentDirectory())
                .AddJsonFile("appsettings.json", optional: true, reloadOnChange: true)
                .AddJsonFile($"appsettings.{env.EnvironmentName}.json", optional: true)
                .AddJsonFile($"secrets/appsettings.json", optional: true)
                .AddEnvironmentVariables();

            Configuration = builder.Build();
        }

        public void ConfigureServices(IServiceCollection services)
        {
            services
                .AddAutoMapper(typeof(Startup).Assembly)
                .AddCors()
                .AddControllers()
                .SetCompatibilityVersion(CompatibilityVersion.Version_3_0)
                .AddFluentValidation(fluentValidator => fluentValidator.RegisterValidatorsFromAssemblyContaining<Startup>());

            services.AddJWTAuthorization(Configuration);

            services.AddHealthChecks();
            //.AddDbContextCheck<PlatformIdentityDbContext>(nameof(PlatformIdentityDbContext), HealthStatus.Unhealthy, tags: new[] { "DB" })
        }

        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            app.UseRouting();

            app.UseCors(builder => builder
                .AllowAnyOrigin()
                .AllowAnyMethod()
                .AllowAnyHeader());

            app.UseMetricServer();

            if (env.EnvironmentName.Equals("Development", StringComparison.InvariantCultureIgnoreCase))
            {
                app.UseDeveloperExceptionPage();
                app.UseDatabaseErrorPage();
            }

            app.UseAuthentication();
            app.UseAuthorization();

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllers();

                endpoints.MapHealthChecks("/ready", new HealthCheckOptions()
                {
                    Predicate = (check) => !check.Tags.Contains("DB"),
                    // The following StatusCodes are the default assignments for
                    // the HealthCheckStatus properties.
                    ResultStatusCodes =
                    {
                        [HealthStatus.Healthy] = StatusCodes.Status200OK,
                        [HealthStatus.Degraded] = StatusCodes.Status200OK,
                        [HealthStatus.Unhealthy] = StatusCodes.Status503ServiceUnavailable
                    },
                    // The default value is false.
                    AllowCachingResponses = false
                }).WithDisplayName("Workspace API Health Check");

                endpoints.MapHealthChecks("/ready/business", new HealthCheckOptions()
                {
                    // The following StatusCodes are the default assignments for
                    // the HealthCheckStatus properties.
                    ResultStatusCodes =
                    {
                        [HealthStatus.Healthy] = StatusCodes.Status200OK,
                        [HealthStatus.Degraded] = StatusCodes.Status503ServiceUnavailable,
                        [HealthStatus.Unhealthy] = StatusCodes.Status503ServiceUnavailable
                    },
                    // The default value is false.
                    AllowCachingResponses = false
                }).WithDisplayName("Workspace API Business Health Check");
            });
        }
    }
}
