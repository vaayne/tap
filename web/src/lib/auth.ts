import { env } from "cloudflare:workers"

/**
 * Verifies the X-Tap-Secret header against TAP_SCRIPTS_SECRET env var.
 * Returns true if the secret is valid, false otherwise.
 */
export function verifySecret(request: Request): boolean {
  const secret = (env as unknown as Record<string, string | undefined>).TAP_SCRIPTS_SECRET
  if (!secret) {
    console.error("TAP_SCRIPTS_SECRET is not configured")
    return false
  }
  const header = request.headers.get("X-Tap-Secret")
  return header === secret
}
