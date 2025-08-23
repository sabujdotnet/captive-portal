1. Installation:
   ```bash
   git clone <repository-url>
   cd captive-portal
   go mod download
   ```
2. Configuration:
   · Update config.yaml with your specific settings
   · Set up your database (MySQL or MikroTik local)
   · Configure MikroTik hotspot settings
   · Set up PHPNuxBill integration
3. Build and Run:
   ```bash
   go build -o portal ./cmd/portal-server
   ./portal
   ```
4. Docker Deployment:
   ```bash
   docker build -t captive-portal .
   docker run -d -p 8080:8080 --name portal captive-portal
   ```

Additional Features

This solution includes:

1. Multi-language Support: Using JSON files in the locales directory
2. User Analytics: Tracks login attempts, bandwidth usage, and session duration
3. Customizable Branding: Easily change logos, colors, and company name
4. Admin Dashboard: Monitor active users, generate reports, and manage vouchers
5. API Endpoints: For integration with other systems
6. Session Management: Secure cookie-based sessions with expiration
7. Error Handling: Comprehensive error handling and logging

Extensibility

The architecture is designed to be easily extensible:

· Add new authentication methods
· Support additional social login providers
· Integrate with other billing systems
· Add new analytics and reporting features
· Customize the user portal with additional functionality

This captive portal solution provides a solid foundation that can be customized and extended based on specific requirements.