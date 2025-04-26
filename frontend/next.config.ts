import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  output: 'standalone', // needed for docker
  rewrites: async () => {
    return [    {
    source: '/api/:path*',
    destination: 'http://localhost:8080/:path*' // Proxy to Backend
  }]}
};

export default nextConfig;
