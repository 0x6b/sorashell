package shell

import "github.com/c-bata/go-prompt"

var Commands = []prompt.Suggest{
	{Text: "audit-logs", Description: "Retrieve audit logs."},
	{Text: "bills", Description: "Show or export billing info."},
	{Text: "coupons", Description: "List or register coupons."},
	{Text: "credentials", Description: "List, create, update or delete credentials sets."},
	{Text: "data", Description: "Get stored data from subscribers."},
	{Text: "devices", Description: "Manage devices."},
	{Text: "event-handlers", Description: "List, create, update or delete event handlers."},
	{Text: "files", Description: "Manage files on Harvest Files."},
	{Text: "gadgets", Description: "Manage gadgets."},
	{Text: "groups", Description: "List, create, update or delete groups."},
	{Text: "lagoon", Description: "Manage Lagoon settings."},
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
	{Text: "roles", Description: "List, create, update or delete roles."},
	{Text: "sandbox", Description: "Sandbox related operations."},
	{Text: "shipping-addresses", Description: "List, create, update or delete shipping addresses."},
	{Text: "sigfox-devices", Description: "Manage Sigfox devices."},
	{Text: "stats", Description: "Show or export statistics."},
	{Text: "subscribers", Description: "Manage subscribers."},
	{Text: "test", Description: "Do diagnostics & testings."},
	{Text: "users", Description: "Manage SAM users."},
	{Text: "version", Description: "Show version info."},
	{Text: "volume-discounts", Description: "Manage volume discounts (long-term discounts)."},
	{Text: "vpg", Description: "List, create, update or delete VPGs."},
}