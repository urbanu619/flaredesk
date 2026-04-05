#!/usr/bin/env node
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";

const CF_API = "https://api.cloudflare.com/client/v4";

function getToken(): string {
  const token = process.env.CF_API_TOKEN;
  if (!token) throw new Error("CF_API_TOKEN environment variable is required");
  return token;
}

async function cfFetch(path: string, options: RequestInit = {}): Promise<any> {
  const res = await fetch(`${CF_API}${path}`, {
    ...options,
    headers: {
      "Authorization": `Bearer ${getToken()}`,
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
  });
  return res.json();
}

// 分页获取所有 zones
async function fetchAllZones(accountId: string): Promise<any[]> {
  let page = 1;
  const all: any[] = [];
  while (true) {
    const data = await cfFetch(`/zones?account.id=${accountId}&per_page=50&page=${page}`);
    const result = data.result || [];
    all.push(...result);
    if (page >= (data.result_info?.total_pages || 1)) break;
    page++;
  }
  return all;
}

const server = new McpServer({
  name: "flaredesk",
  version: "1.0.0",
  description: "Cloudflare DNS batch management — multi-account, bulk operations, proxy toggle",
});

// ── 1. 列出所有 Zone（域名）──────────────────────────────────────────
server.tool(
  "list_zones",
  "List all Cloudflare zones (domains) for an account",
  {
    account_id: z.string().describe("Cloudflare Account ID"),
  },
  async ({ account_id }) => {
    const zones = await fetchAllZones(account_id);
    const list = zones.map((z: any) => ({
      id: z.id,
      name: z.name,
      status: z.status,
      plan: z.plan?.name,
    }));
    return {
      content: [{ type: "text", text: JSON.stringify(list, null, 2) }],
    };
  }
);

// ── 2. 列出 DNS 记录 ─────────────────────────────────────────────────
server.tool(
  "list_dns_records",
  "List all DNS records for a zone",
  {
    account_id: z.string().describe("Cloudflare Account ID"),
    zone_id: z.string().describe("Zone ID"),
    type: z.string().optional().describe("Filter by record type (A, CNAME, TXT, etc.)"),
    name: z.string().optional().describe("Filter by record name"),
  },
  async ({ account_id, zone_id, type, name }) => {
    let url = `/accounts/${account_id}/zones/${zone_id}/dns_records?per_page=100`;
    if (type) url += `&type=${type}`;
    if (name) url += `&name=${name}`;
    const data = await cfFetch(url);
    return {
      content: [{ type: "text", text: JSON.stringify(data.result || [], null, 2) }],
    };
  }
);

// ── 3. 新增 DNS 记录 ─────────────────────────────────────────────────
server.tool(
  "create_dns_record",
  "Create a new DNS record",
  {
    account_id: z.string(),
    zone_id: z.string(),
    type: z.enum(["A", "AAAA", "CNAME", "TXT", "MX", "NS", "SRV", "CAA"]),
    name: z.string().describe("Record name, e.g. @ or www"),
    content: z.string().describe("Record content, e.g. 1.2.3.4"),
    ttl: z.number().default(1).describe("TTL in seconds, 1 = auto"),
    proxied: z.boolean().default(false).describe("Enable Cloudflare proxy (orange cloud)"),
  },
  async ({ account_id, zone_id, type, name, content, ttl, proxied }) => {
    const data = await cfFetch(`/accounts/${account_id}/zones/${zone_id}/dns_records`, {
      method: "POST",
      body: JSON.stringify({ type, name, content, ttl, proxied }),
    });
    return {
      content: [{ type: "text", text: data.success ? `Created: ${data.result?.id}` : JSON.stringify(data.errors) }],
    };
  }
);

// ── 4. 更新 DNS 记录 ─────────────────────────────────────────────────
server.tool(
  "update_dns_record",
  "Update an existing DNS record",
  {
    account_id: z.string(),
    zone_id: z.string(),
    record_id: z.string(),
    type: z.enum(["A", "AAAA", "CNAME", "TXT", "MX", "NS", "SRV", "CAA"]).optional(),
    name: z.string().optional(),
    content: z.string().optional(),
    ttl: z.number().optional(),
    proxied: z.boolean().optional(),
  },
  async ({ account_id, zone_id, record_id, ...fields }) => {
    const data = await cfFetch(`/accounts/${account_id}/zones/${zone_id}/dns_records/${record_id}`, {
      method: "PATCH",
      body: JSON.stringify(fields),
    });
    return {
      content: [{ type: "text", text: data.success ? "Updated" : JSON.stringify(data.errors) }],
    };
  }
);

// ── 5. 删除 DNS 记录 ─────────────────────────────────────────────────
server.tool(
  "delete_dns_record",
  "Delete a DNS record",
  {
    account_id: z.string(),
    zone_id: z.string(),
    record_id: z.string(),
  },
  async ({ account_id, zone_id, record_id }) => {
    const data = await cfFetch(`/accounts/${account_id}/zones/${zone_id}/dns_records/${record_id}`, {
      method: "DELETE",
    });
    return {
      content: [{ type: "text", text: data.success ? "Deleted" : JSON.stringify(data.errors) }],
    };
  }
);

// ── 6. 批量新增 DNS 记录（单域名）───────────────────────────────────
server.tool(
  "batch_create_dns_records",
  "Batch create multiple DNS records for a single zone",
  {
    account_id: z.string(),
    zone_id: z.string(),
    records: z.array(z.object({
      type: z.string(),
      name: z.string(),
      content: z.string(),
      ttl: z.number().default(1),
      proxied: z.boolean().default(false),
    })).describe("List of DNS records to create"),
  },
  async ({ account_id, zone_id, records }) => {
    let success = 0, fail = 0;
    for (const rec of records) {
      const data = await cfFetch(`/accounts/${account_id}/zones/${zone_id}/dns_records`, {
        method: "POST",
        body: JSON.stringify(rec),
      });
      data.success ? success++ : fail++;
    }
    return {
      content: [{ type: "text", text: `Done: ${success} created, ${fail} failed` }],
    };
  }
);

// ── 7. 跨域名批量新增 ────────────────────────────────────────────────
server.tool(
  "cross_zone_create_dns_records",
  "Create the same DNS records across multiple zones",
  {
    account_id: z.string(),
    zone_ids: z.array(z.string()).describe("List of zone IDs"),
    records: z.array(z.object({
      type: z.string(),
      name: z.string(),
      content: z.string(),
      ttl: z.number().default(1),
      proxied: z.boolean().default(false),
    })),
  },
  async ({ account_id, zone_ids, records }) => {
    let success = 0, fail = 0;
    for (const zone_id of zone_ids) {
      for (const rec of records) {
        const data = await cfFetch(`/accounts/${account_id}/zones/${zone_id}/dns_records`, {
          method: "POST",
          body: JSON.stringify(rec),
        });
        data.success ? success++ : fail++;
      }
    }
    return {
      content: [{ type: "text", text: `Done across ${zone_ids.length} zones: ${success} created, ${fail} failed` }],
    };
  }
);

// ── 8. 跨域名批量删除 ────────────────────────────────────────────────
server.tool(
  "cross_zone_delete_records",
  "Delete DNS records matching type/name across multiple zones",
  {
    account_id: z.string(),
    zone_ids: z.array(z.string()).describe("List of zone IDs, or use 'all' to fetch all zones"),
    type: z.string().optional().describe("Record type to match (A, CNAME, etc.)"),
    name: z.string().optional().describe("Record name to match"),
  },
  async ({ account_id, zone_ids, type, name }) => {
    const zones = zone_ids.includes("all") ? (await fetchAllZones(account_id)).map((z: any) => z.id) : zone_ids;
    let deleted = 0, fail = 0;
    for (const zone_id of zones) {
      let url = `/accounts/${account_id}/zones/${zone_id}/dns_records?per_page=100`;
      if (type) url += `&type=${type}`;
      if (name) url += `&name=${name}`;
      const data = await cfFetch(url);
      for (const rec of (data.result || [])) {
        const del = await cfFetch(`/accounts/${account_id}/zones/${zone_id}/dns_records/${rec.id}`, { method: "DELETE" });
        del.success ? deleted++ : fail++;
      }
    }
    return {
      content: [{ type: "text", text: `Done: ${deleted} deleted, ${fail} failed across ${zones.length} zones` }],
    };
  }
);

// ── 9. 跨域名批量切换橙云代理 ───────────────────────────────────────
server.tool(
  "cross_zone_toggle_proxy",
  "Enable or disable Cloudflare proxy (orange cloud) for matching records across multiple zones",
  {
    account_id: z.string(),
    zone_ids: z.array(z.string()).describe("List of zone IDs, or use 'all' to fetch all zones"),
    proxied: z.boolean().describe("true = enable proxy, false = disable"),
    type: z.string().optional().describe("Only match records of this type (A, AAAA, CNAME)"),
    name: z.string().optional().describe("Only match records with this name"),
  },
  async ({ account_id, zone_ids, proxied, type, name }) => {
    const zones = zone_ids.includes("all") ? (await fetchAllZones(account_id)).map((z: any) => z.id) : zone_ids;
    let updated = 0, fail = 0;
    for (const zone_id of zones) {
      let url = `/accounts/${account_id}/zones/${zone_id}/dns_records?per_page=100`;
      if (type) url += `&type=${type}`;
      if (name) url += `&name=${name}`;
      const data = await cfFetch(url);
      for (const rec of (data.result || [])) {
        if (!["A", "AAAA", "CNAME"].includes(rec.type)) continue;
        const patch = await cfFetch(`/accounts/${account_id}/zones/${zone_id}/dns_records/${rec.id}`, {
          method: "PATCH",
          body: JSON.stringify({ proxied }),
        });
        patch.success ? updated++ : fail++;
      }
    }
    return {
      content: [{ type: "text", text: `Done: ${updated} updated, ${fail} failed across ${zones.length} zones` }],
    };
  }
);

// ── 启动 ─────────────────────────────────────────────────────────────
const transport = new StdioServerTransport();
server.connect(transport).catch(console.error);
