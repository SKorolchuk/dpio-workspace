using System;
using System.Text;
using Deeproxio.Persistence.Identity.Jwt;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.IdentityModel.Tokens;

namespace Deeproxio.UserManagement.API.Extensions
{
    public static class ServiceRegistrationExtensions
    {
        public static void AddJWTAuthorization(this IServiceCollection services, IConfiguration configuration)
        {
            string secretKey = configuration.GetValue("JwtIssuerOptions:SecretKey", "");

            if (string.IsNullOrWhiteSpace(secretKey))
            {
                throw new ArgumentNullException(secretKey);
            }

            var signingKey = new SymmetricSecurityKey(Encoding.ASCII.GetBytes(secretKey));

            var jwtAppSettingOptions = configuration.GetSection(nameof(JwtIssuerOptions));

            // Configure JwtIssuerOptions
            services.Configure<JwtIssuerOptions>(options =>
            {
                options.Issuer = jwtAppSettingOptions[nameof(JwtIssuerOptions.Issuer)];
                options.Audience = jwtAppSettingOptions[nameof(JwtIssuerOptions.Audience)];
                options.SigningCredentials = new SigningCredentials(signingKey, SecurityAlgorithms.HmacSha256);
            });

            var tokenValidationParameters = new TokenValidationParameters
            {
                ValidateIssuer = false,
                ValidIssuer = jwtAppSettingOptions[nameof(JwtIssuerOptions.Issuer)],

                ValidateAudience = false,
                ValidAudience = jwtAppSettingOptions[nameof(JwtIssuerOptions.Audience)],

                ValidateIssuerSigningKey = true,
                IssuerSigningKey = signingKey,

                RequireExpirationTime = false,
                ValidateLifetime = true,
                ClockSkew = TimeSpan.Zero
            };

            services
                .AddAuthentication(options =>
                {
                    options.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
                    options.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
                })
                .AddJwtBearer(options =>
                {
                    options.ClaimsIssuer = jwtAppSettingOptions[nameof(JwtIssuerOptions.Issuer)];
                    options.TokenValidationParameters = tokenValidationParameters;
                    options.SaveToken = true;
                });

            services.AddAuthorization(options =>
            {
                options.AddPolicy("JWT Claims", policy => policy.RequireClaim(
                    Constants.Strings.JwtClaimIdentifiers.Rol,
                    Constants.Strings.JwtClaims.ApiAccess));
            });
        }
    }
}
