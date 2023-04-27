/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    appDir: true,
  },
  images: {
    unoptimized: true,
  },
  trailingSlash: true,
  output: "export",
}

module.exports = nextConfig
