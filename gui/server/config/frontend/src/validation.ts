export function validateIPv4(ip: string): boolean {
  if (!ip) return false;
  const ipv4Regex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
  return ipv4Regex.test(ip.trim());
}

export function validateMulticastIPv4(ip: string): boolean {
  if (!validateIPv4(ip)) return false;
  const firstOctet = Number(ip.trim().split('.')[0]);
  return firstOctet >= 224 && firstOctet <= 239;
}
