// Stdio launcher for Smithery MCPB bundling (runs published GHCR image).
const { spawn } = require("node:child_process");

const verbose = process.env.TINY_GO_MCP_VERBOSE || "";
const child = spawn(
  "docker",
  [
    "run",
    "-i",
    "--rm",
    "-e",
    `TINY_GO_MCP_VERBOSE=${verbose}`,
    "ghcr.io/kioie/tiny-go-mcp:1.2.0",
  ],
  { stdio: "inherit" },
);

child.on("exit", (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }
  process.exit(code ?? 1);
});

child.on("error", (err) => {
  console.error(err.message);
  process.exit(1);
});
