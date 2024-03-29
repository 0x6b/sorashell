package sorashell

import "github.com/c-bata/go-prompt"

// Commands define top-level commands from API definition
var Commands = []prompt.Suggest{
	{Text: "audit-logs", Description: "Retrieve audit logs."},
	{Text: "bills", Description: "Show or export billing info."},
	{Text: "cell-locations", Description: "Retrieves cell tower location information."},
	{Text: "coupons", Description: "List or register coupons."},
	{Text: "credentials", Description: "List, create, update or delete credentials sets."},
	{Text: "data", Description: "Get stored data from subscribers."},
	{Text: "devices", Description: "Manage devices."},
	{Text: "diagnostics", Description: "Do diagnostics and get the reports."},
	{Text: "emails", Description: "Manage email addresses."},
	{Text: "event-handlers", Description: "List, create, update or delete event handlers."},
	{Text: "files", Description: "Manage files on Harvest Files."},
	{Text: "gadgets", Description: "Manage gadgets."},
	{Text: "groups", Description: "List, create, update or delete groups."},
	{Text: "help", Description: "Help about any command"},
	{Text: "lagoon", Description: "Manage Lagoon settings."},
	{Text: "logout", Description: "Revoke API key and API token to access to the SORACOM API."},
	{Text: "logs", Description: "List logs."},
	{Text: "lora-devices", Description: "Manage LoRa devices."},
	{Text: "lora-gateways", Description: "Manage LoRa gateways."},
	{Text: "lora-network-sets", Description: "Manage LoRa network sets."},
	{Text: "operator", Description: "Manage operators."},
	{Text: "orders", Description: "List, create or cancel orders."},
	{Text: "payer-information", Description: "Get or edit payer information."},
	{Text: "payment-history", Description: "List payment history."},
	{Text: "payment-methods", Description: "Create or update payment methods."},
	{Text: "payment-statements", Description: "List or export payment statements."},
	{Text: "port-mappings", Description: "Manage port mappings for on-demand remote access."},
	{Text: "products", Description: "List products."},
	{Text: "query", Description: "Search resources such as subscribers or sigfox devices."},
	{Text: "resource-summaries", Description: "Get resource summary."},
	{Text: "roles", Description: "List, create, update or delete roles."},
	{Text: "sandbox", Description: "Sandbox related operations."},
	{Text: "self-update", Description: "Updates soracom-cli to the latest version."},
	{Text: "shipping-addresses", Description: "List, create, update or delete shipping addresses."},
	{Text: "sigfox-devices", Description: "Manage Sigfox devices."},
	{Text: "sims", Description: "Manage SIMs."},
	{Text: "sora-cam", Description: "Manage Soracom Cloud Camera Services (SoraCam) devices and licenses."},
	{Text: "soralets", Description: "Manage Soralets for Orbit."},
	{Text: "stats", Description: "Show or export statistics."},
	{Text: "subscribers", Description: "Manage subscribers."},
	{Text: "system-notifications", Description: "Manage system notifications."},
	{Text: "unconfigure", Description: "Remove configurations."},
	{Text: "users", Description: "Manage SAM users."},
	{Text: "version", Description: "Show version info."},
	{Text: "volume-discounts", Description: "Manage volume discounts (long-term discounts)."},
	{Text: "vpg", Description: "List, create, update or delete VPGs."},
}
